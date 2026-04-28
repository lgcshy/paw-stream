// PawStream API Type Definitions

// ============= Auth API =============

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  nickname?: string
}

export interface UserInfo {
  id: string
  username: string
  nickname: string
  avatar_url?: string
  disabled: boolean
  created_at: string
  updated_at: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user: UserInfo
}

// ============= Device API =============

export interface DeviceInfo {
  id: string
  name: string
  location: string
  publish_path: string
  disabled: boolean
  created_at: string
  updated_at: string
}

export interface CreateDeviceRequest {
  name: string
  location?: string
}

export interface CreateDeviceResponse {
  device: DeviceInfo
  secret: string // Only returned once!
}

export interface UpdateDeviceRequest {
  name?: string
  location?: string
  disabled?: boolean
}

export interface UpdateDeviceResponse {
  device: DeviceInfo
}

export interface RotateSecretResponse {
  device_id: string
  new_secret: string // Only returned once!
}

// ============= Path API =============

export interface PathInfo {
  publish_path: string
  device_id: string
  device_name: string
  device_location: string
}

// ============= Error Response =============

export interface ApiError {
  error: string
  message: string
  request_id?: string
  details?: Record<string, unknown>
}

// ============= Server Config API =============

export interface ServerConfig {
  mediamtx: {
    webrtc_url: string
    rtsp_url: string
  }
}

// ============= API Response Wrapper =============

export interface ApiResponse<T> {
  data?: T
  error?: ApiError
}
