# Change: Add Device Management UI and User Registration

## Why

Phase 4 已完成 Web UI 与 API 的集成和 WebRTC 播放功能,但用户体验还不完整:

1. **用户注册** - 目前只能通过 API 创建用户,没有注册界面
2. **设备管理** - 用户只能查看设备,无法在 Web UI 中添加、编辑或删除设备
3. **用户体验** - 缺少设备管理的完整流程,需要手动调用 API

为了让 PawStream 成为一个完整的、用户友好的应用,需要实现:
- 用户注册页面 - 让新用户可以自助注册
- 设备管理界面 - 完整的设备 CRUD 操作
- 用户信息管理 - 查看和编辑个人信息
- 改进的导航和布局 - 更好的用户体验

这将使 PawStream 达到 MVP 可用状态,用户无需直接调用 API 即可完成所有操作。

## What Changes

### 1. 用户注册页面
- 创建注册页面 (`RegisterView.vue`)
  - 用户名、密码、昵称输入
  - 密码强度提示
  - 表单验证
  - 注册成功后自动登录
- 登录页面添加"注册账号"链接

### 2. 设备管理界面
- 设备列表页面增强
  - 添加"新增设备"按钮
  - 设备项支持长按操作（编辑/删除）
  - 显示设备状态和最后更新时间
- 创建设备表单页面 (`DeviceFormView.vue`)
  - 设备名称、位置输入
  - 创建成功后显示 device_secret (仅一次)
  - Secret 复制功能
- 设备编辑功能
  - 修改名称、位置
  - 启用/禁用设备
- 设备删除确认
  - 二次确认对话框
  - 警告信息

### 3. 设备详情页面
- 创建设备详情页面 (`DeviceDetailView.vue`)
  - 显示设备完整信息
  - 显示 publish_path
  - 启用/禁用开关
  - 轮换 Secret 功能
  - 删除设备按钮

### 4. 用户信息页面
- 创建用户中心页面 (`ProfileView.vue`)
  - 显示用户信息
  - 登出按钮
  - 统计信息（设备数量等）
- Layout 组件改进
  - 添加底部导航栏（流列表、设备管理、个人中心）
  - 顶部栏添加用户头像/昵称

### 5. API 扩展
- 更新 device API 模块
  - createDevice()
  - updateDevice()
  - deleteDevice()
  - rotateSecret()
- 添加 register API

### 6. 状态管理增强
- device store 扩展
  - createDevice action
  - updateDevice action
  - deleteDevice action
- auth store 扩展
  - register action

### 7. 用户体验改进
- 添加确认对话框组件
- 添加 Secret 显示和复制组件
- 改进错误提示
- 添加操作成功反馈
- 优化加载状态

### 8. 路由扩展
- `/register` - 注册页面
- `/devices/new` - 创建设备
- `/devices/:id/edit` - 编辑设备
- `/devices/:id` - 设备详情
- `/profile` - 个人中心

**No breaking changes** - 所有新增功能,不影响现有代码。

## Impact

### Affected Specs
- **MODIFIED**: `web-ui` - 添加设备管理和用户注册功能

### Affected Code
- `web/src/api/`
  - `device.ts` (新增方法)
  - `auth.ts` (新增 register)
- `web/src/stores/`
  - `device.ts` (新增 actions)
  - `auth.ts` (新增 register)
- `web/src/views/`
  - `RegisterView.vue` (新建)
  - `DeviceFormView.vue` (新建)
  - `DeviceDetailView.vue` (新建)
  - `ProfileView.vue` (新建)
  - `StreamListView.vue` (修改 - 改进 UI)
- `web/src/components/`
  - `Layout.vue` (修改 - 添加底部导航)
  - `SecretDisplay.vue` (新建 - Secret 显示组件)
  - `ConfirmDialog.vue` (新建 - 确认对话框)
- `web/src/router/index.ts` (添加新路由)

### Dependencies
无新增外部依赖,使用现有 Vant 组件。

## Success Criteria

- ✅ 用户可以通过 Web UI 注册账号
- ✅ 新用户注册后自动登录
- ✅ 用户可以在 Web UI 中创建设备
- ✅ 创建设备时显示 device_secret 并支持复制
- ✅ 用户可以编辑设备信息（名称、位置、启用/禁用）
- ✅ 用户可以删除设备（带确认）
- ✅ 用户可以轮换设备 secret
- ✅ 用户可以查看个人信息和统计
- ✅ 用户可以登出
- ✅ 底部导航栏提供便捷访问
- ✅ 所有操作有友好的反馈和确认
- ✅ 移动端体验优秀
- ✅ TypeScript 类型检查通过
- ✅ 生产构建成功

## Out of Scope (Future)

- 密码修改功能
- 头像上传
- 用户角色和权限
- 设备分组
- 设备共享
- 数据统计图表
- 设备在线状态实时更新
- 推送通知设置
