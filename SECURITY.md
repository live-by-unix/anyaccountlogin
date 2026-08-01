# Security Documentation

This document describes the security model, threat analysis, and security considerations for AnyAccountLogin.

## Table of Contents

- [Security Overview](#security-overview)
- [Threat Model](#threat-model)
- [Architecture Security](#architecture-security)
- [Identifier Binding](#identifier-binding)
- [Flash Drive Protection](#flash-drive-protection)
- [Cryptographic Implementation](#cryptographic-implementation)
- [Platform-Specific Security](#platform-specific-security)
- [Security Best Practices](#security-best-practices)
- [Vulnerability Reporting](#vulnerability-reporting)

## Security Overview

AnyAccountLogin provides hardware-based authentication using flash drives and device identification. The system is designed with defense-in-depth principles:

1. **Something you have**: Flash drive with cryptographic keys
2. **Something you are**: Device binding (UUID/TPM)
3. **Something you know**: Password (optional)

### Security Goals

- **Confidentiality**: Protect authentication secrets and keys
- **Integrity**: Ensure authentication data is not tampered with
- **Availability**: Ensure authentication services are available when needed
- **Authenticity**: Verify the identity of users and devices

## Threat Model

### Threat Actors

#### 1. Physical Attackers

**Threat**: An attacker gains physical access to the device or flash drive.

**Mitigations**:
- Flash drive keys are encrypted (optional feature)
- Device binding prevents using flash drive on unauthorized devices
- Password adds second factor of authentication
- Keys are stored with restrictive file permissions (0600)

#### 2. Malicious Software

**Threat**: Malware attempts to steal authentication credentials or keys.

**Mitigations**:
- Keys are stored in protected directories
- Daemon runs with minimal privileges
- No keys are stored in environment variables
- Memory is cleared after use

#### 3. Network Attackers

**Threat**: Attackers attempt to intercept authentication over the network.

**Mitigations**:
- SSH integration uses standard SSH encryption
- Local daemon communication uses Unix sockets (Linux/macOS)
- Windows uses named pipes with proper ACLs
- No authentication data transmitted over network without encryption

#### 4. Insider Threats

**Threat**: Authorized users attempt to bypass authentication controls.

**Mitigations**:
- Device binding prevents unauthorized device registration
- Audit logging of authentication attempts
- PAM/Credential Provider integration for system-wide enforcement
- Regular security reviews

### Attack Vectors

#### 1. Flash Drive Cloning

**Attack**: Attacker copies flash drive contents to another drive.

**Mitigations**:
- Device binding requires matching system UUID/TPM
- TPM provides hardware-rooted trust
- Password adds second factor
- Optional encryption of flash drive contents

#### 2. Key Extraction

**Attack**: Attacker extracts private keys from flash drive or system.

**Mitigations**:
- Private keys stored with 0600 permissions
- Daemon runs with minimal privileges
- Keys are not logged or exposed in error messages
- Optional encrypted key storage

#### 3. Replay Attacks

**Attack**: Attacker captures and replays authentication requests.

**Mitigations**:
- Challenge-response authentication (future enhancement)
- Time-based token expiration (future enhancement)
- Nonce-based authentication (future enhancement)

#### 4. Man-in-the-Middle

**Attack**: Attacker intercepts communication between components.

**Mitigations**:
- Unix sockets for local communication (Linux/macOS)
- Named pipes with ACLs (Windows)
- Mutual authentication between components
- Integrity checks on all data

## Architecture Security

### Component Isolation

```
┌─────────────────────────────────────────────────┐
│               User Interface                    │
│  (CLI / LoginWindow / Credential Provider)      │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│              Authentication Daemon              │
│  (HTTP Server / Unix Socket / Named Pipe)      │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│           Authentication Modules                 │
│  (Flash Drive Validator / Device Registry)      │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│              Storage Layer                       │
│  (Keys / Device IDs / Config)                    │
└─────────────────────────────────────────────────┘
```

### Privilege Separation

- **Daemon**: Runs as root/administrator for system integration
- **CLI**: Runs as user for key management
- **Authentication Modules**: Minimal privileges for specific tasks

### Secure Communication

**Linux/macOS**:
- Unix domain sockets with filesystem permissions
- No network exposure by default
- Optional TLS for remote management

**Windows**:
- Named pipes with proper ACLs
- Local-only communication
- Windows security descriptors

## Identifier Binding

### System UUID

The system UUID provides hardware-based device identification:

**Linux**:
- `/sys/class/dmi/id/product_uuid`
- `/etc/machine-id`
- `/var/lib/dbus/machine-id`

**macOS**:
- `ioreg -rd1 -c IOPlatformExpertDevice`

**Windows**:
- `wmic csproduct get uuid`

**Security Considerations**:
- UUIDs can be spoofed in virtual environments
- Physical hardware provides stronger assurance
- TPM provides additional hardware-rooted trust

### TPM Integration

Trusted Platform Module (TPM) provides hardware-rooted security:

**Benefits**:
- Hardware-based key storage
- Attestation capabilities
- Anti-rollback protection
- Secure boot integration

**Implementation**:
- TPM 2.0 support
- EK (Endorsement Key) certificate validation
- Platform Configuration Registers (PCR) validation

**Limitations**:
- Not available on all systems
- May require BIOS configuration
- Virtual machine TPMs are less secure

## Flash Drive Protection

### Key Storage

Flash drives contain:
- `PasswordAuth.pem` - RSA private key (4096-bit)
- `PasswordAuthCode.txt` - Authentication code

**Security Measures**:
- Keys generated with high entropy
- Private key permissions set to 0600
- Optional encryption of flash drive contents
- Recommended use of encrypted flash drives

### Recommendations

1. **Use Encrypted Flash Drives**
   - BitLocker To Go (Windows)
   - FileVault (macOS)
   - LUKS (Linux)

2. **Physical Security**
   - Keep flash drive secure when not in use
   - Use tamper-evident cases
   - Consider hardware-encrypted drives

3. **Backup Strategy**
   - Keep secure backup of keys
   - Store in separate secure location
   - Document recovery procedures

4. **Regular Rotation**
   - Rotate keys periodically
   - Update device registration
   - Revoke compromised keys

## Cryptographic Implementation

### Key Generation

**RSA Key Parameters**:
- Key size: 4096 bits
- Public exponent: 65537 (standard)
- Randomness: CSPRNG (crypto/rand or OpenSSL)

**Key Storage Format**:
- PKCS#1 for private keys
- PKIX for public keys
- PEM encoding for transport

### Random Number Generation

**Go Implementation**:
```go
import "crypto/rand"

b := make([]byte, 32)
_, err := rand.Read(b)
```

**OpenSSL Implementation**:
```bash
openssl rand -hex 32
```

### Cross-Platform Security

**Architecture Support**:
- **amd64/x86_64**: Standard 64-bit Intel/AMD architecture
- **arm64**: ARM 64-bit architecture (Apple Silicon, ARM servers, Windows ARM)

**Platform-Specific Considerations**:
- All binaries are built with consistent security properties across platforms
- Cross-compiled binaries maintain the same cryptographic implementations
- Architecture-specific optimizations do not compromise security
- ARM64 binaries are tested on Apple Silicon and ARM Linux servers

### Hash Functions

- SHA-256 for integrity checks
- SHA-512 for future enhancements
- Constant-time comparison for security-sensitive data

### Future Enhancements

- Elliptic Curve Cryptography (ECC) support
- Key derivation functions (KDF)
- Hardware security module (HSM) integration

## Platform-Specific Security

### Linux (PAM Integration)

**Security Features**:
- PAM module runs with user context
- System-wide authentication enforcement
- Integration with sudo, login, sshd

**Hardening**:
```bash
# Restrict PAM module permissions
chmod 644 /lib/security/pam_anyaccountlogin.so
chown root:root /lib/security/pam_anyaccountlogin.so

# Secure PAM configuration
chmod 644 /etc/pam.d/anyaccountlogin
```

### macOS (LoginWindow Plugin)

**Security Features**:
- Runs in loginwindow context
- System Authorization (Authorization Services)
- Code signing requirement

**Hardening**:
- Code sign all binaries
- Enable SIP (System Integrity Protection)
- Restrict plugin permissions

### Windows (Credential Provider)

**Security Features**:
- Runs in secure desktop (Logon UI)
- Windows security checks
- Credential Provider Framework isolation

**Hardening**:
- Code sign DLL and executables
- Enable Windows Defender
- Configure UAC properly

## Security Best Practices

### For Users

1. **Key Management**
   - Never share your flash drive
   - Keep backup copies secure
   - Rotate keys regularly
   - Use strong passwords

2. **Device Security**
   - Keep systems updated
   - Enable disk encryption
   - Use secure boot
   - Configure firewall

3. **Operational Security**
   - Monitor authentication logs
   - Report suspicious activity
   - Test recovery procedures
   - Document security incidents

### For Administrators

1. **Deployment**
   - Test in non-production first
   - Roll out gradually
   - Monitor for issues
   - Have rollback plan

2. **Monitoring**
   - Enable audit logging
   - Set up alerts
   - Regular security reviews
   - Compliance checks

3. **Maintenance**
   - Keep software updated
   - Rotate credentials
   - Review access logs
   - Test disaster recovery

### For Developers

1. **Code Security**
   - Follow secure coding practices
   - Use static analysis tools
   - Regular code reviews
   - Dependency scanning

2. **Testing**
   - Security testing in CI/CD
   - Penetration testing
   - Threat modeling
   - Incident response testing

## Vulnerability Reporting

### Responsible Disclosure

If you discover a security vulnerability, please report it responsibly:

1. **Email**: security@anyaccountlogin.com
2. **PGP Key**: Available on request
3. **Response Time**: Within 48 hours
4. **Disclosure Policy**: Coordinated disclosure

### Bug Bounty Program

We offer a bug bounty program for security vulnerabilities:

- **Critical**: $1,000 - $5,000
- **High**: $500 - $1,000
- **Medium**: $100 - $500
- **Low**: $50 - $100

### Security Announcements

Security announcements will be:
- Published on GitHub Security Advisories
- Emailed to security mailing list
- Posted on project website
- Included in release notes

## Compliance

AnyAccountLogin is designed to support:

- **GDPR**: Data protection and privacy
- **SOC 2**: Security controls
- **PCI DSS**: Payment card industry (with proper configuration)
- **HIPAA**: Healthcare (with proper configuration)

## References

- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CIS Controls](https://www.cisecurity.org/controls/)
- [Platform Security Guidelines](https://developer.apple.com/security/)

## Contact

- **Security Email**: security@anyaccountlogin.com
- **PGP Key**: Available on request
- **Security Issues**: [GitHub Security Advisories](https://github.com/live-by-unix/anyaccountlogin/security/advisories)

---

This document is regularly updated to reflect the latest security practices and threat landscape. Last updated: 2026-08-01
