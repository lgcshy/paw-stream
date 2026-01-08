# Design Document: API Server Initialization

## Context

PawStream 需要一个控制面 API 服务器来管理用户、设备和权限，并为 MediaMTX 媒体服务器提供鉴权回调。这是一个典型的 REST API 服务，需要支持：

1. 业务用户的注册和登录
2. 设备的注册、绑定和凭证管理
3. MediaMTX 的 publish/read 鉴权回调
4. 未来可能的扩展：录像管理、AI 功能等

**约束条件：**
- 单人开发项目，需要保持简单和可维护
- 初期无实体硬件，需要支持模拟数据进行开发
- 必须可部署在低成本服务器上
- 需要与 MediaMTX 集成（外部依赖）

**参考文档：**
- `docs/backend_project_layout.md` - 详细的项目结构设计
- `openspec/project.md` - 项目约定和原则

---

## Goals / Non-Goals

### Goals
- 搭建清晰、可扩展的 Go 项目结构
- 实现快速验证的内存存储方案
- 准备生产级 PostgreSQL 存储方案
- 为 Phase 3 的功能实现打好基础
- 支持无数据库的本地开发

### Non-Goals
- 完整的业务逻辑实现（留给 Phase 3）
- 高级功能（录像、AI、通知等）
- 性能优化和压测
- Kubernetes 部署配置
- 前端 API 集成（等待 Phase 4）

---

## Decisions

### 1. Web 框架选择：Fiber v2

**决策：** 使用 Fiber v2 作为 HTTP 框架

**理由：**
- 性能优秀，基于 fasthttp
- Express 风格的 API，易于上手
- 内置中间件支持（logger, CORS, recover）
- 轻量级，适合单体应用
- 活跃的社区和文档

**替代方案：**
- **Gin**: 更流行，但 Fiber 性能更好且 API 更现代
- **Echo**: 功能相似，但 Fiber 的中间件生态更丰富
- **标准库 net/http**: 过于底层，需要自己实现太多功能

---

### 2. 项目结构：分层架构 + 依赖注入

**决策：** 采用 `transport → domain → store` 三层架构

**结构说明：**
```
transport/http/    # HTTP 层：handlers, middleware
  → domain/        # 业务逻辑：user, device, acl
    → store/       # 数据持久化：memory, postgres
```

**理由：**
- 清晰的职责分离，便于测试
- 依赖倒置：domain 定义接口，store 实现接口
- 易于替换实现（内存 ↔ PostgreSQL）
- 符合 Go 社区的最佳实践

**关键原则：**
- `internal/` 限制外部包导入，防止 API 泄漏
- Domain 层不依赖具体的 transport 或 store 实现
- 使用接口（Repository）而非具体类型

---

### 3. 存储策略：SQLite（CGO-free）

**决策：** 使用 SQLite 作为唯一数据库，使用 modernc.org/sqlite 驱动（CGO-free）

**SQLite 存储（`internal/store/sqlite/`）：**
- 单文件数据库（默认：`data/pawstream.db`）
- 无需独立数据库服务器
- 使用 modernc.org/sqlite 驱动（纯 Go 实现，无 CGO 依赖）
- 支持连接池和事务
- 启用 WAL 模式提升并发性能
- 通过迁移脚本管理 schema

**CGO-free 的优势：**
- 静态编译：`CGO_ENABLED=0 go build` 生成完全独立的二进制
- 跨平台编译：无需交叉编译 C 库
- 简化部署：单个二进制文件即可运行
- 容器镜像更小：无需 libc 等依赖

**性能考虑：**
- SQLite 在 4-8 路并发流场景下性能充足
- 单机部署，无网络延迟
- WAL 模式支持读写并发
- 对于小型应用（< 100 并发用户）完全够用

**理由：**
- 简化部署，无需配置数据库服务器
- 数据持久化（相比内存存储）
- 易于备份（复制单个文件）
- 适合边缘设备和小型服务器
- CGO-free 简化编译和分发

---

### 4. 配置管理：Viper + YAML + 环境变量

**决策：** 使用 Viper 读取 YAML 配置文件，支持环境变量覆盖

