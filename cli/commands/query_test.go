package commands

import (
	"testing"
)

func TestQueryCmd(t *testing.T) {
	// Test that the command is properly defined
	if QueryCmd.Use != "query [prompt]" {
		t.Errorf("Expected Use 'query [prompt]', got '%s'", QueryCmd.Use)
	}

	if QueryCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if QueryCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test that the command requires exactly one argument
	if QueryCmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}
