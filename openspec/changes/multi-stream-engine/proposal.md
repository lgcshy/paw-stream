# Proposal: 多推流引擎支持（FFmpeg + GStreamer）

**状态**: 提议中  
**创建日期**: 2026-01-09  
**负责人**: Edge Client 团队  

## 📋 概述

为 Edge Client 添加多推流引擎支持，允许用户在 FFmpeg 和 GStreamer 之间选择，满足不同场景的需求。

## 🎯 目标

### 核心目标
1. **架构重构**：将推流引擎抽象为可插拔接口
2. **FFmpeg 引擎**：保持现有功能，作为默认引擎
3. **GStreamer 引擎**：提供低延迟推流选项
4. **预设配置**：提供常见场景的优化预设
5. **向后兼容**：现有配置无需修改

### 非目标
- 不支持引擎热切换（需要重启）
- 不支持同时使用多个引擎
- 第一版不支持自定义 GStreamer pipeline

## 💡 动机

### 业务需求
1. **延迟敏感场景**：监控、直播需要更低延迟
2. **硬件加速**：充分利用硬件编码能力
3. **灵活性**：不同场景选择最优方案

### 技术优势
- **FFmpeg**：兼容性好，使用广泛，适合通用场景
- **GStreamer**：延迟更低，插件丰富，专业场景更优

## 🎨 设计方案

### 1. 架构设计

```
┌─────────────────────────────────────────┐
│          Stream Manager                 │
│  (重连逻辑、状态管理、错误处理)         │
└──────────────┬──────────────────────────┘
               │
               ├─── StreamEngine 接口
               │     │
               │     ├─── FFmpegEngine    (默认)
               │     └─── GStreamerEngine (可选)
               │
               └─── InputSource 接口
                     ├─── V4L2Source
                     ├─── RTSPSource
                     ├─── FileSource
                     └─── TestSource
```

### 2. StreamEngine 接口

```go
type StreamEngine interface {
    // Start 启动推流
    Start(ctx context.Context) error
    
    // Stop 停止推流
    Stop() error
    
    // IsRunning 是否正在运行
    IsRunning() bool
    
    // ErrorCh 返回错误通道
    ErrorCh() <-chan error
    
    // Name 返回引擎名称
    Name() string
    
    // Stats 返回统计信息
    Stats() EngineStats
}

type EngineStats struct {
    FPS       float64
    Bitrate   int
    DroppedFrames int
}
```

### 3. 配置结构

```yaml
stream:
  # 引擎选择
  engine: ffmpeg  # ffmpeg | gstreamer
  
  # 预设配置（可选）
  preset: low-latency  # low-latency | high-quality | balanced | power-save
  
  # 输出配置
  url: rtsp://localhost:8554/stream
  reconnect_interval: 5s
  max_reconnect_attempts: -1
  
  # FFmpeg 特有配置
  ffmpeg:
    preset: ultrafast      # FFmpeg preset
    tune: zerolatency      # FFmpeg tune
    hwaccel: auto          # none | auto | vaapi | nvenc | qsv
    extra_args: []         # 自定义参数
    
  # GStreamer 特有配置
  gstreamer:
    latency_ms: 100        # 管道延迟
    use_hardware: true     # 使用硬件编码
    buffer_size: 1000      # 缓冲区大小（微秒）
```

### 4. 预设配置

#### low-latency（低延迟）
- 引擎：GStreamer
- 延迟：100-200ms
- 场景：监控、实时互动

#### high-quality（高质量）
- 引擎：FFmpeg
- 编码：slower preset
- 场景：录制、存档

#### balanced（平衡）
- 引擎：FFmpeg
- 编码：medium preset
- 场景：通用场景

#### power-save（省电）
- 引擎：优先硬件加速
- 编码：hardware encoder
- 场景：边缘设备、移动设备

### 5. GStreamer Pipeline 设计

#### V4L2 输入
```gstreamer
v4l2src device=/dev/video0 ! 
  video/x-raw,width=1280,height=720,framerate=30/1 ! 
  vaapih264enc bitrate=2000 ! 
  rtspclientsink location=rtsp://...
```

#### RTSP 输入
```gstreamer
rtspsrc location=rtsp://... latency=100 ! 
  rtph264depay ! 
  h264parse ! 
  vaapih264enc bitrate=2000 ! 
  rtspclientsink location=rtsp://...
```

#### Test 输入
```gstreamer
videotestsrc pattern=smpte ! 
  video/x-raw,width=1280,height=720,framerate=30/1 ! 
  vaapih264enc bitrate=2000 ! 
  rtspclientsink location=rtsp://...
```

## 📋 任务分解

### Phase 1: 核心架构（3-4 天）
- [ ] 创建 `StreamEngine` 接口
- [ ] 重构 `FFmpegManager` 为 `FFmpegEngine`
- [ ] 修改 `Manager` 支持引擎工厂
- [ ] 更新配置结构
- [ ] 添加引擎统计信息接口

### Phase 2: GStreamer 引擎（4-5 天）
- [ ] 实现 `GStreamerEngine` 基础框架
- [ ] 实现 V4L2 输入 pipeline
- [ ] 实现 RTSP 输入 pipeline
- [ ] 实现 Test 输入 pipeline
- [ ] 实现 File 输入 pipeline
- [ ] 添加硬件加速检测（VAAPI/NVENC）
- [ ] 实现错误处理和重连逻辑
- [ ] 添加性能统计（FPS、码率）

