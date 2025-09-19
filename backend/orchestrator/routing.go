package orchestrator

import (
	"context"
	"fmt"
	"math"
	"time"
)

type Provider struct {
	Name         string   `json:"name"`
	Model        string   `json:"model"`
	CostPer1K    float64  `json:"cost_per_1k"`
	AvgLatency   int      `json:"avg_latency_ms"`
	Reliability  float64  `json:"reliability"` // 0-1 score
	MaxTokens    int      `json:"max_tokens"`
	Capabilities []string `json:"capabilities"`
}

type RoutingResult struct {
	Provider   string     `json:"provider"`
	Model      string     `json:"model"`
	Confidence float64    `json:"confidence"`
	Reasoning  string     `json:"reasoning"`
	Fallbacks  []Provider `json:"fallbacks"`
}

type RoutingCriteria struct {
	CostWeight        float64 `json:"cost_weight"`
	LatencyWeight     float64 `json:"latency_weight"`
	ReliabilityWeight float64 `json:"reliability_weight"`
	QualityWeight     float64 `json:"quality_weight"`
}

type ProviderStatus struct {
	Name       string    `json:"name"`
	Model      string    `json:"model"`
	Healthy    bool      `json:"healthy"`
	LastCheck  time.Time `json:"last_check"`
	ErrorCount int       `json:"error_count"`
	LastError  string    `json:"last_error"`
}

type FallbackStrategy struct {
	MaxRetries     int           `json:"max_retries"`
	RetryDelay     time.Duration `json:"retry_delay"`
	Timeout        time.Duration `json:"timeout"`
	CircuitBreaker bool          `json:"circuit_breaker"`
}

var (
	// Available providers with their characteristics
	availableProviders = map[string][]Provider{
		"text-generation": {
			{
				Name:         "openai",
				Model:        "gpt-4",
				CostPer1K:    0.03,
				AvgLatency:   2000,
				Reliability:  0.95,
				MaxTokens:    8000,
				Capabilities: []string{"text-generation", "reasoning", "code"},
			},
			{
				Name:         "openai",
				Model:        "gpt-3.5-turbo",
				CostPer1K:    0.002,
				AvgLatency:   1000,
				Reliability:  0.90,
				MaxTokens:    4000,
				Capabilities: []string{"text-generation", "conversation"},
			},
			{
				Name:         "claude",
				Model:        "claude-3-opus",
				CostPer1K:    0.015,
				AvgLatency:   3000,
				Reliability:  0.98,
				MaxTokens:    200000,
				Capabilities: []string{"text-generation", "reasoning", "analysis"},
			},
			{
				Name:         "claude",
				Model:        "claude-3-sonnet",
				CostPer1K:    0.003,
				AvgLatency:   1500,
				Reliability:  0.95,
				MaxTokens:    200000,
				Capabilities: []string{"text-generation", "conversation"},
			},
			{
				Name:         "gemini",
				Model:        "gemini-pro",
				CostPer1K:    0.001,
				AvgLatency:   2500,
				Reliability:  0.85,
				MaxTokens:    30000,
				Capabilities: []string{"text-generation", "multimodal"},
			},
		},
		"summarization": {
			{
				Name:         "claude",
				Model:        "claude-3-sonnet",
				CostPer1K:    0.003,
				AvgLatency:   1500,
				Reliability:  0.95,
				MaxTokens:    200000,
				Capabilities: []string{"summarization", "analysis"},
			},
			{
				Name:         "openai",
				Model:        "gpt-3.5-turbo",
				CostPer1K:    0.002,
				AvgLatency:   1000,
				Reliability:  0.90,
				MaxTokens:    4000,
				Capabilities: []string{"summarization"},
			},
		},
		"code-generation": {
			{
				Name:         "openai",
				Model:        "gpt-4",
				CostPer1K:    0.03,
				AvgLatency:   2000,
				Reliability:  0.95,
				MaxTokens:    8000,
				Capabilities: []string{"code-generation", "debugging"},
			},
			{
				Name:         "claude",
				Model:        "claude-3-opus",
				CostPer1K:    0.015,
				AvgLatency:   3000,
				Reliability:  0.98,
				MaxTokens:    200000,
				Capabilities: []string{"code-generation", "code-review"},
			},
		},
	}

	// Default routing criteria - can be overridden per user
	defaultCriteria = RoutingCriteria{
		CostWeight:        0.3,
		LatencyWeight:     0.3,
		ReliabilityWeight: 0.3,
		QualityWeight:     0.1,
	}

	// Global cache instance
	cache Cache

	// Provider status tracking
	providerStatus = make(map[string]*ProviderStatus)

	// Default fallback strategy
	defaultFallbackStrategy = FallbackStrategy{
		MaxRetries:     2,
		RetryDelay:     1 * time.Second,
		Timeout:        30 * time.Second,
		CircuitBreaker: true,
	}
)

