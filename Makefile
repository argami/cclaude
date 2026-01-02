# cclaude Makefile

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build install test clean release install-script

# Build the binary
build:
	go build $(LDFLAGS) -o cclaude ./cmd/cclaude/

# Install the binary
install: build
	cp cclaude /usr/local/bin/
	chmod +x /usr/local/bin/cclaude

# Run tests
test:
	go test ./... -v -cover

# Clean build artifacts
clean:
	rm -f cclaude
	rm -f cclaude-*

# Release builds for all platforms
release: clean
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o cclaude-linux-amd64 ./cmd/cclaude/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o cclaude-darwin-amd64 ./cmd/cclaude/
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o cclaude-darwin-arm64 ./cmd/cclaude/

# Print installation instructions
install-script:
	@echo "=== cclaude Installation ==="
	@echo ""
	@echo "Option 1: Using make"
	@echo "  make install"
	@echo ""
	@echo "Option 2: Manual"
	@echo "  go build -o cclaude ./cmd/cclaude/"
	@echo "  cp cclaude /usr/local/bin/"
	@echo "  chmod +x /usr/local/bin/cclaude"
	@echo ""
	@echo "Option 3: From release binaries"
	@echo "  curl -L https://github.com/argami/cclaude/releases/download/$(VERSION)/cclaude-$$(uname -s | tr '[:upper:]' '[:lower:]')-$$(uname -m) -o cclaude"
	@echo "  chmod +x cclaude"
	@echo "  sudo mv cclaude /usr/local/bin/"
