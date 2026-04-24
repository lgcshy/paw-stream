<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { List, Cell, Tag, Empty, Loading, PullRefresh, Icon, showFailToast, showSuccessToast } from 'vant'
import Layout from '@/components/Layout.vue'
import StreamPreviewCard from '@/components/StreamPreviewCard.vue'
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

        <!-- 图片/卡片模式（实时预览） -->
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
</style>
