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

#### macOS: Aviso de Seguranca / Security Warning

O app nao e assinado com certificado Apple Developer. O macOS bloqueia a primeira execucao.
The app is not signed with an Apple Developer certificate. macOS will block first launch.

**Metodo 1 / Method 1 — Ajustes do Sistema / System Settings (recomendado / recommended):**
1. Tente abrir normalmente / Try opening normally
2. Ajustes do Sistema → Privacidade e Seguranca → "Abrir Mesmo Assim"
3. System Settings → Privacy & Security → "Open Anyway"

**Metodo 2 / Method 2 — Terminal:**
```
xattr -cr /Applications/CFS\ Spool.app
```

Para mais detalhes, veja o [README](https://github.com/robertocorreajr/cfs_spool#macos-permitir-execucao-gatekeeper).
For detailed instructions, see the [README](https://github.com/robertocorreajr/cfs_spool#macos-permitir-execucao-gatekeeper).

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
