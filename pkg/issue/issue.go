package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Status represents the status of an issue
type Status string

const (
	// StatusOpen represents an open issue
	StatusOpen Status = "open"
	// StatusClosed represents a closed issue
	StatusClosed Status = "closed"
)

// Issue represents a tracked issue in the repository
type Issue struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ClosedAt    time.Time `json:"closed_at,omitempty"`
	AssignedTo  string    `json:"assigned_to,omitempty"`
	CreatedBy   string    `json:"created_by"`
}

// IssueManager manages issues in a repository
type IssueManager struct {
	RepoPath string
}

// NewIssueManager creates a new issue manager for a repository
func NewIssueManager(repoPath string) *IssueManager {
	return &IssueManager{
		RepoPath: repoPath,
	}
}

// getIssuesDir returns the path to the issues directory
func (im *IssueManager) getIssuesDir() string {
	return filepath.Join(im.RepoPath, ".snap", "issues")
}

// getIssueFilePath returns the path to an issue file
func (im *IssueManager) getIssueFilePath(id int) string {
	return filepath.Join(im.getIssuesDir(), fmt.Sprintf("%d.json", id))
}

// getNextIssueID returns the next available issue ID
func (im *IssueManager) getNextIssueID() (int, error) {
	// Read all issue files
	issuesDir := im.getIssuesDir()
	files, err := os.ReadDir(issuesDir)
	if err != nil {
		if os.IsNotExist(err) {
			// If the directory doesn't exist, create it
			if err := os.MkdirAll(issuesDir, 0755); err != nil {
				return 0, fmt.Errorf("failed to create issues directory: %w", err)
			}
			return 1, nil
		}
		return 0, fmt.Errorf("failed to read issues directory: %w", err)
	}

	// Find the highest issue ID
	maxID := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var id int
		_, err := fmt.Sscanf(file.Name(), "%d.json", &id)
		if err != nil {
			continue
		}

		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1, nil
}

// CreateIssue creates a new issue
func (im *IssueManager) CreateIssue(title, description, createdBy string) (*Issue, error) {
	// Get next issue ID
	id, err := im.getNextIssueID()
	if err != nil {
		return nil, err
	}

	// Create issue
	issue := &Issue{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
	}

	// Save issue
	if err := im.SaveIssue(issue); err != nil {
		return nil, err
	}

	return issue, nil
}

// SaveIssue saves an issue to disk
func (im *IssueManager) SaveIssue(issue *Issue) error {
	// Create issues directory if it doesn't exist
	issuesDir := im.getIssuesDir()
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		return fmt.Errorf("failed to create issues directory: %w", err)
	}

	// Marshal issue to JSON
	data, err := json.MarshalIndent(issue, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal issue: %w", err)
	}

	// Write issue to file
	filePath := im.getIssueFilePath(issue.ID)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write issue file: %w", err)
	}

	return nil
}

// GetIssue gets an issue by ID
func (im *IssueManager) GetIssue(id int) (*Issue, error) {
	// Read issue file
	filePath := im.getIssueFilePath(id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read issue file: %w", err)
	}

	// Unmarshal issue from JSON
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue: %w", err)
	}

	return &issue, nil
}

// ListIssues lists all issues
func (im *IssueManager) ListIssues(showClosed bool) ([]*Issue, error) {
	// Read all issue files
	issuesDir := im.getIssuesDir()
	files, err := os.ReadDir(issuesDir)
	if err != nil {
		if os.IsNotExist(err) {
			// If the directory doesn't exist, return an empty list
			return []*Issue{}, nil
		}
		return nil, fmt.Errorf("failed to read issues directory: %w", err)
	}

	// Read each issue file
	var issues []*Issue
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var id int
		_, err := fmt.Sscanf(file.Name(), "%d.json", &id)
		if err != nil {
			continue
		}

		issue, err := im.GetIssue(id)
		if err != nil {
			continue
		}

		// Skip closed issues if showClosed is false
		if !showClosed && issue.Status == StatusClosed {
			continue
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// CloseIssue closes an issue
func (im *IssueManager) CloseIssue(id int) error {
	// Get issue
	issue, err := im.GetIssue(id)
	if err != nil {
		return err
	}

	// Update issue
	issue.Status = StatusClosed
	issue.ClosedAt = time.Now()

	// Save issue
	if err := im.SaveIssue(issue); err != nil {
		return err
	}

	return nil
}

// AssignIssue assigns an issue to a user
func (im *IssueManager) AssignIssue(id int, assignee string) error {
	// Get issue
	issue, err := im.GetIssue(id)
	if err != nil {
		return err
	}

	// Update issue
	issue.AssignedTo = assignee

	// Save issue
	if err := im.SaveIssue(issue); err != nil {
		return err
	}

	return nil
}
