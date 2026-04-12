# CFS Spool Makefile — Wails v2
VERSION ?= dev

.PHONY: all dev build build-all test clean install-frontend install-wails install-deps-ubuntu install-deps-macos help

# Default target
all: build

# Desenvolvimento com hot-reload
dev:
	wails dev

# Build para plataforma atual
build:
	wails build -ldflags "-X main.version=$(VERSION)"

# Build para todas as plataformas
build-all:
	wails build -platform darwin/arm64 -ldflags "-X main.version=$(VERSION)"
	wails build -platform darwin/amd64 -ldflags "-X main.version=$(VERSION)"
	wails build -platform linux/amd64 -ldflags "-X main.version=$(VERSION)"
	wails build -platform windows/amd64 -ldflags "-X main.version=$(VERSION)"

# Testes
test:
	go test -v ./...

# Limpar artefatos
clean:
	rm -rf build/bin/ frontend/dist/ frontend/node_modules/

# Instalar dependências do frontend
install-frontend:
	cd frontend && npm install

# Instalar Wails CLI
install-wails:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Dependências do sistema (Ubuntu/Debian)
install-deps-ubuntu:
	sudo apt-get update
	sudo apt-get install -y pcscd libpcsclite-dev libgtk-3-dev libwebkit2gtk-4.1-dev

# Dependências do sistema (macOS)
install-deps-macos:
	brew install pcsc-lite

# Ajuda
help:
	@echo "CFS Spool — Wails v2 + React + shadcn/ui"
	@echo ""
	@echo "Targets:"
	@echo "  dev              - Desenvolvimento com hot-reload"
	@echo "  build            - Build para plataforma atual"
	@echo "  build-all        - Build para todas as plataformas"
	@echo "  test             - Executar testes"
	@echo "  clean            - Limpar artefatos"
	@echo "  install-frontend - Instalar dependências do frontend"
	@echo "  install-wails    - Instalar Wails CLI"
	@echo "  install-deps-*   - Instalar dependências do sistema"
