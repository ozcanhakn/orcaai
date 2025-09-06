package commands

import (
	"testing"
)

func TestVersionCmd(t *testing.T) {
	// Test that the command is properly defined
	if VersionCmd.Use != "version" {
		t.Errorf("Expected Use 'version', got '%s'", VersionCmd.Use)
	}

	if VersionCmd.Short == "" {
		t.Error("Expected Short description to be non-empty")
	}

	if VersionCmd.Long == "" {
		t.Error("Expected Long description to be non-empty")
	}
}
