package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"orcaai/database"
	"orcaai/models"
	"orcaai/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AIQueryRequest struct {
	Prompt    string                 `json:"prompt" binding:"required"`
	TaskType  string                 `json:"task_type,omitempty"`
	Model     string                 `json:"model,omitempty"`
	MaxTokens int                    `json:"max_tokens,omitempty"`
	Provider  string                 `json:"provider,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

type AIQueryResponse struct {
	ID         string                 `json:"id"`
	Content    string                 `json:"content"`
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model"`
	TokensUsed models.TokenUsage      `json:"tokens_used"`
	Cost       float64                `json:"cost"`
	Latency    int                    `json:"latency_ms"`
	CacheHit   bool                   `json:"cache_hit"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func AIQuery(c *gin.Context) {
	startTime := time.Now()

	// Parse request
	var req AIQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user info from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Route to best AI provider with fallback support
	routingResult, err := orchestrator.RouteRequestWithFallback(req.Prompt, req.TaskType, req.Provider, req.Model, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to route request",
			"details": err.Error(),
		})
		return
	}

	// Check if this is a cache hit
	if routingResult.Reasoning == "Cache hit - returning cached result" {
		// For cache hits, we need to simulate a response since we don't have the actual content
		latency := int(time.Since(startTime).Milliseconds())

		// Record metrics for cache hit
		metrics := &orchestrator.RequestMetrics{
			TaskType: req.TaskType,
			Provider: routingResult.Provider,
			Model:    routingResult.Model,
			Latency:  time.Since(startTime),
			Cost:     0,
			CacheHit: true,
			Success:  true,
		}
		orchestrator.LogRequest(metrics)

		c.JSON(http.StatusOK, AIQueryResponse{
			ID:         requestID,
			Content:    "Cached response for: " + req.Prompt, // This is a placeholder
			Provider:   routingResult.Provider,
			Model:      routingResult.Model,
			TokensUsed: models.TokenUsage{Input: 0, Output: 0},
			Cost:       0, // Cache hits are free
			Latency:    latency,
			CacheHit:   true,
			Timestamp:  time.Now(),
		})
		return
	}

	// Make AI request with enhanced fallback
	aiResponse, err := makeAIRequestWithEnhancedFallback(routingResult, req)
	if err != nil {
		latency := int(time.Since(startTime).Milliseconds())
		logRequest(userID.(string), requestID, "error", req, nil, latency, false, err)

		// Record metrics for failed request
		metrics := &orchestrator.RequestMetrics{
			TaskType:  req.TaskType,
			Provider:  routingResult.Provider,
			Model:     routingResult.Model,
			Latency:   time.Since(startTime),
			Cost:      0,
			CacheHit:  false,
			Success:   false,
			ErrorType: "provider_error",
		}
		orchestrator.LogRequest(metrics)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "AI request failed",
			"details": err.Error(),
		})
		return
	}

	// Calculate final metrics
	latency := int(time.Since(startTime).Milliseconds())

	// Cache the response using the orchestrator's cache function
	orchestrator.CacheResult(req.Prompt, req.TaskType, aiResponse.Provider, aiResponse.Model, aiResponse)

	// Log successful request
	logRequest(userID.(string), requestID, aiResponse.Provider, req, aiResponse, latency, false, nil)

	// Record metrics for successful request
	metrics := &orchestrator.RequestMetrics{
		TaskType: req.TaskType,
		Provider: aiResponse.Provider,
		Model:    aiResponse.Model,
		Latency:  time.Since(startTime),
		Cost:     aiResponse.Cost,
		CacheHit: false,
		Success:  true,
	}
	orchestrator.LogRequest(metrics)

	// Return response
	c.JSON(http.StatusOK, AIQueryResponse{
		ID:         requestID,
		Content:    aiResponse.Content,
		Provider:   aiResponse.Provider,
		Model:      aiResponse.Model,
		TokensUsed: aiResponse.TokensUsed,
		Cost:       aiResponse.Cost,
		Latency:    latency,
		CacheHit:   false,
		Timestamp:  time.Now(),
		Metadata:   aiResponse.Metadata,
	})
}

// makeAIRequestWithEnhancedFallback makes an AI request with enhanced fallback support
func makeAIRequestWithEnhancedFallback(routing *orchestrator.RoutingResult, req AIQueryRequest) (*models.AIResponse, error) {
	// Try primary provider
	response, err := makeAIRequest(routing.Provider, routing.Model, req)
	if err == nil {
		// Update provider status as healthy
		orchestrator.UpdateProviderStatus(routing.Provider, routing.Model, true, nil)
		return response, nil
	}

	// Update provider status as unhealthy
	orchestrator.UpdateProviderStatus(routing.Provider, routing.Model, false, err)

	// Try fallback providers
	for _, fallback := range routing.Fallbacks {
		response, err = makeAIRequest(fallback.Name, fallback.Model, req)
		if err == nil {
			// Update fallback provider status as healthy
			orchestrator.UpdateProviderStatus(fallback.Name, fallback.Model, true, nil)
			return response, nil
		}

		// Update fallback provider status as unhealthy
		orchestrator.UpdateProviderStatus(fallback.Name, fallback.Model, false, err)
	}

	return nil, fmt.Errorf("all providers failed after fallback attempts")
}

func makeAIRequest(provider, model string, req AIQueryRequest) (*models.AIResponse, error) {
	switch provider {
	case "openai":
		return makeOpenAIRequest(model, req)
	case "claude":
		return makeClaudeRequest(model, req)
	case "gemini":
		return makeGeminiRequest(model, req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func makeOpenAIRequest(model string, req AIQueryRequest) (*models.AIResponse, error) {
	// OpenAI API request implementation
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
		"max_tokens": req.MaxTokens,
	}

	jsonBody, _ := json.Marshal(requestBody)

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+getAPIKey("openai"))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse OpenAI response
	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Calculate cost (example rates)
	cost := float64(openAIResp.Usage.PromptTokens)*0.001 + float64(openAIResp.Usage.CompletionTokens)*0.002

	return &models.AIResponse{
		Content:  openAIResp.Choices[0].Message.Content,
		Provider: "openai",
		Model:    model,
		TokensUsed: models.TokenUsage{
			Input:  openAIResp.Usage.PromptTokens,
			Output: openAIResp.Usage.CompletionTokens,
		},
		Cost: cost,
	}, nil
}

func makeClaudeRequest(model string, req AIQueryRequest) (*models.AIResponse, error) {
	// Claude API request implementation (similar structure)
	// This is a placeholder - implement based on Claude's actual API
	return &models.AIResponse{
		Content:    "Claude response placeholder",
		Provider:   "claude",
		Model:      model,
		TokensUsed: models.TokenUsage{Input: 100, Output: 50},
		Cost:       0.01,
	}, nil
}

func makeGeminiRequest(model string, req AIQueryRequest) (*models.AIResponse, error) {
	// Gemini API request implementation (similar structure)
	// This is a placeholder - implement based on Gemini's actual API
	return &models.AIResponse{
		Content:    "Gemini response placeholder",
		Provider:   "gemini",
		Model:      model,
		TokensUsed: models.TokenUsage{Input: 100, Output: 50},
		Cost:       0.005,
	}, nil
}

func generateCacheKey(prompt, taskType, model string) string {
	data := fmt.Sprintf("%s:%s:%s", prompt, taskType, model)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func checkCache(key string) (*models.AIResponse, bool) {
	ctx := context.Background()
	data, err := database.Redis.Get(ctx, "cache:"+key).Result()
	if err != nil {
		return nil, false
	}

	var response models.AIResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, false
	}

	return &response, true
}

func cacheResponse(key string, response *models.AIResponse) {
	ctx := context.Background()
	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	database.Redis.Set(ctx, "cache:"+key, data, 24*time.Hour)
}

func getAPIKey(provider string) string {
	// In production, get from secure config or environment
	switch provider {
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "claude":
		return os.Getenv("CLAUDE_API_KEY")
	case "gemini":
		return os.Getenv("GEMINI_API_KEY")
	default:
		return ""
	}
}

func logRequest(userID, requestID, provider string, req AIQueryRequest, response *models.AIResponse, latency int, cacheHit bool, err error) {
	var tokensInput, tokensOutput int
	var cost float64
	var status, errorMsg string

	if response != nil {
		tokensInput = response.TokensUsed.Input
		tokensOutput = response.TokensUsed.Output
		cost = response.Cost
		status = "success"
	} else if err != nil {
		status = "error"
		errorMsg = err.Error()
	}

	query := `
		INSERT INTO request_logs 
		(user_id, provider, model, prompt_tokens, completion_tokens, cost_usd, latency_ms, cache_hit, status, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, dbErr := database.DB.Exec(query, userID, provider, req.Model, tokensInput, tokensOutput, cost, latency, cacheHit, status, errorMsg)
	if dbErr != nil {
		fmt.Printf("Failed to log request: %v\n", dbErr)
	}
}

