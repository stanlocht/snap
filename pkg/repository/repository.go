package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// SnapDirName is the name of the directory where Snap stores its data
	SnapDirName = ".snap"
)

// Repository represents a Snap repository
type Repository struct {
	Path string // Path to the repository root
}

// IsInitialized checks if the current directory is a Snap repository
func IsInitialized(path string) bool {
	snapDir := filepath.Join(path, SnapDirName)
	_, err := os.Stat(snapDir)
	return err == nil
}

// Init initializes a new Snap repository
func Init(path string) (*Repository, error) {
	if IsInitialized(path) {
		return nil, errors.New("repository already initialized")
	}

	// Create .snap directory
	snapDir := filepath.Join(path, SnapDirName)
	if err := os.MkdirAll(snapDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .snap directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		filepath.Join(snapDir, "objects"),
		filepath.Join(snapDir, "refs"),
		filepath.Join(snapDir, "refs", "heads"),
		filepath.Join(snapDir, "issues"),
		filepath.Join(snapDir, "users"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create initial HEAD file pointing to master branch
	headFile := filepath.Join(snapDir, "HEAD")
	if err := os.WriteFile(headFile, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to create HEAD file: %w", err)
	}

	// Create config file
	configFile := filepath.Join(snapDir, "config")
	defaultConfig := `[core]
	repositoryformatversion = 0
	filemode = true
[user]
	name = 
	email = 
`
	if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
		return nil, fmt.Errorf("failed to create config file: %w", err)
	}

	return &Repository{
		Path: path,
	}, nil
}

// Find locates the repository from the current directory or any parent directory
func Find(startPath string) (*Repository, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}

	for {
		if IsInitialized(path) {
			return &Repository{Path: path}, nil
		}

		// Move to parent directory
		parent := filepath.Dir(path)
		if parent == path {
			// We've reached the root directory
			return nil, errors.New("not a snap repository (or any of the parent directories)")
		}
		path = parent
	}
}
