# CFS Spool - Creality RFID Tag Reader/Writer

🏷️ **Complete system for reading and writing Creality File System (CFS) RFID tags**

[![English](https://img.shields.io/badge/lang-en-blue.svg)](README.md){:target="_blank"} [![Portuguese](https://img.shields.io/badge/lang-pt--BR-green.svg)](README.pt-BR.md){:target="_blank"}

## 📋 Description

CFS Spool is a complete Go application that provides both command-line and web interfaces for interacting with MIFARE Classic RFID tags used in Creality's filament system. The tool allows reading and writing filament spool information such as material, color, batch, manufacturing date, and other metadata stored encrypted on the tags.

## ✨ Features

### 🖥️ Web Interface (Main)
- 🎨 **Visual color selector**: 35 predefined colors palette with real-time preview
- 🧠 **Smart logic**: Auto-selection of supplier based on chosen material
- 📝 **Auto-fill**: Optional fields with automatic padding
- 📖 **Visual reading**: Preview colors from existing tags
- 🔄 **Responsive interface**: Works on desktop and mobile

### 📟 CLI (Command Line)
- 📖 **CFS tag reading**: Decodes complete filament information
- ✏️ **Tag writing**: Programs new tags with custom data
- 🔐 **AES-ECB decryption**: Full support for Creality encryption system
- 🔧 **Debug mode**: Technical data display for development
- 🔄 **Robust authentication**: Multiple fallback methods for reading

### 🛠️ Advanced Features
- 🎯 **Key derivation**: Complete algorithm based on tag UID
- 🔒 **Compatibility**: Works with new tags (FFFFFFFFFFFF) and used tags (derived key)
- 🧪 **Diagnostic tools**: Complete troubleshooting suite
- 📦 **Native installers**: DMG for macOS, AppImage for Linux, executable for Windows

## 🚀 Installation

### 📥 Ready Downloads (Recommended)

Download the latest native installers:

**[⬇️ Releases - GitHub](https://github.com/robertocorreajr/cfs_spool/releases/latest)**

- 🍎 **macOS**: `CFS-Spool-macOS.dmg` (drag-and-drop installer)
- 🐧 **Linux**: `CFS-Spool-Linux.AppImage` (portable)
- 🪟 **Windows**: `CFS-Spool-Windows.exe` (installer)

### 🛠️ Manual Compilation

#### Prerequisites

- **Go 1.21+**
- **Compatible RFID reader** (tested with ACR122U)
- **PC/SC Smart Card Daemon** 
  - macOS: already included
  - Linux: `sudo apt install pcscd libpcsclite-dev`
  - Windows: RFID reader driver

#### Compilation

```bash
git clone https://github.com/robertocorreajr/cfs_spool.git
cd cfs_spool

# Web Interface (main)
go build -o cfs-spool-web ./cmd/app

# Traditional CLI
go build -o cfs-spool-cli ./cmd/cfs-spool
```

## 📱 Usage

### 🖥️ Web Interface (Recommended)

1. **Run application**:
   ```bash
   ./cfs-spool-web
   # or on Windows: CFS-Spool.exe
   ```

2. **Access interface**: Browser opens automatically at `http://localhost:8080`

3. **Use interface**:
   - **"Read Tag" tab**: Place tag on reader and click "Read Tag"
   - **"Write Tag" tab**: Fill in data and click "Write Tag"

#### 🎨 Web Interface Features

- **Color palette**: 35 predefined colors with visual preview
- **Smart auto-fill**: 
  - Empty batch → `000`
  - Empty serial → `000001`
  - Auto-padding with leading zeros
- **Smart logic**:
  - Generic material → Generic supplier (automatic)
  - Creality material → 1B3D supplier (automatic)
  - Material filtering by supplier

### 📟 CLI Interface

```bash
# Basic reading
./cfs-spool-cli read-tag

# Debug mode (technical data)
./cfs-spool-cli read-tag -debug

# Tag writing
./cfs-spool-cli write-tag -batch "1A5" -material "04001" -color "FF40130"
```

### Output Example

```
╔══════════════════════════════════════════╗
║           TAG INFORMATION                ║
╚══════════════════════════════════════════╝
📦 Batch:       1A5
📅 Date:        January 20, 2024
🏭 Supplier:    1B3D
🧪 Material:    CR-PLA (standard)
🎨 Color:       #77BB41 (hex)
📏 Length:      330cm (1kg filament)
🔢 Serial:      000001
```

## 🛠️ Supported Hardware

### 🛒 Recommended Hardware (Affiliate Links)

- **🏷️ [ACR122U RFID Reader](https://s.click.aliexpress.com/e/_ok8qAl9){:target="_blank"}** - Reader used in development (compatibility guaranteed)
- **📇 [MIFARE Classic 1K Tags](https://s.click.aliexpress.com/e/_oBPVnEb){:target="_blank"}** - Compatible tags tested in the project

### Tested RFID Readers
- **ACR122U** ✅ (recommended)
- **Other PC/SC readers** (compatibility not guaranteed)

### Supported Tags
- **MIFARE Classic 1K** ✅
- **MIFARE Classic 4K** ✅
- **Creality CFS Tags** ✅

## 🔧 Development

### Project Structure

```
cfs-spool/
├── cmd/
│   ├── app/                # 🖥️ Web Interface (main)
│   │   └── main.go         # Web server with REST API
│   ├── cfs-spool/          # 📟 Traditional CLI
│   │   ├── main.go         # Command line interface
│   │   └── write_tag.go    # Read/write commands
│   └── web-server/         # (deprecated)
├── internal/
│   ├── creality/           # Creality-specific logic
│   │   ├── crypto.go       # AES-ECB cryptography
│   │   └── fields.go       # Field parsing and formatting
│   └── rfid/               # RFID communication
│       └── reader.go       # PC/SC interface
├── web/                    # 🎨 Web interface frontend
│   ├── index.html          # HTML/CSS/JS interface
│   └── favicon.svg         # Application icon
├── tests/                  # 🧪 Test tools
│   ├── test_auth_read.go   # Authentication test
│   ├── test_basic_read.go  # Basic reading test
│   ├── test_decode_cfs.go  # Decoding test
│   └── test_read_diagnosis.go # Complete diagnosis
├── assets/                 # 🎨 Visual resources
│   ├── icons/              # Icons for installers
│   └── dmg-background.svg  # macOS installer background
├── .github/workflows/      # 🚀 CI/CD
│   └── build.yml           # Automatic build pipeline
├── scripts/                # 📦 Release scripts
│   └── release.sh          # Packaging script
└── Dockerfile              # 🐳 Docker container
```

### REST API (Web Interface)

The web interface exposes a simple REST API:

- `GET /api/status` - Application status
- `GET /api/options` - Options for dropdowns (materials, suppliers, etc.)
- `POST /api/read-tag` - RFID tag reading
- `POST /api/write` - RFID tag writing

### Dependencies

- `github.com/ebfe/scard` - PC/SC interface for RFID communication
- `crypto/aes` - AES cryptography (standard library)
- Native web interface (no external dependencies)

### 🧪 Diagnostic Tools

```bash
# Complete RFID reading diagnosis
go run tests/test_read_diagnosis.go

# Authentication test
go run tests/test_auth_read.go

# CFS decoding test
go run tests/test_decode_cfs.go
```

## 📊 Technical Reference

### Known Vendors

| **Vendor Code** | **Brand / Notes**                                  |
|:---------------:|:--------------------------------------------------:|
|  0x0276         | Creality • Hyper • Ender • HP (official lines)    |
|  0xFFFF         | Generic (any non-official manufacturer)            |

### Known Materials

| **Material Code** | **Description**       |
|:-----------------:|:---------------------:|
|  00001            | Generic PLA           |
|  00002            | Generic PLA-Silk      |
|  00003            | Generic PETG          |
|  00004            | Generic ABS           |
|  00005            | Generic TPU           |
|  00006            | Generic PLA-CF        |
|  00007            | Generic ASA           |
|  00008            | Generic PA            |
|  00009            | Generic PA-CF         |
|  00010            | Generic BVOH          |
|  00011            | Generic PVA           |
|  00012            | Generic HIPS          |
|  00013            | Generic PET-CF        |
|  00014            | Generic PETG-CF       |
|  00015            | Generic PA6-CF        |
|  00016            | Generic PAHT-CF       |
|  00017            | Generic PPS           |
|  00018            | Generic PPS-CF        |
|  00019            | Generic PP            |
|  00020            | Generic PET           |
|  00021            | Generic PC            |
|  01001            | Hyper PLA             |
|  02001            | Hyper PLA-CF          |
|  03001            | Hyper ABS             |
|  04001            | CR-PLA                |
|  05001            | CR-Silk               |
|  06001            | CR-PETG               |
|  06002            | Hyper PETG            |
|  07001            | CR-ABS                |
|  08001            | Ender-PLA             |
|  09001            | EN-PLA+               |
|  09002            | Ender Fast PLA        |
|  10001            | HP-TPU                |
|  11001            | CR-Nylon              |
|  13001            | CR-PLA Carbon         |
|  14001            | CR-PLA Matte          |
|  15001            | CR-PLA Fluo           |
|  16001            | CR-TPU                |
|  17001            | CR-Wood               |
|  18001            | HP Ultra PLA          |
|  19001            | HP-ASA                |

### CFS Tag Format

The Creality CFS system stores data in sectors 1-2 of MIFARE Classic tags:

- **Sector 1 (Blocks 4-6)**: Encrypted filament data
- **Encryption**: AES-ECB with UID-derived keys
- **S1 Key**: Derived from UID using key "q3bu^t1nqfZ(pf$1"
- **Payload**: Decrypted with key "H@CFkRnz@KAtBJp2"

#### Authentication Algorithm

1. **New tags**: Key A = `FFFFFFFFFFFF` (MIFARE default)
2. **Used tags**: Key A = derived from UID using AES algorithm
3. **Fallback**: Multiple attempts with different methods

## 🎨 Predefined Color Palette

The web interface includes 35 predefined colors based on the Creality system:

| Category | Colors |
|----------|--------|
| **Blues** | #25C4DA, #0099A7, #0B359A, #0A4AB6, #11B6EE, #90C6F5 |
| **Oranges/Yellows** | #FA7C0C, #F7B30F, #E5C20F, #B18F2E, #F8E911, #F6D311 |
| **Browns** | #8D766D, #6C4E43 |
| **Reds/Pinks** | #E62E2E, #EE2862, #EA2A2B, #E83D89, #AE2E65 |
| **Purples** | #611C8B, #8D60C7, #B287C9 |
| **Greens** | #006764, #018D80, #42B5AE, #1D822D, #54B351, #72E115 |
| **Grays** | #474747, #668798, #B1BEC6, #58636E |
| **Special** | #F2EFCE, #FFFFFF, #000000 |

## 🚀 Releases and Versioning

- **v1.2.0+**: Complete web interface with color palette
- **v1.1.1**: Critical fix in key derivation
- **v1.1.0**: First version with native installers
- **v1.0.x**: Basic CLI versions

### 📦 Automatic Build System

Each `v*` tag automatically generates:
- 🍎 DMG installer for macOS (with custom icon)
- 🐧 Portable AppImage for Linux
- 🪟 Windows executable with installer
- 🐳 Multi-architecture Docker image

## ❓ FAQ

### How to choose between CLI and Web Interface?

- **Web Interface**: Recommended for general use, more intuitive
- **CLI**: Ideal for automation, scripts, and development

### Is the color palette limited?

No! You can:
- Choose one of the 35 predefined colors (click on palette)
- Type any hex code manually in the text field
- Use the color picker (click on the colored square)

### Optional fields don't work?

The **Batch** and **Serial** fields are optional:
- Empty batch → automatically `000`
- Empty serial → automatically `000001`
- Automatic padding with leading zeros

### How to diagnose reading problems?

```bash
go run tests/test_read_diagnosis.go
```

This command systematically tests all authentication methods.

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

### 🔧 Local Development

```bash
# Web Interface
go run cmd/app/main.go

# CLI
go run cmd/cfs-spool/main.go read-tag

# Tests
go run tests/test_read_diagnosis.go
```

## 📄 License

This project is under MIT license. See details in each source file.

## ⚠️ Disclaimer

This project is developed for educational and interoperability purposes. It is not affiliated with Creality 3D Technology Co., Ltd.

---

**🏷️ CFS Spool v1.2.0+** - Complete system for Creality RFID tags  
*Developed with ❤️ in Go*
