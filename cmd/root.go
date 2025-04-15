package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snap is a playful, community-driven version control system with a fresh take on collaboration",
	Long: `Snap is a playful, community-driven version control system with a fresh take on collaboration.
It supports all the essentials—like initializing a repo, committing changes, and branching—but adds its own flair:

- Commits begin with a Snapmoji to keep things expressive and fun (e.g., :sparkles:, ✨)
- Built-in issue tracking makes managing tasks simple and seamless
- Contributions are tracked in a rewarding way, celebrating your work and progress over time`,
	Run: func(cmd *cobra.Command, args []string) {
		// If version flag is set, print version and exit
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			PrintVersion()
			os.Exit(0)
		}

		// Otherwise, show help
		cmd.Help()
	},
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
