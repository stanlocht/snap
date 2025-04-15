package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stanlocht/snap/pkg/repository"
	"github.com/stanlocht/snap/pkg/snapmoji"
	"github.com/stanlocht/snap/pkg/storage"
	"github.com/stanlocht/snap/pkg/user"
	"github.com/spf13/cobra"
)

// boomCmd represents the boom command
var boomCmd = &cobra.Command{
	Use:   "boom [message]",
	Short: "Shortcut for quick commit",
	Long: `Shortcut for quick commit (like git commit -am).
This command adds all modified files and commits them with the given message.
The message must start with a Snapmoji (e.g., :sparkles:, ✨).`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		selectEmoji, _ := cmd.Flags().GetBool("select-emoji")
		autoConvert, _ := cmd.Flags().GetBool("auto-convert")

		// Get commit message
		var message string
		if len(args) > 0 {
			message = args[0]
		}

		// Handle emoji selection
		if selectEmoji {
			// Display numbered list of emojis
			fmt.Println(snapmoji.GetNumberedSnapmojiList())
			fmt.Print("Enter emoji number: ")

			// Read user input
			var emojiNumber int
			_, err := fmt.Scanf("%d", &emojiNumber)
			if err != nil || emojiNumber < 1 || emojiNumber > len(snapmoji.Snapmojis) {
				fmt.Fprintf(os.Stderr, "Error: invalid emoji number\n")
				os.Exit(1)
			}

			// Get selected emoji
			selectedSnapmoji, _ := snapmoji.GetSnapmojiByNumber(emojiNumber)

			// Get commit message text
			if message == "" {
				fmt.Print("Enter commit message (without emoji): ")
				var messageText string
				// Clear the input buffer
				fmt.Scanln()
				messageText, _ = bufio.NewReader(os.Stdin).ReadString('\n')
				messageText = strings.TrimSpace(messageText)
				if messageText == "" {
					fmt.Fprintf(os.Stderr, "Error: commit message cannot be empty\n")
					os.Exit(1)
				}
				message = selectedSnapmoji.Emoji + " " + messageText
			} else {
				// Prepend selected emoji to existing message
				message = selectedSnapmoji.Emoji + " " + message
			}
		} else if message == "" {
			// No message provided and no emoji selection
			fmt.Fprintln(os.Stderr, "Error: commit message is required")
			fmt.Fprintln(os.Stderr, "Provide a message as an argument or use --select-emoji (-s) to select an emoji from a list")
			fmt.Fprintln(os.Stderr, "\nAvailable Snapmojis:")
			fmt.Fprintln(os.Stderr, snapmoji.GetSnapmojiList())
			os.Exit(1)
		} else if autoConvert {
			// Auto-convert keywords to emojis
			oldMessage := message
			message = snapmoji.AutoConvertKeywordsToEmoji(message)
			if oldMessage != message {
				fmt.Printf("Auto-converted: '%s' to '%s'\n", oldMessage, message)
			}
		}

		// Validate commit message (must start with a snapmoji)
		if err := snapmoji.ValidateCommitMessage(message); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprintln(os.Stderr, "Try using --select-emoji (-s) to select an emoji from a list")
			fmt.Fprintln(os.Stderr, "Or use --auto-convert (-c) to auto-convert keywords to emojis")
			fmt.Fprintln(os.Stderr, "\nAvailable Snapmojis:")
			fmt.Fprintln(os.Stderr, snapmoji.GetSnapmojiList())
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

		// Find all modified files
		modifiedFiles, err := findModifiedFiles(repo.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding modified files: %v\n", err)
			os.Exit(1)
		}

		if len(modifiedFiles) == 0 {
			fmt.Println("No modified files found")
			return
		}

		// Load index
		index, err := storage.LoadIndex(repo.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading index: %v\n", err)
			os.Exit(1)
		}

		// Add each file to the index
		for _, filePath := range modifiedFiles {
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

		// Record user action
		userManager := user.NewUserManager(repo.Path)
		timestamp := time.Now().Format(time.RFC3339)
		err = userManager.RecordAction(authorName, user.ActionCommit, message, timestamp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error recording user action: %v\n", err)
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

// findModifiedFiles finds all modified files in the repository
func findModifiedFiles(repoPath string) ([]string, error) {
	var modifiedFiles []string

	// Walk through all files in the repository
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .snap directory
		if info.IsDir() && info.Name() == ".snap" {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		// Add file to list
		modifiedFiles = append(modifiedFiles, relPath)

		return nil
	})

	return modifiedFiles, err
}

func init() {
	rootCmd.AddCommand(boomCmd)
	boomCmd.Flags().BoolP("select-emoji", "s", false, "Select a snapmoji from a list")
	boomCmd.Flags().BoolP("auto-convert", "c", false, "Auto-convert keywords to snapmojis (e.g., 'feature:' to '✨')")
}
