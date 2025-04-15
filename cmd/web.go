package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/stanlocht/snap/pkg/repository"
	"github.com/stanlocht/snap/pkg/web"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start a web interface for the repository",
	Long: `Start a lightweight web interface for browsing the repository.
The web interface provides access to commits, issues, and user statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get port from flag
		port, _ := cmd.Flags().GetInt("port")

		// Get open flag
		openBrowser, _ := cmd.Flags().GetBool("open")

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

		// Start web server
		serverURL := fmt.Sprintf("http://localhost:%d", port)
		fmt.Printf("Starting Snap web interface at %s\n", serverURL)
		fmt.Println("Press Ctrl+C to stop the server")

		// Open browser if requested
		if openBrowser {
			go func() {
				// Wait a moment for the server to start
				openURL(serverURL)
			}()
		}

		// Start the web server
		if err := web.StartServer(repo, port); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
			os.Exit(1)
		}
	},
}

// openURL opens the specified URL in the default browser
func openURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening browser: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntP("port", "p", 8123, "Port to run the web server on")
	webCmd.Flags().BoolP("open", "o", false, "Open the web interface in the default browser")
}
