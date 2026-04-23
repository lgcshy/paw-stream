<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { List, Cell, Tag, Empty, Loading, PullRefresh, Icon, showFailToast, showSuccessToast } from 'vant'
import Layout from '@/components/Layout.vue'
import { useDeviceStore } from '@/stores/device'

const VIEW_MODE_KEY = 'pawstream_view_mode'
type ViewMode = 'list' | 'grid'

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
    showFailToast('加载设备列表失败')
  }
}

async function onRefresh() {
  try {
    await deviceStore.refreshPaths()
    showSuccessToast('刷新成功')
  } catch (error) {
    showFailToast('刷新失败')
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
          <h2>直播流</h2>
          <Icon
            :name="viewMode === 'list' ? 'photo' : 'bars'"
            size="22"
            class="view-toggle"
            @click="toggleViewMode"
          />
        </div>
        <p class="subtitle">选择一个设备观看实时画面</p>
      </div>

      <PullRefresh v-model="loading" @refresh="onRefresh">
        <div v-if="loading && streams.length === 0" class="loading-container">
          <Loading type="spinner" size="48px">加载中...</Loading>
        </div>
        <Empty v-else-if="streams.length === 0" description="暂无在线设备">
          <template #description>
            <p>暂无在线设备</p>
            <p style="font-size: 12px; color: #969799; margin-top: 8px">请先在"设备管理"中创建并启用设备</p>
          </template>
        </Empty>
        <!-- 列表模式 -->
        <List v-else-if="viewMode === 'list'">
          <Cell
            v-for="stream in streams"
            :key="stream.id"
            :title="stream.name"
            :label="`位置: ${stream.location || '未设置'} | 路径: ${stream.id}`"
            is-link
            @click="goToPlayer(stream.id)"
          >
            <template #right-icon>
              <Tag :type="stream.status === 'online' ? 'success' : 'default'">
                {{ stream.status === 'online' ? '在线' : '离线' }}
              </Tag>
            </template>
          </Cell>
        </List>

        <!-- 图片/卡片模式 -->
        <div v-else class="grid-view">
          <div
            v-for="stream in streams"
            :key="stream.id"
            class="stream-card"
            @click="goToPlayer(stream.id)"
          >
            <div class="card-thumbnail">
              <Icon name="video" size="36" color="rgba(255,255,255,0.6)" />
            </div>
            <div class="card-info">
              <div class="card-title">{{ stream.name }}</div>
              <div class="card-location">{{ stream.location || '未设置' }}</div>
              <Tag :type="stream.status === 'online' ? 'success' : 'default'" size="small">
                {{ stream.status === 'online' ? '在线' : '离线' }}
              </Tag>
            </div>
          </div>
        </div>
      </PullRefresh>
    </div>
  </Layout>
</template>

<style scoped>
.stream-list-view {
  width: 100%;
  /* 移除 min-height: 100%，让内容自然增长，避免滚动问题 */
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

/* 图片/卡片模式 */
.grid-view {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 12px;
}

.stream-card {
  background-color: var(--bg-card);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 8px var(--shadow-light);
  cursor: pointer;
  transition: transform 0.2s;
}

.stream-card:active {
  transform: scale(0.97);
}

.card-thumbnail {
  aspect-ratio: 16 / 9;
  background: linear-gradient(135deg, #2c3e50, #3498db);
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-info {
  padding: 10px;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-location {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
