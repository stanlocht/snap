package cmd

import (
	"fmt"
	"os"

	"github.com/stanlocht/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Undo last commit",
	Long: `Undo last commit.
This command undoes the last commit, keeping the changes in the working directory.`,
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

		// Get current HEAD commit ID before undoing
		currentCommitID, err := repo.GetHEADCommitID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting HEAD commit ID: %v\n", err)
			os.Exit(1)
		}

		// If there are no commits, print a message and exit
		if currentCommitID == "" {
			fmt.Println("No commits to undo")
			return
		}

		// Get the current commit
		currentCommit, err := repo.GetCommit(currentCommitID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current commit: %v\n", err)
			os.Exit(1)
		}

		// Undo the last commit
		if err := repo.UndoLastCommit(); err != nil {
			fmt.Fprintf(os.Stderr, "Error undoing last commit: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Undid last commit %s\n", currentCommitID[:7])
		fmt.Printf("Message: %s\n", currentCommit.Message)
		fmt.Println("Changes from the commit are still in your working directory")
	},
}

func init() {
	rootCmd.AddCommand(popCmd)
}