// InitializeCache initializes the cache system
func InitializeCache(cacheType, addr, password string, db int) error {
	var err error
	switch cacheType {
	case "redis":
		cache, err = NewRedisCache(addr, password, db)
	case "memory":
		cache = NewInMemoryCache()
	default:
		cache = NewInMemoryCache()
	}
	return err
}

// CloseCache closes the cache connection
func CloseCache() error {
	if cache != nil {
		return cache.Close()
	}
	return nil
}

// UpdateProviderStatus updates the health status of a provider
func UpdateProviderStatus(provider, model string, healthy bool, err error) {
	key := fmt.Sprintf("%s:%s", provider, model)

	status, exists := providerStatus[key]
	if !exists {
		status = &ProviderStatus{
			Name:      provider,
			Model:     model,
			Healthy:   true,
			LastCheck: time.Now(),
		}
		providerStatus[key] = status
	}

	status.LastCheck = time.Now()

	if healthy {
		status.ErrorCount = 0
		status.LastError = ""
		status.Healthy = true
	} else {
		status.ErrorCount++
		if err != nil {
			status.LastError = err.Error()
		}

		// Mark as unhealthy if too many errors
		if status.ErrorCount > 5 {
			status.Healthy = false
		}
	}
}

// IsProviderHealthy checks if a provider is healthy
func IsProviderHealthy(provider, model string) bool {
	key := fmt.Sprintf("%s:%s", provider, model)

	status, exists := providerStatus[key]
	if !exists {
		return true // Assume healthy if no status recorded
	}

	// If status is old, assume healthy
	if time.Since(status.LastCheck) > 5*time.Minute {
		return true
	}

	return status.Healthy
}

// GetFallbackProviders returns a list of healthy fallback providers
func GetFallbackProviders(taskType, excludeProvider, excludeModel string) []Provider {
	providers, exists := availableProviders[taskType]
	if !exists {
		providers = availableProviders["text-generation"]
	}

	fallbacks := []Provider{}
	for _, provider := range providers {
		// Skip excluded provider
		if provider.Name == excludeProvider && provider.Model == excludeModel {
			continue
		}

		// Only include healthy providers
		if IsProviderHealthy(provider.Name, provider.Model) {
			fallbacks = append(fallbacks, provider)
		}
	}

	// Sort by reliability (most reliable first)
	for i := 0; i < len(fallbacks)-1; i++ {
		for j := i + 1; j < len(fallbacks); j++ {
			if fallbacks[i].Reliability < fallbacks[j].Reliability {
				fallbacks[i], fallbacks[j] = fallbacks[j], fallbacks[i]
			}
		}
	}

	return fallbacks
}

// RouteRequestWithFallback routes a request with fallback support
func RouteRequestWithFallback(prompt, taskType, preferredProvider, preferredModel string, strategy *FallbackStrategy) (*RoutingResult, error) {
	if strategy == nil {
		strategy = &defaultFallbackStrategy
	}

	// Try primary routing
	result, err := RouteRequest(prompt, taskType, preferredProvider, preferredModel)
	if err != nil {
		return nil, err
	}

	// Add fallback providers
	result.Fallbacks = GetFallbackProviders(taskType, result.Provider, result.Model)

	return result, nil
}

