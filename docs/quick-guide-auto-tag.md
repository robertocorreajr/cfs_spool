# Quick Guide: Automatic Tagging System

This document is a quick reference guide for the automatic tagging system implemented in the CFS Spool project.

## How It Works

The system uses GitHub Actions workflows to detect commit messages with special suffixes and automatically create tags following Semantic Versioning (SemVer).

## Types of Increments

When committing to the `main` branch, you can add one of the following suffixes to your message:

- **`#patch`**: Increments the patch version (z) in x.y.z
  - *Example*: `"Fix UI bug #patch"`
  - *Result*: v1.0.0 → v1.0.1

- **`#minor`**: Increments the minor version (y) in x.y.z
  - *Example*: `"Add new feature #minor"`
  - *Result*: v1.0.0 → v1.1.0

- **`#major`**: Increments the major version (x) in x.y.z
  - *Example*: `"Breaking API change #major"`
  - *Result*: v1.0.0 → v2.0.0

## Workflow

1. Make your changes in a feature/fix branch
2. When finished, merge the branch with `main`
3. In the merge commit message, add the appropriate suffix (#patch, #minor, #major)
4. The workflow will:
   - Detect the suffix in the message
   - Automatically increment the version
   - Create a new tag
   - Trigger the build and release workflow

## Manual Verification

To check if the system is working correctly:

1. Run the script `./scripts/teste-auto-tag.sh`
2. Check the latest tag created
3. Follow the indicated steps to test the workflow

## Manual Triggering

You can also trigger the workflow manually:

1. Go to the Actions tab on GitHub
2. Click on "Auto Tag"
3. Click on "Run workflow"
4. Select the desired increment type
5. Click on "Run workflow"

## Troubleshooting

If the system is not automatically creating tags:

1. Check if the commit was made to the `main` branch
2. Make sure the message contains the correct suffix (#patch, #minor, #major)
3. Check the workflow logs in the GitHub Actions tab
4. Make sure the workflow has permissions to create tags
