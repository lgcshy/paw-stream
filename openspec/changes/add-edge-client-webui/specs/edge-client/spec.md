# edge-client Spec Delta

## ADDED Requirements

### Requirement: Web UI 配置界面
边缘客户端 SHALL 提供内嵌的 Web UI 配置界面，降低首次配置的技术门槛。

#### Scenario: 首次启动自动进入配置模式
- **WHEN** 客户端首次启动且未检测到有效配置文件
- **THEN** 应自动启动 HTTP 服务器（端口 8080）并输出访问地址日志（如 `Web UI available at http://192.168.1.100:8080`）

#### Scenario: 通过浏览器访问配置向导
- **WHEN** 用户在浏览器中访问 `http://<设备IP>:8080`
- **THEN** 应显示首次配置向导页面，包括欢迎信息和配置步骤导航

#### Scenario: 手动进入配置模式
- **WHEN** 执行 `edge-client setup` 命令
- **THEN** 应启动 Web UI 配置服务器，允许用户重新配置

### Requirement: API 服务器配置
Web UI SHALL 允许用户配置 PawStream API 服务器地址。

#### Scenario: 输入 API 服务器地址
- **WHEN** 用户在配置向导中输入 API 服务器地址（如 `http://192.168.1.100:3000`）并点击"下一步"
- **THEN** 应验证地址格式是否有效（URL 格式）

#### Scenario: 验证 API 服务器连通性
- **WHEN** 用户提交 API 服务器地址
- **THEN** 应调用该服务器的健康检查端点（`/health`），验证连通性并显示结果

#### Scenario: API 服务器不可达
- **WHEN** API 服务器地址无法连接
- **THEN** 应显示明确的错误提示（如"无法连接到 API 服务器，请检查地址和网络"）并允许用户重新输入

### Requirement: 用户认证和设备绑定
Web UI SHALL 允许用户通过登录账号来绑定设备。

#### Scenario: 用户登录
- **WHEN** 用户输入 PawStream 账号用户名和密码并提交
- **THEN** 应调用 API 登录接口（`POST /api/login`），获取 JWT token

#### Scenario: 登录失败
- **WHEN** 用户输入错误的用户名或密码
- **THEN** 应显示登录失败提示（如"用户名或密码错误"）并允许重试

#### Scenario: 选择现有设备
- **WHEN** 用户登录成功后，查看设备列表
- **THEN** 应显示该用户所有已创建的设备，允许用户选择一个设备进行绑定

#### Scenario: 创建新设备
- **WHEN** 用户选择"创建新设备"并输入设备名称和位置
- **THEN** 应调用 API 创建设备接口（`POST /api/devices`），获取设备 ID 和密钥

#### Scenario: 设备绑定成功
- **WHEN** 用户选择或创建设备成功
- **THEN** 应显示设备信息（名称、ID）并自动保存设备 ID 和密钥到配置中

### Requirement: 视频输入源配置
Web UI SHALL 允许用户配置视频输入源。

#### Scenario: 选择输入源类型
- **WHEN** 用户在配置向导中选择视频输入源类型（V4L2 摄像头、RTSP、文件、测试模式）
- **THEN** 应根据选择显示对应的配置字段

#### Scenario: V4L2 摄像头配置
- **WHEN** 用户选择 V4L2 摄像头并输入设备路径（如 `/dev/video0`）
- **THEN** 应验证设备路径格式并保存到配置中

#### Scenario: RTSP 输入配置
- **WHEN** 用户选择 RTSP 输入并输入 RTSP URL（如 `rtsp://192.168.1.100:554/stream`）
- **THEN** 应验证 URL 格式并保存到配置中

#### Scenario: 测试模式配置
- **WHEN** 用户选择测试模式
- **THEN** 应自动配置为生成测试画面，无需额外输入

#### Scenario: 自动检测可用输入源（可选）
- **WHEN** 用户进入视频源配置页面
- **THEN** 应尝试检测系统中可用的 V4L2 设备并在下拉列表中显示

### Requirement: 配置文件生成和保存
Web UI SHALL 根据用户输入生成配置文件并保存。

#### Scenario: 生成 config.yaml
- **WHEN** 用户完成所有配置步骤并点击"完成配置"
- **THEN** 应根据用户输入生成完整的 `config.yaml` 文件，包括设备信息、API 地址、输入源等

#### Scenario: 保存配置文件
- **WHEN** 配置文件生成完成
- **THEN** 应将配置文件保存到默认位置（如 `./config.yaml` 或 `/etc/pawstream/config.yaml`）

#### Scenario: 配置保存成功提示
- **WHEN** 配置文件成功保存
- **THEN** 应显示成功提示和下一步操作指引（如"配置完成！客户端将自动重启并开始推流"）

#### Scenario: 配置保存失败
- **WHEN** 配置文件保存失败（如权限不足）
- **THEN** 应显示错误提示和解决建议（如"无法保存配置文件，请检查文件权限"）

### Requirement: 配置完成后自动启动
边缘客户端 SHALL 在配置完成后自动进入正常运行模式。

#### Scenario: 配置完成自动重启
- **WHEN** 用户完成 Web UI 配置
- **THEN** Web UI 应关闭，客户端应重新启动并进入正常推流模式

#### Scenario: 配置完成后 Web UI 关闭
- **WHEN** 配置完成并保存
- **THEN** Web UI HTTP 服务器应停止监听，释放端口 8080

### Requirement: 状态仪表盘（可选）
边缘客户端 SHALL 提供可选的状态仪表盘，显示运行状态。

#### Scenario: 启动状态仪表盘
- **WHEN** 执行 `edge-client dashboard` 命令
- **THEN** 应启动 HTTP 服务器（端口 8080）并显示状态仪表盘页面

#### Scenario: 查看推流状态
- **WHEN** 用户在浏览器中访问仪表盘
- **THEN** 应显示当前推流状态（在线/离线）、连接时长、推流路径等信息

#### Scenario: 查看系统信息
- **WHEN** 用户在仪表盘中查看系统信息
- **THEN** 应显示 CPU 使用率、内存使用、运行时长等系统指标

#### Scenario: 查看最近日志
- **WHEN** 用户在仪表盘中查看日志
- **THEN** 应显示最近 50-100 条日志记录，包括时间、级别、消息

### Requirement: 静态资源嵌入
Web UI 的所有静态资源 SHALL 嵌入到 Go 二进制文件中。

#### Scenario: 单一二进制部署
- **WHEN** 用户下载或编译 edge-client 二进制文件
- **THEN** 该文件应包含所有必需的 Web UI 资源（HTML、CSS、JavaScript），无需额外文件

#### Scenario: 离线运行
- **WHEN** 边缘客户端在没有互联网连接的环境中运行
- **THEN** Web UI 应正常工作，所有资源从嵌入的静态文件加载

### Requirement: 安全性
Web UI SHALL 在局域网环境下安全运行。

#### Scenario: 仅监听局域网地址
- **WHEN** Web UI HTTP 服务器启动
- **THEN** 应默认监听所有接口（`0.0.0.0:8080`），但不应对外网开放（通过防火墙配置）

#### Scenario: 配置模式访问控制（可选）
- **WHEN** 客户端首次启动进入配置模式
- **THEN** 可选择要求输入临时访问码（6 位数字，显示在日志中）

#### Scenario: Dashboard 访问认证
- **WHEN** 用户访问状态仪表盘（`edge-client dashboard`）
- **THEN** 应要求输入设备密钥或简单密码进行认证
