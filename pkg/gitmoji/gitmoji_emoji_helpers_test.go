package gitmoji

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetNumberedGitmojiList(t *testing.T) {
	// Get numbered gitmoji list
	list := GetNumberedGitmojiList()

	// Check if list is not empty
	if list == "" {
		t.Errorf("Expected numbered gitmoji list to not be empty")
	}

	// Check if list contains "Select a Gitmoji by number:"
	if !strings.Contains(list, "Select a Gitmoji by number:") {
		t.Errorf("Expected numbered gitmoji list to contain 'Select a Gitmoji by number:'")
	}

	// Check if list contains numbered entries
	for i := 1; i <= len(Gitmojis); i++ {
		if !strings.Contains(list, fmt.Sprintf("%d.", i)) {
			t.Errorf("Expected numbered gitmoji list to contain entry number %d", i)
		}
	}
}

func TestGetGitmojiByNumber(t *testing.T) {
	// Test cases
	testCases := []struct {
		number      int
		expectError bool
	}{
		{0, true},                 // Invalid number (too low)
		{len(Gitmojis) + 1, true}, // Invalid number (too high)
		{1, false},                // Valid number (first)
		{len(Gitmojis), false},    // Valid number (last)
	}

	// Run tests
	for _, tc := range testCases {
		gitmoji, err := GetGitmojiByNumber(tc.number)
		if tc.expectError && err == nil {
			t.Errorf("Expected error for number %d, but got none", tc.number)
		} else if !tc.expectError && err != nil {
			t.Errorf("Expected no error for number %d, but got: %v", tc.number, err)
		}

		if !tc.expectError {
			// Check if gitmoji is correct
			if tc.number >= 1 && tc.number <= len(Gitmojis) {
				expectedGitmoji := Gitmojis[tc.number-1]
				if gitmoji.Emoji != expectedGitmoji.Emoji {
					t.Errorf("Expected emoji to be '%s', got '%s'", expectedGitmoji.Emoji, gitmoji.Emoji)
				}
				if gitmoji.Code != expectedGitmoji.Code {
					t.Errorf("Expected code to be '%s', got '%s'", expectedGitmoji.Code, gitmoji.Code)
				}
				if gitmoji.Description != expectedGitmoji.Description {
					t.Errorf("Expected description to be '%s', got '%s'", expectedGitmoji.Description, gitmoji.Description)
				}
			}
		}
	}
}

func TestAutoConvertKeywordsToEmoji(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},                     // Empty message
		{"No keyword", "No keyword"}, // No keyword
		{"feature: Add new feature", ":sparkles: Add new feature"},                 // feature: keyword
		{"feat: Add new feature", ":sparkles: Add new feature"},                    // feat: keyword
		{"fix: Fix bug", ":bug: Fix bug"},                                          // fix: keyword
		{"docs: Update docs", ":books: Update docs"},                               // docs: keyword
		{"refactor: Refactor code", ":recycle: Refactor code"},                     // refactor: keyword
		{"FEATURE: Add new feature", ":sparkles: Add new feature"},                 // Case insensitive
		{"Feature: Add new feature", ":sparkles: Add new feature"},                 // Case insensitive
		{"✨ Already has emoji", "✨ Already has emoji"},                             // Already has emoji
		{":sparkles: Already has emoji code", ":sparkles: Already has emoji code"}, // Already has emoji code
	}

	// Run tests
	for _, tc := range testCases {
		result := AutoConvertKeywordsToEmoji(tc.input)
		if result != tc.expected {
			t.Errorf("Expected '%s' to convert to '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}
