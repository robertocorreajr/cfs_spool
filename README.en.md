# CFS Spool - Creality Filament Spool RFID Manager

Native desktop application (Wails v2 + React + shadcn/ui) for reading and writing Creality File System (CFS) RFID tags used in Creality 3D printer filament spools.

[![Portugues](https://img.shields.io/badge/lang-pt--BR-green.svg)](README.md)
[![English](https://img.shields.io/badge/lang-en-blue.svg)](README.en.md)

## Features

- **Native desktop app** -- no browser or server needed, runs directly on your OS
- **Enhanced color selector**: 35 predefined colors based on the Creality system + color picker for any custom hex color
- **Smart logic**: Auto-selection of supplier based on chosen material
- **Auto-fill**: Optional fields (Batch, Serial) with automatic padding
- **Visual reading**: Preview colors from existing tags
- **AES-ECB encryption/decryption**: Full support for Creality encryption system
- **Robust authentication**: Multiple fallback methods for RFID reading
- **Key derivation**: Complete algorithm based on tag UID
- **Compatibility**: Works with new tags (FFFFFFFFFFFF) and used tags (derived key)

## Screenshots

| | |
|:---:|:---:|
| ![Main screen](docs/screenshots/app-main-screen.png) | ![Color picker](docs/screenshots/app-color-picker.png) |
| *Main screen* | *Color picker* |

## Installation

### Ready Downloads (Recommended)

Download the latest release for your platform:

**[Releases - GitHub](https://github.com/robertocorreajr/cfs_spool/releases/latest)**

| Platform | File | Format |
|:---:|:---:|:---:|
| macOS (Apple Silicon) | `cfs-spool-darwin-arm64.dmg` | DMG (drag to Applications) |
| Linux (x86_64) | `cfs-spool-linux-amd64.zip` | ZIP (extract and run) |
| Windows (x86_64) | `cfs-spool-windows-amd64.zip` | ZIP (extract and run) |

### macOS: Bypass Gatekeeper

The app is not signed with an Apple Developer certificate, so macOS will block it on first launch. Use one of the methods below to allow it:

#### Method 1: System Settings (recommended)

1. Open the DMG and drag **CFS Spool** to the **Applications** folder
2. Try opening the app normally (double-click) — it will be blocked
3. **Important**: In the blocking dialog, click **"OK"**. **Do not click "Move to Trash"**, as this will delete the app and you'll need to drag it from the DMG again
4. Open **System Settings** → **Privacy & Security**
5. Under the "Security" section, you'll see a message about CFS Spool
6. Click **"Open Anyway"**
7. Confirm in the next dialog

![Gatekeeper blocked](docs/screenshots/macos-gatekeeper-blocked.png)
![Privacy & Security](docs/screenshots/macos-privacy-security.png)

#### Method 2: Terminal

Run the following command in Terminal to remove the quarantine attribute:

```bash
xattr -cr /Applications/CFS\ Spool.app
```

Then open the app normally from Launchpad or the Applications folder.

### Build from Source

#### Prerequisites

- **Go 1.24+**
- **Node.js 18+**
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Compatible RFID reader** (tested with ACR122U)
- **PC/SC headers**:
  - macOS: built-in (no action needed)
  - Linux: `sudo apt install pcscd libpcsclite-dev libgtk-3-dev libwebkit2gtk-4.1-dev`
  - Windows: built-in (winscard)

#### Build Commands

```bash
git clone https://github.com/robertocorreajr/cfs_spool.git
cd cfs_spool

# Development mode (hot-reload)
wails dev

# Production build
wails build

# Cross-compile all platforms
make build-all
```

The production binary is output to `build/bin/`.

## Usage

1. Connect your ACR122U RFID reader
2. Launch the application (double-click or run from terminal)
3. **Read Tag**: Place a tag on the reader and click "Read Tag"
4. **Write Tag**: Fill in the fields and click "Write Tag"

### Color Selection

You have full flexibility for choosing colors:
- Click one of the **35 predefined colors** from the palette
- Type any **6-digit hex code** manually in the text field
- Click the **color picker** square to open a full color spectrum selector

### Smart Auto-fill

- Empty batch automatically becomes `000`
- Empty serial automatically becomes `000001`
- Automatic padding with leading zeros

### Smart Logic

- Generic material selects Generic supplier automatically
- Creality material selects 0276 supplier (Creality) automatically
- Material list is filtered by selected supplier

## Supported Hardware

### Recommended Hardware (Affiliate Links)

#### AliExpress (International)
- **[ACR122U RFID Reader](https://s.click.aliexpress.com/e/_ok8qAl9)** -- Reader used in development (compatibility guaranteed)
- **[MIFARE Classic 1K Tags](https://s.click.aliexpress.com/e/_oBPVnEb)** -- Compatible tags tested in the project

#### Mercado Livre (Brazil)
- **[RFID Reader/Writer](https://meli.la/13HiRy2)** -- National option for purchasing the RFID reader

### Tested RFID Readers
- **ACR122U** (recommended)
- **Other PC/SC readers** (compatibility not guaranteed)

### Supported Tags
- **MIFARE Classic 1K**
- **MIFARE Classic 4K**
- **Creality CFS Tags**

## Development

### Project Structure

```
cfs_spool/
├── main.go                 # Wails app entry point
├── app.go                  # App struct with Wails-bound methods
├── app_options.go          # Options handler (materials, vendors)
├── wails.json              # Wails configuration
├── frontend/               # React + shadcn/ui frontend
│   ├── src/                # React components and pages
│   ├── package.json        # Node.js dependencies
│   ├── tailwind.config.js  # Tailwind CSS config
│   └── vite.config.ts      # Vite bundler config
├── internal/
│   ├── creality/           # Creality-specific logic
│   │   ├── crypto.go       # AES-ECB cryptography
│   │   └── fields.go       # Field parsing and formatting
│   └── rfid/               # RFID communication
│       └── reader.go       # PC/SC interface
├── build/                  # Wails build output and assets
├── tests/                  # Diagnostic tools
├── assets/                 # Visual resources
├── .github/workflows/      # CI/CD pipelines
├── Makefile                # Build shortcuts
└── go.mod                  # Go module definition
```

### Local Development

```bash
# Run in dev mode with hot-reload
wails dev

# Run Go tests
go test -v ./internal/... ./...

# RFID diagnostic tools (require hardware)
go run tests/test_read_diagnosis.go
go run tests/test_auth_read.go
go run tests/test_decode_cfs.go
```

### Dependencies

- `github.com/wailsapp/wails/v2` -- Desktop application framework
- `github.com/ebfe/scard` -- PC/SC interface for RFID communication
- `crypto/aes` -- AES cryptography (Go standard library)
- React + shadcn/ui + Tailwind CSS (frontend)

## Technical Reference

### Known Vendors

| **Vendor Code** | **Brand / Notes** |
|:---:|:---:|
| 0x0276 | Creality, Hyper, Ender, HP (official lines) |
| 0xFFFF | Generic (any non-official manufacturer) |

### Known Materials

| **Material Code** | **Description** |
|:---:|:---:|
| 00001 | Generic PLA |
| 00002 | Generic PLA-Silk |
| 00003 | Generic PETG |
| 00004 | Generic ABS |
| 00005 | Generic TPU |
| 00006 | Generic PLA-CF |
| 00007 | Generic ASA |
| 00008 | Generic PA |
| 00009 | Generic PA-CF |
| 00010 | Generic BVOH |
| 00011 | Generic PVA |
| 00012 | Generic HIPS |
| 00013 | Generic PET-CF |
| 00014 | Generic PETG-CF |
| 00015 | Generic PA6-CF |
| 00016 | Generic PAHT-CF |
| 00017 | Generic PPS |
| 00018 | Generic PPS-CF |
| 00019 | Generic PP |
| 00020 | Generic PET |
| 00021 | Generic PC |
| 01001 | Hyper PLA |
| 02001 | Hyper PLA-CF |
| 03001 | Hyper ABS |
| 04001 | CR-PLA |
| 05001 | CR-Silk |
| 06001 | CR-PETG |
| 06002 | Hyper PETG |
| 07001 | CR-ABS |
| 08001 | Ender-PLA |
| 09001 | EN-PLA+ |
| 09002 | Ender Fast PLA |
| 10001 | HP-TPU |
| 11001 | CR-Nylon |
| 13001 | CR-PLA Carbon |
| 14001 | CR-PLA Matte |
| 15001 | CR-PLA Fluo |
| 16001 | CR-TPU |
| 17001 | CR-Wood |
| 18001 | HP Ultra PLA |
| 19001 | HP-ASA |

### CFS Tag Format

The Creality CFS system stores data in sectors 1-2 of MIFARE Classic tags:

- **Sector 1 (Blocks 4-6)**: Encrypted filament data
- **Encryption**: AES-ECB with UID-derived keys
- **S1 Key**: Derived from UID using key `q3bu^t1nqfZ(pf$1`
- **Payload**: Decrypted with key `H@CFkRnz@KAtBJp2`

#### Field Layout (38 bytes)

```
Date(5) + Supplier(4) + Batch(2) + Material(5) + Color(7) + Length(4) + Serial(6) + Reserve(4)
```

- Color format: `"0" + 6-char hex` (e.g., `"077BB41"`)
- Batch defaults to `"A2"`, Reserve defaults to `"0000"`

#### Authentication Algorithm

1. **New tags**: Key A = `FFFFFFFFFFFF` (MIFARE default)
2. **Used tags**: Key A = derived from UID using AES algorithm
3. **Fallback**: Multiple attempts with different methods

### Predefined Color Palette

The interface includes 35 predefined colors based on the Creality system:

| Category | Colors |
|---|---|
| **Blues** | #25C4DA, #0099A7, #0B359A, #0A4AB6, #11B6EE, #90C6F5 |
| **Oranges/Yellows** | #FA7C0C, #F7B30F, #E5C20F, #B18F2E, #F8E911, #F6D311 |
| **Browns** | #8D766D, #6C4E43 |
| **Reds/Pinks** | #E62E2E, #EE2862, #EA2A2B, #E83D89, #AE2E65 |
| **Purples** | #611C8B, #8D60C7, #B287C9 |
| **Greens** | #006764, #018D80, #42B5AE, #1D822D, #54B351, #72E115 |
| **Grays** | #474747, #668798, #B1BEC6, #58636E |
| **Special** | #F2EFCE, #FFFFFF, #000000 |

## Releases and Versioning

Version tagging is automatic via `.github/workflows/auto-tag.yml`:

- **Auto-incrementing**: Patch version increases automatically on push to main
- **Semantic Versioning**: Control version type using commit message flags:
  - `git commit -m "Message #patch"` -- increment patch (v1.0.0 -> v1.0.1)
  - `git commit -m "Message #minor"` -- increment minor (v1.0.0 -> v1.1.0)
  - `git commit -m "Message #major"` -- increment major (v1.0.0 -> v2.0.0)
- **Manual trigger**: Available through GitHub Actions interface

Each `v*` tag automatically builds native binaries for macOS, Linux, and Windows via Wails.

## FAQ

### How to choose custom colors?

You have full flexibility:
- Choose one of the 35 predefined colors by clicking on the palette
- Type any hex code manually in the text field (6 digits)
- Click the colored square to open a color picker and choose any color from the spectrum

### Optional fields don't work?

The **Batch** and **Serial** fields are optional:
- Empty batch automatically becomes `000`
- Empty serial automatically becomes `000001`
- Automatic padding with leading zeros

### How to diagnose reading problems?

```bash
go run tests/test_read_diagnosis.go
```

This command systematically tests all authentication methods.

## Contributing

Contributions are welcome! Please:

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## License

This project is under the MIT license. See details in each source file.

## Disclaimer

This project is developed for educational and interoperability purposes. It is not affiliated with Creality 3D Technology Co., Ltd.
