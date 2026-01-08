# PawStream Edge Client

PawStream 边缘推流客户端 - 部署在摄像头设备上的轻量级 Go 应用，负责视频采集和推流。

## 功能特性

- ✅ 多种视频输入源支持（V4L2、RTSP、文件、测试画面）
- ✅ RTSP 推流到 MediaMTX，带设备密钥认证
- ✅ 与 PawStream API 服务器集成，自动获取推流配置
- ✅ **Web UI 管理界面**（配置、状态监控、日志查看）
- ✅ **配置文件热重载**（无需重启即可更新配置）
- ✅ **实时 SSE 推送**（状态和日志实时更新）
- ✅ 配置文件和环境变量管理
- ✅ 结构化日志记录（zerolog）
- ✅ 优雅启动和关闭
- ✅ **Daemon 模式支持**（前台/后台运行）
- ✅ Systemd 服务支持（生产环境推荐）
- ✅ 跨平台支持（Linux、Windows、macOS）
- ✅ 可选的健康检查 HTTP 端点

## 快速开始

### 安装

```bash
# 从源码编译
git clone https://github.com/lgc/pawstream.git
cd pawstream/client/edge
make build

# 或下载预编译二进制
# wget https://github.com/lgc/pawstream/releases/latest/download/edge-client-linux-arm64
```

### 首次使用（零配置）

**小白用户 - 只需一条命令：**

```bash
# 启动客户端（自动进入设置向导）
./edge-client start

# 或后台运行
./edge-client start --daemon
```

**设置向导会自动：**
1. 🔍 检测配置文件
2. 🌐 启动 Web UI (http://localhost:8088/setup)
3. 🔗 打开浏览器
4. 📝 引导您完成配置
5. 💾 自动生成配置文件
6. 🚀 启动推流

### 配置

创建配置文件 `config.yaml`：

```yaml
# 设备信息
device:
  id: "your-device-id"          # 从 PawStream Web UI 获取
  secret: "your-device-secret"  # 从 PawStream Web UI 获取

# API 服务器
api:
  url: "http://api.pawstream.example.com"
  timeout: 10s

# 视频输入源
input:
  type: "v4l2"          # v4l2 | rtsp | file | test
  source: "/dev/video0" # 根据 type 不同，source 格式不同
  # type: rtsp -> source: "rtsp://192.168.1.100:554/stream"
  # type: file -> source: "/path/to/video.mp4"
  # type: test -> source: "" (生成测试画面)

# 视频编码参数（可选）
video:
  codec: "h264"
  width: 1280
  height: 720
  framerate: 30
  bitrate: 2000000  # 2 Mbps

# MediaMTX 推流
stream:
  url: ""  # 留空，将从 API 获取
  reconnect_interval: 5s
  max_reconnect_attempts: 0  # 0 = 无限重试

# 日志
log:
  level: "info"  # debug | info | warn | error
  file: ""       # 留空输出到 stdout，否则输出到文件

# 健康检查（可选）
health:
  enabled: false
  address: ":9090"

# Web UI（可选）
webui:
  enabled: true
  host: "0.0.0.0"          # 监听地址：0.0.0.0 = 所有接口
  port: 8088
  auth:                    # HTTP Basic Auth（可选）
    enabled: false
    username: "admin"
    password: "secret"
```

或使用环境变量：

```bash
export PAWSTREAM_DEVICE_ID="your-device-id"
export PAWSTREAM_DEVICE_SECRET="your-device-secret"
export PAWSTREAM_API_URL="http://api.pawstream.example.com"
export PAWSTREAM_INPUT_TYPE="v4l2"
export PAWSTREAM_INPUT_SOURCE="/dev/video0"
```

### 运行

#### 小白用户（推荐）

```bash
# 首次运行 - 自动进入设置向导
./edge-client start

# 后续运行 - 使用已保存的配置
./edge-client start

# 后台运行
./edge-client start --daemon

# 查看状态
./edge-client status

# 停止
./edge-client stop
```

#### 开发人员（高级配置）

```bash
# 查看帮助
./edge-client

# 查看版本
./edge-client version

# 使用自定义配置文件
./edge-client start --config /path/to/config.yaml

# 使用命令行参数（覆盖配置文件）
./edge-client start \
  --device-id d3cf4e12-a4f0-4066-a7b6-74b22ff8cffd \
  --device-secret "your-secret" \
  --api-url http://api.example.com

# 组合使用配置文件和参数
./edge-client start --config config.yaml --input-type test

# 完全命令行配置（无需配置文件）
./edge-client start \
  --device-id xxx \
  --device-secret yyy \
  --api-url http://api.example.com \
  --input-type v4l2 \
  --input-source /dev/video0

# 调试模式
./edge-client start --log-level debug

# 重启
./edge-client restart
```

#### 配置文件自动查找

如果不指定 `--config`，将按以下顺序自动查找：

1. `./config.yaml`
2. `./configs/config.yaml`
3. `~/.pawstream/config.yaml`

找不到配置文件时，自动进入设置向导模式。
```

### Web UI 管理界面

启动客户端后，访问 Web UI 管理界面：

```
http://localhost:8088
```

Web UI 功能：

- **配置标签页**：编辑设备配置、输入源、视频参数等，保存后自动热重载
- **状态标签页**：实时查看客户端状态、系统资源使用情况
- **日志标签页**：实时查看应用日志，支持自动滚动

特性：

- ✅ 配置在线编辑，无需重启
- ✅ 实时状态监控（SSE 推送）
- ✅ 系统资源监控（CPU、内存、磁盘）
- ✅ 实时日志查看
- ✅ 响应式设计，支持移动端访问
- ✅ 可选的 HTTP Basic Auth 保护

适用场景：

- 边缘设备调试和配置
- 分布式部署的设备监控
- 无技术背景用户的简化配置
- 树莓派等嵌入式设备的本地管理
```

### Systemd 服务

```bash
# 复制服务文件
sudo cp configs/systemd/pawstream-edge.service /etc/systemd/system/

# 编辑配置文件路径
sudo nano /etc/systemd/system/pawstream-edge.service

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable pawstream-edge
sudo systemctl start pawstream-edge

# 查看状态
sudo systemctl status pawstream-edge

# 查看日志
sudo journalctl -u pawstream-edge -f
```

## 开发

### 构建

```bash
# 本地平台
make build

# 跨平台编译
make build-linux-amd64
make build-linux-arm64
make build-windows
make build-macos

# 构建所有平台
make build-all
```

### 测试

```bash
# 运行测试
make test

# 带覆盖率
make test-coverage
```

### 开发调试

```bash
# 启用 debug 日志
./edge-client start --config config.yaml --log-level debug

# 使用测试输入源
./edge-client start --config config.yaml --input-type test
```

## 目录结构

```
client/edge/
├── cmd/
│   └── edge-client/     # 主程序入口
│       └── main.go
├── internal/            # 内部包（不对外暴露）
│   ├── config/          # 配置管理
│   ├── auth/            # API 认证
│   ├── capture/         # 视频采集
│   ├── stream/          # 推流
│   └── health/          # 健康检查
├── configs/             # 配置文件示例和 systemd 服务
├── docs/                # 文档
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 依赖项

- Go 1.21+
- FFmpeg（用于视频采集和推流）
- V4L2（Linux 摄像头支持，仅 Linux）

## 许可证

MIT License

## 支持

- 文档: https://pawstream.example.com/docs
- Issues: https://github.com/lgc/pawstream/issues
