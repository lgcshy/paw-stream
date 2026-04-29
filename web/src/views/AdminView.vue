<template>
  <div class="admin-page">
    <van-nav-bar :title="$t('admin.title')" left-arrow @click-left="router.back()" fixed placeholder />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- Loading -->
      <div v-if="loading && !dashboard" class="loading-container">
        <van-loading type="spinner" size="40px">{{ $t('common.loading') }}</van-loading>
      </div>

      <!-- No permission -->
      <van-empty v-else-if="noPermission" image="error" :description="$t('admin.noPermission')" />

      <!-- Dashboard -->
      <div v-else-if="dashboard" class="dashboard-container">
        <!-- Overview Cards -->
        <van-cell-group inset :title="$t('admin.overview')">
          <van-grid :column-num="3" :border="false">
            <van-grid-item>
              <div class="stat-card">
                <div class="stat-value">{{ dashboard.total_devices }}</div>
                <div class="stat-label">{{ $t('admin.totalDevices') }}</div>
              </div>
            </van-grid-item>
            <van-grid-item>
              <div class="stat-card stat-online">
                <div class="stat-value">{{ dashboard.online_devices }}</div>
                <div class="stat-label">{{ $t('admin.onlineDevices') }}</div>
              </div>
            </van-grid-item>
            <van-grid-item>
              <div class="stat-card stat-offline">
                <div class="stat-value">{{ dashboard.total_devices - dashboard.online_devices }}</div>
                <div class="stat-label">{{ $t('admin.offlineDevices') }}</div>
              </div>
            </van-grid-item>
          </van-grid>
        </van-cell-group>

        <!-- Device List -->
        <van-cell-group inset :title="$t('admin.allDevices')">
          <van-cell
            v-for="device in dashboard.devices"
            :key="device.id"
            :title="device.name"
            :label="deviceLabel(device)"
            is-link
            @click="goToDevice(device.id)"
          >
            <template #icon>
              <span class="status-indicator" :class="device.is_online ? 'online' : 'offline'" />
            </template>
            <template #right-icon>
              <van-tag :type="device.is_online ? 'success' : 'default'" size="medium">
                {{ device.is_online ? $t('admin.online') : $t('admin.offline') }}
              </van-tag>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { deviceApi } from '@/api/device'
import type { AdminDashboard, AdminDeviceInfo } from '@/types/api'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const refreshing = ref(false)
const noPermission = ref(false)
const dashboard = ref<AdminDashboard | null>(null)

function deviceLabel(device: AdminDeviceInfo) {
  const parts = [device.location || t('common.notSet')]
  if (device.last_seen_at) {
    const date = new Date(device.last_seen_at)
    parts.push(`${t('admin.lastSeen')}: ${date.toLocaleString()}`)
  } else {
    parts.push(`${t('admin.lastSeen')}: ${t('admin.never')}`)
  }
  return parts.join(' · ')
}

function goToDevice(id: string) {
  router.push(`/devices/${id}`)
}

async function loadDashboard() {
  loading.value = true
  try {
    dashboard.value = await deviceApi.getAdminDashboard()
    noPermission.value = false
  } catch (error: any) {
    if (error.statusCode === 403) {
      noPermission.value = true
    } else {
      showFailToast(t('admin.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  refreshing.value = true
  try {
    dashboard.value = await deviceApi.getAdminDashboard()
    showSuccessToast(t('admin.refreshSuccess'))
  } catch (error: any) {
    showFailToast(t('admin.refreshFailed'))
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.admin-page {
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

.dashboard-container {
  padding-top: 8px;
}

.stat-card {
  text-align: center;
  padding: 12px 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-online .stat-value {
  color: #07c160;
}

.stat-offline .stat-value {
  color: #969799;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.van-cell-group {
  margin-bottom: 16px;
}

.status-indicator {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 12px;
  flex-shrink: 0;
}

.status-indicator.online {
  background-color: #07c160;
  box-shadow: 0 0 6px rgba(7, 193, 96, 0.6);
}

.status-indicator.offline {
  background-color: #c8c9cc;
}
</style>
