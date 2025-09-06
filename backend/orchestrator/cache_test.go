package orchestrator

import (
	"testing"
	"time"
)

func TestCacheResult(t *testing.T) {
	// Reset cache for testing
	cache = make(map[string]CacheEntry)

	// Create a test AI response
	response := &AIResponse{
		Content:  "Test response",
		Provider: "openai",
		Model:    "gpt-3.5-turbo",
		Cost:     0.001,
	}

	// Cache the result
	CacheResult("Test prompt", "text-generation", "openai", "gpt-3.5-turbo", response)

	// Check that the result was cached
	key := generateCacheKey("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")
	if _, exists := cache[key]; !exists {
		t.Error("Expected result to be cached")
	}
}

func TestGetCachedResult(t *testing.T) {
	// Reset cache for testing
	cache = make(map[string]CacheEntry)

	// Create a test AI response
	response := &AIResponse{
		Content:  "Test response",
		Provider: "openai",
		Model:    "gpt-3.5-turbo",
		Cost:     0.001,
	}

	// Cache the result
	CacheResult("Test prompt", "text-generation", "openai", "gpt-3.5-turbo", response)

	// Retrieve the cached result
	cachedResponse, found := GetCachedResult("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")
	if !found {
		t.Fatal("Expected cached result to be found")
	}

	if cachedResponse.Content != "Test response" {
		t.Errorf("Expected content 'Test response', got '%s'", cachedResponse.Content)
	}

	if cachedResponse.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", cachedResponse.Provider)
	}

	if cachedResponse.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", cachedResponse.Model)
	}

	if cachedResponse.Cost != 0.001 {
		t.Errorf("Expected cost 0.001, got %f", cachedResponse.Cost)
	}
}

func TestCacheExpiry(t *testing.T) {
	// Reset cache for testing
	cache = make(map[string]CacheEntry)

	// Create a test AI response with a short expiry
	response := &AIResponse{
		Content:  "Test response",
		Provider: "openai",
		Model:    "gpt-3.5-turbo",
		Cost:     0.001,
	}

	// Cache the result with a short expiry
	key := generateCacheKey("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")
	cache[key] = CacheEntry{
		Response:  response,
		Timestamp: time.Now().Add(-2 * time.Hour), // Expired 2 hours ago
	}

	// Try to retrieve the expired result
	_, found := GetCachedResult("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")
	if found {
		t.Error("Expected expired result to not be found")
	}

	// Check that the expired entry was removed
	if _, exists := cache[key]; exists {
		t.Error("Expected expired entry to be removed")
	}
}

func TestGenerateCacheKey(t *testing.T) {
	// Test that the cache key generation is consistent
	key1 := generateCacheKey("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")
	key2 := generateCacheKey("Test prompt", "text-generation", "openai", "gpt-3.5-turbo")

	if key1 != key2 {
		t.Error("Expected cache keys to be identical for identical inputs")
	}

	// Test that different inputs produce different keys
	key3 := generateCacheKey("Different prompt", "text-generation", "openai", "gpt-3.5-turbo")
	if key1 == key3 {
		t.Error("Expected cache keys to be different for different prompts")
	}
}
