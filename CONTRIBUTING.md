# Contributing to AnyAccountLogin

Thank you for your interest in contributing to AnyAccountLogin! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and constructive in all interactions.

### Our Standards

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Platform-specific build tools (see INSTALL.md)

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/anyaccountlogin.git
   cd anyaccountlogin
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/live-by-unix/anyaccountlogin.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the project:
   ```bash
   go build -o anyaccountlogin ./cmd/anyaccountlogin
   go build -o anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon
   ```

## Development Workflow

### Branching Strategy

- `main` - Stable release branch
- `develop` - Development branch for next release
- `feature/*` - Feature branches
- `bugfix/*` - Bug fix branches
- `hotfix/*` - Emergency fixes for production

### Platform Testing

Test on all supported platforms before submitting:
- macOS (Intel amd64, Apple Silicon arm64)
- Linux (amd64, arm64)
- Windows (amd64, arm64)

Use cross-compilation to build for all platforms from your development machine.

### Creating a Feature Branch

1. Ensure your `develop` branch is up to date:
   ```bash
   git checkout develop
   git pull upstream develop
   ```

2. Create a new feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes and commit them (see [Commit Messages](#commit-messages))

4. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

5. Create a pull request to `develop`

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(auth): add TPM support for device registration

Implement TPM ID detection and registration for enhanced
security on systems with TPM 2.0 chips.

Closes #123
```

```
fix(cli): resolve password prompt not showing in interactive mode

The prompt was being written to stdout instead of stderr,
causing issues with password masking.
```

## Coding Standards

### Go Code Style

Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines and use `gofmt`:

```bash
gofmt -s -w .
```

### Code Organization

- Package names should be lowercase, single words
- Exported functions should have documentation comments
- Keep functions focused and small (< 50 lines when possible)
- Use interfaces for abstraction

### Error Handling

- Always handle errors explicitly
- Use error wrapping for context:
  ```go
  return fmt.Errorf("failed to generate key: %w", err)
  ```
- Create custom error types for common errors

### Security Considerations

- Never hardcode credentials or secrets
- Validate all user inputs
- Use constant-time comparison for sensitive data
- Follow security best practices in [SECURITY.md](SECURITY.md)

### Platform-Specific Code

Use build tags for platform-specific code:

```go
// +build linux

package device

func getLinuxUUID() (string, error) {
    // Linux-specific implementation
}
```

## Testing

### Unit Tests

Write unit tests for all new functionality:

```go
func TestGenerateKeyPair(t *testing.T) {
    err := GenerateKeyPair("test_key")
    if err != nil {
        t.Fatalf("GenerateKeyPair failed: %v", err)
    }
    
    // Verify files exist
    if _, err := os.Stat("anyaccountloginkeys/test_key.pem"); os.IsNotExist(err) {
        t.Error("Private key file not created")
    }
}
```

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

### Integration Tests

Integration tests should be placed in `tests/` directory and test cross-platform functionality.

### Manual Testing

Test on all supported platforms before submitting:
- macOS (Intel and Apple Silicon)
- Linux (Ubuntu, Fedora, CentOS)
- Windows (10, 11)

## Documentation

### Code Documentation

- Document all exported packages, functions, and types
- Use godoc format:
  ```go
  // GenerateKeyPair generates an RSA key pair and saves it to disk.
  // The private key is saved with 0600 permissions for security.
  //
  // Parameters:
  //   baseName - The base name for the key files (without extension)
  //
  // Returns:
  //   error - An error if key generation or saving fails
  func GenerateKeyPair(baseName string) error {
      // ...
  }
  ```

### README Updates

Update the README when:
- Adding new features
- Changing installation instructions
- Modifying CLI commands
- Updating security information

### Changelog

Add entries to CHANGELOG.md for:
- New features
- Bug fixes
- Breaking changes
- Security updates

## Submitting Changes

### Pull Request Process

1. Ensure your branch is up to date with `develop`
2. Run all tests and ensure they pass
3. Update documentation as needed
4. Create a pull request with:
   - Clear title and description
   - Reference to related issues
   - Screenshots for UI changes
   - Testing instructions

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Commit messages follow conventional commits
- [ ] No merge conflicts with target branch

### Code Review

- Be respectful and constructive in reviews
- Address review comments promptly
- Ask for clarification if needed
- Keep discussions focused on the code

## Release Process

### Versioning

Follow [Semantic Versioning](https://semver.org/):
- MAJOR: Breaking changes
- MINOR: New features (backwards compatible)
- PATCH: Bug fixes (backwards compatible)

### Release Steps

1. Update version in `internal/version/version.go`
2. Update CHANGELOG.md
3. Create release branch: `release/v1.0.0`
4. Run full test suite on all platforms
5. Build cross-platform binaries:
   - macOS (amd64, arm64)
   - Linux (amd64, arm64)
   - Windows (amd64, arm64)
6. Build installers using CI/CD or packaging scripts
7. Create GitHub release with all platform artifacts
8. Merge release branch to `main`
9. Tag the release: `git tag v1.0.0`
10. Push tag: `git push origin v1.0.0`

### CI/CD

The project uses GitHub Actions for:
- Automated testing
- Cross-platform builds
- Package generation
- Release automation

## Getting Help

- GitHub Issues: [https://github.com/live-by-unix/anyaccountlogin/issues](https://github.com/live-by-unix/anyaccountlogin/issues)
- Discussions: [https://github.com/live-by-unix/anyaccountlogin/discussions](https://github.com/live-by-unix/anyaccountlogin/discussions)
- Email: dev@anyaccountlogin.com

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to AnyAccountLogin!
