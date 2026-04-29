<template>
  <div class="register-page">
    <van-nav-bar :title="$t('register.title')" left-arrow @click-left="router.back()" />

    <div class="register-container">
      <div class="register-header">
        <h1>🐾 PawStream</h1>
        <p>{{ $t('register.subtitle') }}</p>
      </div>

      <van-form @submit="handleRegister">
        <!-- Username -->
        <van-cell-group inset>
          <van-field
            v-model="form.username"
            name="username"
            :label="$t('register.username')"
            :placeholder="$t('register.usernamePlaceholder')"
            clearable
            :rules="[
              { required: true, message: $t('register.usernameRequired') },
              { pattern: /^[a-zA-Z0-9_]{3,20}$/, message: $t('register.usernamePattern') },
            ]"
          />
        </van-cell-group>

        <!-- Nickname -->
        <van-cell-group inset>
          <van-field
            v-model="form.nickname"
            name="nickname"
            :label="$t('register.nickname')"
            :placeholder="$t('register.nicknamePlaceholder')"
            clearable
            :rules="[{ pattern: /^.{0,50}$/, message: $t('register.nicknameMaxLength') }]"
          />
        </van-cell-group>

        <!-- Password -->
        <van-cell-group inset>
          <van-field
            v-model="form.password"
            name="password"
            type="password"
            :label="$t('register.password')"
            :placeholder="$t('register.passwordPlaceholder')"
            clearable
            :rules="[
              { required: true, message: $t('register.passwordRequired') },
              { validator: validatePassword, message: $t('register.passwordPattern') },
            ]"
          />
        </van-cell-group>

        <!-- Confirm Password -->
        <van-cell-group inset>
          <van-field
            v-model="form.confirmPassword"
            name="confirmPassword"
            type="password"
            :label="$t('register.confirmPassword')"
            :placeholder="$t('register.confirmPasswordPlaceholder')"
            clearable
            :rules="[
              { required: true, message: $t('register.confirmPasswordRequired') },
              { validator: validateConfirmPassword, message: $t('register.confirmPasswordMismatch') },
            ]"
          />
        </van-cell-group>

        <!-- Password Strength Indicator -->
        <div v-if="form.password" class="password-strength">
          <div class="strength-label">{{ $t('register.passwordStrength') }}</div>
          <div class="strength-bar">
            <div
              class="strength-fill"
              :class="`strength-${passwordStrength}`"
              :style="{ width: passwordStrengthWidth }"
            ></div>
          </div>
          <div class="strength-text" :class="`strength-${passwordStrength}`">
            {{ passwordStrengthText }}
          </div>
        </div>

        <!-- Submit Button -->
        <div class="register-actions">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            {{ $t('register.submit') }}
          </van-button>
        </div>
      </van-form>

      <!-- Login Link -->
      <div class="register-footer">
        <span>{{ $t('register.hasAccount') }}</span>
        <router-link to="/login" class="login-link">{{ $t('register.goLogin') }}</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { showSuccessToast, showFailToast } from 'vant'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = ref({
  username: '',
  nickname: '',
  password: '',
  confirmPassword: '',
})

// Password validation
function validatePassword(value: string) {
  if (value.length < 8) return false
  // Must contain letters and numbers
  const hasLetter = /[a-zA-Z]/.test(value)
  const hasNumber = /\d/.test(value)
  return hasLetter && hasNumber
}

function validateConfirmPassword(value: string) {
  return value === form.value.password
}

// Password strength calculation
const passwordStrength = computed(() => {
  const pwd = form.value.password
  if (!pwd) return 'none'

  let score = 0

  // Length
  if (pwd.length >= 8) score++
  if (pwd.length >= 12) score++

  // Character types
  if (/[a-z]/.test(pwd)) score++
  if (/[A-Z]/.test(pwd)) score++
  if (/\d/.test(pwd)) score++
  if (/[^a-zA-Z0-9]/.test(pwd)) score++

  if (score <= 2) return 'weak'
  if (score <= 4) return 'medium'
  return 'strong'
})

