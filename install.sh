#!/bin/bash

# Finton Text Editor Installation Script

set -e

echo "Building Finton..."
go build -o finton cmd/finton/main.go

echo "Installing to /usr/local/bin..."
sudo mv finton /usr/local/bin/

echo ""
echo "✓ Finton installed successfully!"
echo ""
echo "Usage: finton <filename>"
echo ""
echo "Example:"
echo "  finton myfile.txt"
echo ""
echo "Configuration: ~/.finton.yaml (auto-created on first run)"
echo ""
echo "Key bindings:"
echo "  Ctrl+S - Save"
echo "  Ctrl+Q - Quit"
echo "  Ctrl+A - AI prompt"
echo "  Ctrl+H - Toggle theme"
echo "  Ctrl+F - Format (JSON)"
echo ""
