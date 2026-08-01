package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// KeysDir is the directory where keys are stored
	KeysDir = "anyaccountloginkeys"
	// PasswordAuthPEM is the name of the password auth PEM file
	PasswordAuthPEM = "PasswordAuth.pem"
	// PasswordAuthCodeFile is the name of the password auth code file
	PasswordAuthCodeFile = "PasswordAuthCode.txt"
)

// AuthManager handles authentication operations
type AuthManager struct {
	keysDir string
}

// NewAuthManager creates a new AuthManager
func NewAuthManager() *AuthManager {
	return &AuthManager{
		keysDir: KeysDir,
	}
}

// SetupFlashDrive sets up a flash drive with authentication files
func (am *AuthManager) SetupFlashDrive(flashDrivePath string) error {
	// Ensure flash drive path exists
	if err := os.MkdirAll(flashDrivePath, 0700); err != nil {
		return fmt.Errorf("failed to create flash drive directory: %w", err)
	}

	// Generate auth code
	authCode, err := generateAuthCode()
	if err != nil {
		return fmt.Errorf("failed to generate auth code: %w", err)
	}

	// Save auth code to flash drive
	authCodePath := filepath.Join(flashDrivePath, PasswordAuthCodeFile)
	if err := os.WriteFile(authCodePath, []byte(authCode), 0600); err != nil {
		return fmt.Errorf("failed to write auth code: %w", err)
	}

	// Generate key pair
	keyPairPath := filepath.Join(flashDrivePath, PasswordAuthPEM)
	if err := generateKeyPair(keyPairPath); err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	return nil
}

// ValidateFlashDrive validates the flash drive authentication
func (am *AuthManager) ValidateFlashDrive(flashDrivePath, password string) error {
	// Check if flash drive path exists
	if _, err := os.Stat(flashDrivePath); os.IsNotExist(err) {
		return fmt.Errorf("flash drive path does not exist: %s", flashDrivePath)
	}

	// Read auth code from flash drive
	authCodePath := filepath.Join(flashDrivePath, PasswordAuthCodeFile)
	authCode, err := os.ReadFile(authCodePath)
	if err != nil {
		return fmt.Errorf("failed to read auth code: %w", err)
	}

	// Validate password matches auth code
	if password != string(authCode) {
		return fmt.Errorf("invalid password")
	}

	// Verify PEM key exists
	pemPath := filepath.Join(flashDrivePath, PasswordAuthPEM)
	if _, err := os.Stat(pemPath); os.IsNotExist(err) {
		return fmt.Errorf("PEM key not found on flash drive")
	}

	return nil
}

// LoadPrivateKey loads a private key from a PEM file
func (am *AuthManager) LoadPrivateKey(pemPath string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// GetCompositePassword generates a composite password for SSH login
func (am *AuthManager) GetCompositePassword(flashDrivePath, pemPath string) (string, error) {
	// Read auth code from flash drive
	authCodePath := filepath.Join(flashDrivePath, PasswordAuthCodeFile)
	authCode, err := os.ReadFile(authCodePath)
	if err != nil {
		return "", fmt.Errorf("failed to read auth code: %w", err)
	}

	// Get absolute path of PEM key
	absPemPath, err := filepath.Abs(pemPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create composite password
	compositePassword := fmt.Sprintf("%s-%s", strings.TrimSpace(string(authCode)), absPemPath)
	return compositePassword, nil
}

// RegisterDevice registers a device with the given identifier
func (am *AuthManager) RegisterDevice(identifier string) error {
	// Ensure keys directory exists
	if err := os.MkdirAll(am.keysDir, 0700); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Save device identifier
	deviceIDPath := filepath.Join(am.keysDir, "device_id.txt")
	if err := os.WriteFile(deviceIDPath, []byte(identifier), 0600); err != nil {
		return fmt.Errorf("failed to write device ID: %w", err)
	}

	return nil
}

// GetRegisteredDevice retrieves the registered device identifier
func (am *AuthManager) GetRegisteredDevice() (string, error) {
	deviceIDPath := filepath.Join(am.keysDir, "device_id.txt")
	data, err := os.ReadFile(deviceIDPath)
	if err != nil {
		return "", fmt.Errorf("failed to read device ID: %w", err)
	}
	return string(data), nil
}

// ValidateDevice validates that the current device matches the registered device
func (am *AuthManager) ValidateDevice(currentID string) error {
	registeredID, err := am.GetRegisteredDevice()
	if err != nil {
		return fmt.Errorf("failed to get registered device: %w", err)
	}

	if currentID != registeredID {
		return fmt.Errorf("device mismatch: current=%s, registered=%s", currentID, registeredID)
	}

	return nil
}

// generateAuthCode generates a random auth code
func generateAuthCode() (string, error) {
	// This would use crypto/rand in a real implementation
	// For now, return a placeholder
	return "generated-auth-code", nil
}

// generateKeyPair generates an RSA key pair
func generateKeyPair(basePath string) error {
	// This would use the crypto package in a real implementation
	// For now, create a placeholder file
	privateKeyPath := basePath
	publicKeyPath := basePath + ".pub"

	// Create placeholder files
	if err := os.WriteFile(privateKeyPath, []byte("placeholder-private-key"), 0600); err != nil {
		return err
	}
	return os.WriteFile(publicKeyPath, []byte("placeholder-public-key"), 0644)
}
