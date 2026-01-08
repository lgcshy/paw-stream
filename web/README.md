# PawStream Web UI

PawStream 的移动优先 Web 前端应用，基于 Vue 3 + TypeScript + Vite 7 + Vant 4 构建。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript 超集
- **Vite 7** - 下一代前端构建工具
- **Vant 4** - 移动端 Vue 组件库（Vue 3 兼容）
- **Vue Router 4** - Vue.js 官方路由
- **Pinia** - Vue 状态管理库
- **Prettier** - 代码格式化工具

## 项目结构

```
web/
├── src/
│   ├── api/             # API 客户端层 ✨
│   │   ├── client.ts    # HTTP 客户端封装
│   │   ├── auth.ts      # 认证 API
│   │   ├── device.ts    # 设备管理 API 🆕
│   │   ├── path.ts      # 路径查询 API
│   │   └── index.ts     # API 导出
│   ├── assets/          # 静态资源
│   │   └── styles/      # 样式文件
│   │       └── variables.css  # 主题 CSS 变量 🌙
│   ├── components/      # 可复用组件
│   │   ├── Layout.vue         # 布局组件 + 底部导航 ✨
│   │   ├── SecretDisplay.vue  # Secret 显示组件 🆕
│   │   └── ConfirmDialog.vue  # 确认对话框 🆕
│   ├── router/          # 路由配置
│   │   └── index.ts     # 路由定义 + 认证守卫 ✨
│   ├── stores/          # Pinia 状态管理 ✨
│   │   ├── auth.ts      # 用户认证状态 + 注册 ✨
│   │   ├── device.ts    # 设备数据状态 + CRUD ✨
│   │   └── theme.ts     # 主题状态管理 🌙
│   ├── types/           # TypeScript 类型定义
│   │   ├── api.ts       # API 类型定义 ✨
│   │   └── stream.ts    # 视频流相关类型
│   ├── utils/           # 工具函数 ✨
│   │   └── webrtc.ts    # WebRTC 播放器工具
│   ├── views/           # 页面组件
│   │   ├── HomeView.vue          # 首页
│   │   ├── LoginView.vue         # 登录页 ✅
│   │   ├── RegisterView.vue      # 注册页 🆕
│   │   ├── StreamListView.vue    # 直播流列表 ✅
│   │   ├── StreamPlayerView.vue  # 播放器页面 ✅
│   │   ├── DeviceListView.vue    # 设备列表 🆕
│   │   ├── DeviceFormView.vue    # 设备表单 🆕
│   │   ├── DeviceDetailView.vue  # 设备详情 🆕
│   │   └── ProfileView.vue       # 用户中心 🆕
│   ├── App.vue          # 根组件
│   └── main.ts          # 应用入口
├── public/              # 公共静态资源
├── .env.development     # 开发环境配置 ✨
├── .env.production      # 生产环境配置 ✨
├── index.html           # HTML 入口文件
├── vite.config.ts       # Vite 配置
├── tsconfig.json        # TypeScript 配置
├── .prettierrc          # Prettier 配置
├── TEST_PLAN.md         # 测试计划 ✨
└── package.json         # 项目依赖
```

✨ = Phase 4 新增/更新
✅ = Phase 4 完全实现
🆕 = Phase 5 新增
🌙 = 暗黑模式支持

## 快速开始

### 前置要求

1. **API 服务器**: 确保 PawStream API 服务器正在运行
   ```bash
   cd ../server/api
   ./bin/api
   # 默认运行在 http://localhost:3000
   ```

2. **MediaMTX** (可选,仅用于测试实际流播放):
   ```bash
   cd ../server/mediamtx
   docker-compose up -d
   # 默认运行在 http://localhost:8889
   ```

## 路由结构

```
/login              - 登录页（无需认证）
/register           - 注册页（无需认证）🆕
/streams            - 直播流列表（需要认证，显示底部导航）
/stream/:id         - 播放器页面（需要认证）
/devices            - 设备列表（需要认证，显示底部导航）🆕
/devices/new        - 创建设备（需要认证）🆕
/devices/:id        - 设备详情（需要认证）🆕
/devices/:id/edit   - 编辑设备（需要认证）🆕
/profile            - 用户中心（需要认证，显示底部导航）🆕
```

