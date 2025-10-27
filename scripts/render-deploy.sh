#!/bin/bash

# Render Deployment Script for FlexFume E-Commerce Backend
# This script helps prepare and deploy the application to Render

set -e

echo "================================"
echo "FlexFume Render Deployment"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if git is available
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed${NC}"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo -e "${YELLOW}Current branch: ${CURRENT_BRANCH}${NC}"
echo ""

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}⚠️  You have uncommitted changes${NC}"
    echo "Commit your changes before deploying:"
    echo "  git add ."
    echo "  git commit -m 'Your message'"
    exit 1
fi

echo -e "${GREEN}✓ All changes committed${NC}"
echo ""

# Verify build locally
echo "Building application locally to verify..."
if go build -o bin/server cmd/api/main.go; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

echo ""
echo "================================"
echo "Deployment Instructions"
echo "================================"
echo ""
echo "1. Ensure all environment variables are set in Render Dashboard:"
echo "   - Go to https://dashboard.render.com"
echo "   - Select your web service"
echo "   - Click 'Environment'"
echo "   - Verify all required variables are set"
echo ""
echo "2. Push to main branch to trigger automatic deployment:"
echo "   git push origin main"
echo ""
echo "3. Monitor deployment:"
echo "   - Go to https://dashboard.render.com"
echo "   - Select flexfume-ecom-backend"
echo "   - Click 'Logs' to view deployment progress"
echo ""
echo "4. After deployment, run migrations:"
echo "   - Use Render Shell or SSH"
echo "   - Run: go run cmd/migrate/main.go -cmd=up"
echo ""
echo -e "${GREEN}Ready to deploy!${NC}"
echo ""

# Ask for confirmation
read -p "Push to main branch now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Pushing to main..."
    git push origin main
    echo -e "${GREEN}✓ Pushed successfully${NC}"
    echo ""
    echo "Deployment started! Monitor at: https://dashboard.render.com"
else
    echo "Deployment cancelled."
fi
