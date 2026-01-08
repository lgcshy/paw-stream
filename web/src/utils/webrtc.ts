// WebRTC Player Utility for MediaMTX

export interface WebRTCPlayerOptions {
  path: string
  token: string
  videoElement: HTMLVideoElement
  mediamtxURL: string // MediaMTX WebRTC URL from config
  onConnectionStateChange?: (state: RTCPeerConnectionState) => void
  onError?: (error: Error) => void
}

export class WebRTCPlayer {
  private pc: RTCPeerConnection | null = null
  private mediaStream: MediaStream | null = null
  private options: WebRTCPlayerOptions

  constructor(options: WebRTCPlayerOptions) {
    this.options = options
  }

  async start() {
    try {
      // Create peer connection
      this.pc = new RTCPeerConnection({
        iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
      })

      // Handle connection state changes
      this.pc.onconnectionstatechange = () => {
        if (this.pc && this.options.onConnectionStateChange) {
          this.options.onConnectionStateChange(this.pc.connectionState)
        }

        // Auto-cleanup on failed/closed
        if (this.pc && (this.pc.connectionState === 'failed' || this.pc.connectionState === 'closed')) {
          this.stop()
        }
      }

      // Handle incoming tracks
      this.pc.ontrack = (event) => {
        console.log('Received track:', event.track.kind)
        if (!this.mediaStream) {
          this.mediaStream = new MediaStream()
        }
        this.mediaStream.addTrack(event.track)
        this.options.videoElement.srcObject = this.mediaStream
      }

      // Add transceiver for receiving video
      this.pc.addTransceiver('video', { direction: 'recvonly' })

      // Add transceiver for receiving audio (if available)
      this.pc.addTransceiver('audio', { direction: 'recvonly' })

      // Create offer
      const offer = await this.pc.createOffer()
      await this.pc.setLocalDescription(offer)

      // Wait for ICE gathering to complete
      if (this.pc.iceGatheringState !== 'complete') {
        await new Promise<void>((resolve) => {
          const checkState = () => {
            if (this.pc?.iceGatheringState === 'complete') {
              this.pc.removeEventListener('icegatheringstatechange', checkState)
              resolve()
            }
          }
          this.pc?.addEventListener('icegatheringstatechange', checkState)
          checkState()
        })
      }

      // Send offer to MediaMTX
      const response = await fetch(`${this.options.mediamtxURL}/${this.options.path}/whep`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/sdp',
          'Authorization': `Bearer ${this.options.token}`,
        },
        body: this.pc.localDescription?.sdp,
      })

      if (!response.ok) {
        throw new Error(`MediaMTX WebRTC connection failed: ${response.status} ${response.statusText}`)
      }

      // Set remote description from answer
      const answerSDP = await response.text()
      await this.pc.setRemoteDescription(new RTCSessionDescription({
        type: 'answer',
        sdp: answerSDP,
      }))

      console.log('WebRTC connection established')
    } catch (error) {
      console.error('WebRTC connection error:', error)
      if (this.options.onError) {
        this.options.onError(error instanceof Error ? error : new Error(String(error)))
      }
      this.stop()
      throw error
    }
  }

  stop() {
    if (this.mediaStream) {
      this.mediaStream.getTracks().forEach((track) => track.stop())
      this.mediaStream = null
    }

    if (this.pc) {
      this.pc.close()
      this.pc = null
    }

    if (this.options.videoElement) {
      this.options.videoElement.srcObject = null
    }
  }

  getConnectionState(): RTCPeerConnectionState | null {
    return this.pc?.connectionState || null
  }
}
