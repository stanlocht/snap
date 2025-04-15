package cmd

import (
	"fmt"
	"os"

	"github.com/stanlocht/snap/pkg/repository"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Snap repository",
	Long: `Initialize a new Snap repository in the current directory.
This creates a .snap directory with all necessary files and directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Initialize repository
		repo, err := repository.Init(currentDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Initialized empty Snap repository in %s\n", repo.Path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
