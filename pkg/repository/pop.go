package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UndoLastCommit undoes the last commit, keeping the changes in the working tree
func (r *Repository) UndoLastCommit() error {
	// Get current HEAD commit ID
	currentCommitID, err := r.GetHEADCommitID()
	if err != nil {
		return fmt.Errorf("failed to get HEAD commit ID: %w", err)
	}

	// If there are no commits, return an error
	if currentCommitID == "" {
		return fmt.Errorf("no commits to undo")
	}

	// Get the current commit
	currentCommit, err := r.GetCommit(currentCommitID)
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	// If there is no parent commit, we're undoing the first commit
	if currentCommit.ParentID == "" {
		// Update HEAD to empty
		headPath := filepath.Join(r.Path, SnapDirName, "HEAD")
		headContent, err := os.ReadFile(headPath)
		if err != nil {
			return fmt.Errorf("failed to read HEAD file: %w", err)
		}

		head := string(headContent)
		if strings.HasPrefix(head, "ref: ") {
			// HEAD is a reference (e.g., "ref: refs/heads/master")
			refPath := head[5:] // Remove "ref: " prefix
			refPath = filepath.Join(r.Path, SnapDirName, refPath)
			
			// Remove the reference file
			if err := os.Remove(refPath); err != nil {
				return fmt.Errorf("failed to remove reference file: %w", err)
			}
		} else {
			// HEAD is a commit ID, update it to empty
			if err := os.WriteFile(headPath, []byte(""), 0644); err != nil {
				return fmt.Errorf("failed to write HEAD file: %w", err)
			}
		}

		return nil
	}

	// Update HEAD to point to the parent commit
	if err := r.UpdateHEAD(currentCommit.ParentID); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	return nil
}
