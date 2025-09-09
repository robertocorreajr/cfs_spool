# Automatic Tagging System

This document describes in detail the automatic tagging system implemented for the CFS Spool project, which uses GitHub Actions to create automatic tags following Semantic Versioning (SemVer).

## Overview

The system automates the process of creating project versions (tags) using a GitHub Actions workflow that:

1. Monitors commits on the main branch (`main`)
2. Detects specific keywords in commit messages (`#patch`, `#minor`, `#major`)
3. Increments the version according to the keyword found
4. Creates a new tag with the incremented version
5. Triggers the build workflow to automatically generate a new release

## Semantic Versioning

The system follows the [Semantic Versioning 2.0.0](https://semver.org/) standard, where:

- **Patch Version (z)**: Fixes bugs while maintaining backward compatibility
- **Minor Version (y)**: Adds functionality while maintaining backward compatibility
- **Major Version (x)**: Makes changes incompatible with previous versions

## How to Use

### Automatic Increment via Commit Messages

To trigger version increments, add one of the following hashtags at the end of the commit message:

- `#patch`: Increments the patch version
  ```bash
  git commit -m "Fix RFID reading issue #patch"
  ```

- `#minor`: Increments the minor version
  ```bash
  git commit -m "Add support for new RFID types #minor"
  ```

- `#major`: Increments the major version
  ```bash
  git commit -m "Refactor integration API for new readers #major"
  ```

### Recommended Workflow

1. Develop on feature/bugfix branches
2. When complete, create a Pull Request to the `main` branch
3. In the PR merge message, include the appropriate hashtag (#patch, #minor, #major)
4. After merging, the workflow will be triggered automatically

### Manual Triggering

You can also trigger the workflow manually through the GitHub interface:

1. Go to the "Actions" tab in the repository
2. Select the "Auto Tag" workflow
3. Click on "Run workflow"
4. Select the `main` branch and the desired increment type
5. Click on "Run workflow"

## How It Works

The system consists of two main workflows:

### 1. Automatic Tagging Workflow (auto-tag.yml)

This workflow is responsible for:
- Monitoring commits on the `main` branch
- Detecting keywords in commit messages
- Calculating the new version
- Creating a new tag with the incremented version
- Triggering the build workflow

### 2. Build Workflow (build.yml)

This workflow is responsible for:
- Monitoring the creation of new tags
- Compiling the project
- Creating a new release on GitHub
- Attaching compiled artifacts to the release

## Troubleshooting

### Workflow Not Being Triggered

- Check if the commit was made on the `main` branch
- Make sure the message contains exactly one of the hashtags: `#patch`, `#minor`, or `#major`
- Check the workflow logs in GitHub Actions

### Tag Not Being Created

- Check if the workflow has permissions to create tags
- Make sure there are no tags with the same version
- Check the workflow logs to understand the reason for the failure

### Build Not Being Triggered

- Check if the tag was created correctly
- Make sure the build workflow is configured to be triggered by tag creation
- Check the build workflow logs

## Utility Scripts

### Tagging Test Script

The `scripts/teste-auto-tag.sh` script is available to:
- Check the current tag
- Simulate version increments
- Provide commands to test the system

To use the script:
```bash
./scripts/teste-auto-tag.sh
```

## Customization

The workflows are defined in the files:
- `.github/workflows/auto-tag.yml`
- `.github/workflows/build.yml`

If necessary, you can modify these files to:
- Change the events that trigger the workflows
- Customize the tag format
- Modify the build process
- Add notifications or integrations
