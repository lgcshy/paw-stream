<template>
  <Layout>
    <div class="profile-page">
      <!-- User Info Header -->
      <div class="user-header">
        <div class="avatar-wrapper" @click="triggerAvatarUpload">
          <img v-if="avatarURL" :src="avatarURL" class="avatar-img" alt="头像" />
          <van-icon v-else name="user-circle-o" size="60" color="#1989fa" />
          <div class="avatar-edit-hint">
            <van-icon name="photograph" size="14" color="#fff" />
          </div>
        </div>
        <input
          ref="fileInputRef"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          class="hidden-input"
          @change="onFileSelected"
        />
        <div class="user-info">
          <h2>{{ user?.nickname || user?.username || '用户' }}</h2>
          <p class="username">@{{ user?.username }}</p>
        </div>
      </div>

      <!-- Statistics -->
      <van-cell-group inset title="统计信息">
        <van-cell title="设备总数" :value="deviceCount" />
        <van-cell title="在线设备" :value="enabledDeviceCount" />
        <van-cell title="注册时间" :value="formatDate(user?.created_at || '')" />
      </van-cell-group>

      <!-- Account Info -->
      <van-cell-group inset title="账号信息">
        <van-cell title="用户ID" :value="user?.id" />
        <van-cell title="账号状态" :value="user?.disabled ? '已禁用' : '正常'" />
      </van-cell-group>

      <!-- Settings -->
      <van-cell-group inset title="设置">
        <van-cell title="主题设置" is-link :value="themeLabel" @click="showThemePicker = true">
          <template #icon>
            <van-icon :name="themeIcon" color="#1989fa" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Actions -->
      <van-cell-group inset title="操作">
        <van-cell title="关于 PawStream" is-link @click="handleAbout">
          <template #icon>
            <van-icon name="info-o" color="#1989fa" />
          </template>
        </van-cell>
        <van-cell title="退出登录" is-link @click="handleLogout">
          <template #icon>
            <van-icon name="warning-o" color="#ee0a24" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Logout Confirm Dialog -->
      <ConfirmDialog
        v-model:show="showLogoutConfirm"
        title="退出登录"
        message="确定要退出登录吗？"
        confirm-text="退出"
        cancel-text="取消"
        @confirm="confirmLogout"
      />

      <!-- About Dialog -->
      <van-dialog v-model:show="showAboutDialog" title="关于 PawStream" :show-cancel-button="false">
        <div class="about-content">
          <div class="about-logo">🐾</div>
          <h3>PawStream</h3>
          <p class="version">Version 1.0.0</p>
          <p class="description">基于 WebRTC 的宠物实时监控系统</p>
          <div class="tech-stack">
            <p><strong>前端技术:</strong></p>
            <p>Vue 3 + Vite + Vant + TypeScript</p>
            <p><strong>后端技术:</strong></p>
            <p>Go + Fiber + SQLite + MediaMTX</p>
          </div>
        </div>
      </van-dialog>

      <!-- Theme Picker -->
      <van-action-sheet
        v-model:show="showThemePicker"
        :actions="themeActions"
        cancel-text="取消"
        close-on-click-action
        @select="onThemeSelect"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDeviceStore } from '@/stores/device'
import { useThemeStore } from '@/stores/theme'
import { showSuccessToast, showFailToast, showLoadingToast, closeToast } from 'vant'
import type { ActionSheetAction } from 'vant'
import Layout from '@/components/Layout.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

const router = useRouter()
const authStore = useAuthStore()
const deviceStore = useDeviceStore()
const themeStore = useThemeStore()

const showLogoutConfirm = ref(false)
const showAboutDialog = ref(false)
const showThemePicker = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)

const avatarURL = computed(() => {
  if (!user.value?.avatar_url) return ''
  return `${API_BASE_URL}${user.value.avatar_url}?t=${user.value.updated_at}`
})

const user = computed(() => authStore.user)
const deviceCount = computed(() => deviceStore.devices?.length ?? 0)
const enabledDeviceCount = computed(() => deviceStore.enabledDevices?.length ?? 0)

// Theme related
const themeLabel = computed(() => {
  switch (themeStore.mode) {
    case 'light':
      return '浅色模式'
    case 'dark':
      return '深色模式'
    case 'auto':
      return '跟随系统'
    default:
      return '跟随系统'
  }
})

const themeIcon = computed(() => {
  if (themeStore.mode === 'auto') {
    return 'setting-o'
  }
  return themeStore.isDark ? 'moon-o' : 'sunny'
})

