// Device API

import { apiClient } from './client'
import type {
  DeviceInfo,
  CreateDeviceRequest,
  CreateDeviceResponse,
  UpdateDeviceRequest,
  UpdateDeviceResponse,
  RotateSecretResponse,
} from '@/types/api'

export const deviceApi = {
  /**
   * List all user devices
   */
  async listDevices(): Promise<DeviceInfo[]> {
    // API returns array directly, not wrapped in { devices: [...] }
    return apiClient.get<DeviceInfo[]>('/api/devices')
  },

  /**
   * Get device by ID
   */
  async getDevice(id: string): Promise<DeviceInfo> {
    return apiClient.get<DeviceInfo>(`/api/devices/${id}`)
  },

  /**
   * Create a new device
   */
  async createDevice(data: CreateDeviceRequest): Promise<CreateDeviceResponse> {
    return apiClient.post<CreateDeviceResponse>('/api/devices', data)
  },

  /**
   * Update device information
   */
  async updateDevice(id: string, data: UpdateDeviceRequest): Promise<UpdateDeviceResponse> {
    return apiClient.put<UpdateDeviceResponse>(`/api/devices/${id}`, data)
  },

  /**
   * Delete a device
   */
  async deleteDevice(id: string): Promise<void> {
    await apiClient.delete(`/api/devices/${id}`)
  },

  /**
   * Rotate device secret (generates new secret)
   */
  async rotateSecret(id: string): Promise<RotateSecretResponse> {
    return apiClient.post<RotateSecretResponse>(`/api/devices/${id}/rotate-secret`)
  },
}
