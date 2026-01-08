// Device Store

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { pathApi, deviceApi } from '@/api'
import type { PathInfo, DeviceInfo, CreateDeviceRequest, UpdateDeviceRequest } from '@/types/api'
import type { Stream } from '@/types/stream'

export const useDeviceStore = defineStore('device', () => {
  // State
  const paths = ref<PathInfo[]>([])
  const devices = ref<DeviceInfo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const streams = computed<Stream[]>(() => {
    return paths.value?.map((path) => ({
      id: path.publish_path,
      name: path.device_name,
      deviceId: path.device_id,
      location: path.device_location,
      status: 'online' as const, // All enabled devices are considered online
    })) ?? []
  })

  const enabledStreams = computed(() => {
    return streams.value?.filter((s) => s.status === 'online') ?? []
  })

  const enabledDevices = computed(() => {
    return devices.value?.filter((d) => !d.disabled) ?? []
  })

  // Actions - Paths (for streaming)
  async function fetchPaths() {
    loading.value = true
    error.value = null
    try {
      paths.value = await pathApi.listPaths()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取设备列表失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function refreshPaths() {
    return fetchPaths()
  }

  function getStreamById(id: string): Stream | undefined {
    return streams.value.find((s) => s.id === id)
  }

  // Actions - Device Management
  async function fetchDevices() {
    loading.value = true
    error.value = null
    try {
      devices.value = await deviceApi.listDevices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取设备列表失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getDevice(id: string) {
    loading.value = true
    error.value = null
    try {
      return await deviceApi.getDevice(id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取设备详情失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function getDeviceById(id: string): DeviceInfo | undefined {
    return devices.value?.find((d) => d.id === id)
  }

  async function createDevice(data: CreateDeviceRequest) {
    loading.value = true
    error.value = null
    try {
      const response = await deviceApi.createDevice(data)
      // Add to local state
      if (!devices.value) {
        devices.value = []
      }
      devices.value.push(response.device)
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建设备失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateDevice(id: string, data: UpdateDeviceRequest) {
    loading.value = true
    error.value = null
    try {
      const response = await deviceApi.updateDevice(id, data)
      // Update local state
      if (devices.value) {
        const index = devices.value.findIndex((d) => d.id === id)
        if (index !== -1) {
          devices.value[index] = response.device
        }
      }
      return response.device
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新设备失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteDevice(id: string) {
    loading.value = true
    error.value = null
    try {
      await deviceApi.deleteDevice(id)
      // Remove from local state
      if (devices.value) {
        devices.value = devices.value.filter((d) => d.id !== id)
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '删除设备失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function rotateSecret(id: string) {
    loading.value = true
    error.value = null
    try {
      return await deviceApi.rotateSecret(id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '轮换密钥失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function refreshDevices() {
    return fetchDevices()
  }

  return {
    // State
    paths,
    devices,
    loading,
    error,
    // Getters
    streams,
    enabledStreams,
    enabledDevices,
    // Actions - Paths
    fetchPaths,
    refreshPaths,
    getStreamById,
    // Actions - Device Management
    fetchDevices,
    getDevice,
    getDeviceById,
    createDevice,
    updateDevice,
    deleteDevice,
    rotateSecret,
    refreshDevices,
  }
})
