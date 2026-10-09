# Maintainer: jotalea <main@jotalea.com.ar>
# Build from a local source directory (no download):
#   cp packaging/local/PKGBUILD /path/to/jotasrv-source/
#   cd /path/to/jotasrv-source && makepkg -si

pkgname=jotasrv
pkgver=2.1.0
pkgrel=1
pkgdesc="A lightweight file server with directory listing, file operations, upload, and a web UI"
arch=('x86_64')
url="https://github.com/jotalea/jotasrv"
license=('BSD-3-Clause')
backup=('etc/jotasrv.conf')
makedepends=('go')
source=()
sha256sums=()

build() {
	cd "$startdir"
	CGO_ENABLED=0 go build -ldflags="-s -w" -o "${srcdir}/jotasrv"
}

package() {
	install -Dm755 "${srcdir}/jotasrv" "${pkgdir}/usr/bin/${pkgname}"
	install -Dm644 "${startdir}/LICENSE" "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"

	install -Dm644 /dev/stdin "${pkgdir}/usr/lib/systemd/system/${pkgname}.service" <<'EOF2'
[Unit]
Description=jotasrv file server
Documentation=https://github.com/jotalea/jotasrv
After=network-online.target
Wants=network-online.target

[Service]
User=jotasrv
Group=jotasrv
EnvironmentFile=-/etc/jotasrv.conf
ExecStart=/usr/bin/jotasrv $JOTASRV_ARGS $JOTASRV_ROOT
Restart=on-failure
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes
PrivateDevices=yes
ReadWritePaths=/srv/jotasrv

[Install]
WantedBy=multi-user.target
EOF2

	install -Dm644 /dev/stdin "${pkgdir}/etc/jotasrv.conf" <<'EOF2'
# Configuration for the jotasrv systemd service.
# Apply changes with: systemctl restart jotasrv

# Command-line flags (see `jotasrv -h`), e.g. "-p 8080 -log -persist-dl"
# or "-cert /etc/ssl/cert.pem -key /etc/ssl/key.pem" for HTTPS.
JOTASRV_ARGS="-p 1725 -log"

# Directory to serve. The service user needs access to it, and write access
# if uploads/file operations are used (add it to ReadWritePaths via
# `systemctl edit jotasrv` if it is outside /srv/jotasrv).
JOTASRV_ROOT="/srv/jotasrv"
EOF2

	install -Dm644 /dev/stdin "${pkgdir}/usr/lib/sysusers.d/${pkgname}.conf" <<'EOF2'
u jotasrv - "jotasrv file server" /srv/jotasrv
EOF2

	install -Dm644 /dev/stdin "${pkgdir}/usr/lib/tmpfiles.d/${pkgname}.conf" <<'EOF2'
d /srv/jotasrv 0755 jotasrv jotasrv -
EOF2
}
