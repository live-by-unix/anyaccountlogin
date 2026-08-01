# AnyAccountLogin

A cross-platform authentication system that provides secure login using flash drives and device identification. AnyAccountLogin enables passwordless and hardware-based authentication across Linux, macOS, and Windows.

## Features

- **Hardware-Based Authentication**: Use flash drives containing cryptographic keys for secure login
- **Device Binding**: Register devices using System UUID and TPM ID for enhanced security
- **Cross-Platform Support**: Native integration with PAM (Linux), loginwindow (macOS), and Credential Provider (Windows)
- **SSH Integration**: Composite password authentication for SSH connections
- **Flexible Key Management**: Generate RSA keys using OpenSSL or Go's crypto library
- **Boot-Time Daemon**: System service running at startup for authentication requests

## Installation

### macOS

Download and install the `.pkg` installer from the [latest release](https://github.com/live-by-unix/anyaccountlogin/releases).

```bash
# Using Homebrew (if available)
brew install anyaccountlogin
```

### Linux

#### Debian/Ubuntu

```bash
sudo dpkg -i anyaccountlogin_1.0.0_amd64.deb
sudo systemctl enable anyaccountlogin.service
sudo systemctl start anyaccountlogin.service
```

#### RHEL/CentOS/Fedora

```bash
sudo rpm -i anyaccountlogin-1.0.0-1.x86_64.rpm
sudo systemctl enable anyaccountlogin.service
sudo systemctl start anyaccountlogin.service
```

### Windows

Download and run the `.msi` installer from the [latest release](https://github.com/live-by-unix/anyaccountlogin/releases).

## Quick Start

### 1. Generate Keys

```bash
# Generate user keys
anyaccountlogin create key user

# Generate device keys
anyaccountlogin create key device

# Use OpenSSL instead of Go's crypto library
anyaccountlogin create key user --openssl
```

### 2. Set Up Flash Drive

```bash
# Set up a flash drive with authentication files
anyaccountlogin create flash-drive /Volumes/USB
```

This creates:
- `PasswordAuth.pem` - RSA private key
- `PasswordAuthCode.txt` - Authentication code

### 3. Register Device

```bash
# Auto-detect device information
anyaccountlogin register device

# Or specify a custom identifier
anyaccountlogin register device my-custom-device-id
```

### 4. Test Authentication

```bash
# Test login with flash drive
anyaccountlogin login --flash-drive /Volumes/USB --password <auth-code>

# Interactive mode
anyaccountlogin login interactive
```

### 5. SSH Integration

```bash
# SSH with composite password
anyaccountlogin ssh user@hostname --flash-drive /Volumes/USB

# Generate SSH config snippet
anyaccountlogin ssh config --flash-drive /Volumes/USB --host myserver --user myuser
```

## CLI Commands

### `create`
Generate authentication resources.

```bash
anyaccountlogin create key <user|device>    # Generate cryptographic keys
anyaccountlogin create flash-drive <path>   # Set up flash drive
```

### `register`
Register resources with the authentication system.

```bash
anyaccountlogin register device [identifier]  # Register a device
anyaccountlogin register validate-device     # Validate device registration
```

### `login`
Authenticate using flash drive and password.

```bash
anyaccountlogin login --flash-drive <path> --password <password>
anyaccountlogin login interactive
```

### `ssh`
SSH with composite password authentication.

```bash
anyaccountlogin ssh [user@]hostname [command] --flash-drive <path>
anyaccountlogin ssh config --flash-drive <path> --host <host> --user <user>
```

### `version`
Print version information.

```bash
anyaccountlogin version           # Print version
anyaccountlogin version --full    # Print full version with git commit
anyaccountlogin version --json    # Output in JSON format
```

### `help`
Display help information.

```bash
anyaccountlogin help              # Show general help
anyaccountlogin <command> --help  # Show command-specific help
```

## Platform-Specific Integration

### Linux (PAM Module)

AnyAccountLogin includes a PAM module for system-wide authentication:

1. Install the package
2. Enable PAM authentication by adding to `/etc/pam.d/system-auth`:
   ```
   auth    required    pam_anyaccountlogin.so flash_drive=/media/usb
   ```
3. Restart the daemon: `sudo systemctl restart anyaccountlogin.service`

### macOS (LoginWindow Plugin)

The macOS package includes a loginwindow plugin for login screen authentication:

1. Install the `.pkg` installer
2. The daemon is automatically started via launchd
3. Login screen will show AnyAccountLogin authentication option

### Windows (Credential Provider)

The Windows package includes a Credential Provider for login screen integration:

1. Install the `.msi` installer
2. The service is automatically installed and started
3. Login screen will show AnyAccountLogin credential option

## Security Notes

### Flash Drive Protection

- Keep your flash drive secure at all times
- Never share the `PasswordAuthCode.txt` contents
- The flash drive contains sensitive cryptographic material
- Consider using encrypted flash drives for additional protection

### Device Binding

- Device registration binds authentication to specific hardware
- TPM ID is used when available for enhanced security
- Changing hardware may require re-registration

### Key Management

- Keys are stored in `anyaccountloginkeys/` directory
- Private keys have restrictive permissions (0600)
- Use `--openssl` flag for OpenSSL-based key generation if preferred

### Threat Model

See [SECURITY.md](SECURITY.md) for detailed threat analysis and security considerations.

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/live-by-unix/anyaccountlogin.git
cd anyaccountlogin

# Build the CLI (current platform)
go build -o anyaccountlogin ./cmd/anyaccountlogin

# Build the daemon (current platform)
go build -o anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon
```

### Cross-Compilation

The project supports cross-compilation for multiple platforms and architectures:

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/darwin-amd64/anyaccountlogin ./cmd/anyaccountlogin
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/darwin-amd64/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/darwin-arm64/anyaccountlogin ./cmd/anyaccountlogin
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/darwin-arm64/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

# Linux (AMD64)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/linux-amd64/anyaccountlogin ./cmd/anyaccountlogin
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/linux-amd64/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

# Linux (ARM64)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/linux-arm64/anyaccountlogin ./cmd/anyaccountlogin
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/linux-arm64/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

# Windows (AMD64)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/windows-amd64/anyaccountlogin.exe ./cmd/anyaccountlogin
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/windows-amd64/anyaccountlogin-daemon.exe ./cmd/anyaccountlogin-daemon

# Windows (ARM64)
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o build/windows-arm64/anyaccountlogin.exe ./cmd/anyaccountlogin
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o build/windows-arm64/anyaccountlogin-daemon.exe ./cmd/anyaccountlogin-daemon
```

**Supported Platforms:**
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Linux: amd64, arm64
- Windows: amd64, arm64

### Building Installers

```bash
# macOS .pkg
cd packaging/macos
./build-pkg.sh

# Linux .deb
cd packaging/linux
./build-deb.sh

# Linux .rpm
cd packaging/linux
./build-rpm.sh

# Windows .msi (requires WiX Toolset)
cd packaging/windows
./build-msi.sh
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.

## Support

- GitHub Issues: [https://github.com/live-by-unix/anyaccountlogin/issues](https://github.com/live-by-unix/anyaccountlogin/issues)
- Documentation: [https://docs.anyaccountlogin.com](https://docs.anyaccountlogin.com)
- Email: support@anyaccountlogin.com

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Cryptographic operations powered by Go's standard library and OpenSSL
- Platform-specific integrations using native OS APIs

## Version

Current version: 1.0.0

For version history, see [CHANGELOG.md](CHANGELOG.md).