const themeActions = computed<ActionSheetAction[]>(() => [
  {
    name: '跟随系统',
    subname: themeStore.mode === 'auto' ? '当前' : '',
  },
  {
    name: '浅色模式',
    subname: themeStore.mode === 'light' ? '当前' : '',
  },
  {
    name: '深色模式',
    subname: themeStore.mode === 'dark' ? '当前' : '',
  },
])

function onThemeSelect(action: ActionSheetAction) {
  let mode: 'light' | 'dark' | 'auto' = 'auto'
  
  switch (action.name) {
    case '跟随系统':
      mode = 'auto'
      break
    case '浅色模式':
      mode = 'light'
      break
    case '深色模式':
      mode = 'dark'
      break
  }
  
  themeStore.setMode(mode)
  showSuccessToast(`已切换到${action.name}`)
}

function triggerAvatarUpload() {
  fileInputRef.value?.click()
}

async function onFileSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // Reset input so same file can be re-selected
  input.value = ''

  if (file.size > 2 * 1024 * 1024) {
    showFailToast('图片大小不能超过 2MB')
    return
  }

  showLoadingToast({ message: '上传中...', forbidClick: true })
  try {
    await authStore.updateAvatar(file)
    closeToast()
    showSuccessToast('头像更新成功')
  } catch (error) {
    closeToast()
    showFailToast('上传失败')
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function handleLogout() {
  showLogoutConfirm.value = true
}

async function confirmLogout() {
  await authStore.logout()
  showSuccessToast('已退出登录')
  router.push('/login')
}

function handleAbout() {
  showAboutDialog.value = true
}

// Load device count on mount
onMounted(async () => {
  if (deviceStore.devices.length === 0) {
    try {
      await deviceStore.fetchDevices()
    } catch (error) {
      // Silently fail - user can still see other info
      console.error('Failed to load devices:', error)
    }
  }
})
</script>

<style scoped>
.profile-page {
  width: 100%;
  /* 移除 min-height: 100%，让内容自然增长，保持与其他页面一致 */
}

.user-header {
  display: flex;
  align-items: center;
  padding: 48px 24px;
  background: var(--gradient-primary);
  color: white;
  margin-bottom: 16px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 24px var(--shadow-heavy);
}

.user-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 10% 20%, rgba(56, 239, 125, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 90% 80%, rgba(17, 153, 142, 0.1) 0%, transparent 50%);
  animation: headerFloat 15s ease-in-out infinite;
}

.user-header::after {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(128, 222, 234, 0.08) 0%, transparent 70%);
  animation: headerFloat2 20s ease-in-out infinite reverse;
}

@keyframes headerFloat {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(20px, -20px) scale(1.05);
  }
}

@keyframes headerFloat2 {
  0%, 100% {
    transform: translate(0, 0);
  }
  50% {
    transform: translate(-30px, 30px);
  }
}

.user-header > * {
  position: relative;
  z-index: 1;
}

.avatar-wrapper {
  position: relative;
  cursor: pointer;
  width: 64px;
  height: 64px;
  flex-shrink: 0;
}

.avatar-img {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.5);
}

.avatar-edit-hint {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hidden-input {
  display: none;
}

.user-info {
  margin-left: 20px;
  flex: 1;
}

.user-info h2 {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 3px 12px rgba(0, 0, 0, 0.2);
  letter-spacing: 1px;
  color: #ffffff;
}

.username {
  font-size: 15px;
  opacity: 0.85;
  margin: 0;
  font-weight: 300;
  color: #80deea;
  letter-spacing: 0.5px;
}

.van-cell-group {
  margin-bottom: 16px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 12px var(--shadow-medium);
  border: 1px solid var(--border-light);
}

:deep(.van-cell__left-icon) {
  margin-right: 12px;
  font-size: 20px;
}

.about-content {
  padding: 20px;
  text-align: center;
}

.about-logo {
  font-size: 48px;
  margin-bottom: 12px;
}

.about-content h3 {
  font-size: 20px;
  font-weight: bold;
  margin: 8px 0;
}

.version {
  font-size: 12px;
  color: var(--text-disabled);
  margin: 4px 0 12px 0;
}

.description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 12px 0;
}

.tech-stack {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
  text-align: left;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.tech-stack p {
  margin: 4px 0;
}

.tech-stack strong {
  color: var(--text-primary);
}
</style>
