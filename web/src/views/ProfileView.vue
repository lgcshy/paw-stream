<template>
  <Layout>
    <div class="profile-page">
      <!-- User Info Header -->
      <div class="user-header">
        <div class="avatar-wrapper" @click="triggerAvatarUpload">
          <img v-if="avatarURL" :src="avatarURL" class="avatar-img" :alt="$t('profile.defaultNickname')" />
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
          <h2>{{ user?.nickname || user?.username || $t('profile.defaultNickname') }}</h2>
          <p class="username">@{{ user?.username }}</p>
        </div>
      </div>

      <!-- Statistics -->
      <van-cell-group inset :title="$t('profile.statistics')">
        <van-cell :title="$t('profile.totalDevices')" :value="deviceCount" />
        <van-cell :title="$t('profile.onlineDevices')" :value="enabledDeviceCount" />
        <van-cell :title="$t('profile.registeredAt')" :value="formatDate(user?.created_at || '')" />
      </van-cell-group>

      <!-- Account Info -->
      <van-cell-group inset :title="$t('profile.accountInfo')">
        <van-cell :title="$t('profile.userId')" :value="user?.id" />
        <van-cell :title="$t('profile.accountStatus')" :value="user?.disabled ? $t('profile.statusDisabled') : $t('profile.statusNormal')" />
      </van-cell-group>

      <!-- Settings -->
      <van-cell-group inset :title="$t('profile.settings')">
        <van-cell :title="$t('profile.themeSetting')" is-link :value="themeLabel" @click="showThemePicker = true">
          <template #icon>
            <van-icon :name="themeIcon" color="#1989fa" />
          </template>
        </van-cell>
        <van-cell :title="$t('profile.languageSetting')" is-link :value="languageLabel" @click="showLanguagePicker = true">
          <template #icon>
            <van-icon name="globe-o" color="#1989fa" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Actions -->
      <van-cell-group inset :title="$t('profile.actions')">
        <van-cell :title="$t('admin.title')" is-link @click="router.push('/admin')">
          <template #icon>
            <van-icon name="chart-trending-o" color="#1989fa" />
          </template>
        </van-cell>
        <van-cell :title="$t('profile.about')" is-link @click="handleAbout">
          <template #icon>
            <van-icon name="info-o" color="#1989fa" />
          </template>
        </van-cell>
        <van-cell :title="$t('profile.logout')" is-link @click="handleLogout">
          <template #icon>
            <van-icon name="warning-o" color="#ee0a24" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Logout Confirm Dialog -->
      <ConfirmDialog
        v-model:show="showLogoutConfirm"
        :title="$t('profile.logoutConfirmTitle')"
        :message="$t('profile.logoutConfirmMessage')"
        :confirm-text="$t('profile.logoutButton')"
        :cancel-text="$t('common.cancel')"
        @confirm="confirmLogout"
      />

      <!-- About Dialog -->
      <van-dialog v-model:show="showAboutDialog" :title="$t('profile.aboutTitle')" :show-cancel-button="false">
        <div class="about-content">
          <div class="about-logo">🐾</div>
          <h3>PawStream</h3>
          <p class="version">Version 1.0.0</p>
          <p class="description">{{ $t('profile.aboutDescription') }}</p>
          <div class="tech-stack">
            <p><strong>{{ $t('profile.aboutFrontend') }}</strong></p>
            <p>Vue 3 + Vite + Vant + TypeScript</p>
            <p><strong>{{ $t('profile.aboutBackend') }}</strong></p>
            <p>Go + Fiber + SQLite + MediaMTX</p>
          </div>
        </div>
      </van-dialog>

      <!-- Theme Picker -->
      <van-action-sheet
        v-model:show="showThemePicker"
        :actions="themeActions"
        :cancel-text="$t('common.cancel')"
        close-on-click-action
        @select="onThemeSelect"
      />

      <!-- Language Picker -->
      <van-action-sheet
        v-model:show="showLanguagePicker"
        :actions="languageActions"
        :cancel-text="$t('common.cancel')"
        close-on-click-action
        @select="onLanguageSelect"
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
import { useI18n } from 'vue-i18n'
import { setLocale, getLocale } from '@/locales'

const { t } = useI18n()
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

const router = useRouter()
const authStore = useAuthStore()
const deviceStore = useDeviceStore()
const themeStore = useThemeStore()

const showLogoutConfirm = ref(false)
const showAboutDialog = ref(false)
const showThemePicker = ref(false)
const showLanguagePicker = ref(false)
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
      return t('profile.themeLight')
    case 'dark':
      return t('profile.themeDark')
    case 'auto':
      return t('profile.themeAuto')
    default:
      return t('profile.themeAuto')
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
    name: t('profile.themeAuto'),
    subname: themeStore.mode === 'auto' ? t('profile.themeCurrent') : '',
  },
  {
    name: t('profile.themeLight'),
    subname: themeStore.mode === 'light' ? t('profile.themeCurrent') : '',
  },
  {
    name: t('profile.themeDark'),
    subname: themeStore.mode === 'dark' ? t('profile.themeCurrent') : '',
  },
])

function onThemeSelect(action: ActionSheetAction) {
  let mode: 'light' | 'dark' | 'auto' = 'auto'

  if (action.name === t('profile.themeAuto')) {
    mode = 'auto'
  } else if (action.name === t('profile.themeLight')) {
    mode = 'light'
  } else if (action.name === t('profile.themeDark')) {
    mode = 'dark'
  }

  themeStore.setMode(mode)
  showSuccessToast(t('profile.themeSwitched', { name: action.name }))
}

// Language related
const languageLabel = computed(() => {
  return getLocale() === 'zh-CN' ? '中文' : 'English'
})

const languageActions = computed<ActionSheetAction[]>(() => [
  {
    name: '中文',
    subname: getLocale() === 'zh-CN' ? t('profile.themeCurrent') : '',
  },
  {
    name: 'English',
    subname: getLocale() === 'en' ? t('profile.themeCurrent') : '',
  },
])

function onLanguageSelect(action: ActionSheetAction) {
  const locale = action.name === '中文' ? 'zh-CN' : 'en'
  setLocale(locale)
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
    showFailToast(t('profile.avatarSizeLimit'))
    return
  }

  showLoadingToast({ message: t('profile.avatarUpload'), forbidClick: true })
  try {
    await authStore.updateAvatar(file)
    closeToast()
    showSuccessToast(t('profile.avatarSuccess'))
  } catch (error) {
    closeToast()
    showFailToast(t('profile.avatarFailed'))
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const locale = getLocale() === 'zh-CN' ? 'zh-CN' : 'en-US'
  return date.toLocaleDateString(locale, {
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
  showSuccessToast(t('profile.logoutSuccess'))
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
