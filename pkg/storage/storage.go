package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Index represents the staging area
type Index struct {
	Entries map[string]string // Map of file paths to object IDs
}

// NewIndex creates a new empty index
func NewIndex() *Index {
	return &Index{
		Entries: make(map[string]string),
	}
}

// AddFile adds a file to the index
func (idx *Index) AddFile(repoPath, filePath string) (string, error) {
	// Get absolute path of the file
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if file exists
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	if fileInfo.IsDir() {
		return "", fmt.Errorf("%s is a directory, not a file", filePath)
	}

	// Read file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate SHA1 hash of content
	hash := sha1.New()
	hash.Write(content)
	objectID := hex.EncodeToString(hash.Sum(nil))

	// Store object in .snap/objects
	objectsDir := filepath.Join(repoPath, ".snap", "objects")
	objectPath := filepath.Join(objectsDir, objectID[:2], objectID[2:])

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create objects directory: %w", err)
	}

	// Write content to object file
	if err := os.WriteFile(objectPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write object file: %w", err)
	}

	// Add to index
	relPath, err := filepath.Rel(repoPath, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	idx.Entries[relPath] = objectID

	return objectID, nil
}

// SaveIndex saves the index to the .snap/index file
func (idx *Index) SaveIndex(repoPath string) error {
	indexPath := filepath.Join(repoPath, ".snap", "index")
	file, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	// Write each entry to the index file
	for path, objectID := range idx.Entries {
		line := fmt.Sprintf("%s %s\n", objectID, path)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write to index file: %w", err)
		}
	}

	return nil
}

// LoadIndex loads the index from the .snap/index file
func LoadIndex(repoPath string) (*Index, error) {
	indexPath := filepath.Join(repoPath, ".snap", "index")
	
	// If index file doesn't exist, return an empty index
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return NewIndex(), nil
	}

	file, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	idx := NewIndex()
	scanner := io.ReadSeeker(file)
	
	// Read each line from the index file
	buf := make([]byte, 1024)
	for {
		n, err := scanner.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read index file: %w", err)
		}

		// Process the buffer
		// This is a simplified implementation
		// In a real implementation, we would parse the buffer properly
		// For now, we'll just return an empty index
		if n > 0 {
			// Placeholder for parsing logic
		}
	}

	return idx, nil
}
