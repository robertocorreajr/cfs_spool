# CFS Spool Makefile
VERSION ?= dev
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION)"

# Default target
.PHONY: all
all: build-cli build-web

# Build CLI for current platform
.PHONY: build-cli
build-cli:
	CGO_ENABLED=1 go build $(LDFLAGS) -o bin/cfs-spool-cli cmd/cfs-spool/main.go

# Build web server for current platform  
.PHONY: build-web
build-web:
	CGO_ENABLED=1 go build $(LDFLAGS) -o bin/cfs-spool-web cmd/web-server/main.go

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build $(LDFLAGS) -o dist/cfs-spool-cli-windows-amd64.exe cmd/cfs-spool/main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build $(LDFLAGS) -o dist/cfs-spool-web-windows-amd64.exe cmd/web-server/main.go
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-cli-linux-amd64 cmd/cfs-spool/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-web-linux-amd64 cmd/web-server/main.go
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build $(LDFLAGS) -o dist/cfs-spool-cli-linux-arm64 cmd/cfs-spool/main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build $(LDFLAGS) -o dist/cfs-spool-web-linux-arm64 cmd/web-server/main.go
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-cli-darwin-amd64 cmd/cfs-spool/main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-web-darwin-amd64 cmd/web-server/main.go
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-cli-darwin-arm64 cmd/cfs-spool/main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build $(LDFLAGS) -o dist/cfs-spool-web-darwin-arm64 cmd/web-server/main.go

# Create release packages
.PHONY: package
package: build-all
	@echo "Creating release packages..."
	@mkdir -p releases
	
	# Windows AMD64
	@mkdir -p tmp/cfs-spool-windows-amd64
	@cp dist/cfs-spool-*-windows-amd64.exe tmp/cfs-spool-windows-amd64/
	@cp -r web tmp/cfs-spool-windows-amd64/
	@cp README.md tmp/cfs-spool-windows-amd64/
	@echo '@echo off\necho Starting CFS Spool Web Interface...\nstart http://localhost:8080\ncfs-spool-web-windows-amd64.exe' > tmp/cfs-spool-windows-amd64/start.bat
	@cd tmp && zip -r ../releases/cfs-spool-windows-amd64.zip cfs-spool-windows-amd64/
	
	# Linux AMD64
	@mkdir -p tmp/cfs-spool-linux-amd64
	@cp dist/cfs-spool-*-linux-amd64 tmp/cfs-spool-linux-amd64/
	@cp -r web tmp/cfs-spool-linux-amd64/
	@cp README.md tmp/cfs-spool-linux-amd64/
	@echo '#!/bin/bash\necho "Starting CFS Spool Web Interface..."\nxdg-open http://localhost:8080 &\n./cfs-spool-web-linux-amd64' > tmp/cfs-spool-linux-amd64/start.sh
	@chmod +x tmp/cfs-spool-linux-amd64/*.sh tmp/cfs-spool-linux-amd64/cfs-spool-*
	@cd tmp && tar -czf ../releases/cfs-spool-linux-amd64.tar.gz cfs-spool-linux-amd64/
	
	# macOS Universal (both Intel and Apple Silicon)
	@mkdir -p tmp/cfs-spool-macos
	@cp dist/cfs-spool-*-darwin-amd64 tmp/cfs-spool-macos/
	@cp dist/cfs-spool-*-darwin-arm64 tmp/cfs-spool-macos/
	@cp -r web tmp/cfs-spool-macos/
	@cp README.md tmp/cfs-spool-macos/
	@echo '#!/bin/bash\nARCH=$$(uname -m)\nif [[ "$$ARCH" == "arm64" ]]; then\n  echo "Starting CFS Spool Web Interface (Apple Silicon)..."\n  open http://localhost:8080\n  ./cfs-spool-web-darwin-arm64\nelse\n  echo "Starting CFS Spool Web Interface (Intel)..."\n  open http://localhost:8080\n  ./cfs-spool-web-darwin-amd64\nfi' > tmp/cfs-spool-macos/start.sh
	@chmod +x tmp/cfs-spool-macos/*.sh tmp/cfs-spool-macos/cfs-spool-*
	@cd tmp && tar -czf ../releases/cfs-spool-macos.tar.gz cfs-spool-macos/
	
	@rm -rf tmp
	@echo "Release packages created in releases/ directory"

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run web server
.PHONY: run-web
run-web: build-web
	./bin/cfs-spool-web

# Run CLI help
.PHONY: run-cli
run-cli: build-cli
	./bin/cfs-spool-cli

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/ dist/ releases/ tmp/

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t cfs-spool:$(VERSION) .

# Docker run
.PHONY: docker-run
docker-run: docker-build
	docker run --rm -p 8080:8080 --privileged -v /dev:/dev cfs-spool:$(VERSION)

# Install dependencies for cross-compilation (Ubuntu/Debian)
.PHONY: install-deps-ubuntu
install-deps-ubuntu:
	sudo apt-get update
	sudo apt-get install -y libpcsclite-dev gcc-multilib gcc-mingw-w64 gcc-aarch64-linux-gnu

# Install dependencies for cross-compilation (macOS)
.PHONY: install-deps-macos
install-deps-macos:
	brew install pcsc-lite

# Show help
.PHONY: help
help:
	@echo "CFS Spool Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build CLI and web server for current platform"
	@echo "  build-cli        - Build CLI tool for current platform"
	@echo "  build-web        - Build web server for current platform"
	@echo "  build-all        - Build for all supported platforms"
	@echo "  package          - Create release packages for all platforms"
	@echo "  test             - Run tests"
	@echo "  run-web          - Build and run web server"
	@echo "  run-cli          - Build and run CLI tool"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Build and run Docker container"
	@echo "  clean            - Clean build artifacts"
	@echo "  help             - Show this help"
	@echo ""
	@echo "Cross-compilation setup:"
	@echo "  install-deps-ubuntu - Install dependencies on Ubuntu/Debian"
	@echo "  install-deps-macos  - Install dependencies on macOS"
