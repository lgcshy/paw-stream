# Change: 初始化边缘推流客户端

## Why

PawStream 需要一个轻量级的边缘推流客户端，部署在摄像头设备（如树莓派）上，负责采集视频并推流到 MediaMTX。当前缺少一个官方的、易于部署的推流客户端工具。

## What Changes

- 创建新的 Go CLI 应用 `edge-client`
- 实现视频采集功能（支持多种输入源）
- 实现 RTSP/RTMP 推流到 MediaMTX
- 支持通过 API 服务器进行设备认证（使用设备密钥）
- 提供配置文件和环境变量管理
- 跨平台支持（Linux 主要，Windows/macOS 次要）
- 提供 systemd 服务配置示例
- 实现日志记录和健康检查

## Impact

- **Affected specs**: 新增 `edge-client` capability
- **Affected code**: 
  - 新增 `client/edge/` 目录（Go 应用）
  - 可能需要在 `api-server` 中添加客户端心跳接口（可选，Phase 2）
- **Dependencies**: 
  - 视频采集库（如 gstreamer-go 或 ffmpeg bindings）
  - RTSP/RTMP 推流库
  - API 客户端（复用 api-server 的认证逻辑）
