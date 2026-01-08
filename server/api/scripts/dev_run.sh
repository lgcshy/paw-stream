#!/bin/bash
# Development run script with hot reload

set -e

cd "$(dirname "$0")/.."

echo "Starting PawStream API Server in development mode..."

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "Installing air for hot reload..."
    go install github.com/air-verse/air@latest
fi

# Run with air for hot reload
air
