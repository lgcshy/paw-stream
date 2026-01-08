// Server Configuration API

import { apiClient } from './client'
import type { ServerConfig } from '@/types/api'

/**
 * Get server configuration
 * Returns MediaMTX URLs and other client configuration
 */
export async function getConfig(): Promise<ServerConfig> {
  return apiClient.get<ServerConfig>('/api/config')
}
