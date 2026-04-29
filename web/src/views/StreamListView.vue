<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { List, Cell, Tag, Empty, Loading, PullRefresh, Icon, showFailToast, showSuccessToast } from 'vant'
import Layout from '@/components/Layout.vue'
import StreamPreviewCard from '@/components/StreamPreviewCard.vue'
import { useDeviceStore } from '@/stores/device'
import { useI18n } from 'vue-i18n'

const VIEW_MODE_KEY = 'pawstream_view_mode'
type ViewMode = 'list' | 'grid'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()

const streams = computed(() => deviceStore.enabledStreams)
const loading = computed(() => deviceStore.loading)
const viewMode = ref<ViewMode>((localStorage.getItem(VIEW_MODE_KEY) as ViewMode) || 'list')

function toggleViewMode() {
  viewMode.value = viewMode.value === 'list' ? 'grid' : 'list'
  localStorage.setItem(VIEW_MODE_KEY, viewMode.value)
}

onMounted(async () => {
  await loadStreams()
})

async function loadStreams() {
  try {
    await deviceStore.fetchPaths()
  } catch (error) {
    showFailToast(t('streamList.loadFailed'))
  }
}

async function onRefresh() {
  try {
    await deviceStore.refreshPaths()
    showSuccessToast(t('streamList.refreshSuccess'))
  } catch (error) {
    showFailToast(t('streamList.refreshFailed'))
  }
}

const goToPlayer = (streamId: string) => {
  router.push(`/stream/${encodeURIComponent(streamId)}`)
}
</script>

<template>
  <Layout>
    <div class="stream-list-view">
      <div class="header">
        <div class="header-row">
          <h2>{{ $t('streamList.title') }}</h2>
          <Icon
            :name="viewMode === 'list' ? 'photo' : 'bars'"
            size="22"
            class="view-toggle"
            @click="toggleViewMode"
          />
        </div>
        <p class="subtitle">{{ $t('streamList.subtitle') }}</p>
      </div>

      <PullRefresh v-model="loading" @refresh="onRefresh">
        <div v-if="loading && streams.length === 0" class="loading-container">
          <Loading type="spinner" size="48px">{{ $t('common.loading') }}</Loading>
        </div>
        <Empty v-else-if="streams.length === 0" :description="$t('streamList.noDevices')">
          <template #description>
            <p>{{ $t('streamList.noDevices') }}</p>
            <p style="font-size: 12px; color: #969799; margin-top: 8px">{{ $t('streamList.noDevicesHint') }}</p>
          </template>
        </Empty>
        <List v-else-if="viewMode === 'list'">
          <Cell
            v-for="stream in streams"
            :key="stream.id"
            :title="stream.name"
            :label="`${$t('streamList.location')}: ${stream.location || $t('common.notSet')} | ${$t('streamList.path')}: ${stream.id}`"
            is-link
            @click="goToPlayer(stream.id)"
          >
            <template #right-icon>
              <Tag :type="stream.status === 'online' ? 'success' : 'default'">
                {{ stream.status === 'online' ? $t('streamList.online') : $t('streamList.offline') }}
              </Tag>
            </template>
          </Cell>
        </List>

        <div v-else class="grid-view">
          <StreamPreviewCard
            v-for="stream in streams"
            :key="stream.id"
            :stream="stream"
            @click="goToPlayer(stream.id)"
          />
        </div>
      </PullRefresh>
    </div>
  </Layout>
</template>

<style scoped>
.stream-list-view {
  width: 100%;
}

.header {
  padding: 20px 16px;
  background-color: var(--bg-secondary);
  margin-bottom: 10px;
}

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header h2 {
  font-size: 24px;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.view-toggle {
  cursor: pointer;
  color: var(--text-secondary);
  padding: 4px;
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
}

.grid-view {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 12px;
}
</style>
