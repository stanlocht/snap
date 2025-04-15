package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize repository
	repo, err := Init(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Check if repository path is correct
	if repo.Path != tempDir {
		t.Errorf("Expected repository path to be %s, got %s", tempDir, repo.Path)
	}

	// Check if .snap directory was created
	snapDir := filepath.Join(tempDir, SnapDirName)
	if _, err := os.Stat(snapDir); os.IsNotExist(err) {
		t.Errorf(".snap directory was not created")
	}

	// Check if subdirectories were created
	subdirs := []string{
		filepath.Join(snapDir, "objects"),
		filepath.Join(snapDir, "refs"),
		filepath.Join(snapDir, "refs", "heads"),
		filepath.Join(snapDir, "issues"),
		filepath.Join(snapDir, "users"),
	}

	for _, dir := range subdirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s was not created", dir)
		}
	}

	// Check if HEAD file was created
	headFile := filepath.Join(snapDir, "HEAD")
	if _, err := os.Stat(headFile); os.IsNotExist(err) {
		t.Errorf("HEAD file was not created")
	}

	// Check if config file was created
	configFile := filepath.Join(snapDir, "config")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("config file was not created")
	}
}

func TestIsInitialized(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Check if directory is not initialized
	if IsInitialized(tempDir) {
		t.Errorf("Expected directory to not be initialized")
	}

	// Initialize repository
	_, err = Init(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Check if directory is initialized
	if !IsInitialized(tempDir) {
		t.Errorf("Expected directory to be initialized")
	}
}

func TestFind(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Initialize repository in the root directory
	_, err = Init(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Find repository from subdirectory
	repo, err := Find(subDir)
	if err != nil {
		t.Fatalf("Failed to find repository: %v", err)
	}

	// Check if repository path is correct
	if repo.Path != tempDir {
		t.Errorf("Expected repository path to be %s, got %s", tempDir, repo.Path)
	}

	// Create another temporary directory (not a repository)
	anotherTempDir, err := os.MkdirTemp("", "snap-test-another-")
	if err != nil {
		t.Fatalf("Failed to create another temp directory: %v", err)
	}
	defer os.RemoveAll(anotherTempDir)

	// Try to find repository in a directory that is not a repository
	_, err = Find(anotherTempDir)
	if err == nil {
		t.Errorf("Expected error when finding repository in a directory that is not a repository")
	}
}
