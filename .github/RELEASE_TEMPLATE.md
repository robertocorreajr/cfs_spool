# CFS Spool Release Template

## 🎉 What's New in v{VERSION}

### ✨ New Features
- Feature 1 description
- Feature 2 description

### 🐛 Bug Fixes
- Fix 1 description
- Fix 2 description

### 🔧 Improvements
- Improvement 1 description
- Improvement 2 description

### 📦 Download

Choose the appropriate package for your system:

#### Windows
- **Windows 64-bit (Intel/AMD)**: `cfs-spool-windows-amd64.zip`
- **Windows ARM64**: `cfs-spool-windows-arm64.zip`

#### macOS
- **macOS Universal (Intel + Apple Silicon)**: `cfs-spool-macos.tar.gz`
- **macOS Intel only**: `cfs-spool-darwin-amd64.tar.gz`
- **macOS Apple Silicon only**: `cfs-spool-darwin-arm64.tar.gz`

#### Linux
- **Linux 64-bit (Intel/AMD)**: `cfs-spool-linux-amd64.tar.gz`
- **Linux ARM64**: `cfs-spool-linux-arm64.tar.gz`

### 🐳 Docker

```bash
# Pull and run the latest Docker image
docker pull ghcr.io/robertocorreajr/cfs_spool:{VERSION}
docker run --rm -p 8080:8080 --privileged -v /dev:/dev ghcr.io/robertocorreajr/cfs_spool:{VERSION}
```

### 📋 Requirements

- **RFID Reader**: ACR122U or compatible PC/SC reader
- **Operating System**: Windows 10+, macOS 10.15+, or Linux with PC/SC lite
- **Permissions**: May require running as administrator/sudo for RFID access

### 🚀 Quick Start

1. Download the appropriate package for your system
2. Extract the archive
3. Run the installation script:
   - Windows: Double-click `install.bat`
   - macOS/Linux: Run `./install.sh`
4. Connect your RFID reader
5. Run `cfs-spool-web-*` to start the web interface
6. Open http://localhost:8080 in your browser

### 💡 Usage Examples

#### Web Interface
```bash
# Start web server
./cfs-spool-web-{platform}

# Open browser to http://localhost:8080
```

#### Command Line
```bash
# Read a tag
./cfs-spool-cli-{platform} read-tag

# Write a tag
./cfs-spool-cli-{platform} write-tag --material "CR-PLA" --color "FF0000" --length 250
```

### 🆘 Troubleshooting

**RFID Reader Not Detected:**
- Ensure PC/SC service is running
- Check USB connection
- Try running as administrator/sudo

**Permission Denied:**
- Run as administrator (Windows) or with sudo (Linux/macOS)
- Check that your user is in the `scard` group (Linux)

**Web Interface Not Loading:**
- Ensure port 8080 is not in use
- Check firewall settings
- Verify the binary has execute permissions

### 📞 Support

- **Issues**: https://github.com/robertocorreajr/cfs_spool/issues
- **Documentation**: https://github.com/robertocorreajr/cfs_spool#readme
- **Discussions**: https://github.com/robertocorreajr/cfs_spool/discussions
