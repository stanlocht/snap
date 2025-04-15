package cmd

import (
	"fmt"
)

// Version information
var (
	Version   = "0.0.1"
	BuildDate = "2024-06-01"
	GitCommit = "development"
)

// PrintVersion prints the version information
func PrintVersion() {
	fmt.Printf("Snap version %s\n", Version)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Git commit: %s\n", GitCommit)
}
