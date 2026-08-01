#!/bin/bash

# Build script for Linux .deb package

set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build/linux"
DEB_DIR="${BUILD_DIR}/deb"
PACKAGE_NAME="anyaccountlogin"

echo "Building AnyAccountLogin Linux .deb installer v${VERSION}"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${DEB_DIR}"

# Create Debian package structure
mkdir -p "${DEB_DIR}/DEBIAN"
mkdir -p "${DEB_DIR}/usr/local/bin"
mkdir -p "${DEB_DIR}/lib/security"
mkdir -p "${DEB_DIR}/etc/pam.d"
mkdir -p "${DEB_DIR}/etc/systemd/system"
mkdir -p "${DEB_DIR}/var/lib/anyaccountlogin"
mkdir -p "${DEB_DIR}/var/log/anyaccountlogin"

# Build Go binaries
echo "Building Go binaries..."
GOOS=linux GOARCH=amd64 go build -o "${DEB_DIR}/usr/local/bin/anyaccountlogin" ./cmd/anyaccountlogin
GOOS=linux GOARCH=amd64 go build -o "${DEB_DIR}/usr/local/bin/anyaccountlogin-daemon" ./cmd/anyaccountlogin-daemon

# Set permissions
chmod +x "${DEB_DIR}/usr/local/bin/anyaccountlogin"
chmod +x "${DEB_DIR}/usr/local/bin/anyaccountlogin-daemon"

# Build PAM module
echo "Building PAM module..."
cd packaging/linux
make -f Makefile.pam
cd ../../
cp packaging/linux/pam_anyaccountlogin.so "${DEB_DIR}/lib/security/"
cp packaging/linux/anyaccountlogin.pam "${DEB_DIR}/etc/pam.d/anyaccountlogin"

# Copy systemd service
cp packaging/linux/anyaccountlogin.service "${DEB_DIR}/etc/systemd/system/"

# Create control file
cat > "${DEB_DIR}/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: AnyAccountLogin <support@anyaccountlogin.com>
Description: Cross-platform authentication system
 AnyAccountLogin provides secure login using flash drives and device
 identification. It supports PAM integration on Linux for system-wide
 authentication.
Depends: libc6, libpam0g, systemd
EOF

# Create postinst script
cat > "${DEB_DIR}/DEBIAN/postinst" << 'EOF'
#!/bin/bash

set -e

# Enable and start the daemon
systemctl enable anyaccountlogin.service
systemctl start anyaccountlogin.service

# Setup PAM configuration (optional, user must manually enable)
echo "AnyAccountLogin installed. To enable PAM authentication, add:"
echo "auth    required    pam_anyaccountlogin.so"
echo "to /etc/pam.d/system-auth or /etc/pam.d/common-auth"

exit 0
EOF

chmod +x "${DEB_DIR}/DEBIAN/postinst"

# Create prerm script
cat > "${DEB_DIR}/DEBIAN/prerm" << 'EOF'
#!/bin/bash

set -e

# Stop and disable the daemon
if systemctl is-active --quiet anyaccountlogin.service; then
    systemctl stop anyaccountlogin.service
fi
systemctl disable anyaccountlogin.service

exit 0
EOF

chmod +x "${DEB_DIR}/DEBIAN/prerm"

# Calculate installed size
INSTALLED_SIZE=$(du -sk "${DEB_DIR}" | cut -f1)
echo "Installed-Size: ${INSTALLED_SIZE}" >> "${DEB_DIR}/DEBIAN/control"

# Build package
echo "Building .deb package..."
dpkg-deb --build "${DEB_DIR}" "${BUILD_DIR}/${PACKAGE_NAME}_${VERSION}_amd64.deb"

echo "Package built: ${BUILD_DIR}/${PACKAGE_NAME}_${VERSION}_amd64.deb"
