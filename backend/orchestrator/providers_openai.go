package orchestrator

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "database/sql"
    "os"
    "time"
    _ "github.com/lib/pq"
    "orcaai/backend/database"
    "orcaai/backend/utils"
)

type openAIAdapter struct{}

func (o *openAIAdapter) Name() string { return "openai" }

func (o *openAIAdapter) SupportsModel(model string) bool {
    // Basic check; could be expanded or fetched dynamically
    return model != ""
}

func (o *openAIAdapter) ChatCompletion(ctx context.Context, model string, prompt string, maxTokens int, options map[string]interface{}) (*AdapterResponse, error) {
    start := time.Now()
    body := map[string]interface{}{
        "model": model,
        "messages": []map[string]string{{"role": "user", "content": prompt}},
    }
    if maxTokens > 0 {
        body["max_tokens"] = maxTokens
    }
    for k, v := range options {
        // allow simple pass-through options; provider-specific filtering could be added
        body[k] = v
    }

    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    // Resolve API key: prefer encrypted key in DB, fallback to env
    apiKey := os.Getenv("OPENAI_API_KEY")
    if database.DB != nil {
        var enc sql.NullString
        _ = database.DB.QueryRow("SELECT api_key_encrypted FROM ai_providers WHERE name = $1", "openai").Scan(&enc)
        if enc.Valid {
            secret := os.Getenv("PROVIDER_SECRET_KEY")
            if secret != "" {
                if dec, err := utils.DecryptAESGCM(enc.String, secret); err == nil {
                    apiKey = string(dec)
                }
            }
        }
    }
    req.Header.Set("Authorization", "Bearer "+apiKey)

    httpClient := &http.Client{Timeout: 30 * time.Second}
    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var parsed struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
        Usage struct {
            PromptTokens     int `json:"prompt_tokens"`
            CompletionTokens int `json:"completion_tokens"`
        } `json:"usage"`
        Error struct {
            Message string `json:"message"`
        } `json:"error"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
        return nil, err
    }
    if parsed.Error.Message != "" {
        return nil, fmt.Errorf(parsed.Error.Message)
    }
    if len(parsed.Choices) == 0 {
        return nil, fmt.Errorf("no response from OpenAI")
    }

    // naive cost estimate; real pricing may be model-specific
    cost := float64(parsed.Usage.PromptTokens)*0.001 + float64(parsed.Usage.CompletionTokens)*0.002

    return &AdapterResponse{
        Content:          parsed.Choices[0].Message.Content,
        Provider:         o.Name(),
        Model:            model,
        PromptTokens:     parsed.Usage.PromptTokens,
        CompletionTokens: parsed.Usage.CompletionTokens,
        Cost:             cost,
        Metadata:         map[string]interface{}{},
        Latency:          time.Since(start),
    }, nil
}

func init() {
    RegisterProviderAdapter(&openAIAdapter{})
}


