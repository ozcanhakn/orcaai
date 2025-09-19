package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type LoadBalancerConfig struct {
	Strategy          string        // "round-robin", "least-connections", "weighted-response"
	MaxRetries        int
	RetryDelay        time.Duration
	HealthCheckInterval time.Duration
	CircuitBreakerThreshold int
}

type LoadBalancer struct {
	config     *LoadBalancerConfig
	providers  []*Provider
	status     map[string]*ProviderStatus
	metrics    *MetricsCollector
	mu         sync.RWMutex
}

func NewLoadBalancer() *LoadBalancer {
	lb := &LoadBalancer{
		config: &LoadBalancerConfig{
			Strategy:          "weighted-response",
			MaxRetries:        3,
			RetryDelay:        time.Second * 2,
			HealthCheckInterval: time.Minute,
			CircuitBreakerThreshold: 5,
		},
		status:  make(map[string]*ProviderStatus),
		metrics: NewMetricsCollector(),
	}

	// Start health checker
	go lb.healthChecker()

	return lb
}

func (lb *LoadBalancer) SelectProvider(ctx context.Context, task *TaskProfile) (*Provider, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	switch lb.config.Strategy {
	case "round-robin":
		return lb.roundRobin()
	case "least-connections":
		return lb.leastConnections()
	case "weighted-response":
		return lb.weightedResponse(task)
	default:
		return nil, fmt.Errorf("unknown load balancing strategy: %s", lb.config.Strategy)
	}
}

func (lb *LoadBalancer) ExecuteWithRetry(ctx context.Context, provider *Provider, task func() error) error {
	var lastErr error
	
	for i := 0; i < lb.config.MaxRetries; i++ {
		// Check circuit breaker
		if !lb.isHealthy(provider.Name) {
			continue
		}

		err := task()
		if err == nil {
			// Success - update metrics
			lb.metrics.RecordSuccess(provider.Name)
			return nil
		}

		lastErr = err
		lb.metrics.RecordError(provider.Name, err)

		// Check if we should break the circuit
		if lb.shouldTripCircuitBreaker(provider.Name) {
			lb.markUnhealthy(provider.Name)
			continue
		}

		// Wait before retry
		time.Sleep(lb.config.RetryDelay)
	}

	return fmt.Errorf("all retries failed: %v", lastErr)
}

func (lb *LoadBalancer) healthChecker() {
	ticker := time.NewTicker(lb.config.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		lb.checkHealth()
	}
}

func (lb *LoadBalancer) checkHealth() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, provider := range lb.providers {
		status := lb.status[provider.Name]
		if status == nil {
			status = &ProviderStatus{
				Name:      provider.Name,
				Healthy:   true,
				LastCheck: time.Now(),
			}
			lb.status[provider.Name] = status
		}

		// Perform health check
		err := lb.pingProvider(provider)
		if err != nil {
			status.Healthy = false
			status.ErrorCount++
			status.LastError = err.Error()
		} else {
			status.Healthy = true
			status.ErrorCount = 0
			status.LastError = ""
		}

		status.LastCheck = time.Now()
	}
}

func (lb *LoadBalancer) roundRobin() (*Provider, error) {
	// Implement round-robin selection
	for i := 0; i < len(lb.providers); i++ {
		provider := lb.providers[i]
		if lb.isHealthy(provider.Name) {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("no healthy providers available")
}

func (lb *LoadBalancer) leastConnections() (*Provider, error) {
	var selected *Provider
	minConn := -1

	for _, provider := range lb.providers {
		if !lb.isHealthy(provider.Name) {
			continue
		}

		metrics := lb.metrics.GetProviderMetrics(provider.Name)
		if minConn == -1 || metrics.ActiveConnections < minConn {
			selected = provider
			minConn = metrics.ActiveConnections
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no healthy providers available")
	}

	return selected, nil
}

func (lb *LoadBalancer) weightedResponse(task *TaskProfile) (*Provider, error) {
	var selected *Provider
	var bestScore float64 = -1

	for _, provider := range lb.providers {
		if !lb.isHealthy(provider.Name) {
			continue
		}

		metrics := lb.metrics.GetProviderMetrics(provider.Name)
		
		// Calculate weighted score based on response time, error rate, and load
		score := calculateProviderScore(metrics, task)
		
		if bestScore == -1 || score > bestScore {
			selected = provider
			bestScore = score
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no healthy providers available")
	}

	return selected, nil
}

func (lb *LoadBalancer) isHealthy(providerName string) bool {
	status := lb.status[providerName]
	return status != nil && status.Healthy
}

func (lb *LoadBalancer) markUnhealthy(providerName string) {
	status := lb.status[providerName]
	if status != nil {
		status.Healthy = false
		status.LastCheck = time.Now()
	}
}

func (lb *LoadBalancer) shouldTripCircuitBreaker(providerName string) bool {
	status := lb.status[providerName]
	return status != nil && status.ErrorCount >= lb.config.CircuitBreakerThreshold
}

func (lb *LoadBalancer) pingProvider(provider *Provider) error {
	// Implement provider health check
	// This could be a simple API call to check if the provider is responding
	return nil
}

func calculateProviderScore(metrics *ProviderMetrics, task *TaskProfile) float64 {
	// Weighted scoring based on:
	// - Response time (lower is better)
	// - Error rate (lower is better)
	// - Current load (lower is better)
	// - Success rate (higher is better)
	
	responseTimeWeight := 0.3
	errorRateWeight := 0.25
	loadWeight := 0.2
	successRateWeight := 0.25

	responseTimeScore := 1.0 / (float64(metrics.AverageLatency.Milliseconds()) / 1000.0)
	errorRateScore := 1.0 - metrics.ErrorRate
	loadScore := 1.0 - (float64(metrics.ActiveConnections) / float64(metrics.MaxConnections))
	successRateScore := metrics.SuccessRate

	return (responseTimeScore * responseTimeWeight) +
		   (errorRateScore * errorRateWeight) +
		   (loadScore * loadWeight) +
		   (successRateScore * successRateWeight)
}