func RouteRequest(prompt, taskType, preferredProvider, preferredModel string) (*RoutingResult, error) {
	// Check cache first
	ctx := context.Background()
	cacheKey := GenerateCacheKey(prompt, taskType, preferredProvider, preferredModel)

	if cache != nil {
		if cachedResult, err := cache.Get(ctx, cacheKey); err == nil && cachedResult != nil {
			// Return cached result with modified reasoning
			return &RoutingResult{
				Provider:   cachedResult.Provider,
				Model:      cachedResult.Model,
				Confidence: 1.0,
				Reasoning:  "Cache hit - returning cached result",
				Fallbacks:  GetFallbackProviders(taskType, cachedResult.Provider, cachedResult.Model),
			}, nil
		}
	}

	// If user specified a provider/model, use it (if available and healthy)
	if preferredProvider != "" && preferredModel != "" {
		if provider := findSpecificProvider(preferredProvider, preferredModel, taskType); provider != nil {
			// Check if provider is healthy
			if IsProviderHealthy(preferredProvider, preferredModel) {
				return &RoutingResult{
					Provider:   preferredProvider,
					Model:      preferredModel,
					Confidence: 1.0,
					Reasoning:  "User-specified provider and model",
					Fallbacks:  GetFallbackProviders(taskType, preferredProvider, preferredModel),
				}, nil
			}
			// If not healthy, fall through to auto-routing
		}
	}

	// Auto-route based on task type and criteria
	return autoRoute(prompt, taskType)
}

func autoRoute(prompt, taskType string) (*RoutingResult, error) {
	// Default to text-generation if no specific task type
	if taskType == "" {
		taskType = "text-generation"
	}

	providers, exists := availableProviders[taskType]
	if !exists {
		providers = availableProviders["text-generation"]
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers available for task type: %s", taskType)
	}

	// Score each provider
	bestProvider := scoreProviders(providers, prompt, taskType)

	return &RoutingResult{
		Provider:   bestProvider.Name,
		Model:      bestProvider.Model,
		Confidence: calculateConfidence(bestProvider, providers),
		Reasoning:  generateReasoning(bestProvider, taskType),
		Fallbacks:  getFallbacks(taskType, bestProvider.Name, bestProvider.Model),
	}, nil
}

func scoreProviders(providers []Provider, prompt, taskType string) Provider {
	bestScore := -1.0
	var bestProvider Provider

	promptLength := len(prompt)

	for _, provider := range providers {
		score := calculateScore(provider, promptLength, taskType)
		if score > bestScore {
			bestScore = score
			bestProvider = provider
		}
	}

	return bestProvider
}

func calculateScore(provider Provider, promptLength int, taskType string) float64 {
	// Normalize metrics (0-1 scale)
	costScore := 1.0 - math.Min(provider.CostPer1K/0.05, 1.0)                // Lower cost = higher score
	latencyScore := 1.0 - math.Min(float64(provider.AvgLatency)/5000.0, 1.0) // Lower latency = higher score
	reliabilityScore := provider.Reliability

	// Quality score based on model capabilities and task match
    qualityScore := calculateQualityScore(provider, taskType)

	// Token capacity consideration
	tokenCapacityScore := 1.0
	estimatedTokens := promptLength / 4 // Rough estimation: 1 token ≈ 4 characters
	if estimatedTokens > provider.MaxTokens {
		tokenCapacityScore = 0.0 // Cannot handle this request
	}

	// Weighted final score
	score := (costScore * defaultCriteria.CostWeight) +
		(latencyScore * defaultCriteria.LatencyWeight) +
		(reliabilityScore * defaultCriteria.ReliabilityWeight) +
		(qualityScore * defaultCriteria.QualityWeight)

	return score * tokenCapacityScore
}

func calculateQualityScore(provider Provider, taskType string) float64 {
	// Model-specific quality scores for different tasks
	qualityMatrix := map[string]map[string]float64{
		"gpt-4": {
			"text-generation": 0.95,
			"code-generation": 0.90,
			"summarization":   0.85,
			"reasoning":       0.95,
		},
		"gpt-3.5-turbo": {
			"text-generation": 0.80,
			"code-generation": 0.70,
			"summarization":   0.85,
			"conversation":    0.90,
		},
		"claude-3-opus": {
			"text-generation": 0.98,
			"code-generation": 0.95,
			"summarization":   0.95,
			"reasoning":       0.98,
		},
		"claude-3-sonnet": {
			"text-generation": 0.85,
			"code-generation": 0.80,
			"summarization":   0.90,
			"conversation":    0.88,
		},
		"gemini-pro": {
			"text-generation": 0.75,
			"multimodal":      0.90,
			"summarization":   0.70,
		},
	}

	if modelScores, exists := qualityMatrix[provider.Model]; exists {
		if score, exists := modelScores[taskType]; exists {
			return score
		}
	}

	return 0.7 // Default quality score
}

