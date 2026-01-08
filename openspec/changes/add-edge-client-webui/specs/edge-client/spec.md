# edge-client Spec Delta

## ADDED Requirements

### Requirement: Web UI 管理界面
边缘客户端 SHALL 提供内嵌的 Web UI 管理界面，随客户端启动，提供配置、监控、日志查看的统一界面。

#### Scenario: 启动时自动开启 Web UI
- **WHEN** 执行 `edge-client start` 命令
- **THEN** 应同时启动推流进程和 Web UI HTTP 服务器（默认端口 8080），并在日志中输出访问地址（如 `Web UI available at http://192.168.1.100:8080`）

#### Scenario: 通过浏览器访问 Web UI
- **WHEN** 用户在浏览器中访问 `http://<设备IP>:8080`
- **THEN** 应显示 Web UI 管理界面，包含【配置】【状态】【日志】三个标签页

#### Scenario: 禁用 Web UI
- **WHEN** 执行 `edge-client start --no-webui` 命令
- **THEN** 应仅启动推流进程，不启动 Web UI HTTP 服务器

#### Scenario: 自定义 Web UI 端口
- **WHEN** 执行 `edge-client start --webui-port 9090` 命令
- **THEN** Web UI 应监听端口 9090 而非默认的 8080

#### Scenario: Web UI 持续运行
- **WHEN** 客户端正在推流
- **THEN** Web UI 应持续可访问，可随时通过浏览器查看状态和修改配置

### Requirement: 配置热重载
边缘客户端 SHALL 监控配置文件变化，自动重新加载配置并应用，无需重启进程。

#### Scenario: Web UI 修改配置触发热重载
- **WHEN** 用户在 Web UI 的【配置】标签页修改配置并点击"保存"
- **THEN** 应将新配置写入 config.yaml，自动检测到文件变化，优雅停止当前推流，重新加载配置，重启推流，并在日志中记录"Config reloaded, streaming restarted"

#### Scenario: 手动编辑配置文件触发热重载
- **WHEN** 用户通过编辑器手动修改 config.yaml 并保存
- **THEN** 应自动检测到文件变化，执行配置重载和推流重启

#### Scenario: 无效配置不生效
- **WHEN** 配置文件被修改为无效内容（如缺少必需字段）
- **THEN** 应记录错误日志，拒绝重载，继续使用旧配置运行

#### Scenario: 配置重载失败回滚
- **WHEN** 新配置加载成功但推流启动失败（如输入源不可用）
- **THEN** 应记录错误日志，可选择回滚到旧配置或等待用户修复

### Requirement: 可视化配置管理
Web UI SHALL 提供【配置】标签页，允许用户可视化编辑所有配置项。

#### Scenario: 配置 API 服务器
- **WHEN** 用户在【配置】标签页输入 API 服务器地址（如 `http://192.168.1.100:3000`）
- **THEN** 应验证地址格式，并提供"测试连接"按钮，点击后验证服务器连通性

#### Scenario: 用户登录和设备绑定
- **WHEN** 用户输入 PawStream 账号用户名和密码并点击"登录"
- **THEN** 应调用 API 登录接口，获取 JWT token，然后获取用户的设备列表供选择

#### Scenario: 选择现有设备
- **WHEN** 用户从设备列表中选择一个设备
- **THEN** 应自动填充设备 ID 和密钥到配置中

#### Scenario: 创建新设备
- **WHEN** 用户选择"创建新设备"并输入设备名称和位置
- **THEN** 应调用 API 创建设备接口，获取新设备的 ID 和密钥，并自动填充到配置中

#### Scenario: 配置视频输入源
- **WHEN** 用户选择视频输入源类型（V4L2、RTSP、文件、测试）并填写对应参数
- **THEN** 应验证输入格式，并可选择自动检测可用的 V4L2 设备

#### Scenario: 配置高级参数
- **WHEN** 用户展开高级配置区块，修改视频编码参数、推流参数等
- **THEN** 应显示各参数的说明和默认值，并验证输入范围

#### Scenario: 保存配置
- **WHEN** 用户点击"保存配置"按钮
- **THEN** 应验证所有配置项，将配置写入 config.yaml，触发热重载，并显示成功提示（如"配置已保存并自动应用"）

### Requirement: 实时状态监控
Web UI SHALL 提供【状态】标签页，实时显示客户端运行状态。

#### Scenario: 查看推流状态
- **WHEN** 用户切换到【状态】标签页
- **THEN** 应显示当前推流状态（在线/离线）、运行时长、推流地址、视频源信息

#### Scenario: 查看系统信息
- **WHEN** 用户在【状态】标签页查看系统信息
- **THEN** 应显示 CPU 使用率、内存使用、磁盘空间、网络流量等系统指标

#### Scenario: 查看设备信息
- **WHEN** 用户在【状态】标签页查看设备信息
- **THEN** 应显示设备名称、设备 ID、API 服务器地址、客户端版本等信息

