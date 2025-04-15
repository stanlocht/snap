package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Commit represents a commit in the repository
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	ParentID  string    `json:"parent_id,omitempty"`
	TreeID    string    `json:"tree_id"`
}

// Tree represents a tree object in the repository
type Tree struct {
	Entries map[string]string `json:"entries"` // Map of file paths to object IDs
}

// CreateCommit creates a new commit in the repository
func (r *Repository) CreateCommit(message, author, email string, tree *Tree) (*Commit, error) {
	// Get current HEAD commit ID
	parentID, err := r.GetHEADCommitID()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to get HEAD commit ID: %w", err)
	}

	// Create commit object
	commit := &Commit{
		Message:   message,
		Author:    author,
		Email:     email,
		Timestamp: time.Now(),
		ParentID:  parentID,
	}

	// Save tree
	treeID, err := r.SaveTree(tree)
	if err != nil {
		return nil, fmt.Errorf("failed to save tree: %w", err)
	}
	commit.TreeID = treeID

	// Generate commit ID
	commitJSON, err := json.Marshal(commit)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commit: %w", err)
	}
	hash := sha1.New()
	hash.Write(commitJSON)
	commit.ID = hex.EncodeToString(hash.Sum(nil))

	// Save commit
	if err := r.SaveCommit(commit); err != nil {
		return nil, fmt.Errorf("failed to save commit: %w", err)
	}

	// Update HEAD
	if err := r.UpdateHEAD(commit.ID); err != nil {
		return nil, fmt.Errorf("failed to update HEAD: %w", err)
	}

	return commit, nil
}

// SaveCommit saves a commit to the repository
func (r *Repository) SaveCommit(commit *Commit) error {
	// Create commits directory if it doesn't exist
	commitsDir := filepath.Join(r.Path, SnapDirName, "objects", "commits")
	if err := os.MkdirAll(commitsDir, 0755); err != nil {
		return fmt.Errorf("failed to create commits directory: %w", err)
	}

	// Marshal commit to JSON
	commitJSON, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %w", err)
	}

	// Write commit to file
	commitPath := filepath.Join(commitsDir, commit.ID)
	if err := os.WriteFile(commitPath, commitJSON, 0644); err != nil {
		return fmt.Errorf("failed to write commit file: %w", err)
	}

	return nil
}

// GetCommit gets a commit from the repository
func (r *Repository) GetCommit(id string) (*Commit, error) {
	// Read commit file
	commitPath := filepath.Join(r.Path, SnapDirName, "objects", "commits", id)
	commitJSON, err := os.ReadFile(commitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit file: %w", err)
	}

	// Unmarshal commit from JSON
	var commit Commit
	if err := json.Unmarshal(commitJSON, &commit); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commit: %w", err)
	}

	return &commit, nil
}

// SaveTree saves a tree to the repository
func (r *Repository) SaveTree(tree *Tree) (string, error) {
	// Create trees directory if it doesn't exist
	treesDir := filepath.Join(r.Path, SnapDirName, "objects", "trees")
	if err := os.MkdirAll(treesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create trees directory: %w", err)
	}

	// Marshal tree to JSON
	treeJSON, err := json.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tree: %w", err)
	}

	// Generate tree ID
	hash := sha1.New()
	hash.Write(treeJSON)
	treeID := hex.EncodeToString(hash.Sum(nil))

	// Write tree to file
	treePath := filepath.Join(treesDir, treeID)
	if err := os.WriteFile(treePath, treeJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write tree file: %w", err)
	}

	return treeID, nil
}

// GetTree gets a tree from the repository
func (r *Repository) GetTree(id string) (*Tree, error) {
	// Read tree file
	treePath := filepath.Join(r.Path, SnapDirName, "objects", "trees", id)
	treeJSON, err := os.ReadFile(treePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree file: %w", err)
	}

	// Unmarshal tree from JSON
	var tree Tree
	if err := json.Unmarshal(treeJSON, &tree); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tree: %w", err)
	}

	return &tree, nil
}

// GetHEADCommitID gets the current HEAD commit ID
func (r *Repository) GetHEADCommitID() (string, error) {
	// Read HEAD file
	headPath := filepath.Join(r.Path, SnapDirName, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	// Parse HEAD content
	head := string(headContent)
	if len(head) > 0 && head[0] == 'r' {
		// HEAD is a reference (e.g., "ref: refs/heads/master")
		refPath := head[5:] // Remove "ref: " prefix
		refPath = filepath.Join(r.Path, SnapDirName, refPath)
		
		// Read reference file
		refContent, err := os.ReadFile(refPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Reference file doesn't exist yet (no commits)
				return "", nil
			}
			return "", err
		}
		
		return string(refContent), nil
	}
	
	// HEAD is a commit ID
	return head, nil
}

// UpdateHEAD updates the HEAD reference to point to a commit
func (r *Repository) UpdateHEAD(commitID string) error {
	// Read HEAD file to determine if it's a reference or a commit ID
	headPath := filepath.Join(r.Path, SnapDirName, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD file: %w", err)
	}

	head := string(headContent)
	if len(head) > 0 && head[0] == 'r' {
		// HEAD is a reference (e.g., "ref: refs/heads/master")
		refPath := head[5:] // Remove "ref: " prefix
		refPath = filepath.Join(r.Path, SnapDirName, refPath)
		
		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(refPath), 0755); err != nil {
			return fmt.Errorf("failed to create reference directory: %w", err)
		}
		
		// Write commit ID to reference file
		if err := os.WriteFile(refPath, []byte(commitID), 0644); err != nil {
			return fmt.Errorf("failed to write reference file: %w", err)
		}
	} else {
		// HEAD is a commit ID, update it directly
		if err := os.WriteFile(headPath, []byte(commitID), 0644); err != nil {
			return fmt.Errorf("failed to write HEAD file: %w", err)
		}
	}

	return nil
}

// GetCommitHistory gets the commit history starting from the given commit ID
func (r *Repository) GetCommitHistory(startCommitID string) ([]*Commit, error) {
	var history []*Commit
	currentID := startCommitID

	// If no start commit ID is provided, use HEAD
	if currentID == "" {
		var err error
		currentID, err = r.GetHEADCommitID()
		if err != nil {
			return nil, fmt.Errorf("failed to get HEAD commit ID: %w", err)
		}
		
		// If HEAD is empty (no commits yet), return empty history
		if currentID == "" {
			return history, nil
		}
	}

	// Traverse commit history
	for currentID != "" {
		commit, err := r.GetCommit(currentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get commit %s: %w", currentID, err)
		}

		history = append(history, commit)
		currentID = commit.ParentID
	}

	return history, nil
}
