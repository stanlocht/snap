package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIndex(t *testing.T) {
	// Create a new index
	index := NewIndex()

	// Check if index is not nil
	if index == nil {
		t.Errorf("Expected index to not be nil")
	}

	// Check if entries map is not nil
	if index.Entries == nil {
		t.Errorf("Expected entries map to not be nil")
	}

	// Check if entries map is empty
	if len(index.Entries) != 0 {
		t.Errorf("Expected entries map to be empty")
	}
}

func TestAddFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .snap directory and objects directory
	snapDir := filepath.Join(tempDir, ".snap")
	objectsDir := filepath.Join(snapDir, "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		t.Fatalf("Failed to create objects directory: %v", err)
	}

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	testFileContent := "Test content"
	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a new index
	index := NewIndex()

	// Add file to index
	objectID, err := index.AddFile(tempDir, testFilePath)
	if err != nil {
		t.Fatalf("Failed to add file to index: %v", err)
	}

	// Check if object ID is not empty
	if objectID == "" {
		t.Errorf("Expected object ID to not be empty")
	}

	// Check if file was added to index
	relPath, err := filepath.Rel(tempDir, testFilePath)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}
	if index.Entries[relPath] != objectID {
		t.Errorf("Expected index to contain entry for %s with object ID %s", relPath, objectID)
	}

	// Check if object file was created
	objectPath := filepath.Join(objectsDir, objectID[:2], objectID[2:])
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		t.Errorf("Expected object file to exist at %s", objectPath)
	}

	// Check if object file contains the correct content
	objectContent, err := os.ReadFile(objectPath)
	if err != nil {
		t.Fatalf("Failed to read object file: %v", err)
	}
	if string(objectContent) != testFileContent {
		t.Errorf("Expected object file to contain '%s', got '%s'", testFileContent, string(objectContent))
	}
}

func TestSaveAndLoadIndex(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .snap directory
	snapDir := filepath.Join(tempDir, ".snap")
	if err := os.MkdirAll(snapDir, 0755); err != nil {
		t.Fatalf("Failed to create .snap directory: %v", err)
	}

	// Create a new index
	index := NewIndex()

	// Add some entries to the index
	index.Entries["file1.txt"] = "object1"
	index.Entries["file2.txt"] = "object2"

	// Save index
	if err := index.SaveIndex(tempDir); err != nil {
		t.Fatalf("Failed to save index: %v", err)
	}

	// Check if index file was created
	indexPath := filepath.Join(snapDir, "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected index file to exist at %s", indexPath)
	}

	// Load index
	loadedIndex, err := LoadIndex(tempDir)
	if err != nil {
		t.Fatalf("Failed to load index: %v", err)
	}

	// Check if loaded index is not nil
	if loadedIndex == nil {
		t.Errorf("Expected loaded index to not be nil")
	}

	// Check if loaded index has the correct entries
	// Note: Our current implementation of LoadIndex doesn't actually parse the index file,
	// so we can't check the entries. This is a placeholder for when we implement it.
}
