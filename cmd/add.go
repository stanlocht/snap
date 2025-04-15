package cmd

import (
	"fmt"
	"os"

	"github.com/snap/snap/pkg/repository"
	"github.com/snap/snap/pkg/storage"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file1] [file2] ...",
	Short: "Add file contents to the index",
	Long: `Add file contents to the index (staging area).
This command updates the index using the current content found in the working tree,
preparing the content for the next commit.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Find repository
		repo, err := repository.Find(currentDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Load index
		index, err := storage.LoadIndex(repo.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
			os.Exit(1)
		}

		// Add each file to the index
		for _, filePath := range args {
			objectID, err := index.AddFile(repo.Path, filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error adding file %s: %v\n", filePath, err)
				continue
			}
			fmt.Printf("Added %s (object %s)\n", filePath, objectID)
		}

		// Save index
		if err := index.SaveIndex(repo.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving index: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
