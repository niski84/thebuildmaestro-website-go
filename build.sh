#!/bin/bash
set -e

echo "Building static site generator..."
go mod download
go build -o sitegen ./cmd/sitegen

echo "Generating static site..."
./sitegen -content static/content -template templates -static static -out public -domain https://thebuildmaestro.com

echo "Build complete! Output in public/ directory"


