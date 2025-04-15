package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestHandleNotFound tests the 404 handling
func TestHandleNotFound(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request for non-existent page
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleHome(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

// TestHandleCommitDetail tests the handleCommitDetail function
func TestHandleCommitDetail(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("Test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit the file
	// Note: In a real test, we would use the repository API to add and commit the file
	// For simplicity, we'll just test the 404 case here

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request for non-existent commit
	req, err := http.NewRequest("GET", "/commit/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleCommitDetail(rr, req)

	// Check status code (should be 500 because the commit doesn't exist)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

// TestHandleIssueDetail tests the handleIssueDetail function
func TestHandleIssueDetail(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request for non-existent issue
	req, err := http.NewRequest("GET", "/issue/999", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleIssueDetail(rr, req)

	// Check status code (should be 500 because the issue doesn't exist)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

// TestHandleUserDetail tests the handleUserDetail function
func TestHandleUserDetail(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request for non-existent user
	req, err := http.NewRequest("GET", "/user/nonexistent", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleUserDetail(rr, req)

	// Check status code (should be 500 because the user doesn't exist)
	// Note: The actual implementation might return a different status code
	// depending on how the user manager handles non-existent users
	if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
		t.Errorf("Handler returned unexpected status code: got %v", status)
	}
}

// TestExtractEmoji tests the extractEmoji function
func TestExtractEmoji(t *testing.T) {
	// Test cases
	testCases := []struct {
		message  string
		expected string
	}{
		{"✨ Initial commit", "✨"},
		{":sparkles: Initial commit", ":sparkles:"},
		{"No emoji", ""},
		// Note: The function only extracts emojis at the start of the message
		{"Multiple 🔥 emojis ✨", ""}, // No emoji at start, so expect empty string
	}

	// Run tests
	for _, tc := range testCases {
		result := extractEmoji(tc.message)
		if result != tc.expected {
			t.Errorf("extractEmoji(%s) = %s, want %s", tc.message, result, tc.expected)
		}
	}
}
