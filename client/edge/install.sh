#!/bin/bash
# PawStream Edge Client Installation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/pawstream"
LOG_DIR="/var/log/pawstream"
DATA_DIR="/var/lib/pawstream"
SERVICE_FILE="configs/systemd/pawstream-edge.service"
BINARY="build/edge-client"
USER="pawstream"
GROUP="pawstream"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
       log_error "This script must be run as root (use sudo)"
       exit 1
    fi
}

create_user() {
    if id "$USER" &>/dev/null; then
        log_info "User $USER already exists"
    else
        log_info "Creating user $USER..."
        useradd --system --no-create-home --shell /bin/false $USER
    fi
}

install_binary() {
    log_info "Installing binary to $INSTALL_DIR..."
    if [ ! -f "$BINARY" ]; then
        log_error "Binary not found: $BINARY"
        log_error "Please run 'make build' first"
        exit 1
    fi
    cp "$BINARY" "$INSTALL_DIR/edge-client"
    chmod +x "$INSTALL_DIR/edge-client"
}

create_directories() {
    log_info "Creating directories..."
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOG_DIR"
    mkdir -p "$DATA_DIR"
    
    chown -R $USER:$GROUP "$LOG_DIR"
    chown -R $USER:$GROUP "$DATA_DIR"
}

install_config() {
    log_info "Installing configuration..."
    if [ -f "$CONFIG_DIR/config.yaml" ]; then
        log_warn "Config file already exists, backing up to config.yaml.bak"
        cp "$CONFIG_DIR/config.yaml" "$CONFIG_DIR/config.yaml.bak"
    fi
    
    if [ ! -f "configs/config.example.yaml" ]; then
        log_error "Example config not found: configs/config.example.yaml"
        exit 1
    fi
    
    cp "configs/config.example.yaml" "$CONFIG_DIR/config.yaml"
    chmod 600 "$CONFIG_DIR/config.yaml"
    chown $USER:$GROUP "$CONFIG_DIR/config.yaml"
    
    log_warn "Please edit $CONFIG_DIR/config.yaml to configure your device"
}

install_service() {
    log_info "Installing systemd service..."
    if [ ! -f "$SERVICE_FILE" ]; then
        log_error "Service file not found: $SERVICE_FILE"
        exit 1
    fi
    
    cp "$SERVICE_FILE" /etc/systemd/system/pawstream-edge.service
    systemctl daemon-reload
}

main() {
    log_info "Installing PawStream Edge Client..."
    
    check_root
    create_user
    install_binary
    create_directories
    install_config
    install_service
    
    log_info "Installation complete!"
    echo ""
    log_info "Next steps:"
    echo "  1. Edit configuration: sudo nano $CONFIG_DIR/config.yaml"
    echo "  2. Enable service: sudo systemctl enable pawstream-edge"
    echo "  3. Start service: sudo systemctl start pawstream-edge"
    echo "  4. Check status: sudo systemctl status pawstream-edge"
    echo "  5. View logs: sudo journalctl -u pawstream-edge -f"
}

main "$@"
