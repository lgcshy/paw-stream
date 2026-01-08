# 实施任务清单

## 1. 项目初始化
- [x] 1.1 创建 `client/edge/` 目录结构
- [x] 1.2 初始化 Go module (`go.mod`)
- [x] 1.3 设计目录布局（cmd, internal, configs）
- [x] 1.4 创建 README.md 和配置示例文件

## 2. 配置管理
- [x] 2.1 定义配置结构（设备信息、API 地址、推流参数）
- [x] 2.2 实现配置文件加载（YAML）
- [x] 2.3 支持环境变量覆盖
- [x] 2.4 实现配置验证逻辑

## 3. API 认证集成
- [x] 3.1 实现 API 客户端（HTTP client）
- [x] 3.2 实现设备认证逻辑（使用设备 ID 和密钥）
- [x] 3.3 从 API 获取推流配置（publish_path, MediaMTX 地址）
- [x] 3.4 实现设备专用认证端点（POST /api/device/auth）

## 4. 视频采集
- [x] 4.1 调研和选择视频采集方案（FFmpeg）
- [x] 4.2 实现文件输入源（测试用）
- [x] 4.3 实现摄像头输入源（V4L2 on Linux）
- [x] 4.4 实现 RTSP 输入源（支持 IP 摄像头）
- [x] 4.5 实现测试模式（生成测试视频流）

## 5. 推流功能
- [x] 5.1 调研和选择推流方案（FFmpeg）
- [x] 5.2 实现 RTSP 推流到 MediaMTX
- [x] 5.3 实现推流认证（使用设备密钥）
- [x] 5.4 实现推流错误处理和自动重连
- [x] 5.5 实现推流状态监控

## 6. 日志和监控
- [x] 6.1 集成日志库（zerolog）
- [x] 6.2 实现结构化日志输出
- [x] 6.3 支持日志级别配置
- [x] 6.4 实现可选的日志文件输出
- [x] 6.5 实现健康检查 HTTP endpoint（可选）

## 7. 进程管理
- [x] 7.1 实现优雅启动和关闭
- [x] 7.2 实现信号处理（SIGTERM, SIGINT）
- [x] 7.3 创建 systemd service 配置示例
- [x] 7.4 编写部署文档
- [x] 7.5 实现 daemon 模式支持（go-daemon）
- [x] 7.6 实现子命令架构（start, stop, status, restart）

## 8. 构建和打包
- [x] 8.1 配置 Makefile 或 build 脚本
- [x] 8.2 支持跨平台编译（Linux arm64/amd64, Windows, macOS）
- [x] 8.3 编写安装脚本

## 9. 测试和文档
- [x] 9.1 编写用户文档（README.md）
- [x] 9.2 配置文件示例和说明
- [x] 9.3 端到端测试（推流到 MediaMTX）
- [x] 9.4 验证设备认证流程