**配置优先级：**
```
环境变量 > config.yaml > 默认值
```

**示例配置：**
```yaml
server:
  port: 3000
  host: "0.0.0.0"
  mode: "development"  # or "production"

log:
  level: "info"              # debug, info, warn, error
  file: "logs/api.log"       # 日志文件路径
  console: true              # 是否同时输出到控制台
  max_size: 100              # 每个日志文件最大大小（MB）
  max_backups: 7             # 保留旧日志文件数量
  max_age: 30                # 保留旧日志文件天数
  compress: true             # 是否压缩旧日志

db:
  path: "data/pawstream.db"  # SQLite database file path
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "1h"

jwt:
  secret: "change-me-in-production"
  expiry: "24h"

mediamtx:
  url: "http://localhost:8554"
```

**环境变量格式：**
```bash
PAWSTREAM_SERVER_PORT=3000
PAWSTREAM_SERVER_MODE=production
PAWSTREAM_LOG_LEVEL=info
PAWSTREAM_LOG_FILE=/var/log/pawstream/api.log
PAWSTREAM_DB_PATH=/var/lib/pawstream/data.db
PAWSTREAM_JWT_SECRET=my-secret
```

**理由：**
- YAML 易读易写，适合开发配置
- 环境变量适合容器化部署
- Viper 支持热重载（可选）
- 默认值确保零配置启动

---

### 5. 认证方案：JWT + Device Secret

**业务用户认证：**
- 使用 JWT（JSON Web Tokens）
- 登录成功后签发 token
- Token 包含 `user_id`、`username`、`exp` 等 claims
- 存储在 HTTP header：`Authorization: Bearer <token>`
- 中间件验证 token 并注入用户信息到 context

**设备推流认证：**
- 每个设备有唯一的 `device_secret`（高强度随机字符串）
- MediaMTX 调用 `/mediamtx/auth` 时携带 secret
- API 验证 secret 是否与 path 匹配
- Secret 存储：hash 用于验证，加密后存储用于取回

**理由：**
- JWT 无状态，易于水平扩展
- Device secret 适合机器对机器认证
- 分离用户和设备的认证逻辑
- 支持 secret 轮换（安全性）

---

### 6. 密码和 Secret 处理

**用户密码：**
- 使用 bcrypt 哈希（cost=12）
- 仅存储 hash，不存储明文
- 登录时验证输入与 hash

**设备 Secret：**
- 生成：`crypto/rand` 生成 32 字节随机数，base64 编码
- 存储两份：
  - `secret_hash`: bcrypt hash，用于验证
  - `secret_cipher`: AES 加密，用于用户"复制 secret"功能
- 轮换时递增 `secret_version`

**理由：**
- bcrypt 防止彩虹表攻击
- AES 加密允许用户取回 secret（重新配置设备）
- Secret 轮换提升安全性

---

### 7. MediaMTX 集成方案

**鉴权回调流程：**
```
Device/User → MediaMTX → POST /mediamtx/auth → API Server
                                                    ↓
                                          验证 secret/token
                                                    ↓
                                          返回 200 OK / 403 Forbidden
```

**回调请求格式（MediaMTX 文档）：**
```json
{
  "action": "publish",     // 或 "read", "playback"
  "path": "dogcam/abc123",
  "protocol": "rtsp",
  "ip": "192.168.1.100",
  "user": "",              // basic auth user（可选）
  "password": "",          // basic auth password（可选）
  "token": ""              // JWT token（query param）
}
```

**API 验证逻辑：**
- `action=publish`: 验证 device_secret（从 `password` 或 query param）
- `action=read/playback`: 验证 user_token（从 `token` 或 header）
- 返回 2xx 允许，4xx 拒绝

**配置（mediamtx.yml）：**
```yaml
authHTTPAddress: http://api:3000/mediamtx/auth
```

**理由：**
- MediaMTX 原生支持 HTTP 回调
- API 可完全控制鉴权逻辑
- 日志记录所有鉴权尝试
- 支持动态权限更新

---

### 8. 日志方案：zerolog + lumberjack

