# Guia de Tagueamento Automático de Versões

Este documento explica como funciona o sistema automático de tagueamento semântico do projeto CFS Spool.

## O que é Versionamento Semântico?

O projeto segue o padrão [SemVer (Semantic Versioning)](https://semver.org/), where versions follow the format `vMAJOR.MINOR.PATCH`:

- **MAJOR**: incompatible changes with previous versions
- **MINOR**: addition of features that are compatible with previous versions
- **PATCH**: bug fixes that are compatible with previous versions

## Automated Tag System

The project has a GitHub Actions system that automatically creates version tags when code is pushed to the main branch (`main`).

### How It Works

1. When code is pushed to the `main` branch, the `auto-tag.yml` workflow is triggered
2. The workflow analyzes the most recent commit message to decide the type of version increment
3. A new tag is created and pushed to the repository with the incremented version

### Controlling the Type of Increment

By default, the **patch** version is incremented. To explicitly control the type of increment, add one of the flags below to your commit message:

| Flag | Increment | Example |
|------|------------|---------|
| `#patch` | Increments the patch number | `v1.0.0` → `v1.0.1` |
| `#minor` | Increments the minor number and resets patch | `v1.0.1` → `v1.1.0` |
| `#major` | Increments the major number and resets minor and patch | `v1.1.0` → `v2.0.0` |

#### Examples:

```bash
# Increment patch (default)
git commit -m "Fix color picker issue"

# Increment patch (explicit)
git commit -m "Fix color picker issue #patch"

# Increment minor
git commit -m "Add new color selector #minor"

# Increment major
git commit -m "Completely refactor API #major"
```

### Manual Trigger

It's also possible to trigger the tagging workflow manually through the GitHub interface:

1. Go to the "Actions" tab on GitHub
2. Select the "Auto Tag" workflow
3. Click on "Run workflow"
4. Choose the type of increment: `patch`, `minor`, or `major`
5. Confirm by clicking "Run workflow"

## Build Pipeline

Once the tag is created, the `build.yml` workflow is automatically triggered to:

1. Run automated tests
2. Build binaries for all platforms
3. Create native installers (DMG, AppImage, EXE)
4. Publish a new release on GitHub

## Tips for Commit Messages

- **Bug fixes**: use `#patch` or leave as default
- **New features**: use `#minor`
- **Large/incompatible changes**: use `#major`

## Recommended Flow

1. Develop on feature branches
2. Create Pull Requests to the `main` branch
3. After reviewing and approving the PR, merge to `main`
4. The automatic tagging system will create the new version
5. The CI/CD pipeline will automatically create installers for the new version

## Troubleshooting

If automatic tagging is not triggered or fails:

1. Check if the commit was made directly to the `main` branch
2. Make sure GitHub Actions are enabled in the repository
3. Check the logs in the "Actions" tab on GitHub
4. Use manual tagging through the GitHub interface as an alternative
