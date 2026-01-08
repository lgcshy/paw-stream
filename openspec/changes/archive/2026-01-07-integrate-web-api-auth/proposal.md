# Change: Integrate Web UI with API Authentication and Live Streaming

## Why

Phase 3 已完成 API 服务器的用户认证和设备管理功能,但 Web UI 仍使用占位符实现。为了完成 Phase 4 的目标,需要:

1. **集成真实的 API 认证** - 替换登录页面的占位符,实现真实的用户注册/登录
2. **连接设备管理 API** - 从后端获取真实的设备和流路径数据
3. **实现 WebRTC 播放器** - 集成 MediaMTX 的 WebRTC 流播放功能
4. **状态管理** - 使用 Pinia 管理用户登录态和设备数据
5. **路由守卫** - 实现认证保护,未登录用户重定向到登录页

这将使 PawStream 成为一个端到端可用的宠物监控系统,用户可以通过手机浏览器实时查看宠物摄像头画面。

## What Changes

### 1. API 集成层
- 创建 API 客户端 (`src/api/client.ts`)
  - 封装 fetch 请求
  - 自动添加 JWT token
  - 统一错误处理
- 创建 API 服务模块
  - `src/api/auth.ts` - 用户认证 API
  - `src/api/device.ts` - 设备管理 API
  - `src/api/path.ts` - 路径查询 API
- 定义 TypeScript 类型 (`src/types/api.ts`)

### 2. 状态管理 (Pinia)
- 创建 auth store (`src/stores/auth.ts`)
  - 管理用户登录状态
  - 存储 JWT token (localStorage)
  - 提供登录/登出方法
- 创建 device store (`src/stores/device.ts`)
  - 管理设备列表
  - 管理可访问路径
  - 提供刷新方法

### 3. 登录页面改造
- 替换占位符实现
- 调用真实的 `/api/login` API
- 保存 JWT token
- 登录成功后跳转到流列表
- 错误处理和用户提示

### 4. 流列表页面改造
- 从 `/api/paths` 获取真实数据
- 显示设备名称、位置、在线状态
- 点击跳转到播放器页面
- 空状态提示

### 5. WebRTC 播放器集成
- 集成 MediaMTX WebRTC 客户端
- 实现流播放逻辑
  - 从路径构建 WebRTC URL
  - 附加 JWT token 用于鉴权
  - 处理连接状态
- 播放控制 (播放/暂停/重连)
- 错误处理和重试机制

### 6. 路由守卫
- 实现全局路由守卫
- 检查登录状态
- 未登录重定向到 `/login`
- 已登录访问 `/login` 重定向到 `/streams`

### 7. 环境配置
- 添加 `.env` 配置
  - `VITE_API_BASE_URL` - API 服务器地址
  - `VITE_MEDIAMTX_WEBRTC_URL` - MediaMTX WebRTC 地址
- 开发和生产环境配置

### 8. 用户体验优化
- 加载状态提示
- 错误提示 (Toast)
- 网络请求失败重试
- Token 过期自动登出

**No breaking changes** - 所有改动都是增强现有占位符功能。

## Impact

### Affected Specs
- **MODIFIED**: `web-ui` - 实现占位功能,添加 API 集成和 WebRTC 播放

### Affected Code
- `web/src/api/` (新建目录)
  - `client.ts`, `auth.ts`, `device.ts`, `path.ts`
- `web/src/stores/` (新建目录)
  - `auth.ts`, `device.ts`
- `web/src/types/`
  - `api.ts` (新建)
  - `stream.ts` (修改)
- `web/src/views/`
  - `LoginView.vue` (修改)
  - `StreamListView.vue` (修改)
  - `StreamPlayerView.vue` (修改)
- `web/src/router/index.ts` (修改 - 添加路由守卫)
- `web/.env.development` (新建)
- `web/.env.production` (新建)
- `web/package.json` (可能新增依赖)

### Dependencies
可能新增:
- `pinia` - Vue 3 状态管理 (如果未安装)
- MediaMTX WebRTC 客户端库 (或直接使用原生 WebRTC API)

## Success Criteria

- ✅ 用户可以通过 Web UI 注册和登录
- ✅ 登录后 JWT token 正确保存和使用
- ✅ 流列表页面显示真实的设备数据
- ✅ 点击设备可以打开 WebRTC 播放器
- ✅ WebRTC 播放器成功连接 MediaMTX 并播放视频
- ✅ JWT token 正确附加到 WebRTC 连接用于鉴权
- ✅ 未登录用户无法访问受保护页面
- ✅ Token 过期后自动登出
- ✅ 所有 API 错误有友好的用户提示
- ✅ 移动端体验流畅

## Out of Scope (Phase 5+)

- 用户注册页面 (当前仅登录)
- 设备管理界面 (添加/编辑/删除设备)
- 录像回放功能
- 多设备同时播放
- 推送通知
- 设备在线状态实时更新
- 播放器高级控制 (画质切换、全屏等)
