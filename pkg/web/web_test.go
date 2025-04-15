package web

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestWelcomePages tests the welcome pages when the repository is empty
func TestWelcomePages(t *testing.T) {
	// Setup test repository
	repo, tempDir := setupTestRepo(t)
	defer cleanupTestRepo(tempDir)

	// Create server
	server, err := NewServer(repo)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route requests to the appropriate handler
		switch {
		case r.URL.Path == "/":
			server.handleHome(w, r)
		case r.URL.Path == "/commits":
			server.handleCommits(w, r)
		case r.URL.Path == "/issues":
			server.handleIssues(w, r)
		case r.URL.Path == "/users":
			server.handleUsers(w, r)
		case strings.HasPrefix(r.URL.Path, "/static/"):
			// Set up static file server
			staticContent, err := fs.Sub(staticFS, "static")
			if err != nil {
				t.Fatalf("Error setting up static file server: %v", err)
			}
			http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))).ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	// Test home page
	t.Run("HomePage", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/")
		if err != nil {
			t.Fatalf("Failed to get home page: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Home page returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Print the response body for debugging
		t.Logf("Home page response body: %s", string(body))

		// Check if response contains expected content
		expectedContent := []string{
			"Home",
		}

		for _, content := range expectedContent {
			if !strings.Contains(string(body), content) {
				t.Errorf("Home page does not contain expected content: %s", content)
			}
		}
	})

	// Test issues page
	t.Run("IssuesPage", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/issues")
		if err != nil {
			t.Fatalf("Failed to get issues page: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Issues page returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Print the response body for debugging
		t.Logf("Issues page response body: %s", string(body))

		// Check if response contains expected content
		expectedContent := []string{
			"Issues",
		}

		for _, content := range expectedContent {
			if !strings.Contains(string(body), content) {
				t.Errorf("Issues page does not contain expected content: %s", content)
			}
		}
	})

	// Test users page
	t.Run("UsersPage", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users")
		if err != nil {
			t.Fatalf("Failed to get users page: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Users page returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Print the response body for debugging
		t.Logf("Users page response body: %s", string(body))

		// Check if response contains expected content
		expectedContent := []string{
			"Users",
		}

		for _, content := range expectedContent {
			if !strings.Contains(string(body), content) {
				t.Errorf("Users page does not contain expected content: %s", content)
			}
		}
	})

	// Test static files
	t.Run("StaticFiles", func(t *testing.T) {
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
	})
}
