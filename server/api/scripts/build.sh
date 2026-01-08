#!/bin/bash
# Build script for production binary

set -e

cd "$(dirname "$0")/.."

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

echo "Building PawStream API Server..."
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"

# Build with CGO disabled for static linking
CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
    -o bin/api \
    ./cmd/api

echo "Build complete: bin/api"
