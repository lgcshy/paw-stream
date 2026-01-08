/**
 * Stream related type definitions
 */

// Stream data structure (mapped from API PathInfo)
export interface Stream {
  id: string // publish_path
  name: string // device_name
  deviceId: string // device_id
  location: string // device_location
  status: 'online' | 'offline' // always 'online' for now (enabled devices)
  thumbnail?: string // optional, for future use
}

export interface StreamPlayerProps {
  id: string // publish_path to play
}
