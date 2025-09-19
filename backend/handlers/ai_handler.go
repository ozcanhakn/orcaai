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

	"orcaai/backend/database"
	"orcaai/backend/models"
	"orcaai/backend/orchestrator"

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

// AIQueryStream streams token chunks back to client (basic SSE)
func AIQueryStream(c *gin.Context) {
    c.Writer.Header().Set("Content-Type", "text/event-stream")
    c.Writer.Header().Set("Cache-Control", "no-cache")
    c.Writer.Header().Set("Connection", "keep-alive")
    c.Writer.Flush()

    var req AIQueryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.String(http.StatusBadRequest, "data: %s\n\n", "invalid request")
        return
    }

    // For now, call normal query and emit single chunk
    routingResult, err := orchestrator.RouteRequestWithFallback(req.Prompt, req.TaskType, req.Provider, req.Model, nil)
    if err != nil {
        c.String(http.StatusInternalServerError, "data: %s\n\n", "routing error")
        return
    }
    aiResponse, err := makeAIRequest(routingResult.Provider, routingResult.Model, req)
    if err != nil {
        c.String(http.StatusInternalServerError, "data: %s\n\n", "provider error")
        return
    }
    // Emit once; future: incremental chunks from provider streaming APIs
    payload, _ := json.Marshal(map[string]interface{}{
        "content": aiResponse.Content,
        "provider": aiResponse.Provider,
        "model": aiResponse.Model,
    })
    c.Writer.Write([]byte("data: "))
    c.Writer.Write(payload)
    c.Writer.Write([]byte("\n\n"))
    c.Writer.Flush()
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

	// Initialize orchestrator components
	brain := orchestrator.NewAIBrain()
	cache := orchestrator.NewCache()
	metrics := orchestrator.NewMetricsCollector()

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
	orchestrator.CacheResponseResult(req.Prompt, req.TaskType, aiResponse.Provider, aiResponse.Model, aiResponse)

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
	adapter := orchestrator.GetProviderAdapter(provider)
	if adapter == nil {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := adapter.ChatCompletion(ctx, model, req.Prompt, req.MaxTokens, req.Options)
	if err != nil {
		return nil, err
	}
	return &models.AIResponse{
		Content:  resp.Content,
		Provider: resp.Provider,
		Model:    resp.Model,
		TokensUsed: models.TokenUsage{
			Input:  resp.PromptTokens,
			Output: resp.CompletionTokens,
		},
		Cost:     resp.Cost,
		Metadata: resp.Metadata,
	}, nil
}

// provider-specific direct implementations moved to adapters in orchestrator

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

    // Fetch API keys from database (secure: only hash stored)
    type apiKeyRow struct {
        ID        uuid.UUID
        UserID    uuid.UUID
        Name      string
        KeyHash   string
        IsActive  bool
        CreatedAt time.Time
        LastUsed  *time.Time
    }
    var apiKeys []map[string]interface{}
    query := "SELECT id, user_id, name, key_hash, is_active, created_at, last_used_at FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC"

    rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch API keys",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	// Scan rows into apiKeys slice
    for rows.Next() {
        var r apiKeyRow
        if err := rows.Scan(&r.ID, &r.UserID, &r.Name, &r.KeyHash, &r.IsActive, &r.CreatedAt, &r.LastUsed); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to scan API key",
				"details": err.Error(),
			})
			return
		}

        // expose only prefix/suffix derived from hash (not ideal, but enough for UI masking); typically we'd store prefix separately
        prefix := ""
        suffix := ""
        if len(r.KeyHash) >= 8 {
            prefix = r.KeyHash[:4]
            suffix = r.KeyHash[len(r.KeyHash)-4:]
        }
        apiKeys = append(apiKeys, map[string]interface{}{
            "id": r.ID,
            "user_id": r.UserID,
            "name": r.Name,
            "is_active": r.IsActive,
            "created_at": r.CreatedAt,
            "last_used": r.LastUsed,
            "prefix": prefix,
            "suffix": suffix,
        })
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to iterate API keys",
			"details": err.Error(),
		})
		return
	}

    c.JSON(http.StatusOK, gin.H{"api_keys": apiKeys})
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
    // Hash the key for storage
    sum := sha256.Sum256([]byte(key))
    keyHash := hex.EncodeToString(sum[:])

    id := uuid.New()
    // Insert API key into database (store hash only)
    query := `
        INSERT INTO api_keys (id, user_id, name, key_hash, is_active, created_at)
        VALUES ($1, $2, $3, $4, true, NOW())
    `
    _, err = database.DB.Exec(query, id, userID, req.Name, keyHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create API key",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "API key created successfully",
        "api_key": gin.H{"id": id, "name": req.Name, "key": key},
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

// RotateAPIKey disables old key and creates a new one with same name
func RotateAPIKey(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    keyID := c.Param("id")
    if keyID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "API key ID is required"})
        return
    }
    // Mark old key inactive (soft rotate)
    if _, err := database.DB.Exec("UPDATE api_keys SET is_active = false WHERE id = $1 AND user_id = $2", keyID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate old key"})
        return
    }
    // Create new key with same name suffix
    var name string
    _ = database.DB.QueryRow("SELECT name FROM api_keys WHERE id = $1 AND user_id = $2", keyID, userID).Scan(&name)
    newKey, err := generateAPIKey()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new key"})
        return
    }
    newID := uuid.New()
    _, err = database.DB.Exec(
        `INSERT INTO api_keys (id, user_id, name, key, permissions, is_active, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,true,NOW(),NOW())`,
        newID, userID, name, newKey, fmt.Sprintf("%v", []string{models.PermissionReadMetrics, models.PermissionWriteAPIKeys, models.PermissionReadAPIKeys}),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new key"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "API key rotated", "api_key": gin.H{"id": newID, "name": name, "key": newKey}})
}

// generateAPIKey generates a new random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
