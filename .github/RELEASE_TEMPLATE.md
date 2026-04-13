# CFS Spool Release Template

## What's New in v{VERSION}

### New Features
- Feature 1 description
- Feature 2 description

### Bug Fixes
- Fix 1 description
- Fix 2 description

### Improvements
- Improvement 1 description
- Improvement 2 description

### Download

Choose the appropriate package for your system:

| Platform | File |
|----------|------|
| macOS (Apple Silicon) | `cfs-spool-darwin-arm64.dmg` |
| Linux (x86_64) | `cfs-spool-linux-amd64.zip` |
| Windows (x86_64) | `cfs-spool-windows-amd64.zip` |

### Requirements

- **RFID Reader**: ACR122U or compatible PC/SC reader
- **Operating System**: Windows 10+, macOS 10.15+, or Linux with PC/SC lite
- **Permissions**: May require running as administrator/sudo for RFID access

### Quick Start

1. Download the appropriate package for your system
2. Install:
   - **macOS**: Open the DMG and drag to Applications
   - **Linux/Windows**: Extract the ZIP and run the `cfs-spool` binary
3. Connect your ACR122U RFID reader
4. Launch the application

#### macOS: Security Warning

The app is not signed with an Apple Developer certificate, so macOS may block it.
Run in Terminal to allow:
```
xattr -cr /Applications/CFS\ Spool.app
```
Or right-click the app > "Open" to authorize manually.

### Troubleshooting

**RFID Reader Not Detected:**
- Ensure PC/SC service is running
- Check USB connection
- Try running as administrator/sudo

**Permission Denied:**
- Run as administrator (Windows) or with sudo (Linux/macOS)
- Check that your user is in the `scard` group (Linux)

### Support

- **Issues**: https://github.com/robertocorreajr/cfs_spool/issues
- **Documentation**: https://github.com/robertocorreajr/cfs_spool#readme
