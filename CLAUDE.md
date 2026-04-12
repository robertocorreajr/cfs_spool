# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**CFS Spool** is a native desktop application (Wails v2 + React + shadcn/ui) for reading and writing Creality File System (CFS) RFID tags used in Creality 3D printer filament spools. It manages spool metadata encrypted on MIFARE Classic RFID tags via PC/SC.

## Commands

### Build & Run

```bash
# Desenvolvimento com hot-reload
wails dev

# Build para plataforma atual
wails build

# Build com versão
wails build -ldflags "-X main.version=v3.0.0"

# Makefile shortcuts
make dev          # wails dev (hot-reload)
make build        # wails build (plataforma atual)
make build-all    # Build para macOS, Linux, Windows
make test         # go test -v ./...
make clean        # Limpar artefatos
```

### Tests

```bash
make test                              # go test -v ./...
go run tests/test_read_diagnosis.go    # RFID read diagnostics (requires hardware)
go run tests/test_auth_read.go         # Authentication testing (requires hardware)
go run tests/test_decode_cfs.go        # CFS decoding from known raw bytes
```

### Dependencies

```bash
make install-wails         # go install wails CLI
make install-frontend      # cd frontend && npm install
make install-deps-macos    # brew install pcsc-lite
make install-deps-ubuntu   # apt install pcscd libpcsclite-dev libgtk-3-dev libwebkit2gtk-4.1-dev
```

## Architecture

### Stack

| Layer | Technology |
|-------|-----------|
| Desktop framework | Wails v2 |
| Backend | Go 1.24+ (CGO for PC/SC) |
| Frontend | React 18 + TypeScript + Vite |
| UI Components | shadcn/ui |
| Styling | Tailwind CSS |
| Color picker | react-colorful |
| RFID | github.com/ebfe/scard (PC/SC) |
| Communication | Wails bindings (Go ↔ JS direct, no REST) |

### Data Flow

```
Wails Window (native)
  └─ React + shadcn/ui + Tailwind CSS
       └─ Wails Bindings (direct Go calls)
            ├─ app.go → ReadTag() / WriteTag() / GetOptions() / ValidateColor()
            ├─ internal/rfid/reader.go (PC/SC APDU) → ACR122U → MIFARE Tag
            ├─ internal/creality/crypto.go (AES-ECB)
            └─ internal/creality/fields.go (encode/decode)
```

### Go Backend Bindings (app.go)

| Method | Purpose |
|--------|---------|
| `ReadTag()` | Read + decrypt tag → returns TagData with all fields |
| `WriteTag(data)` | Validate + encrypt + write tag |
| `GetOptions()` | Materials/vendor/length dropdown data |
| `ValidateColor(hex)` | Validate 6-char hex color, return uppercase |
| `GetVersion()` | Return app version string |

### Key Abstractions

**`internal/creality/crypto.go`** — Two-layer cryptography:
1. **S1 Key Derivation**: Tag UID repeated 4×, encrypted with key `q3bu^t1nqfZ(pf$1` → 6-byte MIFARE sector key
2. **Payload Encryption**: 48-byte payload (38 data + 10 padding), AES-ECB with key `H@CFkRnz@KAtBJp2`

**`internal/creality/fields.go`** — Fixed 38-byte field layout:
`Date(5) + Supplier(4) + Batch(2) + Material(5) + Color(7) + Length(4) + Serial(6) + Reserve(4)`
- Color format: `"0" + 6-char hex` (e.g., `"077BB41"`)
- Auto-fix rules: Batch defaults to `"A2"`, Reserve defaults to `"0000"`

**`internal/rfid/reader.go`** — PC/SC abstraction:
- Opens first available reader; sends raw APDU commands to ACR122U
- Read flow: `GET_DATA` (UID) → `LOAD_KEY` → `AUTHENTICATE` sector 1 → read blocks 4–6
- Write flow: authenticate sector 1 → write blocks 4, 5, 6

### Frontend Components (frontend/src/components/)

| Component | Purpose |
|-----------|---------|
| `SpoolForm.tsx` | Main form — read button, all fields, write button |
| `ColorPicker.tsx` | Hex input + react-colorful gradient + 35 preset swatches |
| `MaterialSelect.tsx` | Linked supplier/material dropdowns with auto-filtering |
| `LengthSelect.tsx` | Length dropdown + conditional custom gram input |
| `Header.tsx` | App title + version badge |

## Build Requirements

- **Go 1.24+** with `CGO_ENABLED=1` (required for PC/SC C bindings)
- **Node.js 18+** (for frontend build)
- **Wails CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- C compiler (`gcc` on Linux/macOS)
- PC/SC headers: `libpcsclite-dev` (Linux) or built-in (macOS)
- Linux also needs: `libgtk-3-dev`, `libwebkit2gtk-4.1-dev`
- RFID hardware: ACR122U reader + MIFARE Classic 1K/4K tags

## Release & Versioning

Versioning is automatic via `.github/workflows/auto-tag.yml`:
- Commit keywords: `#major`, `#minor`, default = patch increment
- Merging to `main` triggers auto-tag → triggers build workflow
- Build workflow uses `wails build` for each platform

## Adding Materials/Vendors

Material and vendor lists live in `app_options.go`. Edit the `materials`, `vendors`, or `lengths` slices to add new entries.
