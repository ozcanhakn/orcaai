package commands

import (
	"testing"
)

func TestKeysCmd(t *testing.T) {
	// Test that the command is properly defined
	if KeysCmd.Use != "keys" {
		t.Errorf("Expected Use 'keys', got '%s'", KeysCmd.Use)
	}

	if KeysCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if KeysCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}

	// Test subcommands
	if len(KeysCmd.Commands()) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(KeysCmd.Commands()))
	}
}