func calculateConfidence(bestProvider Provider, allProviders []Provider) float64 {
	if len(allProviders) <= 1 {
		return 1.0
	}

	bestScore := calculateScore(bestProvider, 1000, "text-generation")

	// Find second best score
	secondBestScore := 0.0
	for _, provider := range allProviders {
		if provider.Name != bestProvider.Name || provider.Model != bestProvider.Model {
			score := calculateScore(provider, 1000, "text-generation")
			if score > secondBestScore {
				secondBestScore = score
			}
		}
	}

	// Confidence is based on the gap between best and second best
	gap := bestScore - secondBestScore
	confidence := math.Min(0.5+gap, 1.0)

	return confidence
}

func generateReasoning(provider Provider, taskType string) string {
	reasons := []string{}

	if provider.CostPer1K < 0.01 {
		reasons = append(reasons, "cost-effective")
	}
	if provider.AvgLatency < 2000 {
		reasons = append(reasons, "fast response")
	}
	if provider.Reliability > 0.95 {
		reasons = append(reasons, "high reliability")
	}
	if provider.MaxTokens > 10000 {
		reasons = append(reasons, "large context window")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "balanced performance")
	}

	return fmt.Sprintf("Selected for %s: %s", taskType, fmt.Sprintf("%v", reasons))
}

func getFallbacks(taskType, excludeProvider, excludeModel string) []Provider {
	providers, exists := availableProviders[taskType]
	if !exists {
		providers = availableProviders["text-generation"]
	}

	fallbacks := []Provider{}
	for _, provider := range providers {
		if provider.Name != excludeProvider || provider.Model != excludeModel {
			fallbacks = append(fallbacks, provider)
		}
		if len(fallbacks) >= 2 { // Limit fallbacks
			break
		}
	}

	return fallbacks
}

func findSpecificProvider(name, model, taskType string) *Provider {
	providers, exists := availableProviders[taskType]
	if !exists {
		// Search in all task types
		for _, taskProviders := range availableProviders {
			providers = append(providers, taskProviders...)
		}
	}

	for _, provider := range providers {
		if provider.Name == name && provider.Model == model {
			return &provider
		}
	}

	return nil
}

// GetProviderStats returns current provider statistics
func GetProviderStats() map[string]interface{} {
	stats := map[string]interface{}{
		"total_providers":  0,
		"active_providers": map[string]interface{}{},
		"task_types":       []string{},
	}

	totalProviders := 0
	activeProviders := map[string]interface{}{}

	for taskType, providers := range availableProviders {
		stats["task_types"] = append(stats["task_types"].([]string), taskType)
		totalProviders += len(providers)

		for _, provider := range providers {
			key := fmt.Sprintf("%s:%s", provider.Name, provider.Model)
			activeProviders[key] = map[string]interface{}{
				"cost_per_1k":  provider.CostPer1K,
				"avg_latency":  provider.AvgLatency,
				"reliability":  provider.Reliability,
				"max_tokens":   provider.MaxTokens,
				"capabilities": provider.Capabilities,
			}
		}
	}

	stats["total_providers"] = totalProviders
	stats["active_providers"] = activeProviders

	return stats
}

// CacheResponseResult stores the result in cache
func CacheResponseResult(prompt, taskType, provider, model string, response interface{}) error {
	if cache == nil {
		return nil
	}

	ctx := context.Background()
	cacheKey := GenerateCacheKey(prompt, taskType, provider, model)

	result := &CacheResult{
		Response:  response,
		Provider:  provider,
		Model:     model,
		CachedKey: cacheKey,
	}

	// Cache for 1 hour by default
	return cache.Set(ctx, cacheKey, result, 1*time.Hour)
}
