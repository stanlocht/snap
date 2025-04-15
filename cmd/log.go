package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stanlocht/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long: `Show commit logs.
Displays the commit history of the current branch.`,
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

		// Get commit history
		history, err := repo.GetCommitHistory("")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commit history: %v\n", err)
			os.Exit(1)
		}

		// Check if there are any commits
		if len(history) == 0 {
			fmt.Println("No commits yet")
			return
		}

		// Display commit history
		fmt.Println("Commit History:")
		fmt.Println(strings.Repeat("-", 50))

		for _, commit := range history {
			fmt.Printf("%s %s\n", commit.ID[:7], commit.Message)
			fmt.Printf("Author: %s <%s>\n", commit.Author, commit.Email)
			fmt.Printf("Date:   %s\n\n", commit.Timestamp.Format(time.RFC1123))
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
