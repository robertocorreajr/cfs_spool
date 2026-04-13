#!/bin/bash

# CFS Spool Release Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🏷️  CFS Spool Release Script${NC}"
echo "==============================="

# Check if we're on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${RED}❌ Error: Must be on main branch to create release${NC}"
    echo "Current branch: $CURRENT_BRANCH"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}❌ Error: Uncommitted changes detected${NC}"
    echo "Please commit all changes before creating a release"
    exit 1
fi

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}📋 Current version: ${CURRENT_VERSION}${NC}"

# Ask for new version
echo -e "${YELLOW}📝 Enter new version (e.g., v1.0.0, v1.2.3):${NC}"
read -r NEW_VERSION

# Validate version format
if [[ ! $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Error: Invalid version format. Use vX.Y.Z (e.g., v1.0.0)${NC}"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$NEW_VERSION" >/dev/null 2>&1; then
    echo -e "${RED}❌ Error: Tag $NEW_VERSION already exists${NC}"
    exit 1
fi

echo -e "${YELLOW}📝 Enter release description (optional):${NC}"
read -r RELEASE_DESCRIPTION

# Show summary
echo ""
echo -e "${BLUE}📋 Release Summary:${NC}"
echo "  Version: $NEW_VERSION"
echo "  Description: ${RELEASE_DESCRIPTION:-'(no description)'}"
echo ""

# Confirm
echo -e "${YELLOW}❓ Create release? (y/N):${NC}"
read -r CONFIRM

if [[ $CONFIRM != "y" && $CONFIRM != "Y" ]]; then
    echo -e "${RED}❌ Release cancelled${NC}"
    exit 1
fi

# Update README with new version (if needed)
if grep -q "🏷️ CFS Spool v" README.md; then
    sed -i.bak -E "s/\*\*🏷️ CFS Spool v[0-9]+\.[0-9]+\.[0-9]+[^\*]*\*\*/\*\*🏷️ CFS Spool ${NEW_VERSION}\*\*/" README.md
    rm README.md.bak 2>/dev/null || true
    git add README.md
fi

if grep -q "🏷️ CFS Spool v" README.pt-BR.md; then
    sed -i.bak -E "s/\*\*🏷️ CFS Spool v[0-9]+\.[0-9]+\.[0-9]+[^\*]*\*\*/\*\*🏷️ CFS Spool ${NEW_VERSION}\*\*/" README.pt-BR.md
    rm README.pt-BR.md.bak 2>/dev/null || true
    git add README.pt-BR.md
fi

# Create and push tag
echo -e "${BLUE}🏷️  Creating tag...${NC}"
if [ -n "$RELEASE_DESCRIPTION" ]; then
    git tag -a "$NEW_VERSION" -m "$RELEASE_DESCRIPTION"
else
    git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
fi

echo -e "${BLUE}📤 Pushing tag to remote...${NC}"
git push origin "$NEW_VERSION"

echo -e "${GREEN}✅ Release $NEW_VERSION created successfully!${NC}"
echo ""
echo -e "${BLUE}🔗 Next steps:${NC}"
echo "  1. GitHub Actions will automatically build native binaries for all platforms (Wails)"
echo "  2. Release will be available at: https://github.com/robertocorreajr/cfs_spool/releases/tag/$NEW_VERSION"
echo ""
echo -e "${YELLOW}⏳ Monitor the build at: https://github.com/robertocorreajr/cfs_spool/actions${NC}"
