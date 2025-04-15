package issue

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewIssueManager(t *testing.T) {
	// Create a new issue manager
	manager := NewIssueManager("/path/to/repo")

	// Check if manager is not nil
	if manager == nil {
		t.Errorf("Expected manager to not be nil")
	}

	// Check if repo path is correct
	if manager.RepoPath != "/path/to/repo" {
		t.Errorf("Expected repo path to be /path/to/repo, got %s", manager.RepoPath)
	}
}

func TestCreateAndGetIssue(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new issue manager
	manager := NewIssueManager(tempDir)

	// Create an issue
	issue, err := manager.CreateIssue("Test Issue", "This is a test issue", "testuser")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Check if issue is not nil
	if issue == nil {
		t.Errorf("Expected issue to not be nil")
	}

	// Check if issue ID is 1
	if issue.ID != 1 {
		t.Errorf("Expected issue ID to be 1, got %d", issue.ID)
	}

	// Check if issue title is correct
	if issue.Title != "Test Issue" {
		t.Errorf("Expected issue title to be 'Test Issue', got '%s'", issue.Title)
	}

	// Check if issue description is correct
	if issue.Description != "This is a test issue" {
		t.Errorf("Expected issue description to be 'This is a test issue', got '%s'", issue.Description)
	}

	// Check if issue status is open
	if issue.Status != StatusOpen {
		t.Errorf("Expected issue status to be open, got %s", issue.Status)
	}

	// Check if issue created by is correct
	if issue.CreatedBy != "testuser" {
		t.Errorf("Expected issue created by to be 'testuser', got '%s'", issue.CreatedBy)
	}

	// Check if issue file was created
	issueFilePath := filepath.Join(tempDir, ".snap", "issues", "1.json")
	if _, err := os.Stat(issueFilePath); os.IsNotExist(err) {
		t.Errorf("Expected issue file to exist at %s", issueFilePath)
	}

	// Get the issue
	retrievedIssue, err := manager.GetIssue(1)
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	// Check if retrieved issue is not nil
	if retrievedIssue == nil {
		t.Errorf("Expected retrieved issue to not be nil")
	}

	// Check if retrieved issue ID is correct
	if retrievedIssue.ID != issue.ID {
		t.Errorf("Expected retrieved issue ID to be %d, got %d", issue.ID, retrievedIssue.ID)
	}

	// Check if retrieved issue title is correct
	if retrievedIssue.Title != issue.Title {
		t.Errorf("Expected retrieved issue title to be '%s', got '%s'", issue.Title, retrievedIssue.Title)
	}

	// Check if retrieved issue description is correct
	if retrievedIssue.Description != issue.Description {
		t.Errorf("Expected retrieved issue description to be '%s', got '%s'", issue.Description, retrievedIssue.Description)
	}

	// Check if retrieved issue status is correct
	if retrievedIssue.Status != issue.Status {
		t.Errorf("Expected retrieved issue status to be %s, got %s", issue.Status, retrievedIssue.Status)
	}

	// Check if retrieved issue created by is correct
	if retrievedIssue.CreatedBy != issue.CreatedBy {
		t.Errorf("Expected retrieved issue created by to be '%s', got '%s'", issue.CreatedBy, retrievedIssue.CreatedBy)
	}
}

func TestListIssues(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new issue manager
	manager := NewIssueManager(tempDir)

	// Create some issues
	_, err = manager.CreateIssue("Issue 1", "Description 1", "user1")
	if err != nil {
		t.Fatalf("Failed to create issue 1: %v", err)
	}

	_, err = manager.CreateIssue("Issue 2", "Description 2", "user2")
	if err != nil {
		t.Fatalf("Failed to create issue 2: %v", err)
	}

	// Close issue 2
	err = manager.CloseIssue(2)
	if err != nil {
		t.Fatalf("Failed to close issue 2: %v", err)
	}

	// List all issues
	issues, err := manager.ListIssues(true)
	if err != nil {
		t.Fatalf("Failed to list issues: %v", err)
	}

	// Check if we got 2 issues
	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	// List only open issues
	openIssues, err := manager.ListIssues(false)
	if err != nil {
		t.Fatalf("Failed to list open issues: %v", err)
	}

	// Check if we got 1 open issue
	if len(openIssues) != 1 {
		t.Errorf("Expected 1 open issue, got %d", len(openIssues))
	}

	// Check if the open issue is issue 1
	if openIssues[0].ID != 1 {
		t.Errorf("Expected open issue ID to be 1, got %d", openIssues[0].ID)
	}
}

func TestCloseIssue(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new issue manager
	manager := NewIssueManager(tempDir)

	// Create an issue
	_, err = manager.CreateIssue("Test Issue", "This is a test issue", "testuser")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Close the issue
	err = manager.CloseIssue(1)
	if err != nil {
		t.Fatalf("Failed to close issue: %v", err)
	}

	// Get the issue
	issue, err := manager.GetIssue(1)
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	// Check if issue status is closed
	if issue.Status != StatusClosed {
		t.Errorf("Expected issue status to be closed, got %s", issue.Status)
	}

	// Check if closed at is set
	if issue.ClosedAt.IsZero() {
		t.Errorf("Expected closed at to be set")
	}

	// Check if closed at is within the last minute
	if time.Since(issue.ClosedAt) > time.Minute {
		t.Errorf("Expected closed at to be within the last minute")
	}
}

func TestAssignIssue(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new issue manager
	manager := NewIssueManager(tempDir)

	// Create an issue
	_, err = manager.CreateIssue("Test Issue", "This is a test issue", "testuser")
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Assign the issue
	err = manager.AssignIssue(1, "assignee")
	if err != nil {
		t.Fatalf("Failed to assign issue: %v", err)
	}

	// Get the issue
	issue, err := manager.GetIssue(1)
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	// Check if assigned to is correct
	if issue.AssignedTo != "assignee" {
		t.Errorf("Expected assigned to to be 'assignee', got '%s'", issue.AssignedTo)
	}
}