### 环境配置

创建 `.env.local` 或使用默认的 `.env.development`:

```bash
# API Server URL
VITE_API_BASE_URL=http://localhost:3000

# MediaMTX WebRTC URL
VITE_MEDIAMTX_WEBRTC_URL=http://localhost:8889
```

### 安装依赖

```bash
npm install
```

### 开发模式

启动开发服务器（支持热更新）：

```bash
npm run dev
```

访问 http://localhost:5173

### 类型检查

运行 TypeScript 类型检查：

```bash
npm run type-check
```

### 构建生产版本

```bash
npm run build
```

构建产物将生成在 `dist/` 目录。

### 预览生产构建

```bash
npm run preview
```

### 代码格式化

使用 Prettier 格式化代码：

```bash
npm run format
```

## 开发约定

### 代码风格

- 使用 TypeScript 严格模式进行类型检查
- 遵循 Prettier 默认配置进行代码格式化
- Vue 组件使用 `<script setup lang="ts">` 语法
- 使用组合式 API (Composition API)

### 组件开发

- **页面组件**：放在 `src/views/` 目录，使用 `*View.vue` 命名
- **可复用组件**：放在 `src/components/` 目录
- **类型定义**：放在 `src/types/` 目录，使用 `.ts` 扩展名

### 路径别名

项目配置了 `@` 别名指向 `src/` 目录：

```typescript
import Layout from '@/components/Layout.vue'
import type { Stream } from '@/types/stream'
```

### Vant 组件使用

项目已配置 Vant 4 组件自动导入，无需手动全局注册：

```vue
<script setup lang="ts">
// 直接使用，会自动导入
import { Button, NavBar } from 'vant'
</script>

<template>
  <Button type="primary">按钮</Button>
  <NavBar title="标题" />
</template>
```

## 当前实现状态

### ✅ 已完成

- [x] Vue 3 + TypeScript + Vite 7 项目初始化
- [x] Vant 4 UI 组件库集成和自动导入配置
- [x] Vue Router 4 路由配置
- [x] Pinia 状态管理准备
- [x] TypeScript 严格模式配置
- [x] Prettier 代码格式化配置
- [x] 移动端响应式布局组件
- [x] 登录页面（占位实现）
- [x] 首页
- [x] 视频流列表页面（模拟数据）
- [x] 视频流播放器页面（WebRTC 占位）

### ✅ Phase 4 已完成

- [x] 实际的用户认证和 API 集成
- [x] 从后端 API 获取真实的视频流列表
- [x] WebRTC 播放器集成 (MediaMTX WHEP)
- [x] 路由守卫和 Token 管理
- [x] Pinia 状态管理 (auth, device)
- [x] 错误处理和用户反馈优化
- [x] TypeScript 严格模式验证
- [x] 生产构建成功

### ✅ Phase 5 已完成

- [x] 用户注册界面
- [x] 设备管理界面（完整 CRUD）
- [x] 设备详情和 Secret 管理
- [x] 用户中心（个人信息和统计）
- [x] 底部导航栏
- [x] 可复用组件（SecretDisplay, ConfirmDialog）

### 🌙 暗黑模式（Phase 6）

- [x] 主题系统核心架构（Pinia store + CSS 变量）
- [x] 三种主题模式：浅色 / 深色 / 自动跟随系统
- [x] 主题切换 UI（用户中心）
- [x] 所有页面和组件暗黑模式适配
- [x] Vant 组件暗黑主题配置
- [x] 主题持久化（localStorage）
- [x] 系统主题检测和自动切换

### 🚧 待实现（Phase 7+）

- [ ] 视频流高级控制功能（全屏、画质切换等）
- [ ] 录像回放功能
- [ ] 推送通知
- [ ] 设备分组
- [ ] 设备共享
- [ ] PWA 支持
- [ ] 单元测试和 E2E 测试

## API 集成说明

