# edge-client Specification

## Purpose
TBD - created by archiving change init-edge-client. Update Purpose after archive.
## Requirements
### Requirement: 配置管理
边缘客户端 SHALL 支持通过配置文件和环境变量进行配置管理。

#### Scenario: 从 YAML 配置文件加载配置
- **WHEN** 客户端启动时指定配置文件路径 `--config /etc/pawstream/config.yaml`
- **THEN** 客户端应成功加载配置（设备 ID、密钥、API 地址、输入源等）

#### Scenario: 环境变量覆盖配置文件
- **WHEN** 设置环境变量 `PAWSTREAM_DEVICE_ID` 和 `PAWSTREAM_SECRET`
- **THEN** 环境变量应覆盖配置文件中的对应值

#### Scenario: 配置验证失败
- **WHEN** 缺少必需配置项（如设备 ID 或密钥）
- **THEN** 客户端应输出明确的错误信息并退出，返回非零状态码

### Requirement: 设备认证
边缘客户端 SHALL 使用设备 ID 和密钥向 API 服务器进行认证，获取推流配置。

#### Scenario: 成功认证
- **WHEN** 使用有效的设备 ID 和密钥调用 API `/api/devices/{id}`
- **THEN** 应返回设备信息，包括 `publish_path` 和启用状态

#### Scenario: 认证失败 - 无效密钥
- **WHEN** 使用无效的设备密钥
- **THEN** 应返回 401 错误，客户端应记录错误并退出

#### Scenario: 认证失败 - 设备被禁用
- **WHEN** 设备在 API 服务器中被标记为禁用
- **THEN** 应返回 403 错误，客户端应记录错误并退出

### Requirement: 视频采集
边缘客户端 SHALL 支持多种视频输入源的采集。

#### Scenario: 从本地摄像头采集（V4L2）
- **WHEN** 配置输入源为 `v4l2:///dev/video0`
- **THEN** 客户端应打开该摄像头设备并开始采集视频帧

#### Scenario: 从 RTSP 源采集
- **WHEN** 配置输入源为 `rtsp://192.168.1.100:554/stream`
- **THEN** 客户端应连接到该 RTSP 源并接收视频流

#### Scenario: 从文件读取（测试模式）
- **WHEN** 配置输入源为 `file:///path/to/video.mp4`
- **THEN** 客户端应读取该文件并循环播放作为输入

#### Scenario: 测试模式 - 生成测试画面
- **WHEN** 配置输入源为 `test://`
- **THEN** 客户端应生成测试视频画面（如彩条或计数器）

#### Scenario: 输入源不可用
- **WHEN** 指定的输入源无法打开（如设备不存在）
- **THEN** 客户端应记录错误，等待一段时间后重试或退出

### Requirement: RTSP 推流
边缘客户端 SHALL 将采集的视频流推送到 MediaMTX 服务器。

#### Scenario: 成功推流
- **WHEN** 客户端采集到视频并使用设备密钥推流到 `rtsp://mediamtx:8554/{publish_path}`
- **THEN** MediaMTX 应接受连接并开始接收视频流

#### Scenario: 推流认证
- **WHEN** 推流到 MediaMTX 时，使用设备 ID 作为用户名，密钥作为密码
- **THEN** MediaMTX 应通过回调 API 服务器验证身份，允许推流

#### Scenario: 推流中断 - 自动重连
- **WHEN** 推流连接因网络问题中断
- **THEN** 客户端应记录错误，等待 5 秒后自动重新连接

#### Scenario: 推流失败 - 认证错误
- **WHEN** 推流认证失败（如密钥错误）
- **THEN** 客户端应记录错误并退出，避免无限重试

### Requirement: 日志记录
边缘客户端 SHALL 提供结构化的日志输出，便于问题诊断。

#### Scenario: 控制台日志输出
- **WHEN** 客户端启动且未配置日志文件
- **THEN** 日志应输出到 stdout，格式为 JSON 结构化日志

#### Scenario: 日志文件输出
- **WHEN** 配置日志文件路径 `log_file: /var/log/pawstream/edge-client.log`
- **THEN** 日志应同时输出到控制台和文件

#### Scenario: 日志级别控制
- **WHEN** 配置日志级别为 `info`
- **THEN** 只有 `info`、`warn`、`error` 级别的日志应被输出，`debug` 日志应被过滤

#### Scenario: 关键事件记录
- **WHEN** 发生关键事件（启动、认证成功、推流开始、推流中断、退出）
- **THEN** 应记录包含时间戳、事件类型、详细信息的日志

### Requirement: 进程管理
边缘客户端 SHALL 支持优雅启动和关闭。

#### Scenario: 正常启动
- **WHEN** 执行 `edge-client start --config config.yaml`
- **THEN** 客户端应加载配置、认证、开始采集和推流，并输出启动成功日志

#### Scenario: 优雅关闭 - SIGTERM
- **WHEN** 客户端收到 SIGTERM 信号
- **THEN** 应停止推流、释放资源、输出关闭日志，然后退出

#### Scenario: 优雅关闭 - SIGINT（Ctrl+C）
- **WHEN** 用户按下 Ctrl+C
- **THEN** 应停止推流、释放资源、输出关闭日志，然后退出

#### Scenario: Systemd 服务集成
- **WHEN** 通过 systemd 启动客户端 `systemctl start pawstream-edge`
- **THEN** 客户端应作为后台服务运行，并支持 systemd 的生命周期管理

### Requirement: 跨平台支持
边缘客户端 SHALL 支持在多种操作系统和架构上运行。

#### Scenario: Linux ARM64（树莓派）
- **WHEN** 在 ARM64 Linux 系统上编译和运行客户端
- **THEN** 客户端应正常工作，支持 V4L2 摄像头采集

#### Scenario: Linux AMD64（服务器/PC）
- **WHEN** 在 x86_64 Linux 系统上编译和运行客户端
- **THEN** 客户端应正常工作

#### Scenario: Windows（测试和开发）
- **WHEN** 在 Windows 系统上编译和运行客户端
- **THEN** 客户端应正常工作（V4L2 功能不可用，但支持 RTSP 和文件输入）

#### Scenario: macOS（测试和开发）
- **WHEN** 在 macOS 系统上编译和运行客户端
- **THEN** 客户端应正常工作（V4L2 功能不可用，但支持 RTSP 和文件输入）

### Requirement: 健康检查
边缘客户端 SHALL 支持可选的 HTTP 健康检查端点，可通过配置启用或禁用。

#### Scenario: 健康检查端点 - 正常状态
- **WHEN** 客户端正在推流，访问 `http://localhost:9090/health`
- **THEN** 应返回 200 状态码和 JSON 响应，包含状态信息（如 `{"status": "streaming", "uptime": 3600}`）

#### Scenario: 健康检查端点 - 异常状态
- **WHEN** 客户端推流中断，访问 `http://localhost:9090/health`
- **THEN** 应返回 503 状态码和错误信息

#### Scenario: 禁用健康检查
- **WHEN** 配置中未启用健康检查端点
- **THEN** 客户端应不监听 HTTP 端口

