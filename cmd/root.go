package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snap is a Git-like version control system with a focus on fun and community",
	Long: `Snap is an opinionated, fun, and community-focused version control system.
It supports core version control features like initializing a repo, committing files,
branching, etc., but with some key differences:

- All commits must start with a Gitmoji (e.g., :sparkles:, ✨)
- Built-in issue tracking (create, assign, close issues)
- Gamification system where users earn points for contributions`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringP("author", "a", "", "Author name for commits")
	rootCmd.PersistentFlags().StringP("email", "e", "", "Author email for commits")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "Display version information")
}
