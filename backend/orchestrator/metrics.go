package orchestrator

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics for the orchestrator
type Metrics struct {
	// Request counters
	TotalRequests    *prometheus.CounterVec
	CacheHits        *prometheus.CounterVec
	ProviderRequests *prometheus.CounterVec
	FailedRequests   *prometheus.CounterVec

	// Latency histograms
	RequestLatency  *prometheus.HistogramVec
	ProviderLatency *prometheus.HistogramVec

	// Cost metrics
	TotalCost       *prometheus.CounterVec
	CostPerProvider *prometheus.GaugeVec

	// Cache metrics
	CacheSize      prometheus.Gauge
	CacheEvictions prometheus.Counter

	// Provider health metrics
	ProviderHealth *prometheus.GaugeVec
	ProviderErrors *prometheus.CounterVec
}

// Global metrics instance
var metrics *Metrics

// InitializeMetrics initializes all Prometheus metrics
func InitializeMetrics() {
	metrics = &Metrics{
		TotalRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_requests_total",
			Help: "Total number of AI requests",
		}, []string{"task_type", "provider", "model"}),

		CacheHits: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_cache_hits_total",
			Help: "Total number of cache hits",
		}, []string{"task_type", "provider", "model"}),

		ProviderRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_provider_requests_total",
			Help: "Total number of requests per provider",
		}, []string{"provider", "model"}),

		FailedRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_failed_requests_total",
			Help: "Total number of failed requests",
		}, []string{"provider", "model", "error_type"}),

		RequestLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "orcaai_request_latency_seconds",
			Help:    "Request latency in seconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
		}, []string{"task_type", "provider", "model"}),

		ProviderLatency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "orcaai_provider_latency_seconds",
			Help:    "Provider latency in seconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
		}, []string{"provider", "model"}),

		TotalCost: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_total_cost_dollars",
			Help: "Total cost in USD",
		}, []string{"provider", "model"}),

		CostPerProvider: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "orcaai_cost_per_provider_dollars",
			Help: "Current cost per provider in USD",
		}, []string{"provider", "model"}),

		CacheSize: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "orcaai_cache_size",
			Help: "Current cache size",
		}),

		CacheEvictions: promauto.NewCounter(prometheus.CounterOpts{
			Name: "orcaai_cache_evictions_total",
			Help: "Total number of cache evictions",
		}),

		ProviderHealth: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "orcaai_provider_health",
			Help: "Provider health status (0-1)",
		}, []string{"provider", "model"}),

		ProviderErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "orcaai_provider_errors_total",
			Help: "Total number of provider errors",
		}, []string{"provider", "model", "error_type"}),
	}
}

// RecordRequest records a successful request
func (m *Metrics) RecordRequest(taskType, provider, model string, latency time.Duration, cost float64, cacheHit bool) {
	if m == nil {
		return
	}

	// Record total requests
	m.TotalRequests.WithLabelValues(taskType, provider, model).Inc()

	// Record cache hits
	if cacheHit {
		m.CacheHits.WithLabelValues(taskType, provider, model).Inc()
	}

	// Record provider requests
	m.ProviderRequests.WithLabelValues(provider, model).Inc()

	// Record latency
	m.RequestLatency.WithLabelValues(taskType, provider, model).Observe(latency.Seconds())
	m.ProviderLatency.WithLabelValues(provider, model).Observe(latency.Seconds())

	// Record cost
	m.TotalCost.WithLabelValues(provider, model).Add(cost)
	m.CostPerProvider.WithLabelValues(provider, model).Set(cost)
}

// RecordFailedRequest records a failed request
func (m *Metrics) RecordFailedRequest(provider, model, errorType string) {
	if m == nil {
		return
	}

	m.FailedRequests.WithLabelValues(provider, model, errorType).Inc()
	m.ProviderErrors.WithLabelValues(provider, model, errorType).Inc()
}

// RecordProviderHealth records provider health status
func (m *Metrics) RecordProviderHealth(provider, model string, health float64) {
	if m == nil {
		return
	}

	m.ProviderHealth.WithLabelValues(provider, model).Set(health)
}

// RecordCacheSize records current cache size
func (m *Metrics) RecordCacheSize(size int) {
	if m == nil {
		return
	}

	m.CacheSize.Set(float64(size))
}

// RecordCacheEviction records a cache eviction
func (m *Metrics) RecordCacheEviction() {
	if m == nil {
		return
	}

	m.CacheEvictions.Inc()
}

// GetMetrics returns the current metrics instance
func GetMetrics() *Metrics {
	return metrics
}

// RequestMetrics holds metrics for a single request
type RequestMetrics struct {
	TaskType  string
	Provider  string
	Model     string
	Latency   time.Duration
	Cost      float64
	CacheHit  bool
	Success   bool
	ErrorType string
}

// LogRequest logs metrics for a request
func LogRequest(metrics *RequestMetrics) {
	if metrics == nil || GetMetrics() == nil {
		return
	}

	if metrics.Success {
		GetMetrics().RecordRequest(
			metrics.TaskType,
			metrics.Provider,
			metrics.Model,
			metrics.Latency,
			metrics.Cost,
			metrics.CacheHit,
		)
	} else {
		GetMetrics().RecordFailedRequest(
			metrics.Provider,
			metrics.Model,
			metrics.ErrorType,
		)
	}
}
