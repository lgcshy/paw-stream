# edge-client Specification Delta

## ADDED Requirements

### Requirement: 多推流引擎支持
边缘客户端 SHALL 支持 FFmpeg 和 GStreamer 两种推流引擎，允许用户根据场景选择最优方案。

#### Scenario: 使用 FFmpeg 引擎（默认）
- **WHEN** 配置文件未指定引擎或指定 `stream.engine: ffmpeg`
- **THEN** 客户端应使用 FFmpeg 进行推流，兼容性最佳

#### Scenario: 使用 GStreamer 引擎（低延迟）
- **WHEN** 配置 `stream.engine: gstreamer`
- **THEN** 客户端应使用 GStreamer 进行推流，延迟低于 200ms

#### Scenario: GStreamer 不可用时回退
- **WHEN** 配置使用 GStreamer 但系统未安装
- **THEN** 客户端应记录警告并自动回退到 FFmpeg 引擎

#### Scenario: 引擎状态查询
- **WHEN** 查询客户端状态 `/api/status`
- **THEN** 响应应包含当前使用的引擎名称和版本信息

### Requirement: 预设配置
边缘客户端 SHALL 提供常见场景的优化预设配置，简化用户配置。

#### Scenario: 应用低延迟预设
- **WHEN** 配置 `stream.preset: low-latency`
- **THEN** 客户端应自动选择 GStreamer 引擎，latency 100ms，启用硬件编码

#### Scenario: 应用高质量预设
- **WHEN** 配置 `stream.preset: high-quality`
- **THEN** 客户端应选择 FFmpeg，preset=slow，bitrate=5Mbps，framerate=60fps

#### Scenario: 应用平衡预设
- **WHEN** 配置 `stream.preset: balanced`
- **THEN** 客户端应选择 FFmpeg，preset=medium，bitrate=2Mbps，framerate=30fps

#### Scenario: 应用省电预设
- **WHEN** 配置 `stream.preset: power-save`
- **THEN** 客户端应优先使用硬件编码，降低分辨率和帧率

#### Scenario: 预设与自定义配置冲突
- **WHEN** 同时指定预设和自定义配置参数
- **THEN** 自定义配置应覆盖预设中的对应项

### Requirement: GStreamer Pipeline 支持
边缘客户端 SHALL 支持通过 GStreamer 构建视频处理 pipeline。

#### Scenario: V4L2 输入 + GStreamer
- **WHEN** 配置 `input.type: v4l2` 和 `stream.engine: gstreamer`
- **THEN** 客户端应构建 `v4l2src ! video/x-raw ! vaapih264enc ! rtspclientsink` pipeline

#### Scenario: RTSP 输入 + GStreamer
- **WHEN** 配置 `input.type: rtsp` 和 `stream.engine: gstreamer`
- **THEN** 客户端应构建 `rtspsrc ! rtph264depay ! h264parse ! vaapih264enc ! rtspclientsink` pipeline

#### Scenario: Test 输入 + GStreamer
- **WHEN** 配置 `input.type: test` 和 `stream.engine: gstreamer`
- **THEN** 客户端应构建 `videotestsrc ! video/x-raw ! vaapih264enc ! rtspclientsink` pipeline

### Requirement: 硬件加速检测
边缘客户端 SHALL 自动检测可用的硬件编码器并智能选择。

#### Scenario: 检测到 VAAPI 支持（Intel）
- **WHEN** 系统支持 VAAPI 且配置 `stream.gstreamer.use_hardware: true`
- **THEN** GStreamer pipeline 应使用 `vaapih264enc` 编码器

#### Scenario: 检测到 NVENC 支持（NVIDIA）
- **WHEN** 系统支持 NVENC 且配置 `stream.ffmpeg.hwaccel: auto`
- **THEN** FFmpeg 应使用 `-c:v h264_nvenc` 硬件编码

#### Scenario: 硬件编码失败时降级
- **WHEN** 硬件编码器初始化失败
- **THEN** 客户端应自动回退到软件编码（x264）并记录警告

#### Scenario: 查询可用编码器
- **WHEN** 调用 `/api/engine/encoders`
- **THEN** 应返回当前可用的硬件和软件编码器列表

### Requirement: 引擎性能统计
边缘客户端 SHALL 提供实时的推流性能统计信息。

#### Scenario: 监控实时 FPS
- **WHEN** 推流进行中
- **THEN** 状态 API 应返回当前实际输出帧率（FPS）

