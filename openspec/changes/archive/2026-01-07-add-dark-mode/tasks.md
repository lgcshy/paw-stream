# 实施任务清单 - 暗黑模式

## 1. 基础架构

- [x] 1.1 创建主题 Pinia store (`stores/theme.ts`)
  - [x] 定义主题类型（light / dark / auto）
  - [x] 实现主题状态管理
  - [x] 实现 localStorage 持久化
  - [x] 实现系统主题检测（`prefers-color-scheme`）
  - [x] 实现主题切换逻辑

- [x] 1.2 定义 CSS 变量系统
  - [x] 创建或更新 `assets/styles/variables.css`
  - [x] 定义亮色主题变量（`:root`）
  - [x] 定义暗色主题变量（`.dark-theme`）
  - [x] 定义语义化颜色变量（背景、文字、边框等）

- [x] 1.3 集成主题系统到 App.vue
  - [x] 在 `App.vue` 中使用 ConfigProvider 配置主题
  - [x] 监听主题变化并动态更新 themeVars
  - [x] 监听系统主题变化（`matchMedia`）

## 2. UI 实现

- [x] 2.1 在 ProfileView 添加主题切换 UI
  - [x] 添加"主题设置"单元格
  - [x] 使用 `van-action-sheet` 实现选择器
  - [x] 显示当前主题模式
  - [x] 实现主题切换交互

- [x] 2.2 配置 Vant 暗色主题
  - [x] 在 `App.vue` 中使用 ConfigProvider 配置主题
  - [x] 测试 Vant 组件在暗色模式下的显示效果

## 3. 样式适配

- [x] 3.1 更新所有视图组件样式
  - [x] LoginView.vue - 保持深色背景（已适合暗环境）
  - [x] RegisterView.vue - 保持深色背景
  - [x] StreamListView.vue - 使用 CSS 变量适配
  - [x] StreamPlayerView.vue - 使用 CSS 变量适配
  - [x] DeviceListView.vue - 使用 CSS 变量适配
  - [x] DeviceFormView.vue - 使用 CSS 变量适配
  - [x] DeviceDetailView.vue - 使用 CSS 变量适配
  - [x] ProfileView.vue - 使用 CSS 变量适配

- [x] 3.2 更新布局组件样式
  - [x] Layout.vue - 通过 ConfigProvider 自动适配
  - [x] SecretDisplay.vue - 使用 CSS 变量适配
  - [x] ConfirmDialog.vue - 通过 ConfigProvider 自动适配

- [x] 3.3 特殊元素处理
  - [x] 处理渐变背景（使用 CSS 变量）
  - [x] 处理动画元素（保持兼容性）
  - [x] 确保视频播放器在暗色模式下正常显示

## 4. 测试与验证

- [x] 4.1 功能测试
  - [x] 测试主题切换功能（light / dark / auto）
  - [x] 测试主题持久化（刷新页面后保持）
  - [x] 测试自动模式跟随系统主题
  - [x] 测试系统主题变化时的自动切换

- [x] 4.2 视觉测试
  - [x] 检查所有页面在暗色模式下的可读性
  - [x] 检查对比度（优化后文字清晰可读）
  - [x] 检查渐变和动画效果
  - [x] 检查图标和按钮的可见性

- [x] 4.3 兼容性测试
  - [x] 使用 ConfigProvider 官方方案确保兼容性
  - [x] 修复主题切换时的白色闪烁问题
  - [x] 优化底部导航选中效果

## 5. 文档与完善

- [x] 5.1 更新文档
  - [x] 更新 `web/README.md` 添加暗黑模式说明
  - [x] 创建 IMPLEMENTATION.md 实施总结
  - [x] 创建 DARK_MODE_TEST_REPORT.md 测试清单

- [x] 5.2 代码审查与优化
  - [x] 确保类型定义完整（TypeScript 检查通过）
  - [x] 代码格式化（Prettier）
  - [x] 使用官方 ConfigProvider 方案

- [x] 5.3 用户体验优化
  - [x] 添加主题切换过渡动画（CSS transition）
  - [x] 优化首次加载时的主题闪烁（HTML 预加载脚本）
  - [x] 主题图标清晰（sunny / moon-o / setting-o）
  - [x] 优化底部导航选中效果（Material Design 颜色）

## 验收标准

- ✅ 用户可以在"我的"页面选择亮色、暗色、自动三种主题
- ✅ 选择"自动"时，应用跟随系统主题
- ✅ 主题偏好持久化，刷新页面后保持
- ✅ 所有页面和组件在暗色模式下视觉正常，文字清晰可读
- ✅ Vant 组件在暗色模式下样式正确
- ✅ 主题切换无明显延迟或闪烁
- ✅ 代码通过 TypeScript 类型检查
- ✅ 文档已更新
