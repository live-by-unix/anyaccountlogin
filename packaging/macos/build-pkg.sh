#!/bin/bash

# Build script for macOS .pkg installer

set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build/macos"
PKG_DIR="${BUILD_DIR}/pkg"
ROOT_DIR="${PKG_DIR}/root"
SCRIPTS_DIR="${PKG_DIR}/scripts"

echo "Building AnyAccountLogin macOS installer v${VERSION}"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${ROOT_DIR}" "${SCRIPTS_DIR}"

# Build Go binaries
echo "Building Go binaries..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "${ROOT_DIR}/usr/local/bin/anyaccountlogin" ./cmd/anyaccountlogin
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "${ROOT_DIR}/usr/local/bin/anyaccountlogin-daemon" ./cmd/anyaccountlogin-daemon

# Set permissions
chmod +x "${ROOT_DIR}/usr/local/bin/anyaccountlogin"
chmod +x "${ROOT_DIR}/usr/local/bin/anyaccountlogin-daemon"

# Copy daemon launchd plist
mkdir -p "${ROOT_DIR}/Library/LaunchDaemons"
cp packaging/macos/com.anyaccountlogin.daemon.plist "${ROOT_DIR}/Library/LaunchDaemons/"

# Copy loginwindow plugin (if built)
mkdir -p "${ROOT_DIR}/Library/Security/SecurityAgentPlugins"
# cp packaging/macos/AnyAccountLoginMechanism.bundle "${ROOT_DIR}/Library/Security/SecurityAgentPlugins/"

# Create data directories
mkdir -p "${ROOT_DIR}/var/lib/anyaccountlogin"
mkdir -p "${ROOT_DIR}/var/log/anyaccountlogin"

# Create installer scripts
cat > "${SCRIPTS_DIR}/postinstall" << 'EOF'
#!/bin/bash

# Post-install script for AnyAccountLogin

# Load the daemon
launchctl load /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist

# Start the daemon
launchctl start com.anyaccountlogin.daemon

echo "AnyAccountLogin installed successfully"
exit 0
EOF

chmod +x "${SCRIPTS_DIR}/postinstall"

cat > "${SCRIPTS_DIR}/preinstall" << 'EOF'
#!/bin/bash

# Pre-install script for AnyAccountLogin

# Stop daemon if running
if launchctl list | grep -q com.anyaccountlogin.daemon; then
    launchctl stop com.anyaccountlogin.daemon
    launchctl unload /Library/LaunchDaemons/com.anyaccountlogin.daemon.plist
fi

exit 0
EOF

chmod +x "${SCRIPTS_DIR}/preinstall"

# Build package
echo "Building package..."
pkgbuild \
    --root "${ROOT_DIR}" \
    --scripts "${SCRIPTS_DIR}" \
    --identifier com.anyaccountlogin.anyaccountlogin \
    --version "${VERSION}" \
    --install-location / \
    "${BUILD_DIR}/anyaccountlogin-${VERSION}.pkg"

echo "Package built: ${BUILD_DIR}/anyaccountlogin-${VERSION}.pkg"

# Build product distribution
echo "Building product distribution..."
productbuild \
    --package "${BUILD_DIR}/anyaccountlogin-${VERSION}.pkg" \
    --version "${VERSION}" \
    --identifier com.anyaccountlogin.anyaccountlogin \
    "${BUILD_DIR}/AnyAccountLogin-${VERSION}.pkg"

echo "Installer built: ${BUILD_DIR}/AnyAccountLogin-${VERSION}.pkg"
