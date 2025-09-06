package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// Test with default values
	client := NewClient(Config{
		APIKey: "test-key",
	})

	if client.apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", client.apiKey)
	}

	if client.baseURL != "http://localhost:8080" {
		t.Errorf("Expected base URL 'http://localhost:8080', got '%s'", client.baseURL)
	}

	if client.client.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got '%s'", client.client.Timeout)
	}

	// Test with custom values
	client = NewClient(Config{
		APIKey:  "test-key",
		BaseURL: "https://api.orcaai.com",
		Timeout: 60 * time.Second,
	})

	if client.baseURL != "https://api.orcaai.com" {
		t.Errorf("Expected base URL 'https://api.orcaai.com', got '%s'", client.baseURL)
	}

	if client.client.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got '%s'", client.client.Timeout)
	}
}

func TestQuery(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Send a mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "test-id",
			"content": "Test response",
			"provider": "openai",
			"model": "gpt-3.5-turbo",
			"cost": 0.001,
			"latency_ms": 450,
			"cache_hit": false,
			"timestamp": "2023-01-01T00:00:00Z"
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient(Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	// Make a query
	req := QueryRequest{
		Prompt:   "Test prompt",
		TaskType: "text-generation",
	}

	resp, err := client.Query(req)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	// Verify the response
	if resp.Content != "Test response" {
		t.Errorf("Expected content 'Test response', got '%s'", resp.Content)
	}

	if resp.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", resp.Provider)
	}

	if resp.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", resp.Model)
	}

	if resp.Cost != 0.001 {
		t.Errorf("Expected cost 0.001, got %f", resp.Cost)
	}

	if resp.Latency != 450 {
		t.Errorf("Expected latency 450, got %d", resp.Latency)
	}
}

func TestGetProviders(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", r.Header.Get("Authorization"))
		}

		// Send a mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"providers": [
				{
					"name": "OpenAI",
					"id": "openai",
					"models": ["gpt-3.5-turbo", "gpt-4"],
					"cost_per_1k": 0.002,
					"max_tokens": 4096,
					"status": "active"
				},
				{
					"name": "Claude",
					"id": "claude",
					"models": ["claude-2"],
					"cost_per_1k": 0.001,
					"max_tokens": 100000,
					"status": "active"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient(Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	// Get providers
	resp, err := client.GetProviders()
	if err != nil {
		t.Fatalf("GetProviders failed: %v", err)
	}

	// Verify the response
	if len(resp.Providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(resp.Providers))
	}

	if resp.Providers[0].Name != "OpenAI" {
		t.Errorf("Expected first provider name 'OpenAI', got '%s'", resp.Providers[0].Name)
	}

	if resp.Providers[1].Name != "Claude" {
		t.Errorf("Expected second provider name 'Claude', got '%s'", resp.Providers[1].Name)
	}
}

func TestGetMetrics(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", r.Header.Get("Authorization"))
		}

		// Send a mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"total_requests": 1000,
			"avg_latency": 450,
			"cost_savings": 25.50,
			"uptime": 99.9
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient(Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	// Get metrics
	resp, err := client.GetMetrics()
	if err != nil {
		t.Fatalf("GetMetrics failed: %v", err)
	}

	// Verify the response
	if resp.TotalRequests != 1000 {
		t.Errorf("Expected total requests 1000, got %d", resp.TotalRequests)
	}

	if resp.AvgLatency != 450 {
		t.Errorf("Expected avg latency 450, got %d", resp.AvgLatency)
	}

	if resp.CostSavings != 25.50 {
		t.Errorf("Expected cost savings 25.50, got %f", resp.CostSavings)
	}

	if resp.Uptime != 99.9 {
		t.Errorf("Expected uptime 99.9, got %f", resp.Uptime)
	}
}
