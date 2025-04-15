package cmd

import (
	"fmt"
	"os"

	"github.com/stanlocht/snap/pkg/repository"
	"github.com/stanlocht/snap/pkg/user"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show user stats",
	Long:  `Show user stats and contribution points.`,
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

		// Get author name
		authorName, _ := rootCmd.PersistentFlags().GetString("author")
		if authorName == "" {
			fmt.Fprintln(os.Stderr, "Error: author name is required")
			fmt.Fprintln(os.Stderr, "Use --author (-a) to specify an author name")
			os.Exit(1)
		}

		// Create user manager
		userManager := user.NewUserManager(repo.Path)

		// Get user
		user, err := userManager.GetUser(authorName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting user: %v\n", err)
			os.Exit(1)
		}

		// Print user stats
		fmt.Printf("User: %s\n", user.Name)
		if user.Email != "" {
			fmt.Printf("Email: %s\n", user.Email)
		}
		fmt.Printf("Points: %d\n", user.Points)
		fmt.Printf("Commits: %d\n", user.Commits)
		fmt.Printf("Issues opened: %d\n", user.IssuesOpen)
		fmt.Printf("Issues closed: %d\n", user.IssuesClosed)

		// Print recent actions
		if len(user.ActionLog) > 0 {
			fmt.Println("\nRecent actions:")
			// Show the last 5 actions or all if less than 5
			start := 0
			if len(user.ActionLog) > 5 {
				start = len(user.ActionLog) - 5
			}
			for i := start; i < len(user.ActionLog); i++ {
				action := user.ActionLog[i]
				fmt.Printf("- %s: %s (+%d points)\n", action.Timestamp, action.Description, action.Points)
			}
		}
	},
}

// leaderboardCmd represents the leaderboard command
var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "Show top contributors",
	Long:  `Show top contributors in the repository.`,
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

		// Create user manager
		userManager := user.NewUserManager(repo.Path)

		// Get leaderboard
		users, err := userManager.GetLeaderboard()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting leaderboard: %v\n", err)
			os.Exit(1)
		}

		if len(users) == 0 {
			fmt.Println("No users found")
			return
		}

		// Print leaderboard
		fmt.Println("Leaderboard:")
		fmt.Println("------------")
		for i, user := range users {
			fmt.Printf("%d. %s - %d points\n", i+1, user.Name, user.Points)
		}
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(leaderboardCmd)
}
