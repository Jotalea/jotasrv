# Maintainer: jotalea <you@example.com>
# Build from a local source directory (no download):
#   cp packaging/local/PKGBUILD /path/to/jotasrv-source/
#   cd /path/to/jotasrv-source && makepkg -si

pkgname=jotasrv
pkgver=2.0.0
pkgrel=1
pkgdesc="A lightweight file server with directory listing, file operations, upload, and a web UI"
arch=('x86_64')
url="https://github.com/jotalea/jotasrv"
license=('BSD-3-Clause')
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
}
