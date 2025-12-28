#!/bin/bash
set -e

echo "=========================================="
echo "Cloudflare Pages Deploy Script"
echo "=========================================="
echo "Timestamp: $(date)"
echo "Working directory: $(pwd)"
echo "User: $(whoami)"
echo ""

echo "Checking if public/ directory exists..."
if [ -d "public" ]; then
    echo "✓ public/ directory found"
    echo "Contents of public/:"
    ls -la public/ | head -20
    echo ""
    echo "File count in public/:"
    find public -type f | wc -l
    echo ""
else
    echo "✗ ERROR: public/ directory not found!"
    exit 1
fi

echo "Checking public/index.html..."
if [ -f "public/index.html" ]; then
    echo "✓ public/index.html exists"
    echo "File size: $(du -h public/index.html | cut -f1)"
else
    echo "✗ WARNING: public/index.html not found"
fi

echo ""
echo "Note: Cloudflare Pages automatically deploys the public/ directory"
echo "This script is a no-op to satisfy the required deploy command field"
echo ""
echo "Deploy script completed successfully"
echo "=========================================="

exit 0