### 认证流程
1. 用户在登录页面输入用户名和密码
2. 调用 `POST /api/login` 进行认证
3. 成功后保存 JWT token 到 localStorage
4. 所有 API 请求自动附加 Authorization header
5. Token 过期时自动登出并重定向到登录页

### 设备数据流
1. 登录后自动调用 `GET /api/paths` 获取设备列表
2. 数据存储在 Pinia device store
3. 支持下拉刷新重新加载数据

### WebRTC 播放流程
1. 用户点击设备进入播放器页面
2. 创建 RTCPeerConnection
3. 向 MediaMTX 发送 WHEP 请求: `POST /{path}/whep`
4. 附加 JWT token 用于鉴权
5. 接收视频流并显示

详细测试步骤请参考 [`TEST_PLAN.md`](./TEST_PLAN.md)

## 相关文档

### 项目文档
- [暗黑模式测试报告](./DARK_MODE_TEST_REPORT.md) 🌙
- [Phase 5 实施文档](./PHASE5_IMPLEMENTATION.md) 🆕
- [Phase 4 实施文档](../openspec/changes/archive/2026-01-07-integrate-web-api-auth/IMPLEMENTATION.md)
- [测试计划](./TEST_PLAN.md)
- [API 服务器文档](../server/api/README.md)
- [API 使用指南](../server/api/API_GUIDE.md)

### 技术文档
- [OpenSpec 项目约定](../openspec/project.md)
- [Vue 3 文档](https://vuejs.org/)
- [Vite 文档](https://vitejs.dev/)
- [Vant 4 文档](https://vant-ui.github.io/vant/#/zh-CN)
- [TypeScript 文档](https://www.typescriptlang.org/)

## 🌙 暗黑模式

PawStream Web UI 完整支持暗黑模式，提供更舒适的夜间使用体验。

### 功能特性

- **三种主题模式**：
  - 🌞 **浅色模式** - 经典的明亮主题
  - 🌙 **深色模式** - 护眼的深色主题
  - ⚙️ **自动模式** - 跟随系统主题设置

- **智能切换**：
  - 自动检测系统主题偏好（`prefers-color-scheme`）
  - 系统主题变化时实时响应
  - 主题偏好持久化保存

- **完整适配**：
  - 所有页面和组件支持暗黑模式
  - Vant UI 组件暗色主题配置
  - CSS 变量驱动，易于定制
  - 平滑过渡动画，无闪烁

### 使用方法

1. 登录后进入"我的"页面
2. 点击"主题设置"
3. 选择你喜欢的主题模式：
   - **跟随系统** - 自动跟随系统设置
   - **浅色模式** - 固定使用浅色主题
   - **深色模式** - 固定使用深色主题

主题偏好会自动保存，下次访问时保持你的选择。

### 技术实现

- **状态管理**: Pinia store (`stores/theme.ts`)
- **样式系统**: CSS 变量 (`assets/styles/variables.css`)
- **持久化**: localStorage
- **系统检测**: `matchMedia('prefers-color-scheme: dark')`

详细测试报告请参考 [`DARK_MODE_TEST_REPORT.md`](./DARK_MODE_TEST_REPORT.md)

## 移动端适配

- 所有页面针对移动端优化（320px - 768px）
- 使用 Vant 4 移动端组件确保触控友好
- 支持移动端浏览器（iOS Safari、Android Chrome）
- 视口配置已在 `index.html` 中设置
- 深色模式优化移动设备电池使用（OLED 屏幕）

## 故障排除

### Node.js 版本警告

如果看到 Node.js 版本警告，建议升级到 Node.js 20.19+ 或 22.12+。项目在 Node.js 20.17 上也可以正常运行，但可能会有警告。

### TypeScript 错误

运行 `npm run type-check` 检查类型错误。确保所有组件 props、emits 都有正确的类型定义。

### 样式问题

Vant 样式已在 `src/main.ts` 中全局导入。如果样式缺失，检查 Vite 配置中的 Vant 解析器是否正确配置。

## 许可证

本项目为 PawStream 的一部分，采用相同的许可证。
