#!/bin/bash

# Build script for Linux .rpm package

set -e

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build/linux"
RPM_DIR="${BUILD_DIR}/rpm"
SPEC_DIR="${RPM_DIR}/SPECS"
SOURCE_DIR="${RPM_DIR}/SOURCES"
RPMS_DIR="${RPM_DIR}/RPMS/x86_64"
SRPMS_DIR="${RPM_DIR}/SRPMS"

echo "Building AnyAccountLogin Linux .rpm installer v${VERSION}"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${SPEC_DIR}" "${SOURCE_DIR}" "${RPMS_DIR}" "${SRPMS_DIR}"

# Create source tarball
echo "Creating source tarball..."
mkdir -p "${SOURCE_DIR}/anyaccountlogin-${VERSION}"
cp -r cmd pkg internal "${SOURCE_DIR}/anyaccountlogin-${VERSION}/"
cp go.mod go.sum "${SOURCE_DIR}/anyaccountlogin-${VERSION}/"
cd "${SOURCE_DIR}"
tar czf "anyaccountlogin-${VERSION}.tar.gz" "anyaccountlogin-${VERSION}"
cd - > /dev/null

# Create spec file
cat > "${SPEC_DIR}/anyaccountlogin.spec" << EOF
Name:           anyaccountlogin
Version:        ${VERSION}
Release:        1%{?dist}
Summary:        Cross-platform authentication system
License:        MIT
URL:            https://github.com/live-by-unix/anyaccountlogin
Source0:        anyaccountlogin-%{version}.tar.gz

BuildRequires:  golang
BuildRequires:  pam-devel
BuildRequires:  systemd
Requires:       pam
Requires:       systemd

%description
AnyAccountLogin provides secure login using flash drives and device
identification. It supports PAM integration on Linux for system-wide
authentication.

%prep
%setup -q

%build
# Build CLI
go build -o anyaccountlogin ./cmd/anyaccountlogin
go build -o anyaccountlogin-daemon ./cmd/anyaccountlogin-daemon

# Build PAM module
cd packaging/linux
make -f Makefile.pam
cd ../..

%install
rm -rf %{buildroot}

# Install binaries
install -D -m 755 anyaccountlogin %{buildroot}/usr/local/bin/anyaccountlogin
install -D -m 755 anyaccountlogin-daemon %{buildroot}/usr/local/bin/anyaccountlogin-daemon

# Install PAM module
install -D -m 644 packaging/linux/pam_anyaccountlogin.so %{buildroot}/lib/security/pam_anyaccountlogin.so
install -D -m 644 packaging/linux/anyaccountlogin.pam %{buildroot}/etc/pam.d/anyaccountlogin

# Install systemd service
install -D -m 644 packaging/linux/anyaccountlogin.service %{buildroot}/etc/systemd/system/anyaccountlogin.service

# Create directories
mkdir -p %{buildroot}/var/lib/anyaccountlogin
mkdir -p %{buildroot}/var/log/anyaccountlogin

%post
# Enable and start the daemon
systemctl daemon-reload
systemctl enable anyaccountlogin.service
systemctl start anyaccountlogin.service

echo "AnyAccountLogin installed. To enable PAM authentication, add:"
echo "auth    required    pam_anyaccountlogin.so"
echo "to /etc/pam.d/system-auth or /etc/pam.d/common-auth"

%preun
# Stop and disable the daemon
if [ \$1 -eq 0 ]; then
    systemctl stop anyaccountlogin.service
    systemctl disable anyaccountlogin.service
fi

%files
%defattr(-,root,root,-)
/usr/local/bin/anyaccountlogin
/usr/local/bin/anyaccountlogin-daemon
/lib/security/pam_anyaccountlogin.so
/etc/pam.d/anyaccountlogin
/etc/systemd/system/anyaccountlogin.service
%dir /var/lib/anyaccountlogin
%dir /var/log/anyaccountlogin

%changelog
* $(date +'%a %b %d %Y') AnyAccountLogin <support@anyaccountlogin.com> - ${VERSION}-1
- Initial package release
EOF

# Build RPM
echo "Building .rpm package..."
rpmbuild --define "_topdir ${RPM_DIR}" \
         --define "_sourcedir ${SOURCE_DIR}" \
         --define "_specdir ${SPEC_DIR}" \
         --define "_rpmdir ${RPM_DIR}/RPMS" \
         --define "_srcrpmdir ${SRPMS_DIR}" \
         -ba "${SPEC_DIR}/anyaccountlogin.spec"

echo "Package built: ${RPMS_DIR}/anyaccountlogin-${VERSION}-1.x86_64.rpm"
