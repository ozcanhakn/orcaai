package orchestrator

import (
    "context"
    "time"
)

// AIProviderAdapter defines a unified interface for AI providers
type AIProviderAdapter interface {
    // Name returns the provider identifier, e.g., "openai"
    Name() string
    // SupportsModel returns true if the provider supports the given model
    SupportsModel(model string) bool
    // ChatCompletion executes a chat/completion style request and returns content, token usage, cost and optional metadata
    ChatCompletion(ctx context.Context, model string, prompt string, maxTokens int, options map[string]interface{}) (*AdapterResponse, error)
}

// AdapterResponse is a normalized response from any provider adapter
type AdapterResponse struct {
    Content    string
    Provider   string
    Model      string
    PromptTokens int
    CompletionTokens int
    Cost       float64
    Metadata   map[string]interface{}
    Latency    time.Duration
}

var providerRegistry = map[string]AIProviderAdapter{}

// RegisterProviderAdapter registers a provider adapter by name
func RegisterProviderAdapter(adapter AIProviderAdapter) {
    if adapter == nil {
        return
    }
    providerRegistry[adapter.Name()] = adapter
}

// GetProviderAdapter retrieves a registered adapter by provider name
func GetProviderAdapter(name string) AIProviderAdapter {
    if adapter, ok := providerRegistry[name]; ok {
        return adapter
    }
    return nil
}


