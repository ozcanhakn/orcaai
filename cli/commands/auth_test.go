package commands

import (
	"testing"
)

func TestAuthCmd(t *testing.T) {
	// Test that the command is properly defined
	if AuthCmd.Use != "auth" {
		t.Errorf("Expected Use 'auth', got '%s'", AuthCmd.Use)
	}

	if AuthCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if AuthCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test subcommands
	if len(AuthCmd.Commands()) < 1 {
		t.Error("Expected at least 1 subcommand")
	}
}
