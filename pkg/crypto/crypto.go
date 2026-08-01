package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// KeySize is the RSA key size in bits
	KeySize = 4096
	// KeysDir is the directory where keys are stored
	KeysDir = "anyaccountloginkeys"
)

// GenerateRSAKey generates an RSA private key
func GenerateRSAKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	return privateKey, nil
}

// SavePrivateKeyToPEM saves an RSA private key to a PEM file
func SavePrivateKeyToPEM(privateKey *rsa.PrivateKey, filepath string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	if err := os.WriteFile(filepath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	return nil
}

// SavePublicKeyToPEM saves an RSA public key to a PEM file
func SavePublicKeyToPEM(publicKey *rsa.PublicKey, filepath string) error {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	if err := os.WriteFile(filepath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}

	return nil
}

// GenerateKeyPair generates and saves an RSA key pair
func GenerateKeyPair(baseName string) error {
	// Ensure keys directory exists
	if err := os.MkdirAll(KeysDir, 0700); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	privateKey, err := GenerateRSAKey()
	if err != nil {
		return err
	}

	privateKeyPath := filepath.Join(KeysDir, baseName+".pem")
	publicKeyPath := filepath.Join(KeysDir, baseName+".pub.pem")

	if err := SavePrivateKeyToPEM(privateKey, privateKeyPath); err != nil {
		return err
	}

	if err := SavePublicKeyToPEM(&privateKey.PublicKey, publicKeyPath); err != nil {
		return err
	}

	return nil
}

// GenerateKeyWithOpenSSL generates a key using OpenSSL
func GenerateKeyWithOpenSSL(baseName string) error {
	// Ensure keys directory exists
	if err := os.MkdirAll(KeysDir, 0700); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	privateKeyPath := filepath.Join(KeysDir, baseName+".pem")

	// Generate RSA private key using OpenSSL
	cmd := exec.Command("openssl", "genrsa", "-out", privateKeyPath, fmt.Sprintf("%d", KeySize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate key with OpenSSL: %w", err)
	}

	// Set appropriate permissions
	if err := os.Chmod(privateKeyPath, 0600); err != nil {
		return fmt.Errorf("failed to set permissions on private key: %w", err)
	}

	// Extract public key
	publicKeyPath := filepath.Join(KeysDir, baseName+".pub.pem")
	cmd = exec.Command("openssl", "rsa", "-in", privateKeyPath, "-pubout", "-out", publicKeyPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract public key with OpenSSL: %w", err)
	}

	return nil
}

// GenerateAuthCode generates a random auth code
func GenerateAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate auth code: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}

// SaveAuthCode saves an auth code to a file
func SaveAuthCode(code, filepath string) error {
	if err := os.WriteFile(filepath, []byte(code), 0600); err != nil {
		return fmt.Errorf("failed to write auth code to file: %w", err)
	}
	return nil
}

// LoadAuthCode loads an auth code from a file
func LoadAuthCode(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read auth code from file: %w", err)
	}
	return string(data), nil
}
