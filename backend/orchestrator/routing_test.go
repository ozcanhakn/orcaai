package orchestrator

import (
	"testing"
	"time"
)

func TestSelectBestProvider(t *testing.T) {
	// Reset provider status for testing
	providerStatus = make(map[string]*ProviderStatus)

	// Add test providers
	Providers = []Provider{
		{
			Name:      "openai",
			Model:     "gpt-3.5-turbo",
			Cost:      0.002,
			Latency:   500 * time.Millisecond,
			Quality:   0.9,
			Available: true,
		},
		{
			Name:      "claude",
			Model:     "claude-2",
			Cost:      0.001,
			Latency:   800 * time.Millisecond,
			Quality:   0.85,
			Available: true,
		},
		{
			Name:      "gemini",
			Model:     "gemini-pro",
			Cost:      0.0005,
			Latency:   1000 * time.Millisecond,
			Quality:   0.8,
			Available: true,
		},
	}

	// Test selection based on different priorities
	tests := []struct {
		name         string
		taskType     string
		priority     string
		expectedName string
	}{
		{
			name:         "Cost priority",
			taskType:     "text-generation",
			priority:     "cost",
			expectedName: "gemini",
		},
		{
			name:         "Latency priority",
			taskType:     "text-generation",
			priority:     "latency",
			expectedName: "openai",
		},
		{
			name:         "Quality priority",
			taskType:     "text-generation",
			priority:     "quality",
			expectedName: "openai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectBestProvider(tt.taskType, tt.priority)
			if result.Name != tt.expectedName {
				t.Errorf("Expected provider '%s', got '%s'", tt.expectedName, result.Name)
			}
		})
	}
}

func TestRouteRequest(t *testing.T) {
	// Reset provider status for testing
	providerStatus = make(map[string]*ProviderStatus)

	// Add test providers
	Providers = []Provider{
		{
			Name:      "openai",
			Model:     "gpt-3.5-turbo",
			Cost:      0.002,
			Latency:   500 * time.Millisecond,
			Quality:   0.9,
			Available: true,
		},
		{
			Name:      "claude",
			Model:     "claude-2",
			Cost:      0.001,
			Latency:   800 * time.Millisecond,
			Quality:   0.85,
			Available: true,
		},
	}

	result, err := RouteRequest("Test prompt", "text-generation", "", "")
	if err != nil {
		t.Fatalf("RouteRequest failed: %v", err)
	}

	if result.Provider == "" {
		t.Error("Expected provider to be set")
	}

	if result.Model == "" {
		t.Error("Expected model to be set")
	}

	if result.Reasoning == "" {
		t.Error("Expected reasoning to be set")
	}
}

func TestUpdateProviderStatus(t *testing.T) {
	// Reset provider status for testing
	providerStatus = make(map[string]*ProviderStatus)

	// Test updating provider status
	UpdateProviderStatus("openai", "gpt-3.5-turbo", true, nil)

	status := GetProviderStatus("openai", "gpt-3.5-turbo")
	if status == nil {
		t.Fatal("Expected provider status to be created")
	}

	if !status.Healthy {
		t.Error("Expected provider to be healthy")
	}

	if status.LastChecked.IsZero() {
		t.Error("Expected LastChecked to be set")
	}

	// Test updating with error
	testErr := &ProviderError{Message: "Test error"}
	UpdateProviderStatus("openai", "gpt-3.5-turbo", false, testErr)

	status = GetProviderStatus("openai", "gpt-3.5-turbo")
	if status.Healthy {
		t.Error("Expected provider to be unhealthy")
	}

	if status.LastError == nil {
		t.Error("Expected LastError to be set")
	}
}
