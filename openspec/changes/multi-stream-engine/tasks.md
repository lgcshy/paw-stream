# 任务清单：多推流引擎支持

## Phase 1: 核心架构 (0/5)

- [ ] 1.1 创建 StreamEngine 接口定义
  - internal/stream/engine.go
  - 定义 StreamEngine 接口
  - 定义 EngineStats 结构

- [ ] 1.2 重构 FFmpegManager 为 FFmpegEngine
  - 重命名 internal/stream/ffmpeg.go 中的类型
  - 实现 StreamEngine 接口
  - 添加 Stats() 方法

- [ ] 1.3 创建引擎工厂
  - internal/stream/factory.go
  - NewStreamEngine() 工厂函数
  - 引擎检测和验证

- [ ] 1.4 更新 Manager 支持引擎接口
  - 修改 internal/stream/manager.go
  - 使用 StreamEngine 接口替代直接依赖
  - 支持动态引擎选择

- [ ] 1.5 更新配置结构
  - internal/config/config.go
  - 添加 StreamConfig.Engine 字段
  - 添加 FFmpegConfig 和 GStreamerConfig
  - 添加 Preset 字段

## Phase 2: GStreamer 引擎 (0/10)

- [ ] 2.1 创建 GStreamerEngine 基础框架
  - internal/stream/gstreamer.go
  - 实现 StreamEngine 接口
  - 基础启动/停止逻辑

- [ ] 2.2 实现 pipeline 构建器
  - internal/stream/gstreamer_pipeline.go
  - 根据 InputSource 构建 pipeline
  - 参数化配置

- [ ] 2.3 实现 V4L2 输入支持
  - V4L2Source 的 GStreamer pipeline
  - 测试摄像头捕获

- [ ] 2.4 实现 RTSP 输入支持
  - RTSPSource 的 GStreamer pipeline
  - 测试 RTSP 转发

- [ ] 2.5 实现 Test 输入支持
  - TestSource 的 GStreamer pipeline
  - videotestsrc 配置

- [ ] 2.6 实现 File 输入支持
  - FileSource 的 GStreamer pipeline
  - 文件循环播放

- [ ] 2.7 添加硬件加速检测
  - internal/stream/hwaccel.go
  - 检测 VAAPI 支持
  - 检测 NVENC 支持
  - 自动选择最优编码器

- [ ] 2.8 实现错误处理
  - 解析 GStreamer 错误消息
  - 错误重试逻辑
  - 降级处理（硬件→软件）

- [ ] 2.9 添加性能统计
  - FPS 监控
  - 码率监控
  - 丢帧统计

- [ ] 2.10 集成测试
  - 各种输入源测试
  - 长时间运行测试
  - 错误恢复测试

## Phase 3: 预设配置 (0/6)

- [ ] 3.1 实现预设系统
  - internal/config/preset.go
  - ApplyPreset() 函数
  - 预设验证

- [ ] 3.2 实现 low-latency 预设
  - GStreamer 引擎
  - latency=100ms
  - 硬件编码

- [ ] 3.3 实现 high-quality 预设
  - FFmpeg 引擎
  - slow preset
  - 高码率

- [ ] 3.4 实现 balanced 预设
  - FFmpeg 引擎
  - medium preset
  - 中等码率

- [ ] 3.5 实现 power-save 预设
  - 优先硬件编码
  - 低分辨率/帧率
  - 低码率

- [ ] 3.6 预设文档
  - docs/presets.md
  - 各预设说明
  - 使用场景

## Phase 4: Web UI 支持 (0/7)

- [ ] 4.1 设置向导添加引擎选择
  - web/setup.html 和 setup.js
  - 步骤 4 添加引擎选项
  - FFmpeg/GStreamer 单选

- [ ] 4.2 添加预设配置选择
  - 预设下拉菜单
  - 预设说明和推荐

- [ ] 4.3 动态配置选项
  - 根据引擎显示不同配置
  - FFmpeg 显示 preset/tune
  - GStreamer 显示 latency

- [ ] 4.4 引擎依赖检查 API
  - GET /api/engine/available
  - 返回可用引擎列表
  - 显示缺失依赖提示

- [ ] 4.5 状态页面显示引擎信息
  - web/index.html
  - 显示当前引擎
  - 显示引擎统计

- [ ] 4.6 添加性能统计图表
  - FPS 图表
  - 码率图表
  - 实时更新

- [ ] 4.7 配置导入/导出支持引擎
  - 包含引擎配置
  - 预设信息

## Phase 5: 测试和文档 (0/9)

- [ ] 5.1 单元测试 - StreamEngine 接口
  - stream/engine_test.go
  - 接口行为测试

- [ ] 5.2 单元测试 - FFmpegEngine
  - stream/ffmpeg_test.go
  - 各种输入源测试

- [ ] 5.3 单元测试 - GStreamerEngine
  - stream/gstreamer_test.go
  - Pipeline 构建测试

- [ ] 5.4 集成测试 - 引擎切换
  - 配置切换测试
  - 预设应用测试

- [ ] 5.5 性能对比测试
  - 延迟测试
  - CPU/内存占用对比
  - 测试报告

- [ ] 5.6 更新 README
  - 引擎选择说明
  - 依赖安装指南
  - 配置示例

- [ ] 5.7 添加引擎选择指南
  - docs/engine-selection.md
  - 场景推荐
  - 性能对比

- [ ] 5.8 添加 GStreamer 故障排除
  - docs/troubleshooting-gstreamer.md
  - 常见问题
  - 调试方法

- [ ] 5.9 更新配置文件示例
  - configs/config.yaml
  - 添加引擎配置示例
  - 各预设示例

## 进度统计

- **总任务数**: 37
- **已完成**: 0
- **进行中**: 0
- **待开始**: 37
- **完成率**: 0%

## 里程碑

- [ ] **M1**: Phase 1 完成 - 架构重构（预计 2026-01-11）
- [ ] **M2**: Phase 2 完成 - GStreamer 引擎（预计 2026-01-15）
- [ ] **M3**: Phase 3 完成 - 预设配置（预计 2026-01-17）
- [ ] **M4**: Phase 4 完成 - Web UI 支持（预计 2026-01-20）
- [ ] **M5**: Phase 5 完成 - 测试和文档（预计 2026-01-23）
- [ ] **发布**: v1.1.0 - 多引擎支持（预计 2026-01-24）
