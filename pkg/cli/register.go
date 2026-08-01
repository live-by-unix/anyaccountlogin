package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/live-by-unix/anyaccountlogin/pkg/auth"
	"github.com/live-by-unix/anyaccountlogin/pkg/device"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register resources",
	Long:  `Register resources such as devices with the authentication system.`,
}

var registerDeviceCmd = &cobra.Command{
	Use:   "device <identifier>",
	Short: "Register a device with the authentication system",
	Long: `Register a device using its unique identifier. If no identifier is provided,
the system will attempt to auto-detect the device UUID and TPM ID (if available).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		authManager := auth.NewAuthManager()

		var identifier string
		if len(args) > 0 {
			identifier = args[0]
			fmt.Printf("Registering device with custom identifier: %s\n", identifier)
		} else {
			// Auto-detect device information
			deviceInfo, err := device.GetDeviceInfo()
			if err != nil {
				return fmt.Errorf("failed to get device info: %w", err)
			}

			identifier = device.RegisterDeviceIdentifier(deviceInfo.SystemUUID, deviceInfo.TPMID)
			fmt.Printf("Auto-detected device information:\n")
			fmt.Printf("  System UUID: %s\n", deviceInfo.SystemUUID)
			if deviceInfo.TPMID != "" {
				fmt.Printf("  TPM ID: %s\n", deviceInfo.TPMID)
			}
			fmt.Printf("  Platform: %s\n", deviceInfo.Platform)
			fmt.Printf("\nGenerated device identifier: %s\n", identifier)
		}

		if err := authManager.RegisterDevice(identifier); err != nil {
			return fmt.Errorf("failed to register device: %w", err)
		}

		fmt.Println("Device registered successfully!")
		return nil
	},
}

var validateDeviceCmd = &cobra.Command{
	Use:   "validate-device",
	Short: "Validate the current device against registered device",
	Long: `Validate that the current device matches the registered device identifier.
This is useful for troubleshooting device registration issues.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		authManager := auth.NewAuthManager()

		// Get current device info
		deviceInfo, err := device.GetDeviceInfo()
		if err != nil {
			return fmt.Errorf("failed to get device info: %w", err)
		}

		currentID := device.RegisterDeviceIdentifier(deviceInfo.SystemUUID, deviceInfo.TPMID)

		// Validate against registered device
		if err := authManager.ValidateDevice(currentID); err != nil {
			return fmt.Errorf("device validation failed: %w", err)
		}

		fmt.Println("Device validation successful!")
		fmt.Printf("Current device ID: %s\n", currentID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.AddCommand(registerDeviceCmd)
	registerCmd.AddCommand(validateDeviceCmd)
}
