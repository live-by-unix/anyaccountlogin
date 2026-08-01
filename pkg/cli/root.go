package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anyaccountlogin",
	Short: "AnyAccountLogin - Cross-platform authentication system",
	Long: `AnyAccountLogin is a cross-platform authentication system that provides
secure login using flash drives and device identification.

Features:
- Generate and manage cryptographic keys
- Register devices with unique identifiers
- Secure login using flash drive authentication
- SSH integration with composite passwords
- Cross-platform support (Linux, macOS, Windows)`,
}

// Execute runs the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// GetRootCommand returns the root command for testing
func GetRootCommand() *cobra.Command {
	return rootCmd
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.anyaccountlogin.yaml)")
}
