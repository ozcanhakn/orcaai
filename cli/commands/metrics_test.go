package commands

import (
	"testing"
)

func TestMetricsCmd(t *testing.T) {
	// Test that the command is properly defined
	if MetricsCmd.Use != "metrics" {
		t.Errorf("Expected Use 'metrics', got '%s'", MetricsCmd.Use)
	}

	if MetricsCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if MetricsCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}
