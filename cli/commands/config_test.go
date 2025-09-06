package commands

import (
	"testing"
)

func TestConfigCmd(t *testing.T) {
	// Test that the command is properly defined
	if ConfigCmd.Use != "config" {
		t.Errorf("Expected Use 'config', got '%s'", ConfigCmd.Use)
	}

	if ConfigCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if ConfigCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test subcommands
	if len(ConfigCmd.Commands()) < 1 {
		t.Error("Expected at least 1 subcommand")
	}
}
