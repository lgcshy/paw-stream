// Configuration Store
// Manages server configuration including MediaMTX URLs

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getConfig } from '@/api/config'
import type { ServerConfig } from '@/types/api'

export const useConfigStore = defineStore('config', () => {
  // State
  const config = ref<ServerConfig | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchTime = ref<number>(0)

  // Cache duration: 5 minutes
  const CACHE_DURATION = 5 * 60 * 1000

  // Computed
  const mediamtxWebRTCURL = computed(() => config.value?.mediamtx.webrtc_url || '')
  const mediamtxRTSPURL = computed(() => config.value?.mediamtx.rtsp_url || '')
  
  const isConfigCached = computed(() => {
    if (!config.value) return false
    return Date.now() - lastFetchTime.value < CACHE_DURATION
  })

  // Actions
  async function fetchConfig(force = false): Promise<ServerConfig> {
    // Return cached config if available and not forced
    if (!force && isConfigCached.value && config.value) {
      return config.value
    }

    loading.value = true
    error.value = null

    try {
      const serverConfig = await getConfig()
      config.value = serverConfig
      lastFetchTime.value = Date.now()
      return serverConfig
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取配置失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearConfig() {
    config.value = null
    lastFetchTime.value = 0
    error.value = null
  }

  return {
    // State
    config,
    loading,
    error,
    
    // Computed
    mediamtxWebRTCURL,
    mediamtxRTSPURL,
    isConfigCached,
    
    // Actions
    fetchConfig,
    clearConfig,
  }
})
