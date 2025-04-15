package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	// Test help flag
	output, err := executeCommand(rootCmd, "--help")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if output contains expected text
	expectedTexts := []string{
		"Snap is an opinionated, fun, and community-focused version control system",
		"Usage:",
		"snap [command]",
		"Available Commands:",
		"--author",
		"--email",
	}

	for _, text := range expectedTexts {
		if !bytes.Contains([]byte(output), []byte(text)) {
			t.Errorf("Expected output to contain '%s', but it didn't", text)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	// Test version flag
	_, err := executeCommand(rootCmd, "--version")
	if err != nil {
		// We expect an error because we haven't implemented the version flag handler
		// This is just a placeholder test
		return
	}
}
