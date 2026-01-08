# Change: 为 Web UI 添加暗黑模式支持

## Why

用户在低光环境（夜间查看宠物直播）下使用 PawStream 时，当前的亮色主题会造成眼睛不适，影响用户体验。添加暗黑模式可以：
- 降低在暗环境中的视觉疲劳
- 节省移动设备电量（OLED 屏幕）
- 符合现代应用的用户体验标准
- 尊重用户的系统主题偏好

## What Changes

- 新增主题切换功能，支持亮色、暗色、自动三种模式
- 定义完整的暗色主题 CSS 变量和样式
- 在用户中心（Profile）页面添加主题切换入口
- 使用 Pinia store 管理主题状态和用户偏好
- 持久化主题偏好到 localStorage
- 自动检测系统主题偏好（`prefers-color-scheme`）
- 更新所有页面和组件以支持暗黑模式
- 为 Vant 组件配置暗色主题变量

## Impact

### 受影响的能力（Capabilities）
- **web-ui**: 添加主题系统和暗黑模式支持

### 受影响的代码
- **新增文件**:
  - `web/src/stores/theme.ts` - 主题 store
  - `web/src/composables/useTheme.ts` - 主题 composable（可选）
  
- **修改文件**:
  - `web/src/App.vue` - 添加主题初始化和根元素类切换
  - `web/src/views/ProfileView.vue` - 添加主题切换 UI
  - `web/src/assets/styles/variables.css` - 定义暗色主题变量
  - `web/src/main.ts` - 初始化主题
  - 所有视图和组件 - 确保样式兼容暗黑模式

### 用户体验变化
- 用户可以在"我的"页面选择主题偏好
- 支持自动跟随系统主题
- 主题切换即时生效，无需刷新页面
- 用户偏好在重新访问时保持

### 技术依赖
- 无新增外部依赖
- 利用现有的 Pinia、Vue 3 Composition API
- 使用 CSS 变量实现主题切换
