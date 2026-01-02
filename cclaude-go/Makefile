# Makefile para cclaude-go
BINARY_NAME=cclaude
VERSION=1.0.0
BUILD_DIR=build

.PHONY: build build-linux build-macos build-windows clean test install

build: build-linux build-macos build-windows

build-linux:
	@echo "Building Linux amd64..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-macos:
	@echo "Building macOS amd64..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64
	@echo "Building macOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64

build-windows:
	@echo "Building Windows amd64..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cclaude

clean:
	@echo "Cleaning builds..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test ./...

install:
	@echo "Installing to /usr/local/bin..."
	@go build -o /usr/local/bin/$(BINARY_NAME) ./cmd/cclaude

# Cross-compile all platforms
cross-compile:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cclaude
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 ./cmd/cclaude
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 ./cmd/cclaude
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cclaude