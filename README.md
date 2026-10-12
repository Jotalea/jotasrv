# jotasrv

A lightweight file server with directory listing, file operations, upload, and a web UI.

```
jotasrv [flags] [directory]
  -p 1725        port to listen on
  -local         bind to 127.0.0.1 instead of 0.0.0.0
  -cert/-key     enable HTTPS
  -log           access logging to stdout
  -lite          force lightweight UI (no JavaScript)
  -persist-dl    persist download counts to .dlcounts.json
  -user admin    login username (or env JOTASRV_USER)
  -pass 1234     login password (or env JOTASRV_PASS)
  -public        serve ./public_html at /public/ (read-only, no login)
  -no-auth       disable login
  -version       print version and exit
```

Login uses HTTP basic auth and is **on by default with `admin` / `1234`**,
so change it. Flags can also be written with two dashes (`--user`, `--pass`).
Passing the password via `JOTASRV_PASS` keeps it out of the process list.

```sh
jotasrv --user me --pass 's3cret' /some/dir
JOTASRV_USER=me JOTASRV_PASS='s3cret' jotasrv /some/dir
```

## Building

Requires Go (see `go.mod`).

```sh
# Current platform
go build -trimpath -ldflags="-s -w" -o jotasrv .

# Cross-compile linux amd64, armv7 and arm64 into build/
make            # or: make x86_64 / make arm / make arm64
make clean
```

## Debian package

Native build (on Debian/Ubuntu, with `debhelper` and `golang-go` installed):

```sh
sudo apt install build-essential debhelper devscripts fakeroot golang-go
dpkg-buildpackage -us -uc -b
# -> ../jotasrv_<version>_<arch>.deb
```

Without Debian, using Docker (the `.deb` ends up in `./out`):

```sh
docker build -o out .      # BuildKit

# legacy builder (no buildx):
docker build -t jotasrv-deb . && id=$(docker create jotasrv-deb /) \
  && docker cp "$id":/ out && docker rm "$id"
```

Install with `sudo apt install ./out/jotasrv_*.deb`.

## Arch Linux package

Build from the local source tree:

```sh
makepkg -si
```

Build from the git tag (AUR-style):

```sh
cd aur && makepkg -si
```

## systemd service

Both packages install and enable `jotasrv.service`, running as the `jotasrv`
user and serving `/srv/jotasrv`.

| | Debian | Arch |
|---|---|---|
| Config | `/etc/default/jotasrv` | `/etc/jotasrv.conf` |

Set `JOTASRV_ARGS` (flags), `JOTASRV_ROOT` (directory) and `JOTASRV_USER` /
`JOTASRV_PASS` (login) in the config (mode 600, root-only), then:

```sh
sudo systemctl restart jotasrv
sudo journalctl -u jotasrv -f
```

If `JOTASRV_ROOT` is outside `/srv/jotasrv`, grant write access with
`sudo systemctl edit jotasrv` and set `ReadWritePaths=`.

## License

BSD 3-Clause, see `LICENSE`.
