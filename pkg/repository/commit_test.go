package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndGetCommit(t *testing.T) {
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

	// Create a tree
	tree := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
			"file2.txt": "object2",
		},
	}

	// Create a commit
	commit, err := repo.CreateCommit("✨ Initial commit", "testuser", "test@example.com", tree)
	if err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Check if commit is not nil
	if commit == nil {
		t.Errorf("Expected commit to not be nil")
	}

	// Check if commit ID is not empty
	if commit.ID == "" {
		t.Errorf("Expected commit ID to not be empty")
	}

	// Check if commit message is correct
	if commit.Message != "✨ Initial commit" {
		t.Errorf("Expected commit message to be '✨ Initial commit', got '%s'", commit.Message)
	}

	// Check if commit author is correct
	if commit.Author != "testuser" {
		t.Errorf("Expected commit author to be 'testuser', got '%s'", commit.Author)
	}

	// Check if commit email is correct
	if commit.Email != "test@example.com" {
		t.Errorf("Expected commit email to be 'test@example.com', got '%s'", commit.Email)
	}

	// Check if commit parent ID is empty (first commit)
	if commit.ParentID != "" {
		t.Errorf("Expected commit parent ID to be empty, got '%s'", commit.ParentID)
	}

	// Check if commit tree ID is not empty
	if commit.TreeID == "" {
		t.Errorf("Expected commit tree ID to not be empty")
	}

	// Get the commit
	retrievedCommit, err := repo.GetCommit(commit.ID)
	if err != nil {
		t.Fatalf("Failed to get commit: %v", err)
	}

	// Check if retrieved commit is not nil
	if retrievedCommit == nil {
		t.Errorf("Expected retrieved commit to not be nil")
	}

	// Check if retrieved commit ID is correct
	if retrievedCommit.ID != commit.ID {
		t.Errorf("Expected retrieved commit ID to be '%s', got '%s'", commit.ID, retrievedCommit.ID)
	}

	// Check if retrieved commit message is correct
	if retrievedCommit.Message != commit.Message {
		t.Errorf("Expected retrieved commit message to be '%s', got '%s'", commit.Message, retrievedCommit.Message)
	}

	// Check if retrieved commit author is correct
	if retrievedCommit.Author != commit.Author {
		t.Errorf("Expected retrieved commit author to be '%s', got '%s'", commit.Author, retrievedCommit.Author)
	}

	// Check if retrieved commit email is correct
	if retrievedCommit.Email != commit.Email {
		t.Errorf("Expected retrieved commit email to be '%s', got '%s'", commit.Email, retrievedCommit.Email)
	}

	// Check if retrieved commit parent ID is correct
	if retrievedCommit.ParentID != commit.ParentID {
		t.Errorf("Expected retrieved commit parent ID to be '%s', got '%s'", commit.ParentID, retrievedCommit.ParentID)
	}

	// Check if retrieved commit tree ID is correct
	if retrievedCommit.TreeID != commit.TreeID {
		t.Errorf("Expected retrieved commit tree ID to be '%s', got '%s'", commit.TreeID, retrievedCommit.TreeID)
	}
}

func TestSaveAndGetTree(t *testing.T) {
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

	// Create a tree
	tree := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
			"file2.txt": "object2",
		},
	}

	// Save the tree
	treeID, err := repo.SaveTree(tree)
	if err != nil {
		t.Fatalf("Failed to save tree: %v", err)
	}

	// Check if tree ID is not empty
	if treeID == "" {
		t.Errorf("Expected tree ID to not be empty")
	}

	// Get the tree
	retrievedTree, err := repo.GetTree(treeID)
	if err != nil {
		t.Fatalf("Failed to get tree: %v", err)
	}

	// Check if retrieved tree is not nil
	if retrievedTree == nil {
		t.Errorf("Expected retrieved tree to not be nil")
	}

	// Check if retrieved tree has the correct entries
	if len(retrievedTree.Entries) != len(tree.Entries) {
		t.Errorf("Expected retrieved tree to have %d entries, got %d", len(tree.Entries), len(retrievedTree.Entries))
	}

	for path, objectID := range tree.Entries {
		if retrievedObjectID, ok := retrievedTree.Entries[path]; !ok {
			t.Errorf("Expected retrieved tree to have entry for '%s'", path)
		} else if retrievedObjectID != objectID {
			t.Errorf("Expected retrieved tree entry for '%s' to be '%s', got '%s'", path, objectID, retrievedObjectID)
		}
	}
}

