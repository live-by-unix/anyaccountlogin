# Installation Guide

This guide provides step-by-step installation instructions for AnyAccountLogin on different platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Platform Support](#platform-support)
- [macOS Installation](#macos-installation)
- [Linux Installation](#linux-installation)
- [Windows Installation](#windows-installation)
- [Post-Installation Configuration](#post-installation-configuration)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Common Requirements

- Administrator/root access
- 50MB free disk space
- Flash drive for authentication (USB 2.0 or higher)

### Platform-Specific Requirements

**macOS:**
- macOS 10.15 (Catalina) or later
- Xcode Command Line Tools (for building from source)

**Linux:**
- glibc 2.17 or later
- systemd (for service management)
- PAM libraries
- gcc and make (for building from source)

**Windows:**
- Windows 10 or later
- .NET Framework 4.7.2 or later
- Visual C++ Redistributable

## Platform Support

AnyAccountLogin supports the following platforms and architectures:

### macOS
- **Intel (amd64)**: macOS 10.15 (Catalina) or later
- **Apple Silicon (arm64)**: macOS 11 (Big Sur) or later
- **Universal binaries**: Combined Intel + Apple Silicon

### Linux
- **amd64**: Most modern distributions (Ubuntu 18.04+, Debian 10+, CentOS 7+, Fedora 30+)
- **arm64**: ARM64 distributions (Ubuntu Server ARM, Debian ARM, Fedora ARM)
- **Supported distributions**: Debian, Ubuntu, CentOS, RHEL, Fedora, Arch Linux

### Windows
- **amd64**: Windows 10 or later
- **arm64**: Windows 10/11 on ARM devices (Surface Pro X, etc.)

### Architecture Notes
- All binaries are built with optimization flags for smaller size
- Linux builds use CGO_ENABLED=0 for pure Go binaries
- Windows builds include .exe extensions
- macOS supports both architecture-specific and universal binaries

## macOS Installation

### Method 1: Using .pkg Installer (Recommended)

1. Download the latest `.pkg` installer from [GitHub Releases](https://github.com/live-by-unix/anyaccountlogin/releases)
2. Double-click the downloaded `.pkg` file
3. Follow the installation wizard:
   - Click "Continue" to proceed
   - Review the license and accept
   - Select installation destination (default: Macintosh HD)
   - Click "Install" and enter your password when prompted
4. The installer will:
   - Install binaries to `/usr/local/bin/`
   - Install launchd plist to `/Library/LaunchDaemons/`
   - Start the daemon automatically

### Method 2: Building from Source

1. Install Xcode Command Line Tools:
   ```bash
   xcode-select --install
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/live-by-unix/anyaccountlogin.git
   cd anyaccountlogin
   ```

3. Build the binaries:
   ```bash
   # For current platform
   go build -o /usr/local/bin/anyaccountlogin ./cmd/anyaccountlogin
   go build -o /usr/local/bin/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

   # For cross-compilation (choose your target platform)
   GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o anyaccountlogin.exe ./cmd/anyaccountlogin
   GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o anyaccountlogin.exe ./cmd/anyaccountlogin
   ```

4. Install launchd service:
   ```bash
   sudo cp packaging/macos/com.anyaccountlogin.daemon.plist /Library/LaunchDaemons/
   sudo launchctl load /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist
   sudo launchctl start com.anyaccountlogin.daemon
   ```

5. Create required directories:
   ```bash
   sudo mkdir -p /var/lib/anyaccountlogin
   sudo mkdir -p /var/log/anyaccountlogin
   ```

### Verification

```bash
# Check if daemon is running
sudo launchctl list | grep anyaccountlogin

# Check version
anyaccountlogin version
```

## Linux Installation

### Debian/Ubuntu (Using .deb)

1. Download the latest `.deb` package from [GitHub Releases](https://github.com/live-by-unix/anyaccountlogin/releases)

2. Install the package:
   ```bash
   sudo dpkg -i anyaccountlogin_1.0.0_amd64.deb
   ```

3. If dependencies are missing:
   ```bash
   sudo apt-get install -f
   ```

4. Enable and start the service:
   ```bash
   sudo systemctl enable anyaccountlogin.service
   sudo systemctl start anyaccountlogin.service
   ```

### RHEL/CentOS/Fedora (Using .rpm)

1. Download the latest `.rpm` package from [GitHub Releases](https://github.com/live-by-unix/anyaccountlogin/releases)

2. Install the package:
   ```bash
   sudo rpm -i anyaccountlogin-1.0.0-1.x86_64.rpm
   ```

3. Enable and start the service:
   ```bash
   sudo systemctl enable anyaccountlogin.service
   sudo systemctl start anyaccountlogin.service
   ```

### Building from Source

1. Install build dependencies:
   ```bash
   # Debian/Ubuntu
   sudo apt-get update
   sudo apt-get install -y golang gcc make libpam0g-dev systemd-devel

   # RHEL/CentOS/Fedora
   sudo dnf install -y golang gcc make pam-devel systemd-devel
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/live-by-unix/anyaccountlogin.git
   cd anyaccountlogin
   ```

3. Build the binaries:
   ```bash
   # For current platform
   go build -o /usr/local/bin/anyaccountlogin ./cmd/anyaccountlogin
   go build -o /usr/local/bin/anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

   # For cross-compilation (choose your target platform)
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o anyaccountlogin ./cmd/anyaccountlogin
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o anyaccountlogin.exe ./cmd/anyaccountlogin
   GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o anyaccountlogin.exe ./cmd/anyaccountlogin
   ```

4. Build PAM module:
   ```bash
   cd packaging/linux
   make -f Makefile.pam
   sudo make -f Makefile.pam install
   cd ../..
   ```

5. Install systemd service:
   ```bash
   sudo cp packaging/linux/anyaccountlogin.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable anyaccountlogin.service
   sudo systemctl start anyaccountlogin.service
   ```

6. Create required directories:
   ```bash
   sudo mkdir -p /var/lib/anyaccountlogin
   sudo mkdir -p /var/log/anyaccountlogin
   ```

### Verification

```bash
# Check if service is running
sudo systemctl status anyaccountlogin.service

# Check version
anyaccountlogin version
```

## Windows Installation

### Method 1: Using .msi Installer (Recommended)

1. Download the latest `.msi` installer from [GitHub Releases](https://github.com/live-by-unix/anyaccountlogin/releases)
2. Double-click the downloaded `.msi` file
3. Follow the installation wizard:
   - Click "Next" to proceed
   - Accept the license agreement
   - Choose installation directory (default: `C:\Program Files\AnyAccountLogin`)
   - Click "Install"
4. The installer will:
   - Install binaries to `C:\Program Files\AnyAccountLogin\`
   - Install Credential Provider DLL to `C:\Windows\System32\`
   - Install and start the Windows service
   - Configure registry entries

5. Restart your computer to enable the Credential Provider

### Method 2: Manual Installation

1. Download the latest binaries from [GitHub Releases](https://github.com/live-by-unix/anyaccountlogin/releases)
   - Choose the appropriate architecture:
     - `windows-amd64` for Intel/AMD 64-bit systems
     - `windows-arm64` for ARM-based Windows devices

2. Create installation directory:
   ```cmd
   mkdir "C:\Program Files\AnyAccountLogin"
   ```

3. Copy binaries to installation directory:
   ```cmd
   copy anyaccountlogin.exe "C:\Program Files\AnyAccountLogin\"
   copy anyaccountlogin-daemon.exe "C:\Program Files\AnyAccountLogin\"
   copy anyaccountlogin-service.exe "C:\Program Files\AnyAccountLogin\"
   copy AnyAccountLoginCredentialProvider.dll "C:\Windows\System32\"
   ```

4. Register Credential Provider (requires Administrator):
   ```cmd
   regsvr32 "C:\Windows\System32\AnyAccountLoginCredentialProvider.dll"
   ```

5. Install Windows service:
   ```cmd
   cd "C:\Program Files\AnyAccountLogin"
   anyaccountlogin-service.exe install
   ```

6. Start the service:
   ```cmd
   net start AnyAccountLoginService
   ```

### Verification

```cmd
# Check if service is running
sc query AnyAccountLoginService

# Check version
anyaccountlogin version
```

## Post-Installation Configuration

### 1. Generate Authentication Keys

```bash
# Generate user keys
anyaccountlogin create key user

# Generate device keys
anyaccountlogin create key device
```

### 2. Set Up Flash Drive

```bash
# Insert your flash drive and note its path
# macOS: /Volumes/USB
# Linux: /media/username/USB
# Windows: E:\

# Set up flash drive with authentication files
anyaccountlogin create flash-drive <flash-drive-path>
```

### 3. Register Your Device

```bash
# Auto-detect device information
anyaccountlogin register device

# Or specify custom identifier
anyaccountlogin register device my-device-id
```

### 4. Test Authentication

```bash
# Test with flash drive
anyaccountlogin login --flash-drive <path> --password <auth-code>
```

### 5. Enable System Integration (Optional)

**Linux PAM Integration:**

Add to `/etc/pam.d/system-auth`:
```
auth    required    pam_anyaccountlogin.so flash_drive=/media/usb
```

**macOS LoginWindow:**

The loginwindow plugin is automatically installed. Restart to see it at the login screen.

**Windows Credential Provider:**

The Credential Provider is automatically registered. Restart to see it at the login screen.

## Troubleshooting

### macOS Issues

**Daemon not starting:**
```bash
# Check launchd logs
log show --predicate 'process == "anyaccountlogin-daemon"' --last 1h

# Restart daemon
sudo launchctl unload /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist
sudo launchctl load /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist
```

**Permission denied errors:**
```bash
# Fix permissions
sudo chown root:wheel /usr/local/bin/anyaccountlogin*
sudo chmod 755 /usr/local/bin/anyaccountlogin*
```

### Linux Issues

**Service not starting:**
```bash
# Check service status
sudo systemctl status anyaccountlogin.service

# Check logs
sudo journalctl -u anyaccountlogin.service -n 50

# Restart service
sudo systemctl restart anyaccountlogin.service
```

**PAM module errors:**
```bash
# Check PAM logs
grep anyaccountlogin /var/log/auth.log

# Verify PAM module installation
ls -la /lib/security/pam_anyaccountlogin.so
```

### Windows Issues

**Service not starting:**
```cmd
# Check service status
sc query AnyAccountLoginService

# Check Event Viewer
eventvwr.msc

# Restart service
net stop AnyAccountLoginService
net start AnyAccountLoginService
```

**Credential Provider not showing:**
```cmd
# Verify DLL registration
reg query "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Authentication\Credential Providers"

# Re-register DLL
regsvr32 /u "C:\Windows\System32\AnyAccountLoginCredentialProvider.dll"
regsvr32 "C:\Windows\System32\AnyAccountLoginCredentialProvider.dll"
```

### Common Issues

**Flash drive not detected:**
- Ensure flash drive is properly mounted
- Check file system permissions
- Try different USB port

**Authentication failures:**
- Verify flash drive contains `PasswordAuth.pem` and `PasswordAuthCode.txt`
- Ensure password matches the contents of `PasswordAuthCode.txt`
- Check daemon logs for detailed error messages

**Device registration errors:**
- Run with verbose flag: `anyaccountlogin register device --verbose`
- Check if TPM is available (optional)
- Verify system UUID detection

## Uninstallation

### macOS

```bash
# Stop and unload daemon
sudo launchctl stop com.anyaccountlogin.daemon
sudo launchctl unload /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist

# Remove files
sudo rm /usr/local/bin/anyaccountlogin
sudo rm /usr/local/bin/anyaccountlogin-daemon
sudo rm /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist
sudo rm -rf /var/lib/anyaccountlogin
sudo rm -rf /var/log/anyaccountlogin
```

### Linux

```bash
# Stop and disable service
sudo systemctl stop anyaccountlogin.service
sudo systemctl disable anyaccountlogin.service

# Remove package
sudo apt-get remove anyaccountlogin  # Debian/Ubuntu
sudo dnf remove anyaccountlogin      # RHEL/CentOS/Fedora

# Or manual removal
sudo rm /usr/local/bin/anyaccountlogin*
sudo rm /lib/security/pam_anyaccountlogin.so
sudo rm /etc/systemd/system/anyaccountlogin.service
sudo rm -rf /var/lib/anyaccountlogin
sudo rm -rf /var/log/anyaccountlogin
```

### Windows

```cmd
# Stop and uninstall service
net stop AnyAccountLoginService
"C:\Program Files\AnyAccountLogin\anyaccountlogin-service.exe" uninstall

# Unregister Credential Provider
regsvr32 /u "C:\Windows\System32\AnyAccountLoginCredentialProvider.dll"

# Remove files
rmdir /s "C:\Program Files\AnyAccountLogin"
del "C:\Windows\System32\AnyAccountLoginCredentialProvider.dll"
```

## Getting Help

If you encounter issues not covered in this guide:

1. Check the [GitHub Issues](https://github.com/live-by-unix/anyaccountlogin/issues)
2. Review [SECURITY.md](SECURITY.md) for security-related concerns
3. Contact support at support@anyaccountlogin.com
