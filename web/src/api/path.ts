// Path Query API

import { apiClient } from './client'
import type { PathInfo } from '@/types/api'

export const pathApi = {
  /**
   * List accessible stream paths
   */
  async listPaths(): Promise<PathInfo[]> {
    return apiClient.get<PathInfo[]>('/api/paths')
  },
}