func TestGetCommitHistory(t *testing.T) {
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

	// Create a tree
	tree1 := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
		},
	}

	// Create first commit
	commit1, err := repo.CreateCommit("✨ Initial commit", "testuser", "test@example.com", tree1)
	if err != nil {
		t.Fatalf("Failed to create first commit: %v", err)
	}

	// Create a second tree
	tree2 := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
			"file2.txt": "object2",
		},
	}

	// Create second commit
	commit2, err := repo.CreateCommit("📚 Add file2.txt", "testuser", "test@example.com", tree2)
	if err != nil {
		t.Fatalf("Failed to create second commit: %v", err)
	}

	// Get commit history
	history, err := repo.GetCommitHistory("")
	if err != nil {
		t.Fatalf("Failed to get commit history: %v", err)
	}

	// Check if history has 2 commits
	if len(history) != 2 {
		t.Errorf("Expected history to have 2 commits, got %d", len(history))
	}

	// Check if history is in the correct order (newest first)
	if len(history) >= 2 {
		if history[0].ID != commit2.ID {
			t.Errorf("Expected first commit in history to be '%s', got '%s'", commit2.ID, history[0].ID)
		}

		if history[1].ID != commit1.ID {
			t.Errorf("Expected second commit in history to be '%s', got '%s'", commit1.ID, history[1].ID)
		}
	}
}

func TestUndoLastCommit(t *testing.T) {
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

	// Create a tree
	tree1 := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
		},
	}

	// Create first commit
	commit1, err := repo.CreateCommit("✨ Initial commit", "testuser", "test@example.com", tree1)
	if err != nil {
		t.Fatalf("Failed to create first commit: %v", err)
	}

	// Create a second tree
	tree2 := &Tree{
		Entries: map[string]string{
			"file1.txt": "object1",
			"file2.txt": "object2",
		},
	}

	// Create second commit
	_, err = repo.CreateCommit("📚 Add file2.txt", "testuser", "test@example.com", tree2)
	if err != nil {
		t.Fatalf("Failed to create second commit: %v", err)
	}

	// Get current HEAD commit ID
	currentCommitID, err := repo.GetHEADCommitID()
	if err != nil {
		t.Fatalf("Failed to get HEAD commit ID: %v", err)
	}

	// Undo the last commit
	err = repo.UndoLastCommit()
	if err != nil {
		t.Fatalf("Failed to undo last commit: %v", err)
	}

	// Get new HEAD commit ID
	newCommitID, err := repo.GetHEADCommitID()
	if err != nil {
		t.Fatalf("Failed to get HEAD commit ID after undo: %v", err)
	}

	// Check if HEAD now points to the first commit
	if newCommitID != commit1.ID {
		t.Errorf("Expected HEAD to point to first commit '%s' after undo, got '%s'", commit1.ID, newCommitID)
	}

	// Check if HEAD no longer points to the second commit
	if newCommitID == currentCommitID {
		t.Errorf("Expected HEAD to no longer point to second commit '%s' after undo", currentCommitID)
	}
}

func TestGetStatus(t *testing.T) {
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

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	testFileContent := "Test content"
	if err := os.WriteFile(testFilePath, []byte(testFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create index file with the test file
	indexPath := filepath.Join(tempDir, ".snap", "index")
	indexContent := "abcdef1234567890 test.txt\n"
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	// Create a tree
	tree := &Tree{
		Entries: map[string]string{
			"test.txt": "abcdef1234567890",
		},
	}

	// Save the tree
	_, err = repo.SaveTree(tree)
	if err != nil {
		t.Fatalf("Failed to save tree: %v", err)
	}

	// Create a commit with the tree
	_, err = repo.CreateCommit("✨ Initial commit", "testuser", "test@example.com", tree)
	if err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Get repository status
	status, branch, err := repo.GetStatus()
	if err != nil {
		t.Fatalf("Failed to get repository status: %v", err)
	}

	// Check if branch is correct
	if branch != "master" {
		t.Errorf("Expected branch to be 'master', got '%s'", branch)
	}

	// Check if status is empty (working tree clean)
	if len(status) != 0 {
		t.Errorf("Expected status to be empty, got %d entries", len(status))
	}

	// Modify the index to simulate a modified file
	indexContent = "fedcba0987654321 test.txt\n"
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		t.Fatalf("Failed to modify index file: %v", err)
	}

	// Get repository status again
	status, _, err = repo.GetStatus()
	if err != nil {
		t.Fatalf("Failed to get repository status after modification: %v", err)
	}

	// Check if status has one entry for the modified file
	if len(status) != 1 {
		t.Errorf("Expected status to have 1 entry after modification, got %d", len(status))
	} else {
		if status[0].Path != "test.txt" {
			t.Errorf("Expected modified file path to be 'test.txt', got '%s'", status[0].Path)
		}
		if status[0].Status != "modified" {
			t.Errorf("Expected file status to be 'modified', got '%s'", status[0].Status)
		}
	}
}
