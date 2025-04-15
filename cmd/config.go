package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stanlocht/snap/pkg/config"
	"github.com/stanlocht/snap/pkg/repository"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get or set configuration values",
	Long:  `Get or set configuration values for the Snap repository.`,
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value from the Snap repository.`,
	Args:  cobra.ExactArgs(1),
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

		// Get key
		key := args[0]

		// Read config file
		configPath := filepath.Join(repo.Path, ".snap", "config")
		value, err := config.GetValue(configPath, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		if value == "" {
			fmt.Printf("%s is not set\n", key)
		} else {
			fmt.Printf("%s = %s\n", key, value)
		}
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value in the Snap repository.`,
	Args:  cobra.ExactArgs(2),
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

		// Get key and value
		key := args[0]
		value := args[1]

		// Set config value
		configPath := filepath.Join(repo.Path, ".snap", "config")
		err = config.SetValue(configPath, key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set %s = %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}
