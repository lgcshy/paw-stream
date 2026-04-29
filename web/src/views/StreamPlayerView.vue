<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Loading, Button, showFailToast, showSuccessToast } from 'vant'
import Layout from '@/components/Layout.vue'
import { useAuthStore } from '@/stores/auth'
import { useDeviceStore } from '@/stores/device'
import { useConfigStore } from '@/stores/config'
import { WebRTCPlayer } from '@/utils/webrtc'
import type { StreamPlayerProps } from '@/types/stream'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const props = defineProps<StreamPlayerProps>()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const deviceStore = useDeviceStore()
const configStore = useConfigStore()

const videoRef = ref<HTMLVideoElement | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const connectionState = ref<string>('new')
const streamName = ref('')
const player = ref<WebRTCPlayer | null>(null)

onMounted(async () => {
  const streamPath = props.id || (route.params.id as string)

  if (!streamPath) {
    error.value = t('player.invalidPath')
    loading.value = false
    return
  }

  // Get stream info
  const stream = deviceStore.getStreamById(streamPath)
  if (stream) {
    streamName.value = stream.name
  } else {
    streamName.value = streamPath
  }

  // Start WebRTC connection
  await startStream(streamPath)
})

onUnmounted(() => {
  stopStream()
})

async function startStream(path: string) {
  if (!videoRef.value) {
    error.value = t('player.videoNotReady')
    loading.value = false
    return
  }

  if (!authStore.token) {
    error.value = t('player.notLoggedIn')
    loading.value = false
    return
  }

  loading.value = true
  error.value = null

  try {
    // Fetch MediaMTX URL from config
    await configStore.fetchConfig()

    if (!configStore.mediamtxWebRTCURL) {
      throw new Error(t('player.configError'))
    }

    player.value = new WebRTCPlayer({
      path: path,
      token: authStore.token,
      videoElement: videoRef.value,
      mediamtxURL: configStore.mediamtxWebRTCURL,
      onConnectionStateChange: (state) => {
        connectionState.value = state
        console.log('Connection state:', state)

        if (state === 'connected') {
          loading.value = false
          showSuccessToast(t('player.connected'))
        } else if (state === 'failed') {
          error.value = t('player.connectionFailed')
          loading.value = false
        }
      },
      onError: (err) => {
        console.error('WebRTC error:', err)
        error.value = err.message
        loading.value = false
        showFailToast(t('player.playFailed') + ': ' + err.message)
      },
    })

    await player.value.start()
  } catch (err) {
    console.error('Failed to start stream:', err)
    error.value = err instanceof Error ? err.message : t('player.cannotStart')
    loading.value = false
    showFailToast(t('player.playFailed'))
  }
}

function stopStream() {
  if (player.value) {
    player.value.stop()
    player.value = null
  }
}

async function retry() {
  stopStream()
  const streamPath = props.id || (route.params.id as string)
  await startStream(streamPath)
}

function goBack() {
  router.push('/streams')
}
</script>

<template>
  <Layout>
    <div class="stream-player-view">
      <div class="player-container">
        <video
          ref="videoRef"
          class="video-player"
          autoplay
          playsinline
          muted
          :class="{ hidden: loading || error }"
        />

        <div v-if="loading" class="overlay">
          <Loading type="spinner" size="48px" color="#fff">{{ $t('player.connecting') }}</Loading>
        </div>

        <div v-if="error && !loading" class="overlay error-overlay">
          <div class="error-content">
            <p class="error-icon">⚠️</p>
            <p class="error-message">{{ error }}</p>
            <div class="error-actions">
              <Button type="primary" size="small" @click="retry">{{ $t('common.retry') }}</Button>
              <Button size="small" @click="goBack">{{ $t('common.back') }}</Button>
            </div>
          </div>
        </div>
      </div>

      <div class="controls">
        <div class="stream-info">
          <p class="stream-name">{{ streamName }}</p>
          <p class="stream-status">
            <span :class="['status-dot', connectionState === 'connected' ? 'online' : 'offline']"></span>
            {{ connectionState === 'connected' ? $t('player.playing') : connectionState }}
          </p>
        </div>
        <div class="control-buttons">
          <Button size="small" @click="retry" :disabled="loading">{{ $t('player.reconnect') }}</Button>
          <Button size="small" @click="goBack">{{ $t('player.backToList') }}</Button>
        </div>
      </div>
    </div>
  </Layout>
</template>

<style scoped>
.stream-player-view {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #000;
}

.player-container {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #1a1a1a;
  overflow: hidden;
}

.video-player {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.video-player.hidden {
  display: none;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.8);
}

.error-overlay {
  background-color: rgba(0, 0, 0, 0.9);
}

.error-content {
  text-align: center;
  color: #fff;
  padding: 20px;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  margin-bottom: 24px;
  line-height: 1.5;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.controls {
  padding: 16px;
  background-color: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

.stream-info {
  margin-bottom: 12px;
}

.stream-name {
  color: #fff;
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
}

.stream-status {
  color: var(--text-secondary);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot.online {
  background-color: var(--success-color);
  box-shadow: 0 0 6px var(--success-color);
}

.status-dot.offline {
  background-color: var(--text-disabled);
}

.control-buttons {
  display: flex;
  gap: 12px;
}

.control-buttons button {
  flex: 1;
}
</style>
