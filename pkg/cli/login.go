package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/live-by-unix/anyaccountlogin/pkg/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate using flash drive and password",
	Long: `Authenticate using a flash drive containing authentication files
and a password. This validates both the flash drive presence and
the correct password.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flashDrivePath, _ := cmd.Flags().GetString("flash-drive")
		password, _ := cmd.Flags().GetString("password")

		if flashDrivePath == "" {
			return fmt.Errorf("flash drive path is required (use --flash-drive flag)")
		}

		if password == "" {
			return fmt.Errorf("password is required (use --password flag)")
		}

		authManager := auth.NewAuthManager()

		fmt.Printf("Authenticating with flash drive at %s...\n", flashDrivePath)

		if err := authManager.ValidateFlashDrive(flashDrivePath, password); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Println("Authentication successful!")
		return nil
	},
}

var loginInteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive login with prompts",
	Long: `Interactive login mode that prompts for flash drive path and password
instead of requiring command-line flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prompt for flash drive path
		fmt.Print("Enter flash drive path: ")
		var flashDrivePath string
		fmt.Scanln(&flashDrivePath)

		// Prompt for password
		fmt.Print("Enter password: ")
		var password string
		fmt.Scanln(&password)

		authManager := auth.NewAuthManager()

		if err := authManager.ValidateFlashDrive(flashDrivePath, password); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Println("Authentication successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.AddCommand(loginInteractiveCmd)

	loginCmd.Flags().StringP("flash-drive", "f", "", "Path to flash drive containing authentication files")
	loginCmd.Flags().StringP("password", "p", "", "Password for authentication")
	loginCmd.MarkFlagRequired("flash-drive")
	loginCmd.MarkFlagRequired("password")
}
