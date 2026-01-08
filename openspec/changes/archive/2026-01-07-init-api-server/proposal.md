# Change: Initialize API Server Project

## Why

PawStream 需要一个 Go 语言的 API 服务器作为控制面（Control Plane），负责用户认证、设备管理、权限控制，以及为 MediaMTX 提供鉴权回调接口。当前 `server/api/` 目录为空，需要搭建完整的 Go + Fiber 项目基础架构。

这是 Phase 3（Server API）的基础工作，为后续实现用户注册/登录、设备绑定、流媒体鉴权等功能奠定基础。

## What Changes

- 初始化 Go module 项目结构（基于 `docs/backend_project_layout.md` 设计）
- 搭建 Fiber web 框架基础
- 实现分层架构：transport → domain → store
- 配置 SQLite 数据库（CGO-free 版本）和迁移脚本
- 实现 SQLite 存储层，支持单文件数据库
- 添加基础中间件：request ID、logger、CORS
- 配置日志系统（支持文件输出、日志轮转、可配置路径）
- 创建健康检查端点
- 配置 YAML 文件读取
- 添加 JWT 工具包（业务用户认证）
- 准备 MediaMTX 鉴权回调接口结构

**No breaking changes** - 这是新能力的初始化。

## Impact

### Affected Specs
- **NEW**: `api-server` - Go API 服务器能力

### Affected Code
- `server/api/` 目录 - 将创建完整的 Go 项目结构
- 新文件结构（基于 `docs/backend_project_layout.md`）:
  - `cmd/api/main.go` - 应用入口
  - `internal/app/` - 应用组装和路由
  - `internal/config/` - 配置管理
  - `internal/transport/http/` - HTTP 层（handlers、middleware）
  - `internal/domain/` - 业务领域（user、device、acl）
  - `internal/store/sqlite/` - SQLite 存储层
  - `internal/pkg/` - 工具包（jwtutil、idgen、errors）
  - `internal/integration/mediamtx/` - MediaMTX 集成
  - `deployments/` - Docker Compose 配置
  - `scripts/` - 开发脚本
  - `migrations/` - 数据库迁移脚本

### Dependencies Added
- **Fiber v2** - Web 框架
- **modernc.org/sqlite** - CGO-free SQLite 驱动
- **golang-jwt** - JWT 处理
- **viper** - 配置管理
- **zerolog** - 结构化日志
- **lumberjack.v2** - 日志轮转（按大小和时间）
- **golang-migrate/migrate** - 数据库迁移工具（支持 SQLite）
- **testify** - 测试框架

### Database Schema
- `users` 表 - 业务用户信息
- `devices` 表 - 设备信息和推流凭证

## Success Criteria

- `go run cmd/api/main.go` 成功启动 API 服务器
- 健康检查端点 `GET /health` 返回 200
- SQLite 数据库文件自动创建（默认 `data/pawstream.db`）
- 数据迁移自动执行
- 项目结构符合 `docs/backend_project_layout.md` 设计
- 所有代码通过 `go fmt` 和 `go vet` 检查
- 日志输出结构化且清晰
- 日志文件自动创建并支持轮转
- 编译产物无 CGO 依赖，可静态链接

## Technical Decisions

### 1. 存储策略
- 使用 SQLite 作为唯一数据库（CGO-free 版本：modernc.org/sqlite）
- 单文件数据库，便于备份和迁移
- 无需独立数据库服务器，简化部署
- 支持 4-8 路并发流，性能充足

### 2. 配置管理
- 使用 YAML 配置文件（`config.yaml`）
- 支持环境变量覆盖
- 默认值在 `defaults.go` 中定义
- 日志路径和轮转策略可配置

### 3. 项目结构
- 遵循 Go 标准项目布局
- `internal/` 限制外部包导入
- 分层架构：清晰的职责分离
- 依赖注入：便于测试和替换实现

### 4. 鉴权设计
- 业务用户：JWT token（存储在 HTTP header 或 cookie）
- 设备推流：device_secret（通过 MediaMTX 回调验证）
- MediaMTX 回调接口：`POST /mediamtx/auth`