**决策：** 使用 zerolog 进行结构化 JSON 日志，使用 lumberjack 进行日志轮转

**日志输出：**
- **开发模式**: 控制台输出（美化格式）+ 文件输出（JSON）
- **生产模式**: 文件输出（JSON）

**日志格式：**
```json
{
  "level": "info",
  "time": "2026-01-06T15:30:00Z",
  "request_id": "abc-123",
  "user_id": "user-456",
  "method": "POST",
  "path": "/api/login",
  "status": 200,
  "duration_ms": 45,
  "message": "request completed"
}
```

**日志轮转配置（lumberjack）：**
```go
&lumberjack.Logger{
    Filename:   "logs/api.log",  // 日志文件路径
    MaxSize:    100,              // 每个日志文件最大大小（MB）
    MaxBackups: 7,                // 保留旧日志文件的最大数量
    MaxAge:     30,               // 保留旧日志文件的最大天数
    Compress:   true,             // 是否压缩旧日志文件（gzip）
}
```

**日志级别：**
- **debug**: 详细的调试信息（仅开发模式）
- **info**: 常规信息（默认级别）
- **warn**: 警告信息（需要注意但不影响运行）
- **error**: 错误信息（需要处理）
- **fatal**: 致命错误（导致程序退出）

**理由：**
- zerolog: JSON 格式易于解析和搜索（ELK、Loki），零分配高性能
- lumberjack: 自动日志轮转，防止磁盘占满
- 按大小和时间轮转，灵活可配
- 压缩旧日志，节省磁盘空间
- 支持上下文字段（request_id, user_id）

---

### 9. 错误处理：分层错误类型

**错误分类：**
```go
// Domain 层错误
ErrUserNotFound
ErrDeviceNotFound
ErrDuplicateUsername
ErrInvalidCredentials

// Transport 层映射
ErrUserNotFound → 404 Not Found
ErrDuplicateUsername → 409 Conflict
ErrInvalidCredentials → 401 Unauthorized
```

**错误响应格式：**
```json
{
  "error": "user_not_found",
  "message": "User with ID 'abc' not found",
  "request_id": "xyz-123"
}
```

**理由：**
- 清晰的错误语义
- 便于客户端处理
- 不泄漏内部实现细节
- 统一的错误格式

---

### 10. 数据库迁移：golang-migrate + Auto-migration

**决策：** 使用 golang-migrate 管理数据库 schema，启动时自动执行

**迁移文件命名：**
```
migrations/
├── 001_init_schema.up.sql
├── 001_init_schema.down.sql
├── 002_add_device_location.up.sql
└── 002_add_device_location.down.sql
```

**执行方式：**
1. **自动迁移（推荐）：** 服务器启动时自动检测并执行未应用的迁移
2. **手动迁移：**
   ```bash
   ./scripts/migrate.sh up      # 升级
   ./scripts/migrate.sh down    # 回滚
   ./scripts/migrate.sh version # 查看版本
   ```

**SQLite 特殊配置：**
- 启用 WAL 模式：`PRAGMA journal_mode=WAL;`
- 启用外键约束：`PRAGMA foreign_keys=ON;`
- 初始化迁移表

**理由：**
- 版本化的 schema 管理
- 自动迁移简化部署
- 支持回滚
- SQLite 兼容性好

---

## Risks / Trade-offs

### Risk 1: SQLite 并发写入限制
**风险：** SQLite 对并发写入有限制，可能成为瓶颈

**缓解措施：**
- 启用 WAL 模式，提升并发读写性能
- 对于 4-8 路流的场景，写入压力不大
- 未来如需扩展，可迁移到 PostgreSQL（接口已抽象）
- 监控数据库锁等待时间

### Risk 1.5: SQLite 文件损坏
**风险：** 突然断电或磁盘故障可能导致数据库文件损坏

**缓解措施：**
- 定期备份数据库文件
- WAL 模式提供更好的崩溃恢复
- 记录 SQLite 完整性检查日志
- 提供数据库修复脚本

### Risk 2.5: 日志文件占满磁盘
**风险：** 日志文件持续增长可能占满磁盘空间

