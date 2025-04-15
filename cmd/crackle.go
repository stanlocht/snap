package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/snap/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// crackleCmd represents the crackle command
var crackleCmd = &cobra.Command{
	Use:   "crackle",
	Short: "Stylized commit log view",
	Long: `Stylized commit log view.
This command shows a colorful and stylized view of the commit history.`,
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

		// Display stylized commit history
		fmt.Println("✨ Stylized Commit Log with Snapmojis ✨")
		fmt.Println(strings.Repeat("=", 60))

		// Define some snapmoji prefixes for different types of commits
		emojiMap := map[string]string{
			"sparkles":            "✨",
			"bug":                 "🐛",
			"books":               "📚",
			"recycle":             "♻️",
			"wrench":              "🔧",
			"white_check_mark":    "✅",
			"rocket":              "🚀",
			"lipstick":            "💄",
			"fire":                "🔥",
			"ambulance":           "🚑",
			"art":                 "🎨",
			"zap":                 "⚡️",
			"lock":                "🔒",
			"construction":        "🚧",
			"memo":                "📝",
			"truck":               "🚚",
			"construction_worker": "👷",
			"heavy_plus_sign":     "➕",
			"heavy_minus_sign":    "➖",
			"bookmark":            "🔖",
		}

		// Display each commit with emoji and formatting
		for i, commit := range history {
			// Extract emoji from commit message if present
			emoji := "📦" // Default emoji
			for code, e := range emojiMap {
				if strings.Contains(commit.Message, e) || strings.Contains(commit.Message, ":"+code+":") {
					emoji = e
					break
				}
			}

			// Format date
			date := commit.Timestamp.Format("2006-01-02 15:04:05")

			// Format commit message (remove emoji if present)
			message := commit.Message
			for _, e := range emojiMap {
				message = strings.TrimPrefix(message, e+" ")
				message = strings.TrimPrefix(message, e)
			}
			// Also remove emoji codes
			for code := range emojiMap {
				message = strings.TrimPrefix(message, ":"+code+": ")
				message = strings.TrimPrefix(message, ":"+code+":")
			}

			// Print commit with formatting
			fmt.Printf("%s %s | %s | %s\n", emoji, date, commit.Author, message)

			// Add separator between commits except for the last one
			if i < len(history)-1 {
				fmt.Println(strings.Repeat("-", 60))
			}
		}

		fmt.Println(strings.Repeat("=", 60))
	},
}

func init() {
	rootCmd.AddCommand(crackleCmd)
}
