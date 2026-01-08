# Change: Implement User Authentication and Device Management APIs

## Why

init-api-server 提案已完成基础架构搭建,但业务 API 端点仍是占位符(返回 501 Not Implemented)。为了完成 Phase 3 的目标,需要实现完整的用户认证和设备管理功能,让系统能够:

1. 支持业务用户注册和登录
2. 管理设备的完整生命周期(创建、查询、更新、删除)
3. 为 MediaMTX 提供完整的鉴权能力
4. 让用户能够查看自己可以访问的流路径

这是连接"控制面"和"媒体面"的关键一步,完成后整个系统将具备端到端的流媒体管理能力。

## What Changes

### 用户认证 (Authentication)
- 实现用户注册 API (`POST /api/register`)
  - 用户名唯一性校验
  - 密码强度要求
  - Bcrypt 哈希存储
- 实现用户登录 API (`POST /api/login`)
  - 凭证验证
  - JWT token 签发
  - Token 过期时间配置
- 实现用户信息查询 (`GET /api/me`)
  - 返回当前登录用户信息
  - 需要 JWT 认证

### 设备管理 (Device Management)
- 实现设备创建 API (`POST /api/devices`)
  - 生成唯一 device_id
  - 生成加密随机 device_secret
  - 自动设置 publish_path (格式: `dogcam/<device_id>`)
  - 仅返回一次 secret (安全性)
- 实现设备列表 API (`GET /api/devices`)
  - 仅返回当前用户的设备
  - 支持分页
  - Secret 不暴露在响应中
- 实现设备详情 API (`GET /api/devices/:id`)
  - 权限检查(仅设备所有者可访问)
  - Secret 不暴露
- 实现设备更新 API (`PUT /api/devices/:id`)
  - 更新名称、位置
  - 启用/禁用设备
- 实现设备删除 API (`DELETE /api/devices/:id`)
  - 软删除或硬删除(待定)
  - 权限检查
- 实现 Secret 轮换 API (`POST /api/devices/:id/rotate-secret`)
  - 生成新 secret
  - 递增 secret_version
  - 返回新 secret (仅一次)

### 路径查询 (Path Query)
- 实现路径列表 API (`GET /api/paths`)
  - 返回用户可访问的所有 publish_path
  - 包含设备信息(名称、位置、状态)
  - 用于 Web UI 的流选择

### MediaMTX 集成完善
- 完善 `/mediamtx/auth` 回调处理
  - Publish 鉴权: 验证 device_secret
  - Read/Playback 鉴权: 验证 user_token
  - 详细的鉴权日志
  - 友好的错误消息

### 测试和文档
- 添加 API 集成测试
- 更新 README 的 API 端点文档
- 添加 Postman/curl 示例
- 创建 OpenAPI 规范 (可选)

**No breaking changes** - 所有新增功能,不影响现有代码。

## Impact

### Affected Specs
- **MODIFIED**: `api-server` - 实现占位 API,添加新端点

### Affected Code
- `internal/transport/http/handlers/`
  - `auth_handler.go` - 实现注册/登录逻辑
  - `device_handler.go` - 实现设备 CRUD
  - `path_handler.go` - 实现路径查询 (新建)
  - `mediamtx_auth_handler.go` - 完善鉴权逻辑
- `internal/domain/user/service.go` - 可能需要补充方法
- `internal/domain/device/service.go` - 可能需要补充方法
- `server/api/README.md` - 更新 API 文档
- 新增测试文件

### Dependencies
无新增外部依赖,使用现有技术栈。

## Success Criteria

- ✅ 用户可以通过 API 注册和登录
- ✅ 登录后获得有效的 JWT token
- ✅ 用户可以创建、查看、更新、删除自己的设备
- ✅ 设备创建时获得 device_secret
- ✅ 用户可以查询自己可访问的流路径
- ✅ MediaMTX 回调能正确鉴权 publish 和 read 操作
- ✅ 所有 API 有适当的错误处理和验证
- ✅ 集成测试覆盖主要场景
- ✅ API 文档完整且准确

## Out of Scope (Phase 4+)

- 密码重置功能
- OAuth/第三方登录
- 用户角色和权限系统
- 设备共享功能
- 实时推送通知
- 设备在线状态监控
