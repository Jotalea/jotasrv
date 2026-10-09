# Build the .deb without Debian installed:
#   docker build -o out .          (BuildKit)
#   or, with the legacy builder:
#   docker build -t jotasrv-deb . && id=$(docker create jotasrv-deb /) && docker cp $id:/ out && docker rm $id
# The package(s) end up in ./out
FROM golang:1.26-bookworm AS build
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential debhelper devscripts fakeroot \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY . .
# -d: golang-go from apt isn't needed, the official Go toolchain is on PATH
RUN dpkg-buildpackage -us -uc -b -d && mkdir /out && cp ../*.deb /out/

FROM scratch
COPY --from=build /out/ /
