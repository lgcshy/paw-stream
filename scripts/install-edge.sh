#!/usr/bin/env bash
set -euo pipefail

# PawStream Edge Client Installer
# Usage: curl -sSL https://raw.githubusercontent.com/lgcshy/paw-stream/dev/scripts/install-edge.sh | bash

REPO="lgcshy/paw-stream"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="pawstream-edge"
CONFIG_DIR="/etc/pawstream"
BINARY_NAME="edge-client"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux)  OS="linux" ;;
        darwin) OS="darwin" ;;
        *)      error "Unsupported OS: $OS" ;;
    esac

    case "$ARCH" in
        x86_64|amd64)   ARCH="amd64" ;;
        aarch64|arm64)  ARCH="arm64" ;;
        armv7l|armhf)   ARCH="arm" ;;
        *)              error "Unsupported architecture: $ARCH" ;;
    esac

    info "Detected platform: ${OS}/${ARCH}"
}

# Check for root
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        error "This script must be run as root (use sudo)"
    fi
}

# Get latest release tag
get_latest_version() {
    VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' | head -1 | cut -d'"' -f4) || true

    if [ -z "${VERSION:-}" ]; then
        warn "Could not determine latest version, using 'dev'"
        VERSION="dev"
    fi
    info "Version: ${VERSION}"
}

# Download binary
download_binary() {
    local url="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}"
    local tmp=$(mktemp)

    info "Downloading from: ${url}"
    if ! curl -sSfL -o "$tmp" "$url"; then
        error "Download failed. Check that the release exists for your platform (${OS}/${ARCH})"
    fi

    chmod +x "$tmp"
    mv "$tmp" "${INSTALL_DIR}/${BINARY_NAME}"
    info "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Create config directory
setup_config() {
    if [ ! -d "$CONFIG_DIR" ]; then
        mkdir -p "$CONFIG_DIR"
        info "Created config directory: ${CONFIG_DIR}"
    fi

    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        cat > "${CONFIG_DIR}/config.yaml" <<'YAML'
# PawStream Edge Client Configuration
# Edit this file, then restart the service:
#   sudo systemctl restart pawstream-edge

server:
  # Your PawStream API server URL
  url: "http://localhost:3000"

device:
  # Device ID and secret (from the PawStream web UI)
  id: ""
  secret: ""

stream:
  # Video input source
  input_type: "v4l2"           # v4l2, rtsp, file, test
  input_path: "/dev/video0"    # device path, RTSP URL, or file path

  # Encoding settings
  video_codec: "libx264"
  video_bitrate: 2000          # kbps
  video_width: 1280
  video_height: 720
  video_framerate: 30

  # Reconnection
  reconnect_interval: "5s"
  max_reconnect_attempts: 0    # 0 = unlimited

web:
  # Setup Wizard web UI
  enabled: true
  port: 8088
YAML
        info "Created default config: ${CONFIG_DIR}/config.yaml"
        warn "Edit ${CONFIG_DIR}/config.yaml with your device credentials before starting"
    else
        info "Config already exists: ${CONFIG_DIR}/config.yaml"
    fi
}

# Install systemd service
install_service() {
    if [ ! -d /run/systemd/system ]; then
        warn "systemd not detected, skipping service installation"
        return
    fi

    cat > "/etc/systemd/system/${SERVICE_NAME}.service" <<EOF
[Unit]
Description=PawStream Edge Client
Documentation=https://github.com/${REPO}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/${BINARY_NAME} --config ${CONFIG_DIR}/config.yaml
Restart=always
RestartSec=5
WorkingDirectory=${CONFIG_DIR}

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${CONFIG_DIR}
PrivateTmp=true

# Allow access to video devices
SupplementaryGroups=video

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "${SERVICE_NAME}"
    info "Systemd service installed and enabled"
    info "  Start:   sudo systemctl start ${SERVICE_NAME}"
    info "  Status:  sudo systemctl status ${SERVICE_NAME}"
    info "  Logs:    sudo journalctl -u ${SERVICE_NAME} -f"
}

# Main
main() {
    echo ""
    echo "  🐾 PawStream Edge Client Installer"
    echo "  ===================================="
    echo ""

    check_root
    detect_platform
    get_latest_version
    download_binary
    setup_config
    install_service

    echo ""
    info "Installation complete!"
    echo ""
    echo "  Next steps:"
    echo "  1. Edit config:  sudo nano ${CONFIG_DIR}/config.yaml"
    echo "  2. Set device ID and secret from the PawStream web UI"
    echo "  3. Start:        sudo systemctl start ${SERVICE_NAME}"
    echo "  4. Setup Wizard: http://$(hostname -I | awk '{print $1}'):8088"
    echo ""
}

main "$@"