**缓解措施：**
- 使用 lumberjack 自动轮转日志文件
- 按大小（100MB）和时间（30天）限制
- 压缩旧日志文件（gzip）
- 监控磁盘空间使用
- 提供日志清理脚本

### Risk 3: JWT Secret 泄漏
**风险：** JWT secret 泄漏导致 token 伪造

**缓解措施：**
- 默认 secret 附带警告
- 生产环境通过环境变量配置
- 文档强调 secret 的重要性
- 考虑支持 secret 轮换（Phase 4）

### Risk 4: MediaMTX 回调失败
**风险：** API 服务器不可用时，所有流无法鉴权

**缓解措施：**
- MediaMTX 支持重试
- API 服务器健康检查
- 日志记录所有鉴权失败
- 考虑降级策略（临时允许/拒绝）

### Risk 5: 单体架构扩展性
**风险：** 单体 API 可能成为性能瓶颈

**当前决策：** 暂不优化，理由：
- 项目初期，用户量小
- 4-8 路流 + 少量用户，单体足够
- 提前优化是万恶之源
- 需要时可拆分（User Service、Device Service）

---

## Migration Plan

### 数据库备份和恢复

**备份 SQLite 数据库：**
```bash
# 方法 1: 简单文件复制（服务器停止时）
cp data/pawstream.db data/pawstream.db.backup

# 方法 2: 在线备份（服务器运行时）
sqlite3 data/pawstream.db ".backup data/pawstream.db.backup"
```

**恢复数据库：**
```bash
# 停止服务器
systemctl stop pawstream-api

# 恢复备份
cp data/pawstream.db.backup data/pawstream.db

# 启动服务器
systemctl start pawstream-api
```

### 未来迁移到 PostgreSQL（如果需要）

如果未来需要更高的并发性能：

1. 实现 `internal/store/postgres/` 包
2. Repository 接口保持不变
3. 使用工具导出 SQLite 数据，导入 PostgreSQL
4. 修改配置切换存储实现
5. 重启服务器

**当前不需要，保持简单。**

---

## Open Questions

1. **SQLite 是否足够？**
   - 当前答案：对于 4-8 路流 + 小规模用户，完全足够
   - 重新评估时机：并发用户 > 100 或写入 QPS > 1000
   - 备选方案：PostgreSQL（接口已抽象，易于迁移）

2. **是否需要 Redis 缓存？**
   - 当前答案：Phase 1 不需要
   - 重新评估时机：QPS > 1000 或查询延迟 > 100ms

3. **是否需要 rate limiting？**
   - 当前答案：Phase 1 不需要
   - 后续可添加中间件（例如使用 Fiber 的 limiter）

4. **是否需要 API 版本化？**
   - 当前答案：暂不需要
   - API 路径前缀使用 `/api/` 为未来版本化预留空间

5. **是否需要 OpenAPI 文档？**
   - 当前答案：可选，创建骨架但暂不详细维护
   - Phase 3 实现业务 API 时补充完整

6. **设备离线检测机制？**
   - 当前答案：Phase 1 暂不实现
   - 依赖 MediaMTX 的流状态 API（Phase 3）

---

## Summary

本设计文档定义了 PawStream API 服务器的技术架构和关键决策。通过分层架构、SQLite 存储、清晰的认证方案和结构化日志，我们构建了一个简单但可扩展的系统基础。

**核心原则：**
- 简单优先，按需复杂化
- 清晰的分层和接口
- SQLite + CGO-free：简化部署，保持数据持久化
- 为 Phase 3 的功能实现打好基础

**SQLite 的优势：**
- ✅ 零配置：无需独立数据库服务器
- ✅ 单文件：易于备份和迁移
- ✅ CGO-free：静态编译，简化部署
- ✅ 性能充足：适合 4-8 路流的场景
- ✅ 持久化：相比内存存储，数据不丢失

**下一步：**
按照 `tasks.md` 逐步实现各个组件，优先实现核心路径（健康检查 → SQLite 连接 → 基础中间件 → MediaMTX 回调）。