const passwordStrengthWidth = computed(() => {
  const strength = passwordStrength.value
  if (strength === 'weak') return '33%'
  if (strength === 'medium') return '66%'
  if (strength === 'strong') return '100%'
  return '0%'
})

const passwordStrengthText = computed(() => {
  const strength = passwordStrength.value
  if (strength === 'weak') return t('register.strengthWeak')
  if (strength === 'medium') return t('register.strengthMedium')
  if (strength === 'strong') return t('register.strengthStrong')
  return ''
})

// Form submission
async function handleRegister() {
  loading.value = true
  try {
    await authStore.register(form.value.username, form.value.password, form.value.nickname || undefined)

    showSuccessToast(t('register.success'))

    // Navigate to stream list
    router.push('/streams')
  } catch (error: any) {
    console.error('Register failed:', error)
    showFailToast(error.message || t('register.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(to bottom, #f7f8fa 0%, #ffffff 100%);
}

/* Dark theme now handled by global CSS in variables.css */

.register-container {
  padding: 16px;
}

.register-header {
  text-align: center;
  padding: 48px 24px 40px;
  background: linear-gradient(135deg, rgba(30, 60, 114, 0.08) 0%, rgba(42, 82, 152, 0.08) 100%);
  border-radius: 20px;
  margin: 0 16px 24px;
  box-shadow: 0 4px 16px rgba(30, 60, 114, 0.1);
  border: 1px solid rgba(30, 60, 114, 0.05);
  position: relative;
  overflow: hidden;
}

.register-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at 30% 50%, rgba(56, 239, 125, 0.03) 0%, transparent 70%);
  pointer-events: none;
}

.register-header h1 {
  font-size: 40px;
  font-weight: 700;
  color: #1e3c72;
  margin: 0 0 10px 0;
  letter-spacing: 2px;
  position: relative;
  z-index: 1;
}

.register-header p {
  font-size: 16px;
  color: #546e7a;
  margin: 0;
  font-weight: 300;
  letter-spacing: 0.5px;
  position: relative;
  z-index: 1;
}

.van-cell-group {
  margin-bottom: 12px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(0, 0, 0, 0.02);
}

/* Dark theme handled by global CSS */

.password-strength {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin: 8px 0;
  background-color: white;
  border-radius: 8px;
}

.strength-label {
  font-size: 14px;
  color: #646566;
  margin-right: 12px;
  white-space: nowrap;
}

.strength-bar {
  flex: 1;
  height: 6px;
  background-color: #ebedf0;
  border-radius: 3px;
  overflow: hidden;
  margin-right: 12px;
}

/* Dark theme handled by global CSS */

.strength-fill {
  height: 100%;
  transition: width 0.3s, background-color 0.3s;
  border-radius: 3px;
}

.strength-fill.strength-weak {
  background-color: #ee0a24;
}

.strength-fill.strength-medium {
  background-color: #ff976a;
}

.strength-fill.strength-strong {
  background-color: #07c160;
}

.strength-text {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

.strength-text.strength-weak {
  color: #ee0a24;
}

.strength-text.strength-medium {
  color: #ff976a;
}

.strength-text.strength-strong {
  color: #07c160;
}

.register-actions {
  margin-top: 24px;
  padding: 0 16px;
}

.register-footer {
  text-align: center;
  margin-top: 32px;
  font-size: 14px;
  color: #546e7a;
}

.login-link {
  color: #1e3c72;
  text-decoration: none;
  margin-left: 4px;
  font-weight: 600;
  border-bottom: 1px solid transparent;
  transition: all 0.3s ease;
}

.login-link:hover {
  color: #2a5298;
  border-bottom-color: #2a5298;
}

/* Dark theme handled by global CSS */
</style>
