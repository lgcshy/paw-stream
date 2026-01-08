# PawStream API Server

PawStream 控制面 API 服务器 - 基于 Go + Fiber + SQLite 构建的流媒体管理系统。

## 特性

- ✅ **分层架构**: Transport → Domain → Storage 清晰分离
- ✅ **SQLite 存储**: CGO-free 纯 Go 实现，单文件数据库，易于部署
- ✅ **JWT 认证**: 业务用户身份验证
- ✅ **设备鉴权**: 设备推流 secret 管理
- ✅ **MediaMTX 集成**: 流媒体鉴权回调
- ✅ **结构化日志**: Zerolog JSON 日志 + Lumberjack 自动轮转
- ✅ **静态编译**: 无 CGO 依赖，跨平台编译

## 项目结构

```
server/api/
├── cmd/api/              # 应用入口
├── internal/
│   ├── app/api/          # 应用组装和路由
│   ├── config/           # 配置管理
│   ├── transport/http/   # HTTP 层
│   │   ├── handlers/     # HTTP 处理器
│   │   └── middleware/   # 中间件
│   ├── domain/           # 业务领域
│   │   ├── user/         # 用户管理
│   │   ├── device/       # 设备管理
│   │   └── acl/          # 访问控制
│   ├── store/sqlite/     # SQLite 存储层
│   └── pkg/              # 工具包
├── migrations/           # 数据库迁移
├── deployments/          # Docker 部署配置
└── scripts/              # 开发脚本
```

## 快速开始

### 环境要求

- Go 1.24+
- SQLite (embedded, 无需独立安装)

### 开发运行

```bash
# 1. 安装依赖
go mod download

# 2. 创建配置文件 (可选，使用默认值)
cp config.yaml.example config.yaml

# 3. 运行服务器
go run cmd/api/main.go

# 或使用热重载 (推荐)
./scripts/dev_run.sh
```

服务器将在 `http://localhost:3000` 启动。

### 构建生产二进制

```bash
./scripts/build.sh
# 生成: bin/api
```

### Docker 部署

```bash
cd deployments
docker-compose up -d
```

## 配置

配置优先级: **环境变量 > config.yaml > 默认值**

### 配置文件示例 (config.yaml)

```yaml
server:
  port: "3000"
  host: "0.0.0.0"
  mode: "development"  # or "production"

log:
  level: "info"
  file: "logs/api.log"
  console: true
  max_size: 100       # MB
  max_backups: 7
  max_age: 30         # days
  compress: true

db:
  path: "data/pawstream.db"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime: "1h"

jwt:
  secret: "change-me-in-production"
  expiry: "24h"

mediamtx:
  url: "http://localhost:8554"
```

### 环境变量

```bash
PAWSTREAM_SERVER_PORT=3000
PAWSTREAM_SERVER_MODE=production
PAWSTREAM_LOG_LEVEL=info
PAWSTREAM_LOG_FILE=/var/log/pawstream/api.log
PAWSTREAM_DB_PATH=/var/lib/pawstream/data.db
PAWSTREAM_JWT_SECRET=your-secret-key
```

## API 端点

### 健康检查

```
GET /health
```

返回服务器健康状态。

### MediaMTX 鉴权回调

```
POST /mediamtx/auth
```

MediaMTX 流媒体服务器鉴权回调接口。

### 用户认证

```
POST /api/register   # 用户注册 ✅
POST /api/login      # 用户登录 ✅
GET  /api/me         # 获取当前用户信息 ✅
```

### 设备管理

```
GET    /api/devices                      # 列出设备 ✅
POST   /api/devices                      # 创建设备 ✅
GET    /api/devices/:id                  # 获取设备详情 ✅
PUT    /api/devices/:id                  # 更新设备 ✅
DELETE /api/devices/:id                  # 删除设备 ✅
POST   /api/devices/:id/rotate-secret    # 轮换设备 Secret ✅
```

### 路径查询

```
GET /api/paths   # 列出可访问的流路径 ✅
```

## 数据库

### SQLite 优势

