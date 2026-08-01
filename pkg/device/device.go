package device

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DeviceInfo contains information about the current device
type DeviceInfo struct {
	SystemUUID string
	TPMID      string
	Platform   string
}

// GetDeviceInfo returns device information for the current system
func GetDeviceInfo() (*DeviceInfo, error) {
	info := &DeviceInfo{
		Platform: runtime.GOOS,
	}

	uuid, err := getSystemUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to get system UUID: %w", err)
	}
	info.SystemUUID = uuid

	tpmID, err := getTPMID()
	if err == nil {
		info.TPMID = tpmID
	}

	return info, nil
}

// getSystemUUID returns the system UUID based on the platform
func getSystemUUID() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxUUID()
	case "darwin":
		return getMacUUID()
	case "windows":
		return getWindowsUUID()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// getLinuxUUID retrieves the system UUID on Linux
func getLinuxUUID() (string, error) {
	// Try DMI product UUID first
	if data, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		uuid := strings.TrimSpace(string(data))
		if uuid != "" {
			return uuid, nil
		}
	}

	// Try machine-id from systemd
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		uuid := strings.TrimSpace(string(data))
		if uuid != "" {
			return uuid, nil
		}
	}

	// Try dbus machine-id
	if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		uuid := strings.TrimSpace(string(data))
		if uuid != "" {
			return uuid, nil
		}
	}

	// Fallback to dmidecode
	cmd := exec.Command("dmidecode", "-s", "system-uuid")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		uuid := strings.TrimSpace(out.String())
		if uuid != "" {
			return uuid, nil
		}
	}

	return "", fmt.Errorf("could not determine system UUID on Linux")
}

// getMacUUID retrieves the system UUID on macOS
func getMacUUID() (string, error) {
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run ioreg: %w", err)
	}

	output := out.String()
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.Split(line, `"`)
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}

	return "", fmt.Errorf("could not find IOPlatformUUID")
}

// getWindowsUUID retrieves the system UUID on Windows
func getWindowsUUID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run wmic: %w", err)
	}

	output := strings.TrimSpace(out.String())
	lines := strings.Split(output, "\n")
	if len(lines) >= 2 {
		uuid := strings.TrimSpace(lines[1])
		if uuid != "" {
			return uuid, nil
		}
	}

	return "", fmt.Errorf("could not determine system UUID on Windows")
}

// getTPMID attempts to get the TPM identifier if available
func getTPMID() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxTPMID()
	case "darwin":
		return "", fmt.Errorf("TPM not typically available on macOS")
	case "windows":
		return getWindowsTPMID()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// getLinuxTPMID retrieves the TPM identifier on Linux
func getLinuxTPMID() (string, error) {
	// Check if TPM device exists
	if _, err := os.Stat("/dev/tpm0"); os.IsNotExist(err) {
		return "", fmt.Errorf("TPM device not found")
	}

	// Try to get TPM EK (Endorsement Key) certificate
	cmd := exec.Command("tpm2_getekcertificate", "-out", "/dev/stdout")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		// Use the certificate hash as a TPM identifier
		return fmt.Sprintf("%x", out.Bytes()[:32]), nil
	}

	// Fallback: use TPM version
	cmd = exec.Command("tpm2_getcap", "properties-fixed")
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		output := out.String()
		// Extract TPM version info
		if strings.Contains(output, "TPM2_PT_VENDOR_STRING") {
			return "tpm2-available", nil
		}
	}

	return "", fmt.Errorf("could not determine TPM ID")
}

// getWindowsTPMID retrieves the TPM identifier on Windows
func getWindowsTPMID() (string, error) {
	cmd := exec.Command("powershell", "-Command", "Get-TPM")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run Get-TPM: %w", err)
	}

	output := out.String()
	if strings.Contains(output, "TpmPresent") && strings.Contains(output, "True") {
		// Extract TPM version information
		if strings.Contains(output, "TpmVersion") {
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, "TpmVersion") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						return strings.TrimSpace(parts[1]), nil
					}
				}
			}
		}
		return "tpm-available", nil
	}

	return "", fmt.Errorf("TPM not available on this system")
}

// RegisterDeviceIdentifier generates a device identifier string
func RegisterDeviceIdentifier(uuid, tpmID string) string {
	if tpmID != "" {
		return fmt.Sprintf("%s-%s", uuid, tpmID)
	}
	return uuid
}
