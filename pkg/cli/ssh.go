package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/live-by-unix/anyaccountlogin/pkg/auth"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [user@]hostname [command]",
	Short: "SSH with composite password authentication",
	Long: `SSH to a remote host using composite password authentication.
The composite password is in the format: <PasswordAuthCode>-<absolute-path-to-PEM>.

If using the -i flag for identity file, the composite password is not used.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flashDrivePath, _ := cmd.Flags().GetString("flash-drive")
		pemPath, _ := cmd.Flags().GetString("identity")
		forwardAgent, _ := cmd.Flags().GetBool("forward-agent")
		port, _ := cmd.Flags().GetString("port")

		if flashDrivePath == "" && pemPath == "" {
			return fmt.Errorf("either --flash-drive or --identity must be specified")
		}

		var sshArgs []string

		// Add standard SSH flags
		if port != "" {
			sshArgs = append(sshArgs, "-p", port)
		}
		if forwardAgent {
			sshArgs = append(sshArgs, "-A")
		}
		if pemPath != "" {
			absPemPath, err := filepath.Abs(pemPath)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for identity file: %w", err)
			}
			sshArgs = append(sshArgs, "-i", absPemPath)
		}

		// Add the host and optional command
		sshArgs = append(sshArgs, args...)

		// If using flash drive (not -i flag), set up composite password
		if pemPath == "" && flashDrivePath != "" {
			authManager := auth.NewAuthManager()
			// Get the PEM path from flash drive
			flashDrivePEM := filepath.Join(flashDrivePath, "PasswordAuth.pem")
			compositePassword, err := authManager.GetCompositePassword(flashDrivePath, flashDrivePEM)
			if err != nil {
				return fmt.Errorf("failed to generate composite password: %w", err)
			}

			// Set SSH_ASKPASS environment variable to use a helper script
			// For now, we'll output the password for manual use
			fmt.Println("Composite password generated (use this as your SSH password):")
			fmt.Println(compositePassword)
			fmt.Println("\nLaunching SSH client...")
		}

		// Execute SSH command
		sshExecCmd := exec.Command("ssh", sshArgs...)
		sshExecCmd.Stdin = os.Stdin
		sshExecCmd.Stdout = os.Stdout
		sshExecCmd.Stderr = os.Stderr

		if err := sshExecCmd.Run(); err != nil {
			return fmt.Errorf("SSH command failed: %w", err)
		}

		return nil
	},
}

var sshConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate SSH config snippet",
	Long: `Generate an SSH config snippet that can be added to ~/.ssh/config
to use AnyAccountLogin authentication with specific hosts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flashDrivePath, _ := cmd.Flags().GetString("flash-drive")
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")

		if flashDrivePath == "" {
			return fmt.Errorf("flash drive path is required (use --flash-drive flag)")
		}

		if host == "" {
			return fmt.Errorf("host is required (use --host flag)")
		}

		flashDrivePEM := filepath.Join(flashDrivePath, "PasswordAuth.pem")
		absPemPath, err := filepath.Abs(flashDrivePEM)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		fmt.Println("# Add this to your ~/.ssh/config")
		fmt.Printf("Host %s\n", host)
		fmt.Println("    IdentityFile", absPemPath)
		if user != "" {
			fmt.Printf("    User %s\n", user)
		}
		fmt.Println("    IdentitiesOnly yes")
		fmt.Println("    PasswordAuthentication no")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.AddCommand(sshConfigCmd)

	sshCmd.Flags().StringP("flash-drive", "f", "", "Path to flash drive containing authentication files")
	sshCmd.Flags().StringP("identity", "i", "", "Identity file (PEM key)")
	sshCmd.Flags().BoolP("forward-agent", "A", false, "Enable SSH agent forwarding")
	sshCmd.Flags().StringP("port", "p", "", "Port to connect to on the remote host")

	sshConfigCmd.Flags().StringP("flash-drive", "f", "", "Path to flash drive containing authentication files")
	sshConfigCmd.Flags().StringP("host", "H", "", "Host configuration name")
	sshConfigCmd.Flags().StringP("user", "u", "", "Username for the host")
}
