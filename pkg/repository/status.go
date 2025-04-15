package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileStatus represents the status of a file in the repository
type FileStatus struct {
	Path   string
	Status string // "modified", "new", "deleted", "staged", "untracked"
}

// GetStatus gets the status of the repository
func (r *Repository) GetStatus() ([]FileStatus, string, error) {
	var status []FileStatus
	
	// Get current branch
	branch := "master" // Default branch
	headPath := filepath.Join(r.Path, SnapDirName, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read HEAD file: %w", err)
	}
	
	head := string(headContent)
	if strings.HasPrefix(head, "ref: refs/heads/") {
		branch = strings.TrimPrefix(head, "ref: refs/heads/")
		branch = strings.TrimSpace(branch)
	}
	
	// Get current commit ID
	commitID, err := r.GetHEADCommitID()
	if err != nil && !os.IsNotExist(err) {
		return nil, "", fmt.Errorf("failed to get HEAD commit ID: %w", err)
	}
	
	// If there are no commits yet, return empty status
	if commitID == "" {
		return status, branch, nil
	}
	
	// Get current commit
	commit, err := r.GetCommit(commitID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get commit: %w", err)
	}
	
	// Get current tree
	tree, err := r.GetTree(commit.TreeID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get tree: %w", err)
	}
	
	// Get index
	indexPath := filepath.Join(r.Path, SnapDirName, "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Index doesn't exist yet
		return status, branch, nil
	}
	
	// Read index file
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read index file: %w", err)
	}
	
	// Parse index
	indexEntries := make(map[string]string)
	for _, line := range strings.Split(string(indexContent), "\n") {
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		
		objectID := parts[0]
		path := parts[1]
		indexEntries[path] = objectID
	}
	
	// Compare index with tree
	for path, objectID := range indexEntries {
		if treeObjectID, ok := tree.Entries[path]; ok {
			if objectID != treeObjectID {
				// File is modified
				status = append(status, FileStatus{
					Path:   path,
					Status: "modified",
				})
			}
		} else {
			// File is new
			status = append(status, FileStatus{
				Path:   path,
				Status: "new",
			})
		}
	}
	
	// Check for deleted files
	for path := range tree.Entries {
		if _, ok := indexEntries[path]; !ok {
			// File is deleted
			status = append(status, FileStatus{
				Path:   path,
				Status: "deleted",
			})
		}
	}
	
	return status, branch, nil
}
