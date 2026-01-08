#!/bin/bash
#
# PawStream Edge Client - 一键安装脚本
# 支持 Linux 系统（Ubuntu, Debian, CentOS, Fedora, Arch）
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

# 检测操作系统
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        OS_VERSION=$VERSION_ID
    else
        print_error "无法检测操作系统"
        exit 1
    fi
    
    print_info "检测到操作系统: $OS $OS_VERSION"
}

# 检测系统架构
detect_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="armv7"
            ;;
        *)
            print_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac
    
    print_info "系统架构: $ARCH"
}

# 检查依赖
check_dependencies() {
    print_info "检查系统依赖..."
    
    local deps=("ffmpeg")
    local missing_deps=()
    
    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            missing_deps+=($dep)
        fi
    done
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_warning "缺少依赖: ${missing_deps[*]}"
        print_info "正在安装依赖..."
        
        case $OS in
            ubuntu|debian)
                apt-get update
                apt-get install -y "${missing_deps[@]}"
                ;;
            centos|rhel|fedora)
                yum install -y "${missing_deps[@]}"
                ;;
            arch)
                pacman -S --noconfirm "${missing_deps[@]}"
                ;;
            *)
                print_warning "请手动安装依赖: ${missing_deps[*]}"
                ;;
        esac
    fi
    
    print_success "依赖检查完成"
}

# 创建必要的目录
create_directories() {
    print_info "创建目录结构..."
    
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$LOG_DIR"
    
    print_success "目录创建完成"
}

# 安装二进制文件
install_binary() {
    print_info "安装 Edge Client 二进制文件..."
    
    # 检查是否存在本地编译的二进制文件
    if [ -f "./build/edge-client" ]; then
        print_info "使用本地二进制文件"
        cp "./build/edge-client" "$INSTALL_DIR/$APP_NAME"
    elif [ -f "./edge-client" ]; then
        print_info "使用当前目录的二进制文件"
        cp "./edge-client" "$INSTALL_DIR/$APP_NAME"
    else
        print_error "未找到二进制文件"
        print_info "请先编译: make build"
        exit 1
    fi
    
    chmod +x "$INSTALL_DIR/$APP_NAME"
    
    # 创建符号链接
    ln -sf "$INSTALL_DIR/$APP_NAME" "$BIN_DIR/$APP_NAME"
    
    print_success "二进制文件安装完成"
}

# 复制 Web UI 文件
install_webui() {
    print_info "安装 Web UI 文件..."
    
    if [ -d "./web" ]; then
        cp -r "./web" "$INSTALL_DIR/"
        print_success "Web UI 文件安装完成"
    else
        print_warning "未找到 Web UI 文件"
    fi
}

# 创建示例配置文件
create_config() {
    print_info "创建示例配置文件..."
    
    if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
        cat > "$CONFIG_DIR/config.yaml.example" << 'EOF'
# PawStream Edge Client 配置文件示例
# 复制此文件为 config.yaml 并修改相应配置

device:
  id: "your-device-id"
  secret: "your-device-secret"

api:
  url: "http://your-api-server:3000"

input:
  type: "v4l2"  # v4l2, rtsp, file, test
  source: "/dev/video0"

video:
  width: 1280
  height: 720
  framerate: 30
  bitrate: 2000000

stream:
  url: "rtsp://localhost:8554"
  reconnect_interval: 5
  max_reconnect_attempts: 0

log:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, console

health:
  enabled: true
  address: ":8089"

webui:
  enabled: true
  host: "0.0.0.0"
  port: 8088
EOF
        print_success "示例配置文件创建完成: $CONFIG_DIR/config.yaml.example"
        print_info "首次运行时，客户端将自动启动设置向导"
    else
        print_warning "配置文件已存在，跳过创建"
    fi
}

# 创建 systemd 服务
install_systemd_service() {
    print_info "安装 systemd 服务..."
    
    cat > "$SYSTEMD_DIR/$SERVICE_NAME.service" << EOF
[Unit]
Description=PawStream Edge Client
Documentation=https://github.com/yourusername/pawstream
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/$APP_NAME start --config $CONFIG_DIR/config.yaml
Restart=on-failure
RestartSec=5s
StandardOutput=append:$LOG_DIR/edge-client.log
StandardError=append:$LOG_DIR/edge-client-error.log

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DATA_DIR $LOG_DIR $CONFIG_DIR

[Install]
WantedBy=multi-user.target
EOF
    
    # 重载 systemd
    systemctl daemon-reload
    
    print_success "Systemd 服务安装完成"
}

# 显示安装后的信息
show_post_install_info() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_success "PawStream Edge Client 安装完成！"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "📦 安装位置:"
    echo "   二进制文件: $INSTALL_DIR/$APP_NAME"
    echo "   配置目录:   $CONFIG_DIR"
    echo "   数据目录:   $DATA_DIR"
    echo "   日志目录:   $LOG_DIR"
    echo ""
    echo "🚀 快速开始:"
    echo ""
    echo "   方式1: 使用设置向导（推荐）"
    echo "   $ edge-client start"
    echo "   浏览器将自动打开设置向导"
    echo ""
    echo "   方式2: 手动配置"
    echo "   $ cp $CONFIG_DIR/config.yaml.example $CONFIG_DIR/config.yaml"
    echo "   $ nano $CONFIG_DIR/config.yaml"
    echo "   $ edge-client start"
    echo ""
    echo "🔧 服务管理:"
    echo "   启动服务:   sudo systemctl start $SERVICE_NAME"
    echo "   停止服务:   sudo systemctl stop $SERVICE_NAME"
    echo "   重启服务:   sudo systemctl restart $SERVICE_NAME"
    echo "   查看状态:   sudo systemctl status $SERVICE_NAME"
    echo "   开机自启:   sudo systemctl enable $SERVICE_NAME"
    echo "   禁用自启:   sudo systemctl disable $SERVICE_NAME"
    echo ""
    echo "📝 查看日志:"
    echo "   实时日志:   sudo journalctl -u $SERVICE_NAME -f"
    echo "   全部日志:   sudo journalctl -u $SERVICE_NAME"
    echo "   日志文件:   $LOG_DIR/edge-client.log"
    echo ""
    echo "🌐 Web UI:"
    echo "   访问地址:   http://localhost:8088"
    echo "   设置向导:   http://localhost:8088/setup"
    echo ""
    echo "💡 提示:"
    echo "   - 首次运行会自动打开设置向导"
    echo "   - 配置完成后可在 Web UI 中管理"
    echo "   - 使用 'edge-client --help' 查看所有命令"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# 主安装流程
main() {
    echo ""
    echo "╔════════════════════════════════════════════════════╗"
    echo "║                                                    ║"
    echo "║     PawStream Edge Client - 安装程序              ║"
    echo "║                                                    ║"
    echo "╚════════════════════════════════════════════════════╝"
    echo ""
    
    check_root
    detect_os
    detect_arch
    check_dependencies
    create_directories
    install_binary
    install_webui
    create_config
    install_systemd_service
    show_post_install_info
    
    print_success "安装完成！"
}

# 运行主程序
main
