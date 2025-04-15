package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/snap/snap/pkg/gitmoji"
	"github.com/snap/snap/pkg/repository"
	"github.com/snap/snap/pkg/storage"
	"github.com/snap/snap/pkg/user"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long: `Record changes to the repository with a message.
All commit messages must start with a Gitmoji (e.g., :sparkles:, ✨).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get commit message
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			fmt.Fprintln(os.Stderr, "Error: commit message is required")
			fmt.Fprintln(os.Stderr, "Use --message (-m) to specify a commit message")
			fmt.Fprintln(os.Stderr, "\nAvailable Gitmojis:")
			fmt.Fprintln(os.Stderr, gitmoji.GetGitmojiList())
			os.Exit(1)
		}

		// Validate commit message (must start with a gitmoji)
		if err := gitmoji.ValidateCommitMessage(message); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, gitmoji.GetGitmojiList())
			os.Exit(1)
		}

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

		// Get author name
		authorName, _ := rootCmd.PersistentFlags().GetString("author")
		if authorName == "" {
			fmt.Fprintln(os.Stderr, "Error: author name is required")
			fmt.Fprintln(os.Stderr, "Use --author (-a) to specify an author name")
			os.Exit(1)
		}

		// Record user action
		userManager := user.NewUserManager(repo.Path)
		timestamp := time.Now().Format(time.RFC3339)
		err = userManager.RecordAction(authorName, user.ActionCommit, message, timestamp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error recording user action: %v\n", err)
			os.Exit(1)
		}

		// Load index
		index, err := storage.LoadIndex(repo.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
			os.Exit(1)
		}

		// Create tree from index
		tree := &repository.Tree{
			Entries: index.Entries,
		}

		// Create commit
		email, _ := rootCmd.PersistentFlags().GetString("email")
		commit, err := repo.CreateCommit(message, authorName, email, tree)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating commit: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created commit %s\n", commit.ID[:7])
		fmt.Printf("Message: %s\n", message)
		fmt.Printf("Timestamp: %s\n", timestamp)
		fmt.Printf("Earned %d points for committing!\n", user.PointValues[user.ActionCommit])
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message (must start with a gitmoji)")
}