### Phase 3: 预设配置（1-2 天）
- [ ] 实现预设配置系统
- [ ] 添加 low-latency 预设
- [ ] 添加 high-quality 预设
- [ ] 添加 balanced 预设
- [ ] 添加 power-save 预设
- [ ] 预设配置文档

### Phase 4: Web UI 支持（2-3 天）
- [ ] 设置向导添加引擎选择
- [ ] 添加预设配置选择
- [ ] 根据引擎显示不同配置选项
- [ ] 状态页面显示引擎信息
- [ ] 添加引擎性能统计图表

### Phase 5: 测试和文档（2-3 天）
- [ ] 单元测试（引擎接口）
- [ ] 集成测试（各种输入源）
- [ ] 性能对比测试
- [ ] 延迟测试
- [ ] 更新 README
- [ ] 添加引擎选择指南
- [ ] 添加故障排除文档

## 🔧 技术细节

### 依赖安装

#### GStreamer（Debian/Ubuntu）
```bash
sudo apt-get install -y \
  gstreamer1.0-tools \
  gstreamer1.0-plugins-base \
  gstreamer1.0-plugins-good \
  gstreamer1.0-plugins-bad \
  gstreamer1.0-plugins-ugly \
  gstreamer1.0-rtsp
```

#### 硬件加速（可选）
```bash
# VAAPI (Intel)
sudo apt-get install -y gstreamer1.0-vaapi

# NVENC (NVIDIA)
sudo apt-get install -y gstreamer1.0-plugins-bad
```

### 引擎选择逻辑

```go
func NewStreamEngine(engineType string, input capture.InputSource, config Config) (StreamEngine, error) {
    switch engineType {
    case "ffmpeg", "":
        return NewFFmpegEngine(input, config), nil
    case "gstreamer":
        // 检查 gstreamer 是否安装
        if !isGStreamerInstalled() {
            return nil, fmt.Errorf("gstreamer not installed")
        }
        return NewGStreamerEngine(input, config), nil
    default:
        return nil, fmt.Errorf("unsupported engine: %s", engineType)
    }
}
```

### 预设应用逻辑

```go
func ApplyPreset(cfg *Config, preset string) error {
    switch preset {
    case "low-latency":
        cfg.Stream.Engine = "gstreamer"
        cfg.Stream.GStreamer.LatencyMs = 100
        cfg.Stream.GStreamer.UseHardware = true
        cfg.Video.Bitrate = 2000000
        cfg.Video.Framerate = 30
        
    case "high-quality":
        cfg.Stream.Engine = "ffmpeg"
        cfg.Stream.FFmpeg.Preset = "slow"
        cfg.Video.Bitrate = 5000000
        cfg.Video.Framerate = 60
        
    // ... 其他预设
    }
    return nil
}
```

## 📊 性能预期

### 延迟对比
| 引擎 | 典型延迟 | 最低延迟 |
|------|---------|---------|
| FFmpeg | 1-2s | 500ms |
| GStreamer | 200-500ms | 100ms |

### 资源占用
| 引擎 | CPU | 内存 |
|------|-----|------|
| FFmpeg | 30-50% | 50MB |
| GStreamer | 25-40% | 40MB |

## ⚠️ 风险和限制

### 风险
1. **GStreamer 依赖**：需要系统安装 GStreamer
2. **兼容性**：不同硬件编码器可用性不同
3. **调试复杂度**：GStreamer pipeline 错误较难定位

### 缓解措施
1. **依赖检查**：启动时检测 GStreamer 是否安装
2. **回退机制**：GStreamer 失败时自动回退到 FFmpeg
3. **详细日志**：记录完整的 pipeline 和错误信息
4. **文档完善**：提供故障排除指南

### 限制
1. 第一版不支持自定义 GStreamer pipeline
2. 不支持引擎热切换
3. GStreamer 仅支持 Linux（Windows/macOS 后续）

## 📈 成功指标

1. **功能完整性**：支持所有输入源
2. **延迟改善**：GStreamer 延迟 < 200ms
3. **性能**：CPU 占用 < 50%
4. **稳定性**：24 小时运行无崩溃
5. **用户体验**：Web UI 配置流畅

## 🔄 向后兼容

### 配置兼容
- 现有配置自动使用 FFmpeg 引擎
- 无需修改任何配置文件

### API 兼容
- 状态 API 增加引擎信息字段
- 所有现有 API 保持不变

## 📚 参考资料

- [GStreamer Documentation](https://gstreamer.freedesktop.org/documentation/)
- [GStreamer RTSP Server](https://github.com/GStreamer/gst-rtsp-server)
- [VAAPI Hardware Acceleration](https://wiki.archlinux.org/title/Hardware_video_acceleration)

## ✅ 验收标准

1. ✅ 支持 FFmpeg 和 GStreamer 两种引擎
2. ✅ 默认使用 FFmpeg，配置简单
3. ✅ 提供 4 种预设配置
4. ✅ Web UI 支持引擎选择
5. ✅ GStreamer 延迟 < 200ms
6. ✅ 所有输入源正常工作
7. ✅ 文档完整，包含故障排除
8. ✅ 向后兼容现有配置
