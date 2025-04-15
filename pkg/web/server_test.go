package web

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stanlocht/snap/pkg/repository"
)

// setupTestRepo creates a temporary repository for testing
func setupTestRepo(t *testing.T) (*repository.Repository, string) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "snap-web-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize repository
	repo, err := repository.Init(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	return repo, tempDir
}

// cleanupTestRepo removes the temporary repository
func cleanupTestRepo(tempDir string) {
	os.RemoveAll(tempDir)
}

// TestNewServer tests the NewServer function
func TestNewServer(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Check if server was created
	if server == nil {
		t.Errorf("Expected server to not be nil")
	}

	// Check if templates were parsed
	if server.Templates == nil {
		t.Errorf("Expected templates to not be nil")
	}

	// Check if repository was set
	if server.Repo != repo {
		t.Errorf("Expected server.Repo to be %v, got %v", repo, server.Repo)
	}
}

// TestHandleHome tests the handleHome function
func TestHandleHome(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleHome(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response contains expected content
	if !strings.Contains(rr.Body.String(), "Home") {
		t.Errorf("Handler response does not contain 'Home'")
	}
}

// TestHandleCommits tests the handleCommits function
func TestHandleCommits(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/commits", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleCommits(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response contains expected content
	if !strings.Contains(rr.Body.String(), "Commits") {
		t.Errorf("Handler response does not contain 'Commits'")
	}
}

// TestHandleIssues tests the handleIssues function
func TestHandleIssues(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/issues", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleIssues(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response contains expected content
	if !strings.Contains(rr.Body.String(), "Issues") {
		t.Errorf("Handler response does not contain 'Issues'")
	}
}

// TestHandleUsers tests the handleUsers function
func TestHandleUsers(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	server.handleUsers(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if response contains expected content
	if !strings.Contains(rr.Body.String(), "Users") {
		t.Errorf("Handler response does not contain 'Users'")
	}
}

// TestStaticFiles tests the static file server
func TestStaticFiles(t *testing.T) {
	// Setup test repository
	_, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Start test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set up static file server
		staticContent, err := fs.Sub(staticFS, "static")
		if err != nil {
			t.Fatalf("Error setting up static file server: %v", err)
		}

		// Handle static files
		http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))).ServeHTTP(w, r)
	}))
	defer ts.Close()

	// Test CSS file
	resp, err := http.Get(ts.URL + "/static/style.css")
	if err != nil {
		t.Fatalf("Failed to get CSS file: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Static file server returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check if response contains expected content
	if !strings.Contains(string(body), "Base styles") {
		t.Errorf("Static file response does not contain expected content")
	}
}

// TestFormatTime tests the formatTime function
func TestFormatTime(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected string
	}{
		{"2025-01-01T00:00:00Z", "2025-01-01 00:00:00"},
		{"2025-12-31T23:59:59Z", "2025-12-31 23:59:59"},
	}

	// Run tests
	for _, tc := range testCases {
		// Parse input time
		inputTime, err := time.Parse(time.RFC3339, tc.input)
		if err != nil {
			t.Fatalf("Failed to parse input time: %v", err)
		}

		// Format time
		result := formatTime(inputTime)

		// Check result
		if result != tc.expected {
			t.Errorf("formatTime(%s) = %s, want %s", tc.input, result, tc.expected)
		}
	}
}

// TestTruncateID tests the truncateID function
func TestTruncateID(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected string
	}{
		{"1234567890abcdef", "1234567"},
		{"1234567", "1234567"},
		{"123456", "123456"},
	}

	// Run tests
	for _, tc := range testCases {
		result := truncateID(tc.input)
		if result != tc.expected {
			t.Errorf("truncateID(%s) = %s, want %s", tc.input, result, tc.expected)
		}
	}
}
