// PawStream API Client

import type { ApiError } from '@/types/api'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

export class ApiClient {
  private baseURL: string
  private refreshing: Promise<boolean> | null = null

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL
  }

  /**
   * Get the stored JWT token
   */
  private getToken(): string | null {
    return localStorage.getItem('auth_token')
  }

  /**
   * Attempt to refresh the access token using the stored refresh token
   */
  private async tryRefresh(): Promise<boolean> {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) return false

    try {
      const response = await fetch(`${this.baseURL}/api/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })

      if (!response.ok) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
        return false
      }

      const data = await response.json()
      localStorage.setItem('auth_token', data.token)
      localStorage.setItem('refresh_token', data.refresh_token)
      return true
    } catch {
      return false
    }
  }

  /**
   * Make an authenticated API request
   */
  async request<T>(
    endpoint: string,
    options: RequestInit = {},
    isRetry = false
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const token = this.getToken()

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {}),
    }

    // Add authorization header if token exists
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      })

      // Handle 401 with silent refresh
      if (response.status === 401 && !isRetry) {
        // Deduplicate concurrent refresh attempts
        if (!this.refreshing) {
          this.refreshing = this.tryRefresh().finally(() => { this.refreshing = null })
        }
        const refreshed = await this.refreshing
        if (refreshed) {
          return this.request<T>(endpoint, options, true)
        }
        // Refresh failed - redirect to login
        window.location.href = '/login'
        throw new Error('Authentication required')
      }

      if (response.status === 401) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
        window.location.href = '/login'
        throw new Error('Authentication required')
      }

      if (!response.ok) {
        // Try to parse error response
        const errorData: ApiError = await response.json().catch(() => ({
          error: 'unknown_error',
          message: `HTTP ${response.status}: ${response.statusText}`,
        }))
        throw new ApiClientError(errorData.message, response.status, errorData)
      }

      // Handle 204 No Content
      if (response.status === 204) {
        return null as T
      }

      return await response.json()
    } catch (error) {
      if (error instanceof ApiClientError) {
        throw error
      }

      // Network or other errors
      if (error instanceof Error) {
        throw new ApiClientError(
          error.message || '网络错误,请稍后重试',
          0,
          {
            error: 'network_error',
            message: error.message,
          }
        )
      }

      throw error
    }
  }

  /**
   * GET request
   */
  async get<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'GET',
    })
  }

  /**
   * POST request
   */
  async post<T>(endpoint: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  /**
   * PUT request
   */
  async put<T>(endpoint: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  /**
   * DELETE request
   */
  async delete<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'DELETE',
    })
  }
}

/**
 * Custom error class for API errors
 */
export class ApiClientError extends Error {
  statusCode: number
  errorData?: ApiError

  constructor(
    message: string,
    statusCode: number,
    errorData?: ApiError
  ) {
    super(message)
    this.name = 'ApiClientError'
    this.statusCode = statusCode
    this.errorData = errorData
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
