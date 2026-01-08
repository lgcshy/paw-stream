#!/bin/bash
#
# PawStream Edge Client - 卸载脚本
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
APP_NAME="edge-client"
INSTALL_DIR="/opt/pawstream"
BIN_DIR="/usr/local/bin"
CONFIG_DIR="/etc/pawstream"
DATA_DIR="/var/lib/pawstream"
LOG_DIR="/var/log/pawstream"
SYSTEMD_DIR="/etc/systemd/system"
SERVICE_NAME="pawstream-edge"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "请使用 root 权限运行此脚本"
        echo "尝试: sudo $0"
        exit 1
    fi
}

# 确认卸载
confirm_uninstall() {
    echo ""
    echo "╔════════════════════════════════════════════════════╗"
    echo "║                                                    ║"
    echo "║     PawStream Edge Client - 卸载程序              ║"
    echo "║                                                    ║"
    echo "╚════════════════════════════════════════════════════╝"
    echo ""
    
    print_warning "即将卸载 PawStream Edge Client"
    echo ""
    echo "将要删除:"
    echo "  - 二进制文件: $INSTALL_DIR"
    echo "  - 符号链接:   $BIN_DIR/$APP_NAME"
    echo "  - Systemd 服务: $SERVICE_NAME"
    echo ""
    
    read -p "是否保留配置文件? [y/N]: " keep_config
    read -p "是否保留数据和日志? [y/N]: " keep_data
    echo ""
    
    read -p "确认卸载? [y/N]: " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        print_info "取消卸载"
        exit 0
    fi
}

# 停止并移除 systemd 服务
remove_systemd_service() {
    print_info "停止并移除 systemd 服务..."
    
    if systemctl is-active --quiet $SERVICE_NAME; then
        systemctl stop $SERVICE_NAME
        print_success "服务已停止"
    fi
    
    if systemctl is-enabled --quiet $SERVICE_NAME 2>/dev/null; then
        systemctl disable $SERVICE_NAME
        print_success "禁用开机自启"
    fi
    
    if [ -f "$SYSTEMD_DIR/$SERVICE_NAME.service" ]; then
        rm -f "$SYSTEMD_DIR/$SERVICE_NAME.service"
        systemctl daemon-reload
        print_success "Systemd 服务已移除"
    fi
}

# 移除二进制文件
remove_binary() {
    print_info "移除二进制文件..."
    
    if [ -L "$BIN_DIR/$APP_NAME" ]; then
        rm -f "$BIN_DIR/$APP_NAME"
        print_success "符号链接已移除"
    fi
    
    if [ -d "$INSTALL_DIR" ]; then
        rm -rf "$INSTALL_DIR"
        print_success "安装目录已移除"
    fi
}

# 移除配置文件
remove_config() {
    if [[ $keep_config =~ ^[Yy]$ ]]; then
        print_info "保留配置文件: $CONFIG_DIR"
    else
        if [ -d "$CONFIG_DIR" ]; then
            rm -rf "$CONFIG_DIR"
            print_success "配置文件已移除"
        fi
    fi
}

# 移除数据和日志
remove_data() {
    if [[ $keep_data =~ ^[Yy]$ ]]; then
        print_info "保留数据和日志:"
        print_info "  - $DATA_DIR"
        print_info "  - $LOG_DIR"
    else
        if [ -d "$DATA_DIR" ]; then
            rm -rf "$DATA_DIR"
            print_success "数据目录已移除"
        fi
        
        if [ -d "$LOG_DIR" ]; then
            rm -rf "$LOG_DIR"
            print_success "日志目录已移除"
        fi
    fi
}

# 显示卸载后的信息
show_post_uninstall_info() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_success "PawStream Edge Client 已卸载"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    if [[ $keep_config =~ ^[Yy]$ ]]; then
        echo "📂 保留的配置文件: $CONFIG_DIR"
    fi
    
    if [[ $keep_data =~ ^[Yy]$ ]]; then
        echo "📂 保留的数据: $DATA_DIR"
        echo "📂 保留的日志: $LOG_DIR"
    fi
    
    echo ""
    print_info "感谢使用 PawStream Edge Client！"
    echo ""
}

# 主卸载流程
main() {
    check_root
    confirm_uninstall
    
    remove_systemd_service
    remove_binary
    remove_config
    remove_data
    
    show_post_uninstall_info
}

# 运行主程序
main
