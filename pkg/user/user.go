package user

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Action represents a user action that earns points
type Action string

const (
	// ActionCommit represents a commit action
	ActionCommit Action = "commit"
	// ActionIssueCreate represents an issue creation action
	ActionIssueCreate Action = "issue_create"
	// ActionIssueClose represents an issue closure action
	ActionIssueClose Action = "issue_close"
	// ActionIssueAssign represents an issue assignment action
	ActionIssueAssign Action = "issue_assign"
)

// PointValues defines the point values for different actions
var PointValues = map[Action]int{
	ActionCommit:      10,
	ActionIssueCreate: 5,
	ActionIssueClose:  15,
	ActionIssueAssign: 5,
}

// User represents a user in the repository
type User struct {
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Points       int            `json:"points"`
	ActionLog    []ActionRecord `json:"action_log"`
	Commits      int            `json:"commits"`
	IssuesOpen   int            `json:"issues_open"`
	IssuesClosed int            `json:"issues_closed"`
}

// ActionRecord represents a record of a user action
type ActionRecord struct {
	Action      Action `json:"action"`
	Points      int    `json:"points"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// UserManager manages users in a repository
type UserManager struct {
	RepoPath string
}

// NewUserManager creates a new user manager for a repository
func NewUserManager(repoPath string) *UserManager {
	return &UserManager{
		RepoPath: repoPath,
	}
}

// getUsersDir returns the path to the users directory
func (um *UserManager) getUsersDir() string {
	return filepath.Join(um.RepoPath, ".snap", "users")
}

// getUserFilePath returns the path to a user file
func (um *UserManager) getUserFilePath(name string) string {
	return filepath.Join(um.getUsersDir(), fmt.Sprintf("%s.json", name))
}

// GetUser gets a user by name
func (um *UserManager) GetUser(name string) (*User, error) {
	// Check if user file exists
	filePath := um.getUserFilePath(name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If the file doesn't exist, create a new user
		return &User{
			Name:      name,
			Points:    0,
			ActionLog: []ActionRecord{},
		}, nil
	}

	// Read user file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	// Unmarshal user from JSON
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// SaveUser saves a user to disk
func (um *UserManager) SaveUser(user *User) error {
	// Create users directory if it doesn't exist
	usersDir := um.getUsersDir()
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		return fmt.Errorf("failed to create users directory: %w", err)
	}

	// Marshal user to JSON
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Write user to file
	filePath := um.getUserFilePath(user.Name)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user file: %w", err)
	}

	return nil
}

// RecordAction records a user action and awards points
func (um *UserManager) RecordAction(name string, action Action, description string, timestamp string) error {
	// Get user
	user, err := um.GetUser(name)
	if err != nil {
		return err
	}

	// Get points for action
	points, ok := PointValues[action]
	if !ok {
		return fmt.Errorf("unknown action: %s", action)
	}

	// Update user stats
	user.Points += points
	user.ActionLog = append(user.ActionLog, ActionRecord{
		Action:      action,
		Points:      points,
		Description: description,
		Timestamp:   timestamp,
	})

	// Update specific counters
	switch action {
	case ActionCommit:
		user.Commits++
	case ActionIssueCreate:
		user.IssuesOpen++
	case ActionIssueClose:
		user.IssuesClosed++
	}

	// Save user
	if err := um.SaveUser(user); err != nil {
		return err
	}

	return nil
}

// GetLeaderboard gets the leaderboard of users
func (um *UserManager) GetLeaderboard() ([]*User, error) {
	// Read all user files
	usersDir := um.getUsersDir()
	files, err := os.ReadDir(usersDir)
	if err != nil {
		if os.IsNotExist(err) {
			// If the directory doesn't exist, return an empty list
			return []*User{}, nil
		}
		return nil, fmt.Errorf("failed to read users directory: %w", err)
	}

	// Read each user file
	var users []*User
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Extract name from filename (remove .json extension)
		name := strings.TrimSuffix(file.Name(), ".json")
		if name == file.Name() {
			// If the filename doesn't end with .json, skip it
			continue
		}

		user, err := um.GetUser(name)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	// Sort users by points (descending)
	sort.Slice(users, func(i, j int) bool {
		return users[i].Points > users[j].Points
	})

	return users, nil
}
