package snapmoji

import (
	"strings"
	"testing"
)

func TestValidateCommitMessage(t *testing.T) {
	// Test cases
	testCases := []struct {
		message string
		valid   bool
	}{
		{"", false},                                // Empty message
		{"Initial commit", false},                  // No snapmoji
		{"✨ Initial commit", true},                // With emoji
		{":sparkles: Initial commit", true},        // With emoji code
		{"✨Initial commit", true},                 // No space after emoji
		{":sparkles:Initial commit", true},         // No space after emoji code
		{"Some text ✨ Initial commit", false},     // Emoji not at start
		{"Some text :sparkles: Initial commit", false}, // Emoji code not at start
	}

	// Run tests
	for _, tc := range testCases {
		err := ValidateCommitMessage(tc.message)
		if tc.valid && err != nil {
			t.Errorf("Expected message '%s' to be valid, but got error: %v", tc.message, err)
		} else if !tc.valid && err == nil {
			t.Errorf("Expected message '%s' to be invalid, but got no error", tc.message)
		}
	}
}

func TestGetSnapmojiList(t *testing.T) {
	// Get snapmoji list
	list := GetSnapmojiList()

	// Check if list is not empty
	if list == "" {
		t.Errorf("Expected snapmoji list to not be empty")
	}

	// Check if list contains "Available Snapmojis:"
	if !strings.Contains(list, "Available Snapmojis:") {
		t.Errorf("Expected snapmoji list to contain 'Available Snapmojis:'")
	}

	// Check if list contains at least one snapmoji
	if !strings.Contains(list, "✨") {
		t.Errorf("Expected snapmoji list to contain at least one snapmoji")
	}

	// Check if list contains at least one emoji code
	if !strings.Contains(list, ":sparkles:") {
		t.Errorf("Expected snapmoji list to contain at least one emoji code")
	}

	// Check if list contains at least one description
	if !strings.Contains(list, "Introduce new features") {
		t.Errorf("Expected snapmoji list to contain at least one description")
	}
}