- ✅ **零配置**: 无需独立数据库服务器
- ✅ **单文件**: 易于备份和迁移
- ✅ **CGO-free**: 静态编译，简化部署
- ✅ **性能充足**: 适合 4-8 路流的场景
- ✅ **WAL 模式**: 支持读写并发

### 数据库文件

默认位置: `data/pawstream.db`

### 迁移管理

迁移在服务器启动时自动执行。手动运行:

```bash
./scripts/migrate.sh up      # 升级
./scripts/migrate.sh down    # 回滚
./scripts/migrate.sh version # 查看版本
```

### 备份 SQLite 数据库

```bash
# 简单备份 (服务器停止时)
cp data/pawstream.db data/pawstream.db.backup

# 在线备份 (服务器运行时)
sqlite3 data/pawstream.db ".backup data/pawstream.db.backup"
```

## 日志管理

### 日志输出

- **开发模式**: 控制台 (美化格式) + 文件 (JSON)
- **生产模式**: 文件 (JSON)

### 日志轮转

使用 Lumberjack 自动轮转:
- 按大小: 100MB/文件
- 按时间: 保留 30 天
- 压缩: 旧日志自动 gzip 压缩

### 查看日志

```bash
# 实时查看
tail -f logs/api.log

# 解析 JSON 日志 (使用 jq)
tail -f logs/api.log | jq .
```

## 开发

### 代码检查

```bash
# 格式化
go fmt ./...

# 静态检查
go vet ./...

# 运行测试
go test ./...
```

### 项目约定

参考 `openspec/project.md` 了解项目规范和约定。

## 生产部署

### 安全检查清单

- [ ] 修改 JWT secret (`PAWSTREAM_JWT_SECRET`)
- [ ] 配置 CORS 允许的来源
- [ ] 设置日志级别为 `info` 或 `warn`
- [ ] 配置反向代理 (Nginx/Caddy)
- [ ] 启用 HTTPS
- [ ] 定期备份数据库文件
- [ ] 监控磁盘空间 (日志和数据库)

### 性能考虑

- SQLite 适合 4-8 路并发流 + 少量用户
- WAL 模式提升并发性能
- 如需更高并发,可迁移到 PostgreSQL (接口已抽象)

## 故障排查

### 数据库锁定

如果遇到 "database is locked" 错误:

1. 检查是否有多个服务器实例
2. 确认 WAL 模式已启用
3. 检查磁盘 I/O 性能

### 日志文件占满磁盘

Lumberjack 会自动轮转和清理,如需手动清理:

```bash
# 删除旧日志
find logs/ -name "*.log.*" -mtime +30 -delete
```

### 迁移失败

```bash
# 检查迁移状态
./scripts/migrate.sh version

# 手动回滚
./scripts/migrate.sh down
```

## API 使用指南

完整的 API 使用说明请参考:
- **API_GUIDE.md** - 详细的 API 文档和示例

快速示例:

```bash
# 1. 注册用户
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123","nickname":"Demo User"}'

# 2. 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123"}' | jq -r '.token')

# 3. 创建设备
curl -X POST http://localhost:3000/api/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"我的摄像头","location":"客厅"}'

# 4. 列出设备
curl http://localhost:3000/api/devices \
  -H "Authorization: Bearer $TOKEN"

# 5. 列出可访问路径
curl http://localhost:3000/api/paths \
  -H "Authorization: Bearer $TOKEN"
```

## Phase 3 完成状态

✅ **用户认证**: 注册、登录、获取用户信息  
✅ **设备管理**: 完整的 CRUD + Secret 轮换  
✅ **路径查询**: 获取可访问的流路径  
✅ **MediaMTX 集成**: 发布和读取鉴权完整实现  
✅ **集成测试**: 所有 API 端点测试通过  

## 许可证

(待定)

## 相关文档

- **API 使用指南**: `API_GUIDE.md`
- **验证报告**: `VERIFICATION_REPORT.md`
- 架构设计: `docs/backend_project_layout.md`
- 部署指南: `docs/deployment.md`
- OpenSpec 规范: `openspec/specs/api-server/spec.md`
