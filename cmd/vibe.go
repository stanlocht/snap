package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/snap/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// vibeCmd represents the vibe command
var vibeCmd = &cobra.Command{
	Use:   "vibe",
	Short: "Show mood of repo",
	Long: `Show mood of repo based on recent commits/emojis.
This command analyzes the emojis used in recent commits to determine the mood of the repository.`,
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
			fmt.Println("No commits yet, no vibe to analyze")
			return
		}

		// Define emoji categories for mood analysis
		emojiCategories := map[string]struct {
			emojis  []string
			mood    string
			emoji   string
			message string
		}{
			"productive": {
				emojis:  []string{"✨", "📚", "📝", "🚀"},
				mood:    "Productive",
				emoji:   "📈",
				message: "The team is making great progress!",
			},
			"bugfix": {
				emojis:  []string{"🐛", "🚑", "🔧"},
				mood:    "Bugfixing",
				emoji:   "🔨",
				message: "The team is squashing bugs and improving stability.",
			},
			"refactoring": {
				emojis:  []string{"♻️", "🎨", "✅"},
				mood:    "Refactoring",
				emoji:   "👾",
				message: "The codebase is getting cleaner and more maintainable.",
			},
			"security": {
				emojis:  []string{"🔒", "🔍"},
				mood:    "Security-focused",
				emoji:   "👮",
				message: "The team is prioritizing security and protection.",
			},
			"cleanup": {
				emojis:  []string{"🔥", "🚮"},
				mood:    "Cleaning up",
				emoji:   "🧹",
				message: "The team is removing old code and cleaning up the codebase.",
			},
		}

		// Count emojis in commit messages
		emojiCounts := make(map[string]int)
		for _, commit := range history {
			for category, info := range emojiCategories {
				for _, emoji := range info.emojis {
					if strings.Contains(commit.Message, emoji) {
						emojiCounts[category]++
						break
					}
				}
			}
		}

		// Find the dominant mood
		dominantMood := "Neutral"
		dominantEmoji := "😐"
		dominantMessage := "The repository has a balanced mix of activities."
		maxCount := 0

		for category, count := range emojiCounts {
			if count > maxCount {
				maxCount = count
				dominantMood = emojiCategories[category].mood
				dominantEmoji = emojiCategories[category].emoji
				dominantMessage = emojiCategories[category].message
			}
		}

		// Collect recent emojis
		recentEmojis := []string{}
		for i, commit := range history {
			if i >= 5 { // Only look at the 5 most recent commits
				break
			}

			// Try to extract emoji from commit message
			for _, info := range emojiCategories {
				for _, emoji := range info.emojis {
					if strings.Contains(commit.Message, emoji) {
						recentEmojis = append(recentEmojis, emoji)
						break
					}
				}
			}
		}

		// Display vibe analysis
		fmt.Println("📊 Repository Vibe Analysis 📊")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("Recent emojis: %s\n", strings.Join(recentEmojis, " "))
		fmt.Printf("Dominant mood: %s %s\n", dominantMood, dominantEmoji)
		fmt.Printf("Analysis: %s\n", dominantMessage)
		fmt.Println(strings.Repeat("-", 50))
	},
}

func init() {
	rootCmd.AddCommand(vibeCmd)
}
