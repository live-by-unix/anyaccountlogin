#!/bin/bash

# Build script for Windows .msi installer using WiX

set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build/windows"
WIX_DIR="${BUILD_DIR}/wix"
BINARY_DIR="${BUILD_DIR}/binaries"

echo "Building AnyAccountLogin Windows .msi installer v${VERSION}"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${WIX_DIR}" "${BINARY_DIR}"

# Build Go binaries
echo "Building Go binaries..."
GOOS=windows GOARCH=amd64 go build -o "${BINARY_DIR}/anyaccountlogin.exe" ./cmd/anyaccountlogin
GOOS=windows GOARCH=amd64 go build -o "${BINARY_DIR}/anyaccountlogin-daemon.exe" ./cmd/anyaccountlogin-daemon

# Build Windows service (requires Visual Studio or MinGW)
echo "Building Windows service..."
# This would normally be compiled with Visual Studio or MinGW
# For now, we'll create a placeholder
echo "Windows service build requires Visual Studio or MinGW environment"

# Build Credential Provider (requires Visual Studio)
echo "Building Credential Provider..."
# This would normally be compiled with Visual Studio
# For now, we'll create a placeholder
echo "Credential Provider build requires Visual Studio environment"

# Create config file
cat > "${BINARY_DIR}/config.json" << EOF
{
  "version": "${VERSION}",
  "logLevel": "info",
  "flashDrivePath": "",
  "enableTPM": true
}
EOF

# Build MSI using WiX
echo "Building MSI installer..."
candle packaging/windows/Product.wxs -out "${WIX_DIR}/Product.wixobj"
light "${WIX_DIR}/Product.wixobj" -out "${BUILD_DIR}/AnyAccountLogin-${VERSION}.msi"

echo "Installer built: ${BUILD_DIR}/AnyAccountLogin-${VERSION}.msi"