#### Scenario: 监控码率
- **WHEN** 推流进行中
- **THEN** 状态 API 应返回当前实际输出码率（kbps）

#### Scenario: 监控丢帧数
- **WHEN** 推流进行中
- **THEN** 状态 API 应返回累计丢帧数

#### Scenario: 性能统计通过 SSE 推送
- **WHEN** Web UI 连接到 `/api/events`
- **THEN** 应每秒推送一次引擎性能统计（FPS、码率、丢帧）

## MODIFIED Requirements

### Requirement: 配置管理
边缘客户端 SHALL 支持通过配置文件、环境变量和预设进行配置管理。

#### Scenario: 从 YAML 配置文件加载配置
- **WHEN** 客户端启动时指定配置文件路径 `--config /etc/pawstream/config.yaml`
- **THEN** 客户端应成功加载配置（设备 ID、密钥、API 地址、输入源、引擎选择等）

#### Scenario: 环境变量覆盖配置文件
- **WHEN** 设置环境变量 `PAWSTREAM_DEVICE_ID`、`PAWSTREAM_SECRET` 和 `PAWSTREAM_STREAM_ENGINE`
- **THEN** 环境变量应覆盖配置文件中的对应值

#### Scenario: 配置验证失败
- **WHEN** 缺少必需配置项（如设备 ID 或密钥）或引擎配置无效
- **THEN** 客户端应输出明确的错误信息并退出，返回非零状态码

#### Scenario: 预设配置加载
- **WHEN** 配置文件指定 `stream.preset: low-latency`
- **THEN** 客户端应先应用预设，再用自定义配置覆盖

### Requirement: RTSP 推流
边缘客户端 SHALL 使用选定的引擎将采集的视频流推送到 MediaMTX 服务器。

#### Scenario: 成功推流
- **WHEN** 客户端采集到视频并使用设备密钥通过 FFmpeg 或 GStreamer 推流到 `rtsp://mediamtx:8554/{publish_path}`
- **THEN** MediaMTX 应接受连接并开始接收视频流

#### Scenario: 推流认证
- **WHEN** 推流到 MediaMTX 时，使用设备 ID 作为用户名，密钥作为密码
- **THEN** MediaMTX 应通过回调 API 服务器验证身份，允许推流

#### Scenario: 推流中断 - 自动重连
- **WHEN** 推流连接因网络问题中断
- **THEN** 客户端应记录错误，等待 5 秒后使用相同引擎自动重新连接

#### Scenario: 推流失败 - 认证错误
- **WHEN** 推流认证失败（如密钥错误）
- **THEN** 客户端应记录错误并退出，避免无限重试

#### Scenario: 引擎切换需要重启
- **WHEN** 修改配置文件更改引擎类型（FFmpeg ↔ GStreamer）
- **THEN** 客户端应检测到配置变化，停止当前推流，使用新引擎重新启动

## ADDED Configuration Fields

### StreamConfig
```yaml
stream:
  # 引擎选择（新增）
  engine: ffmpeg  # ffmpeg | gstreamer
  
  # 预设配置（新增）
  preset: ""  # low-latency | high-quality | balanced | power-save
  
  # 原有字段
  url: rtsp://localhost:8554/stream
  reconnect_interval: 5s
  max_reconnect_attempts: -1
  
  # FFmpeg 特有配置（新增）
  ffmpeg:
    preset: ultrafast      # ultrafast | superfast | veryfast | faster | fast | medium | slow | slower | veryslow
    tune: zerolatency      # film | animation | grain | stillimage | fastdecode | zerolatency
    hwaccel: auto          # none | auto | vaapi | nvenc | qsv | videotoolbox
    extra_args: []         # 自定义 FFmpeg 参数
    
  # GStreamer 特有配置（新增）
  gstreamer:
    latency_ms: 100        # Pipeline 延迟（毫秒）
    use_hardware: true     # 使用硬件编码
    buffer_size: 1000      # 缓冲区大小（微秒）
```

### Web UI Setup
在设置向导的步骤 4（输入源配置）之后添加新的步骤 5：

```
步骤 5: 推流引擎和预设

- 推流引擎选择（单选）
  ○ FFmpeg（推荐）- 兼容性好，通用场景
  ○ GStreamer - 低延迟，专业场景
  
- 预设配置（下拉菜单，可选）
  - 低延迟（监控、实时互动）
  - 高质量（录制、存档）
  - 平衡（通用场景）
  - 省电（边缘设备）
  - 自定义（不使用预设）
```
