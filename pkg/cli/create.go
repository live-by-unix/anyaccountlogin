package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/live-by-unix/anyaccountlogin/pkg/crypto"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create authentication resources",
	Long:  `Create authentication resources such as keys and devices.`,
}

var createKeyCmd = &cobra.Command{
	Use:   "key <user|device>",
	Short: "Generate cryptographic keys",
	Long: `Generate cryptographic keys for authentication.

  user  - Generate user keys for personal authentication
  device - Generate device keys for system authentication`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyType := args[0]
		useOpenSSL, _ := cmd.Flags().GetBool("openssl")

		if keyType != "user" && keyType != "device" {
			return fmt.Errorf("invalid key type: %s (must be 'user' or 'device')", keyType)
		}

		baseName := keyType + "_key"
		fmt.Printf("Generating %s key...\n", keyType)

		var err error
		if useOpenSSL {
			err = crypto.GenerateKeyWithOpenSSL(baseName)
		} else {
			err = crypto.GenerateKeyPair(baseName)
		}

		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		fmt.Printf("Successfully generated %s key in %s/\n", keyType, crypto.KeysDir)
		fmt.Printf("  Private key: %s/%s.pem\n", crypto.KeysDir, baseName)
		fmt.Printf("  Public key: %s/%s.pub.pem\n", crypto.KeysDir, baseName)
		return nil
	},
}

var createFlashDriveCmd = &cobra.Command{
	Use:   "flash-drive <path>",
	Short: "Set up a flash drive with authentication files",
	Long: `Set up a flash drive with the necessary authentication files
(PasswordAuth.pem and PasswordAuthCode.txt) for secure login.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flashDrivePath := args[0]

		// Check if path exists
		if _, err := os.Stat(flashDrivePath); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", flashDrivePath)
		}

		fmt.Printf("Setting up flash drive at %s...\n", flashDrivePath)

		// Generate auth code
		authCode, err := crypto.GenerateAuthCode()
		if err != nil {
			return fmt.Errorf("failed to generate auth code: %w", err)
		}

		// Save auth code
		authCodePath := fmt.Sprintf("%s/%s", flashDrivePath, "PasswordAuthCode.txt")
		if err := crypto.SaveAuthCode(authCode, authCodePath); err != nil {
			return fmt.Errorf("failed to save auth code: %w", err)
		}

		// Generate key pair
		baseName := fmt.Sprintf("%s/%s", flashDrivePath, "PasswordAuth")
		if err := crypto.GenerateKeyPair(baseName); err != nil {
			return fmt.Errorf("failed to generate key pair: %w", err)
		}

		fmt.Println("Flash drive setup complete!")
		fmt.Printf("  Auth code: %s\n", authCodePath)
		fmt.Printf("  Private key: %s.pem\n", baseName)
		fmt.Printf("  Public key: %s.pub.pem\n", baseName+".pub.pem")
		fmt.Println("\nIMPORTANT: Keep the flash drive safe and do not share the auth code.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createKeyCmd)
	createCmd.AddCommand(createFlashDriveCmd)

	createKeyCmd.Flags().Bool("openssl", false, "Use OpenSSL instead of Go's crypto library")
}
