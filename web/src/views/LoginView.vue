<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Form, Field, Button, CellGroup, showFailToast, showSuccessToast } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { ApiClientError } from '@/api'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

const onSubmit = async () => {
  if (!username.value || !password.value) {
    showFailToast(t('login.inputRequired'))
    return
  }

  loading.value = true
  try {
    await authStore.login(username.value, password.value)
    showSuccessToast(t('login.success'))

    // Redirect to original destination or streams page
    const redirect = (route.query.redirect as string) || '/streams'
    router.push(redirect)
  } catch (error) {
    if (error instanceof ApiClientError) {
      if (error.statusCode === 401) {
        showFailToast(t('login.invalidCredentials'))
      } else if (error.statusCode === 403) {
        showFailToast(t('login.accountDisabled'))
      } else {
        showFailToast(error.message || t('login.failed'))
      }
    } else {
      showFailToast(t('common.networkError'))
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-view">
    <div class="login-container">
      <div class="logo">
        <h1>🐾 PawStream</h1>
        <p>{{ $t('app.subtitle') }}</p>
      </div>

      <Form @submit="onSubmit">
        <CellGroup inset>
          <Field
            v-model="username"
            name="username"
            :label="$t('login.username')"
            :placeholder="$t('login.usernamePlaceholder')"
            :rules="[{ required: true, message: $t('login.usernameRequired') }]"
          />
          <Field
            v-model="password"
            type="password"
            name="password"
            :label="$t('login.password')"
            :placeholder="$t('login.passwordPlaceholder')"
            :rules="[{ required: true, message: $t('login.passwordRequired') }]"
          />
        </CellGroup>
        <div class="button-container">
          <Button round block type="primary" native-type="submit" :loading="loading" :disabled="loading">
            {{ $t('login.submit') }}
          </Button>
        </div>
      </Form>

      <!-- Register Link -->
      <div class="register-link">
        <span>{{ $t('login.noAccount') }}</span>
        <router-link to="/register" class="link-text">{{ $t('login.goRegister') }}</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
  position: relative;
  padding: 20px;
  overflow: hidden;
}

.login-view::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background:
    radial-gradient(circle at 20% 50%, rgba(56, 239, 125, 0.08) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(17, 153, 142, 0.08) 0%, transparent 50%);
  animation: float 20s ease-in-out infinite;
}

.login-view::after {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.05) 0%, transparent 70%);
  animation: float2 25s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(30px, -30px) scale(1.1);
  }
}

@keyframes float2 {
  0%, 100% {
    transform: translate(0, 0);
  }
  50% {
    transform: translate(-40px, 40px);
  }
}

.login-container {
  width: 100%;
  max-width: 400px;
  position: relative;
  z-index: 1;
}

.logo {
  text-align: center;
  margin-bottom: 50px;
  color: white;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.logo h1 {
  font-size: 52px;
  margin-bottom: 12px;
  font-weight: 700;
  letter-spacing: 3px;
  background: linear-gradient(135deg, #ffffff 0%, #e0f7fa 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: drop-shadow(0 2px 8px rgba(255, 255, 255, 0.3));
}

.logo p {
  font-size: 16px;
  opacity: 0.9;
  font-weight: 300;
  letter-spacing: 1px;
  color: #b3e5fc;
}

/* Form styling for better visibility */
:deep(.van-cell-group) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

:deep(.van-cell) {
  background: transparent;
}

/* Dark theme is now handled by global CSS in variables.css */

.button-container {
  margin-top: 20px;
  padding: 0 16px;
}

.register-link {
  text-align: center;
  margin-top: 24px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
}

.link-text {
  color: #80deea;
  text-decoration: none;
  margin-left: 4px;
  font-weight: 500;
  border-bottom: 1px solid transparent;
  transition: all 0.3s ease;
}

.link-text:hover {
  color: #b3e5fc;
  border-bottom-color: #b3e5fc;
}
</style>
