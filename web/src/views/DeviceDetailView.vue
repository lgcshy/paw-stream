<template>
  <div class="device-detail-page">
    <van-nav-bar :title="$t('deviceDetail.title')" left-arrow @click-left="router.back()" fixed placeholder />

    <!-- Loading State -->
    <div v-if="loading && !device" class="loading-container">
      <van-loading type="spinner" size="40px">{{ $t('common.loading') }}</van-loading>
    </div>

    <!-- Device Details -->
    <div v-else-if="device" class="detail-container">
      <!-- Rotated Secret Display -->
      <div v-if="rotatedSecret" class="secret-section">
        <van-notice-bar type="success" :scrollable="false">
          <template #default> {{ $t('deviceDetail.secretRotatedNotice') }} </template>
        </van-notice-bar>
        <SecretDisplay :secret="rotatedSecret" />
        <van-button type="default" block round @click="rotatedSecret = null">{{ $t('common.close') }}</van-button>
      </div>

      <!-- Device Info -->
      <van-cell-group inset :title="$t('deviceDetail.basicInfo')">
        <van-cell :title="$t('deviceDetail.deviceName')" :value="device.name" />
        <van-cell :title="$t('deviceDetail.location')" :value="device.location || $t('common.notSet')" />
        <van-cell :title="$t('deviceDetail.publishPath')" :value="device.publish_path" />
        <van-cell :title="$t('deviceDetail.deviceId')" :value="device.id" />
        <van-cell :title="$t('deviceDetail.createdTime')" :value="formatDateTime(device.created_at)" />
        <van-cell :title="$t('deviceDetail.updatedTime')" :value="formatDateTime(device.updated_at)" />
      </van-cell-group>

      <!-- Device Status -->
      <van-cell-group inset :title="$t('deviceDetail.status')">
        <van-cell :title="$t('deviceDetail.enableStatus')">
          <template #right-icon>
            <van-switch v-model="deviceEnabled" size="24" @change="handleToggleStatus" :loading="toggleLoading" />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Actions -->
      <van-cell-group inset :title="$t('deviceDetail.actions')">
        <van-cell :title="$t('deviceDetail.editDevice')" is-link @click="handleEdit">
          <template #icon>
            <van-icon name="edit" color="#1989fa" />
          </template>
        </van-cell>
        <van-cell :title="$t('deviceDetail.rotateSecret')" is-link @click="handleRotateSecret">
          <template #icon>
            <van-icon name="replay" color="#ff976a" />
          </template>
        </van-cell>
        <van-cell v-if="deviceEnabled" :title="$t('deviceDetail.watchStream')" is-link @click="handlePlayStream">
          <template #icon>
            <van-icon name="play-circle" color="#07c160" />
          </template>
        </van-cell>
        <van-cell :title="$t('deviceDetail.deleteDevice')" is-link @click="handleDelete">
          <template #icon>
            <van-icon name="delete" color="#ee0a24" />
          </template>
        </van-cell>
      </van-cell-group>
    </div>

    <!-- Confirm Dialog for Delete -->
    <ConfirmDialog
      v-model:show="showDeleteConfirm"
      :title="$t('deviceDetail.deleteConfirmTitle')"
      :message="$t('deviceDetail.deleteConfirmMessage')"
      :confirm-text="$t('common.delete')"
      :cancel-text="$t('common.cancel')"
      danger
      @confirm="confirmDelete"
    />

    <!-- Confirm Dialog for Rotate Secret -->
    <ConfirmDialog
      v-model:show="showRotateConfirm"
      :title="$t('deviceDetail.rotateConfirmTitle')"
      :message="$t('deviceDetail.rotateConfirmMessage')"
      :confirm-text="$t('deviceDetail.rotateButton')"
      :cancel-text="$t('common.cancel')"
      @confirm="confirmRotateSecret"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useDeviceStore } from '@/stores/device'
import { showSuccessToast, showFailToast } from 'vant'
import SecretDisplay from '@/components/SecretDisplay.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { DeviceInfo } from '@/types/api'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const deviceStore = useDeviceStore()

const loading = ref(false)
const toggleLoading = ref(false)
const device = ref<DeviceInfo | null>(null)
const deviceEnabled = ref(true)
const rotatedSecret = ref<string | null>(null)

const showDeleteConfirm = ref(false)
const showRotateConfirm = ref(false)

// Device ID from route
const deviceId = computed(() => route.params.id as string)

// Format helpers
function formatDateTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Actions
function handleEdit() {
  router.push(`/devices/${deviceId.value}/edit`)
}

function handlePlayStream() {
  if (device.value) {
    router.push(`/stream/${device.value.publish_path}`)
  }
}

function handleDelete() {
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  loading.value = true
  try {
    await deviceStore.deleteDevice(deviceId.value)
    showSuccessToast(t('deviceDetail.deleteSuccess'))
    router.push('/devices')
  } catch (error: any) {
    console.error('Delete failed:', error)
    showFailToast(error.message || t('deviceDetail.deleteFailed'))
  } finally {
    loading.value = false
  }
}

function handleRotateSecret() {
  showRotateConfirm.value = true
}

async function confirmRotateSecret() {
  loading.value = true
  try {
    const response = await deviceStore.rotateSecret(deviceId.value)
    rotatedSecret.value = response.new_secret
    showSuccessToast(t('deviceDetail.rotateSuccess'))
  } catch (error: any) {
    console.error('Rotate secret failed:', error)
    showFailToast(error.message || t('deviceDetail.rotateFailed'))
  } finally {
    loading.value = false
  }
}

async function handleToggleStatus() {
  toggleLoading.value = true
  try {
    await deviceStore.updateDevice(deviceId.value, {
      disabled: !deviceEnabled.value,
    })
    if (device.value) {
      device.value.disabled = !deviceEnabled.value
    }
    showSuccessToast(deviceEnabled.value ? t('deviceDetail.deviceEnabled') : t('deviceDetail.deviceDisabled'))
  } catch (error: any) {
    console.error('Toggle status failed:', error)
    // Revert the switch
    deviceEnabled.value = !deviceEnabled.value
    showFailToast(error.message || t('deviceDetail.statusUpdateFailed'))
  } finally {
    toggleLoading.value = false
  }
}

// Load device data
onMounted(async () => {
  loading.value = true
  try {
    device.value = await deviceStore.getDevice(deviceId.value)
    deviceEnabled.value = !device.value.disabled
  } catch (error: any) {
    console.error('Load device failed:', error)
    showFailToast(error.message || t('deviceDetail.loadFailed'))
    router.back()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.device-detail-page {
  min-height: 100vh;
  background-color: var(--bg-primary);
  padding-bottom: 20px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.detail-container {
  padding: 16px 0;
}

.secret-section {
  margin: 16px;
  padding: 16px;
  background-color: var(--bg-card);
  border-radius: 8px;
}

.secret-section .van-notice-bar {
  margin-bottom: 16px;
}

.secret-section .van-button {
  margin-top: 16px;
}

.van-cell-group {
  margin-bottom: 16px;
}

.van-cell {
  align-items: center;
}

:deep(.van-cell__left-icon) {
  margin-right: 12px;
  font-size: 20px;
}
</style>
