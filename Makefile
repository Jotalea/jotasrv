APP := jotasrv
BUILD_DIR := build

.PHONY: all clean x86_64 arm arm64

all: x86_64 arm arm64

x86_64:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(APP)-linux-amd64 .

arm:
	GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(APP)-linux-armv7 .

arm64:
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(APP)-linux-arm64 .

clean:
	rm -rf $(BUILD_DIR)
