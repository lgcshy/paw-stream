// Auth Store

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import type { UserInfo } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string | null>(null)
  const user = ref<UserInfo | null>(null)
  const loading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)

  // Actions
  async function login(username: string, password: string) {
    loading.value = true
    try {
      const response = await authApi.login({ username, password })
      token.value = response.token
      user.value = response.user

      // Persist tokens to localStorage
      localStorage.setItem('auth_token', response.token)
      localStorage.setItem('refresh_token', response.refresh_token)

      return response
    } finally {
      loading.value = false
    }
  }

  async function register(username: string, password: string, nickname?: string) {
    loading.value = true
    try {
      // Register and get user info
      await authApi.register({ username, password, nickname })

      // Auto-login after registration
      const response = await authApi.login({ username, password })
      token.value = response.token
      user.value = response.user

      // Persist tokens to localStorage
      localStorage.setItem('auth_token', response.token)
      localStorage.setItem('refresh_token', response.refresh_token)

      return response
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('refresh_token')
  }

  async function loadToken() {
    const storedToken = localStorage.getItem('auth_token')
    if (!storedToken) {
      return
    }

    token.value = storedToken

    // Try to fetch current user info
    try {
      user.value = await authApi.getCurrentUser()
    } catch (error) {
      // Token invalid, clear it
      await logout()
      throw error
    }
  }

  async function updateAvatar(file: File) {
    const updatedUser = await authApi.uploadAvatar(file)
    user.value = updatedUser
    return updatedUser
  }

  async function checkAuth() {
    if (!token.value) {
      await loadToken()
    }
    return isAuthenticated.value
  }

  return {
    // State
    token,
    user,
    loading,
    // Getters
    isAuthenticated,
    // Actions
    login,
    register,
    logout,
    loadToken,
    updateAvatar,
    checkAuth,
  }
})
