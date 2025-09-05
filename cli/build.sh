#!/bin/bash

# Build script for OrcaAI CLI

# Set the output directory
OUTPUT_DIR="./bin"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Build for different platforms
echo "Building OrcaAI CLI..."

# Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/orcaai-windows-amd64.exe" .

# Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/orcaai-linux-amd64" .

# macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o "$OUTPUT_DIR/orcaai-darwin-amd64" .

echo "Build complete! Binaries are located in the $OUTPUT_DIR directory."