// Authentication API

import { apiClient } from './client'
import type { LoginRequest, LoginResponse, UserInfo } from '@/types/api'

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
}
