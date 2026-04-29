<template>
  <Layout>
    <div class="device-list-page">
      <div class="header">
        <h2>{{ $t('deviceList.title') }}</h2>
        <van-button icon="plus" type="primary" size="small" @click="handleAddDevice">
          {{ $t('deviceList.addDevice') }}
        </van-button>
      </div>

      <!-- Pull to Refresh -->
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Loading State -->
      <div v-if="deviceStore.loading && !refreshing" class="loading-container">
        <van-loading type="spinner" size="40px">{{ $t('common.loading') }}</van-loading>
      </div>

      <!-- Empty State -->
      <van-empty v-else-if="!deviceStore.loading && devices.length === 0" :description="$t('deviceList.noDevices')">
        <van-button type="primary" round @click="handleAddDevice">{{ $t('deviceList.createFirst') }}</van-button>
      </van-empty>

      <!-- Device List -->
      <div v-else class="device-list">
        <van-cell
          v-for="device in devices"
          :key="device.id"
          :title="device.name"
          :label="deviceLabel(device)"
          :value="deviceStatus(device)"
          :value-class="deviceStatusClass(device)"
          is-link
          @click="handleDeviceClick(device.id)"
        >
          <template #icon>
            <van-icon :name="deviceIcon(device)" size="24" :color="deviceIconColor(device)" />
          </template>
        </van-cell>
      </div>
    </van-pull-refresh>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDeviceStore } from '@/stores/device'
import { showFailToast } from 'vant'
import Layout from '@/components/Layout.vue'
import type { DeviceInfo } from '@/types/api'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()
const refreshing = ref(false)

const devices = computed(() => deviceStore.devices ?? [])

// Device display helpers
function deviceLabel(device: DeviceInfo) {
  const parts = []
  if (device.location) {
    parts.push(device.location)
  }
  parts.push(t('deviceList.createdAt', { date: formatDate(device.created_at) }))
  return parts.join(' · ')
}

function deviceStatus(device: DeviceInfo) {
  return device.disabled ? t('deviceList.disabled') : t('deviceList.enabled')
}

function deviceStatusClass(device: DeviceInfo) {
  return device.disabled ? 'status-disabled' : 'status-enabled'
}

function deviceIcon(device: DeviceInfo) {
  return device.disabled ? 'close' : 'success'
}

function deviceIconColor(device: DeviceInfo) {
  return device.disabled ? '#c8c9cc' : '#07c160'
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

// Actions
function handleAddDevice() {
  router.push('/devices/new')
}

function handleDeviceClick(deviceId: string) {
  router.push(`/devices/${deviceId}`)
}

async function onRefresh() {
  refreshing.value = true
  try {
    await deviceStore.refreshDevices()
  } catch (error: any) {
    console.error('Refresh failed:', error)
    showFailToast(error.message || t('deviceList.refreshFailed'))
  } finally {
    refreshing.value = false
  }
}

// Lifecycle
onMounted(async () => {
  // Always fetch devices when mounting
  try {
    await deviceStore.fetchDevices()
  } catch (error: any) {
    console.error('Fetch devices failed:', error)
    showFailToast(error.message || t('deviceList.fetchFailed'))
  }
})
</script>

<style scoped>
.device-list-page {
  width: 100%;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 16px;
  background-color: var(--bg-secondary);
  margin-bottom: 10px;
}

.header h2 {
  font-size: 24px;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.device-list {
  padding-bottom: 20px;
}

.van-cell {
  margin-bottom: 8px;
}

:deep(.van-cell__left-icon) {
  margin-right: 12px;
}

:deep(.status-enabled) {
  color: var(--success-color);
  font-weight: 500;
}

:deep(.status-disabled) {
  color: var(--text-disabled);
}
</style>