func GetProviders(c *gin.Context) {
	providers := []map[string]interface{}{
		{
			"name":        "OpenAI",
			"id":          "openai",
			"models":      []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo"},
			"cost_per_1k": 0.03,
			"max_tokens":  4000,
			"status":      "active",
		},
		{
			"name":        "Claude",
			"id":          "claude",
			"models":      []string{"claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
			"cost_per_1k": 0.015,
			"max_tokens":  100000,
			"status":      "active",
		},
		{
			"name":        "Gemini",
			"id":          "gemini",
			"models":      []string{"gemini-pro", "gemini-pro-vision"},
			"cost_per_1k": 0.001,
			"max_tokens":  30720,
			"status":      "active",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
	})
}

// GetAPIKeys returns all API keys for the authenticated user
func GetAPIKeys(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch API keys from database
	var apiKeys []models.APIKey
	query := "SELECT id, user_id, name, key, permissions, is_active, created_at, updated_at FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC"
	err := database.DB.Select(&apiKeys, query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch API keys",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"api_keys": apiKeys,
	})
}

// CreateAPIKey creates a new API key for the authenticated user
func CreateAPIKey(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Generate a new API key
	key, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Create API key record
	apiKey := models.APIKey{
		ID:          uuid.New(),
		UserID:      userID.(uuid.UUID),
		Name:        req.Name,
		Key:         key,
		Permissions: []string{models.PermissionReadMetrics, models.PermissionWriteAPIKeys, models.PermissionReadAPIKeys},
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert API key into database
	query := `
		INSERT INTO api_keys (id, user_id, name, key, permissions, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = database.DB.Exec(query, apiKey.ID, apiKey.UserID, apiKey.Name, apiKey.Key,
		fmt.Sprintf("%v", apiKey.Permissions), apiKey.IsActive, apiKey.CreatedAt, apiKey.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create API key",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "API key created successfully",
		"api_key": apiKey,
	})
}

// DeleteAPIKey deletes an API key for the authenticated user
func DeleteAPIKey(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get API key ID from URL parameter
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
		return
	}

	// Delete API key from database (only if it belongs to the user)
	query := "DELETE FROM api_keys WHERE id = $1 AND user_id = $2"
	result, err := database.DB.Exec(query, keyID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete API key",
			"details": err.Error(),
		})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check deletion result",
			"details": err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found or not owned by user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key deleted successfully",
	})
}

// generateAPIKey generates a new random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
