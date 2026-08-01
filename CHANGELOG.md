# Changelog

All notable changes to AnyAccountLogin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Windows ARM64 support for ARM-based Windows devices
- Enhanced cross-compilation documentation for all platforms
- Platform support section in installation guide
- Updated GitHub organization to live-by-unix
- Updated copyright to live-by-unix (2026)

### Changed
- Improved cross-compilation build instructions
- Updated all GitHub links to live-by-unix organization
- Enhanced platform testing requirements in contribution guide

### Planned Features
- Challenge-response authentication
- Time-based token expiration
- Elliptic Curve Cryptography (ECC) support
- Hardware Security Module (HSM) integration
- Mobile app for remote authentication
- Biometric integration (fingerprint, Face ID)
- Cloud backup and sync of keys
- Multi-factor authentication options

### Security Enhancements
- Encrypted flash drive support
- Hardware key storage (TPM, Secure Enclave)
- Certificate-based authentication
- Zero-knowledge password proof

## [1.0.0] - 2026-08-01

### Added
- Initial release of AnyAccountLogin
- Cross-platform CLI tool (Linux, macOS, Windows)
- RSA key generation (4096-bit) using Go crypto and OpenSSL
- Flash drive authentication with PEM keys and auth codes
- Device registration with System UUID detection
- TPM ID detection for enhanced security
- SSH integration with composite password authentication
- Boot-time daemon for authentication services
- PAM module for Linux system authentication
- LoginWindow plugin for macOS login screen integration
- Credential Provider for Windows login screen integration
- Systemd service for Linux
- Launchd service for macOS
- Windows Service for Windows
- Comprehensive CLI commands:
  - `create key <user|device>` - Generate cryptographic keys
  - `create flash-drive <path>` - Set up flash drive
  - `register device [identifier]` - Register devices
  - `login` - Authenticate with flash drive
  - `ssh` - SSH with composite password
  - `version` - Display version information
  - `help` - Display help information
- Packaging scripts:
  - macOS .pkg installer (pkgbuild/productbuild)
  - Linux .deb package
  - Linux .rpm package
  - Windows .msi installer (WiX Toolset)
- GitHub Actions CI/CD workflow
- Cross-compilation support for multiple architectures
- Comprehensive documentation:
  - README.md with overview and usage
  - INSTALL.md with platform-specific setup
  - CONTRIBUTING.md with development guidelines
  - SECURITY.md with threat model and security considerations
  - LICENSE (MIT)

### Security Features
- Device binding using System UUID and TPM ID
- Flash drive-based authentication
- RSA 4096-bit key generation
- Secure key storage with proper permissions
- Platform-specific security integrations
- Defense-in-depth architecture

### Platform Support
- Linux (amd64, arm64)
- macOS (amd64, arm64 - universal binaries)
- Windows (amd64, arm64)
- Cross-compilation support for all platforms

### Dependencies
- Go 1.21+
- Cobra CLI framework
- Viper configuration management
- Platform-specific libraries (PAM, Security APIs, Credential Provider Framework)

### Documentation
- Complete installation guides for all platforms
- Security documentation with threat model
- Developer contribution guidelines
- CLI command reference
- Architecture documentation

## [0.1.0] - 2026-07-15

### Added
- Initial project structure
- Basic CLI framework
- Key generation functionality
- Device detection
- Basic authentication logic

### Changed
- Project initialization and planning

---

## Version Format

Version numbers follow the format: MAJOR.MINOR.PATCH

- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality additions
- **PATCH**: Backwards-compatible bug fixes

## Release Process

1. Update version in `internal/version/version.go`
2. Update CHANGELOG.md
3. Create release branch
4. Run full test suite
5. Build installers via CI/CD
6. Create GitHub release
7. Tag release
8. Merge to main branch

## Support

For support and questions:
- GitHub Issues: [https://github.com/live-by-unix/anyaccountlogin/issues](https://github.com/live-by-unix/anyaccountlogin/issues)
- Email: support@anyaccountlogin.com
