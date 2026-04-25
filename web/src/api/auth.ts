// Authentication API

import { apiClient } from './client'
import type { LoginRequest, LoginResponse, UserInfo } from '@/types/api'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

export const authApi = {
  /**
   * User login
   */
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    return apiClient.post<LoginResponse>('/api/login', credentials)
  },

  /**
   * Get current user info
   */
  async getCurrentUser(): Promise<UserInfo> {
    return apiClient.get<UserInfo>('/api/me')
  },

  /**
   * User registration (optional for Phase 4)
   */
  async register(data: { username: string; password: string; nickname?: string }): Promise<UserInfo> {
    return apiClient.post<UserInfo>('/api/register', data)
  },

  /**
   * Upload avatar
   */
  async uploadAvatar(file: File): Promise<UserInfo> {
    const formData = new FormData()
    formData.append('avatar', file)

    const token = localStorage.getItem('auth_token')
    const response = await fetch(`${API_BASE_URL}/api/me/avatar`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: formData,
    })

    if (!response.ok) {
      const err = await response.json().catch(() => ({ message: 'Upload failed' }))
      throw new Error(err.message || 'Upload failed')
    }

    return response.json()
  },
}
