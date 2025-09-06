package commands

import (
	"testing"
)

func TestProvidersCmd(t *testing.T) {
	// Test that the command is properly defined
	if ProvidersCmd.Use != "providers" {
		t.Errorf("Expected Use 'providers', got '%s'", ProvidersCmd.Use)
	}

	if ProvidersCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if ProvidersCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}
