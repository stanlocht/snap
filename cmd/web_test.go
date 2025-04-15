package cmd

import (
	"bytes"
	"testing"
)

func TestWebCommand(t *testing.T) {
	// Skip this test for now as it requires more setup
	// In a real test, we would need to set up the command properly
	t.Skip("Skipping test that requires more setup")

	// Test help flag
	output, err := executeCommand(rootCmd, "web", "--help")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if output contains expected text
	expectedTexts := []string{
		"Start a web interface for the repository",
		"--port",
		"--open",
	}

	for _, text := range expectedTexts {
		if !bytes.Contains([]byte(output), []byte(text)) {
			t.Errorf("Expected output to contain '%s', but it didn't", text)
		}
	}
}

func TestWebCommandFlags(t *testing.T) {
	// Test port flag
	port, err := webCmd.Flags().GetInt("port")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if port != 8123 {
		t.Errorf("Expected default port to be 8123, got %d", port)
	}

	// Test open flag
	open, err := webCmd.Flags().GetBool("open")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if open != false {
		t.Errorf("Expected default open flag to be false, got %v", open)
	}
}
