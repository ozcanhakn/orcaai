package orchestrator

import (
	"testing"
	"time"
)

func TestLogRequest(t *testing.T) {
	// Reset metrics for testing
	requestMetrics = make(chan *RequestMetrics, 100)
	close(requestMetrics)
	requestMetrics = make(chan *RequestMetrics, 100)

	// Create a test metrics entry
	metrics := &RequestMetrics{
		TaskType: "text-generation",
		Provider: "openai",
		Model:    "gpt-3.5-turbo",
		Latency:  500 * time.Millisecond,
		Cost:     0.001,
		CacheHit: false,
		Success:  true,
	}

	// Log the request
	LogRequest(metrics)

	// Since we're not running the metrics processor, we can't directly test
	// that the metrics were processed, but we can test that the function
	// doesn't panic or cause errors
}

func TestGetMetrics(t *testing.T) {
	// Reset metrics for testing
	providerMetrics = make(map[string]*ProviderMetrics)
	providerMetrics["openai"] = &ProviderMetrics{
		TotalRequests: 100,
		TotalCost:     0.5,
		AvgLatency:    450 * time.Millisecond,
		CacheHits:     20,
		Errors:        5,
	}

	// Get metrics
	metrics := GetMetrics()

	if len(metrics) != 1 {
		t.Errorf("Expected 1 provider metric, got %d", len(metrics))
	}

	openaiMetrics := metrics["openai"]
	if openaiMetrics == nil {
		t.Fatal("Expected openai metrics to exist")
	}

	if openaiMetrics.TotalRequests != 100 {
		t.Errorf("Expected 100 total requests, got %d", openaiMetrics.TotalRequests)
	}

	if openaiMetrics.TotalCost != 0.5 {
		t.Errorf("Expected total cost 0.5, got %f", openaiMetrics.TotalCost)
	}

	if openaiMetrics.CacheHits != 20 {
		t.Errorf("Expected 20 cache hits, got %d", openaiMetrics.CacheHits)
	}

	if openaiMetrics.Errors != 5 {
		t.Errorf("Expected 5 errors, got %d", openaiMetrics.Errors)
	}
}

func TestCalculateUptime(t *testing.T) {
	// Reset metrics for testing
	providerMetrics = make(map[string]*ProviderMetrics)
	providerMetrics["openai"] = &ProviderMetrics{
		TotalRequests: 100,
		Errors:        5,
	}

	// Calculate uptime
	uptime := calculateUptime("openai")

	// Expected uptime: (100 - 5) / 100 = 0.95 = 95%
	expected := 95.0
	if uptime != expected {
		t.Errorf("Expected uptime %f, got %f", expected, uptime)
	}

	// Test with no requests
	providerMetrics["claude"] = &ProviderMetrics{
		TotalRequests: 0,
		Errors:        0,
	}

	uptime = calculateUptime("claude")
	if uptime != 100.0 {
		t.Errorf("Expected uptime 100.0 for no requests, got %f", uptime)
	}
}

func TestCalculateCostSavings(t *testing.T) {
	// Reset metrics for testing
	providerMetrics = make(map[string]*ProviderMetrics)
	providerMetrics["openai"] = &ProviderMetrics{
		TotalRequests: 100,
		TotalCost:     1.0, // $1.00 without caching
		CacheHits:     20,  // 20% cache hit rate
	}

	// Calculate cost savings
	// With 20% cache hits, we save 20% of the cost
	// Savings = 1.0 * 0.2 = 0.2
	savings := calculateCostSavings("openai")

	if savings != 0.2 {
		t.Errorf("Expected cost savings 0.2, got %f", savings)
	}

	// Test with no requests
	providerMetrics["claude"] = &ProviderMetrics{
		TotalRequests: 0,
		TotalCost:     0,
		CacheHits:     0,
	}

	savings = calculateCostSavings("claude")
	if savings != 0.0 {
		t.Errorf("Expected cost savings 0.0 for no requests, got %f", savings)
	}
}
