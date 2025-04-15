package cmd

import (
	"fmt"
	"os"

	"github.com/snap/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long: `Show the working tree status.
Displays paths that have differences between the index file and the current HEAD commit,
paths that have differences between the working tree and the index file,
and paths in the working tree that are not tracked by Snap.`,
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

		// Get repository status
		status, branch, err := repo.GetStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting repository status: %v\n", err)
			os.Exit(1)
		}

		// Get current commit ID
		commitID, err := repo.GetHEADCommitID()
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error getting HEAD commit ID: %v\n", err)
			os.Exit(1)
		}

		// Print branch information
		fmt.Printf("On branch %s\n", branch)

		// Print commit information
		if commitID == "" {
			fmt.Println("No commits yet")
		} else {
			commit, err := repo.GetCommit(commitID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting commit: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("HEAD -> %s\n", commitID[:7])
			fmt.Printf("Last commit: %s\n", commit.Message)
		}

		fmt.Println()

		// Print status information
		if len(status) == 0 {
			fmt.Println("Nothing to commit, working tree clean")
		} else {
			// Group files by status
			statusGroups := make(map[string][]string)
			for _, file := range status {
				statusGroups[file.Status] = append(statusGroups[file.Status], file.Path)
			}

			// Print modified files
			if modified, ok := statusGroups["modified"]; ok && len(modified) > 0 {
				fmt.Println("Changes to be committed:")
				fmt.Println("  (use \"snap pop\" to undo the last commit)")
				fmt.Println()
				for _, file := range modified {
					fmt.Printf("\tmodified:   %s\n", file)
				}
				fmt.Println()
			}

			// Print new files
			if new, ok := statusGroups["new"]; ok && len(new) > 0 {
				fmt.Println("New files to be committed:")
				fmt.Println()
				for _, file := range new {
					fmt.Printf("\tnew file:   %s\n", file)
				}
				fmt.Println()
			}

			// Print deleted files
			if deleted, ok := statusGroups["deleted"]; ok && len(deleted) > 0 {
				fmt.Println("Deleted files to be committed:")
				fmt.Println()
				for _, file := range deleted {
					fmt.Printf("\tdeleted:    %s\n", file)
				}
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
