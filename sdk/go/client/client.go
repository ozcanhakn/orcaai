package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the OrcaAI client
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Config represents the client configuration
type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// NewClient creates a new OrcaAI client
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// QueryRequest represents a query request
type QueryRequest struct {
	Prompt    string `json:"prompt"`
	TaskType  string `json:"task_type,omitempty"`
	Provider  string `json:"provider,omitempty"`
	Model     string `json:"model,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty"`
}

// QueryResponse represents a query response
type QueryResponse struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Provider  string                 `json:"provider"`
	Model     string                 `json:"model"`
	Cost      float64                `json:"cost"`
	Latency   int                    `json:"latency_ms"`
	CacheHit  bool                   `json:"cache_hit"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Query sends a query to the OrcaAI platform
func (c *Client) Query(req QueryRequest) (*QueryResponse, error) {
	url := c.baseURL + "/api/v1/ai/query"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &queryResp, nil
}

// Provider represents an AI provider
type Provider struct {
	Name      string   `json:"name"`
	ID        string   `json:"id"`
	Models    []string `json:"models"`
	CostPer1K float64  `json:"cost_per_1k"`
	MaxTokens int      `json:"max_tokens"`
	Status    string   `json:"status"`
}

// ProvidersResponse represents the providers response
type ProvidersResponse struct {
	Providers []Provider `json:"providers"`
}

// GetProviders gets the list of available providers
func (c *Client) GetProviders() (*ProvidersResponse, error) {
	url := c.baseURL + "/api/v1/ai/providers"

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var providersResp ProvidersResponse
	if err := json.NewDecoder(resp.Body).Decode(&providersResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &providersResp, nil
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	TotalRequests int     `json:"total_requests"`
	AvgLatency    int     `json:"avg_latency"`
	CostSavings   float64 `json:"cost_savings"`
	Uptime        float64 `json:"uptime"`
}

// GetMetrics gets usage metrics
func (c *Client) GetMetrics() (*MetricsResponse, error) {
	url := c.baseURL + "/api/v1/metrics"

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var metricsResp MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metricsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metricsResp, nil
}
