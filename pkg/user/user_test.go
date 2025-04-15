package user

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewUserManager(t *testing.T) {
	// Create a new user manager
	manager := NewUserManager("/path/to/repo")

	// Check if manager is not nil
	if manager == nil {
		t.Errorf("Expected manager to not be nil")
	}

	// Check if repo path is correct
	if manager.RepoPath != "/path/to/repo" {
		t.Errorf("Expected repo path to be /path/to/repo, got %s", manager.RepoPath)
	}
}

func TestGetUser(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new user manager
	manager := NewUserManager(tempDir)

	// Get a user that doesn't exist yet
	user, err := manager.GetUser("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check if user is not nil
	if user == nil {
		t.Errorf("Expected user to not be nil")
	}

	// Check if user name is correct
	if user.Name != "testuser" {
		t.Errorf("Expected user name to be 'testuser', got '%s'", user.Name)
	}

	// Check if user points is 0
	if user.Points != 0 {
		t.Errorf("Expected user points to be 0, got %d", user.Points)
	}

	// Check if user action log is empty
	if len(user.ActionLog) != 0 {
		t.Errorf("Expected user action log to be empty")
	}
}

func TestSaveUser(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new user manager
	manager := NewUserManager(tempDir)

	// Create a user
	user := &User{
		Name:      "testuser",
		Email:     "test@example.com",
		Points:    10,
		ActionLog: []ActionRecord{},
	}

	// Save user
	err = manager.SaveUser(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Check if user file was created
	userFilePath := filepath.Join(tempDir, ".snap", "users", "testuser.json")
	if _, err := os.Stat(userFilePath); os.IsNotExist(err) {
		t.Errorf("Expected user file to exist at %s", userFilePath)
	}

	// Get the user
	retrievedUser, err := manager.GetUser("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check if retrieved user is not nil
	if retrievedUser == nil {
		t.Errorf("Expected retrieved user to not be nil")
	}

	// Check if retrieved user name is correct
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected retrieved user name to be '%s', got '%s'", user.Name, retrievedUser.Name)
	}

	// Check if retrieved user email is correct
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected retrieved user email to be '%s', got '%s'", user.Email, retrievedUser.Email)
	}

	// Check if retrieved user points is correct
	if retrievedUser.Points != user.Points {
		t.Errorf("Expected retrieved user points to be %d, got %d", user.Points, retrievedUser.Points)
	}
}

func TestRecordAction(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new user manager
	manager := NewUserManager(tempDir)

	// Record an action
	err = manager.RecordAction("testuser", ActionCommit, "Test commit", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("Failed to record action: %v", err)
	}

	// Get the user
	user, err := manager.GetUser("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check if user points is correct
	expectedPoints := PointValues[ActionCommit]
	if user.Points != expectedPoints {
		t.Errorf("Expected user points to be %d, got %d", expectedPoints, user.Points)
	}

	// Check if user action log has one entry
	if len(user.ActionLog) != 1 {
		t.Errorf("Expected user action log to have 1 entry, got %d", len(user.ActionLog))
	}

	// Check if action log entry is correct
	if user.ActionLog[0].Action != ActionCommit {
		t.Errorf("Expected action to be %s, got %s", ActionCommit, user.ActionLog[0].Action)
	}
	if user.ActionLog[0].Points != expectedPoints {
		t.Errorf("Expected points to be %d, got %d", expectedPoints, user.ActionLog[0].Points)
	}
	if user.ActionLog[0].Description != "Test commit" {
		t.Errorf("Expected description to be 'Test commit', got '%s'", user.ActionLog[0].Description)
	}
	if user.ActionLog[0].Timestamp != "2025-01-01T00:00:00Z" {
		t.Errorf("Expected timestamp to be '2025-01-01T00:00:00Z', got '%s'", user.ActionLog[0].Timestamp)
	}

	// Record another action
	err = manager.RecordAction("testuser", ActionIssueClose, "Closed issue", "2025-01-02T00:00:00Z")
	if err != nil {
		t.Fatalf("Failed to record action: %v", err)
	}

	// Get the user again
	user, err = manager.GetUser("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Check if user points is correct
	expectedPoints += PointValues[ActionIssueClose]
	if user.Points != expectedPoints {
		t.Errorf("Expected user points to be %d, got %d", expectedPoints, user.Points)
	}

	// Check if user action log has two entries
	if len(user.ActionLog) != 2 {
		t.Errorf("Expected user action log to have 2 entries, got %d", len(user.ActionLog))
	}
}

func TestGetLeaderboard(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "snap-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new user manager
	manager := NewUserManager(tempDir)

	// Create users directory
	usersDir := filepath.Join(tempDir, ".snap", "users")
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		t.Fatalf("Failed to create users directory: %v", err)
	}

	// Create user1
	user1 := &User{
		Name:   "user1",
		Points: 10,
		ActionLog: []ActionRecord{{
			Action:      ActionCommit,
			Points:      10,
			Description: "Commit 1",
			Timestamp:   "2025-01-01T00:00:00Z",
		}},
		Commits: 1,
	}

	// Create user2
	user2 := &User{
		Name:   "user2",
		Points: 25,
		ActionLog: []ActionRecord{{
			Action:      ActionCommit,
			Points:      10,
			Description: "Commit 2",
			Timestamp:   "2025-01-02T00:00:00Z",
		}, {
			Action:      ActionIssueClose,
			Points:      15,
			Description: "Closed issue",
			Timestamp:   "2025-01-03T00:00:00Z",
		}},
		Commits:      1,
		IssuesClosed: 1,
	}

	// Save users
	if err := manager.SaveUser(user1); err != nil {
		t.Fatalf("Failed to save user1: %v", err)
	}

	if err := manager.SaveUser(user2); err != nil {
		t.Fatalf("Failed to save user2: %v", err)
	}

	// Get leaderboard
	leaderboard, err := manager.GetLeaderboard()
	if err != nil {
		t.Fatalf("Failed to get leaderboard: %v", err)
	}

	// Check if leaderboard has 2 users
	if len(leaderboard) != 2 {
		t.Errorf("Expected leaderboard to have 2 users, got %d", len(leaderboard))
		return // Avoid panic if leaderboard is empty
	}

	// Check if users are sorted by points (descending)
	if leaderboard[0].Points < leaderboard[1].Points {
		t.Errorf("Expected users to be sorted by points (descending)")
	}

	// Check if user2 is first (has more points)
	if leaderboard[0].Name != "user2" {
		t.Errorf("Expected user2 to be first, got %s", leaderboard[0].Name)
	}

	// Check if user1 is second
	if leaderboard[1].Name != "user1" {
		t.Errorf("Expected user1 to be second, got %s", leaderboard[1].Name)
	}
}