#### Scenario: 状态实时更新
- **WHEN** 推流状态发生变化（如从离线变为在线，或配置重载）
- **THEN** 【状态】标签页应通过 WebSocket 自动更新显示，无需刷新页面

#### Scenario: WebSocket 断线重连
- **WHEN** WebSocket 连接意外断开
- **THEN** 应自动尝试重新连接，并在连接恢复后继续接收状态更新

### Requirement: 日志查看
Web UI SHALL 提供【日志】标签页，显示客户端日志信息。

#### Scenario: 查看最近日志
- **WHEN** 用户切换到【日志】标签页
- **THEN** 应显示最近 100 条日志记录，包括时间戳、日志级别、消息内容

#### Scenario: 日志级别颜色区分
- **WHEN** 显示日志列表
- **THEN** 不同级别的日志应使用不同颜色标识（如 INFO-蓝色、WARN-黄色、ERROR-红色）

#### Scenario: 日志实时推送
- **WHEN** 客户端产生新日志
- **THEN** 应通过 WebSocket 实时推送到【日志】标签页，自动添加到日志列表并滚动到底部

#### Scenario: 日志过滤
- **WHEN** 用户选择日志级别过滤（如仅显示 WARN 和 ERROR）
- **THEN** 应仅显示符合条件的日志记录

#### Scenario: 日志搜索
- **WHEN** 用户在搜索框输入关键词
- **THEN** 应高亮显示包含该关键词的日志行

#### Scenario: 清空日志显示
- **WHEN** 用户点击"清空日志"按钮
- **THEN** 应清空前端显示的日志列表（不影响服务器端日志）

### Requirement: 首次配置引导
Web UI SHALL 在客户端首次启动时提供配置引导。

#### Scenario: 首次启动检测
- **WHEN** 客户端首次启动且无有效配置文件
- **THEN** Web UI 应显示欢迎提示和配置引导信息，引导用户完成首次配置

#### Scenario: 配置向导步骤提示
- **WHEN** 用户在首次配置过程中
- **THEN** 应显示当前步骤（如"步骤 1/4：配置 API 服务器"）和进度指示

#### Scenario: 配置完成自动生效
- **WHEN** 用户完成首次配置并保存
- **THEN** 应自动加载配置，开始推流，并在 Web UI 中显示成功提示和推流状态

### Requirement: 静态资源嵌入
Web UI 的所有静态资源 SHALL 嵌入到 Go 二进制文件中。

#### Scenario: 单一二进制部署
- **WHEN** 用户下载或编译 edge-client 二进制文件
- **THEN** 该文件应包含所有必需的 Web UI 资源（HTML、CSS、JavaScript），无需额外文件

#### Scenario: 离线运行
- **WHEN** 边缘客户端在没有互联网连接的环境中运行
- **THEN** Web UI 应正常工作，所有资源从嵌入的静态文件加载

### Requirement: 响应式设计
Web UI SHALL 支持多种设备和屏幕尺寸。

#### Scenario: 桌面浏览器访问
- **WHEN** 用户在桌面浏览器（1920x1080）中访问 Web UI
- **THEN** 应正常显示所有功能，布局清晰易用

#### Scenario: 移动设备访问
- **WHEN** 用户在手机浏览器（375x667）中访问 Web UI
- **THEN** 应自动调整布局，适配小屏幕，保持功能可用

#### Scenario: 平板设备访问
- **WHEN** 用户在平板浏览器（768x1024）中访问 Web UI
- **THEN** 应自动调整布局，充分利用屏幕空间

### Requirement: 安全性（可选）
Web UI SHALL 支持可选的访问认证。

#### Scenario: 启用 HTTP Basic Auth
- **WHEN** 配置文件中启用 webui.auth.enabled 并设置用户名密码
- **THEN** 访问 Web UI 时应要求输入用户名和密码

#### Scenario: 认证失败
- **WHEN** 用户输入错误的用户名或密码
- **THEN** 应拒绝访问并显示"401 Unauthorized"

#### Scenario: 审计日志
- **WHEN** 用户通过 Web UI 修改配置
- **THEN** 应在日志中记录操作详情（时间、操作类型、修改内容）

### Requirement: 性能要求
Web UI SHALL 保持低资源占用。

#### Scenario: HTTP 服务器资源消耗
- **WHEN** Web UI HTTP 服务器空闲时（无客户端连接）
- **THEN** 应占用不超过 10MB 内存，CPU 使用率接近 0%

#### Scenario: 多客户端并发访问
- **WHEN** 最多 5 个客户端同时访问 Web UI
- **THEN** 应保持正常响应速度，不影响推流性能

#### Scenario: WebSocket 连接稳定性
- **WHEN** WebSocket 连接建立后持续 24 小时
- **THEN** 连接应保持稳定，不发生内存泄漏或性能下降
