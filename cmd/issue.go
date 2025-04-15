package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/stanlocht/snap/pkg/issue"
	"github.com/stanlocht/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  `Manage issues in the repository.`,
}

// issueNewCmd represents the issue new command
var issueNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new issue",
	Long:  `Create a new issue in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get issue title and description
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")

		if title == "" {
			fmt.Fprintln(os.Stderr, "Error: issue title is required")
			fmt.Fprintln(os.Stderr, "Use --title to specify an issue title")
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
			authorName = "unknown"
		}

		// Create issue manager
		issueManager := issue.NewIssueManager(repo.Path)

		// Create issue
		newIssue, err := issueManager.CreateIssue(title, description, authorName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created issue #%d: %s\n", newIssue.ID, newIssue.Title)
	},
}

// issueListCmd represents the issue list command
var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `List issues in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get show-closed flag
		showClosed, _ := cmd.Flags().GetBool("show-closed")

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

		// Create issue manager
		issueManager := issue.NewIssueManager(repo.Path)

		// List issues
		issues, err := issueManager.ListIssues(showClosed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing issues: %v\n", err)
			os.Exit(1)
		}

		if len(issues) == 0 {
			fmt.Println("No issues found")
			return
		}

		// Print issues
		for _, issue := range issues {
			statusStr := "OPEN"
			if issue.Status == "closed" {
				statusStr = "CLOSED"
			}

			fmt.Printf("#%d [%s] %s", issue.ID, statusStr, issue.Title)
			if issue.AssignedTo != "" {
				fmt.Printf(" (assigned to %s)", issue.AssignedTo)
			}
			fmt.Println()
		}
	},
}

// issueShowCmd represents the issue show command
var issueShowCmd = &cobra.Command{
	Use:   "show [issue-id]",
	Short: "Show issue details",
	Long:  `Show details of a specific issue.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse issue ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid issue ID: %s\n", args[0])
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

		// Create issue manager
		issueManager := issue.NewIssueManager(repo.Path)

		// Get issue
		issue, err := issueManager.GetIssue(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: issue #%d not found\n", id)
			os.Exit(1)
		}

		// Print issue details
		fmt.Printf("Issue #%d: %s\n", issue.ID, issue.Title)
		fmt.Printf("Status: %s\n", issue.Status)
		fmt.Printf("Created by: %s at %s\n", issue.CreatedBy, issue.CreatedAt.Format("2006-01-02 15:04:05"))
		if issue.Status == "closed" {
			fmt.Printf("Closed at: %s\n", issue.ClosedAt.Format("2006-01-02 15:04:05"))
		}
		if issue.AssignedTo != "" {
			fmt.Printf("Assigned to: %s\n", issue.AssignedTo)
		}
		fmt.Println()
		fmt.Println("Description:")
		fmt.Println(issue.Description)
	},
}

// issueCloseCmd represents the issue close command
var issueCloseCmd = &cobra.Command{
	Use:   "close [issue-id]",
	Short: "Close an issue",
	Long:  `Close an issue in the repository.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse issue ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid issue ID: %s\n", args[0])
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

		// Create issue manager
		issueManager := issue.NewIssueManager(repo.Path)

		// Close issue
		if err := issueManager.CloseIssue(id); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Closed issue #%d\n", id)
	},
}

// issueAssignCmd represents the issue assign command
var issueAssignCmd = &cobra.Command{
	Use:   "assign [issue-id] [assignee]",
	Short: "Assign an issue to a user",
	Long:  `Assign an issue to a user in the repository.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse issue ID
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid issue ID: %s\n", args[0])
			os.Exit(1)
		}

		// Get assignee
		assignee := args[1]

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

		// Create issue manager
		issueManager := issue.NewIssueManager(repo.Path)

		// Assign issue
		if err := issueManager.AssignIssue(id, assignee); err != nil {
			fmt.Fprintf(os.Stderr, "Error assigning issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Assigned issue #%d to %s\n", id, assignee)
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issueNewCmd)
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueShowCmd)
	issueCmd.AddCommand(issueCloseCmd)
	issueCmd.AddCommand(issueAssignCmd)

	// Add flags for issue new command
	issueNewCmd.Flags().StringP("title", "t", "", "Issue title")
	issueNewCmd.Flags().StringP("description", "d", "", "Issue description")

	// Add flags for issue list command
	issueListCmd.Flags().BoolP("show-closed", "c", false, "Show closed issues")
}
