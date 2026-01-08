# 实施任务清单

## Phase 1: 基础架构 (MVP)

### 1. Web UI 后端模块
- [ ] 1.1 创建 `internal/webui/` 模块目录结构
- [ ] 1.2 实现 HTTP 服务器（`server.go`）
- [ ] 1.3 实现静态资源嵌入（`embedded.go`，使用 `embed`）
- [ ] 1.4 实现 API 处理器（`handlers.go`）
  - [ ] 获取当前配置状态（GET /api/status）
  - [ ] 验证 API 服务器连接（POST /api/validate-server）
  - [ ] 用户登录（POST /api/login）
  - [ ] 获取设备列表（GET /api/devices）
  - [ ] 创建设备（POST /api/devices）
  - [ ] 保存配置（POST /api/save-config）
  - [ ] 检测可用输入源（GET /api/input-sources）

### 2. 前端页面开发
- [ ] 2.1 创建 `web/` 目录存放前端源文件
- [ ] 2.2 实现首次配置向导页面（`setup.html`）
  - [ ] 欢迎页面
  - [ ] API 服务器配置
  - [ ] 用户登录
  - [ ] 设备选择/创建
  - [ ] 视频源配置
  - [ ] 完成页面
- [ ] 2.3 实现样式文件（`style.css`）
- [ ] 2.4 实现前端逻辑（`app.js`）
  - [ ] 表单验证
  - [ ] API 调用封装
  - [ ] 步骤导航控制
  - [ ] 错误处理和显示

### 3. 主程序集成
- [ ] 3.1 修改 `cmd/edge-client/main.go` 支持 `setup` 子命令
- [ ] 3.2 实现配置状态检测逻辑
- [ ] 3.3 在 `start` 命令中，检测未配置时自动进入 setup 模式
- [ ] 3.4 实现配置完成后的自动重启逻辑

### 4. 配置管理增强
- [ ] 4.1 添加配置状态标记（`configured: true/false`）
- [ ] 4.2 实现从 API 响应生成 `config.yaml` 的逻辑
- [ ] 4.3 实现配置文件写入和权限设置

## Phase 2: 状态仪表盘（可选）

### 5. Dashboard 实现
- [ ] 5.1 创建 `dashboard.html` 页面
- [ ] 5.2 实现状态 API
  - [ ] 推流状态（GET /api/stream/status）
  - [ ] 系统信息（GET /api/system/info）
  - [ ] 最近日志（GET /api/logs/recent）
- [ ] 5.3 实现 `dashboard` 子命令
- [ ] 5.4 添加简单的 HTTP Basic Auth 认证

## Phase 3: 增强功能（未来）

### 6. 高级特性
- [ ] 6.1 实现 mDNS 设备发现（`pawstream-edge.local`）
- [ ] 6.2 添加临时访问码验证
- [ ] 6.3 实现扫码绑定功能
- [ ] 6.4 添加配置导入/导出

## Phase 4: 文档和测试

### 7. 文档
- [ ] 7.1 更新 README.md，添加 Web UI 配置说明
- [ ] 7.2 添加 Web UI 使用截图
- [ ] 7.3 编写 Web UI 开发文档

### 8. 测试
- [ ] 8.1 测试首次配置流程
- [ ] 8.2 测试配置文件生成
- [ ] 8.3 测试各种浏览器兼容性
- [ ] 8.4 测试错误处理（API 不可用、认证失败等）
