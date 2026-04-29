<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Icon, Tag, Loading } from 'vant'
import { WebRTCPlayer } from '@/utils/webrtc'
import { useAuthStore } from '@/stores/auth'
import { useConfigStore } from '@/stores/config'
import type { Stream } from '@/types/stream'

const props = defineProps<{
  stream: Stream
}>()

const emit = defineEmits<{
  click: []
}>()

const authStore = useAuthStore()
const configStore = useConfigStore()

const videoRef = ref<HTMLVideoElement | null>(null)
const player = ref<WebRTCPlayer | null>(null)
const previewState = ref<'loading' | 'playing' | 'error'>('loading')

onMounted(async () => {
  if (props.stream.status !== 'online') {
    previewState.value = 'error'
    return
  }
  await startPreview()
})

onUnmounted(() => {
  stopPreview()
})

async function startPreview() {
  if (!videoRef.value || !authStore.token) {
    previewState.value = 'error'
    return
  }

  try {
    await configStore.fetchConfig()
    if (!configStore.mediamtxWebRTCURL) {
      previewState.value = 'error'
      return
    }

    player.value = new WebRTCPlayer({
      path: props.stream.id,
      token: authStore.token,
      videoElement: videoRef.value,
      mediamtxURL: configStore.mediamtxWebRTCURL,
      onConnectionStateChange: (state) => {
        if (state === 'connected') {
          previewState.value = 'playing'
        } else if (state === 'failed' || state === 'closed') {
          previewState.value = 'error'
        }
      },
      onError: () => {
        previewState.value = 'error'
      },
    })

    await player.value.start()
  } catch {
    previewState.value = 'error'
  }
}

function stopPreview() {
  if (player.value) {
    player.value.stop()
    player.value = null
  }
}
</script>

<template>
  <div class="preview-card" @click="emit('click')">
    <div class="preview-thumbnail">
      <video
        ref="videoRef"
        class="preview-video"
        :class="{ visible: previewState === 'playing' }"
        autoplay
        playsinline
        muted
      />
      <div v-if="previewState === 'loading'" class="preview-overlay">
        <Loading type="spinner" size="24px" color="rgba(255,255,255,0.8)" />
      </div>
      <div v-else-if="previewState === 'error'" class="preview-overlay">
        <Icon name="video" size="36" color="rgba(255,255,255,0.5)" />
      </div>
    </div>
    <div class="preview-info">
      <div class="preview-title">{{ stream.name }}</div>
      <div class="preview-location">{{ stream.location || $t('common.notSet') }}</div>
      <Tag :type="stream.status === 'online' ? 'success' : 'default'" size="medium">
        {{ stream.status === 'online' ? $t('streamList.online') : $t('streamList.offline') }}
      </Tag>
    </div>
  </div>
</template>

<style scoped>
.preview-card {
  background-color: var(--bg-card);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 8px var(--shadow-light);
  cursor: pointer;
  transition: transform 0.2s;
}

.preview-card:active {
  transform: scale(0.97);
}

.preview-thumbnail {
  position: relative;
  aspect-ratio: 16 / 9;
  background: linear-gradient(135deg, #2c3e50, #3498db);
  overflow: hidden;
}

.preview-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: none;
}

.preview-video.visible {
  display: block;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-info {
  padding: 10px;
}

.preview-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-location {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
