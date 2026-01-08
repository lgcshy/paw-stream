<template>
  <div class="secret-display">
    <!-- Warning Notice -->
    <van-notice-bar left-icon="warning-o" color="#ed6a0c" background="#fffbe8" :scrollable="false">
      <template #default>
        此密钥仅显示一次，请立即保存！关闭后将无法再次查看
      </template>
    </van-notice-bar>

    <!-- Secret Content -->
    <div class="secret-content">
      <div class="secret-label">设备密钥 (Device Secret)</div>
      <div class="secret-value">{{ secret }}</div>
    </div>

    <!-- Copy Button -->
    <van-button type="primary" block @click="copyToClipboard" :loading="copying">
      <van-icon name="copy" />
      {{ copied ? '已复制' : '复制密钥' }}
    </van-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { showSuccessToast, showFailToast } from 'vant'

interface Props {
  secret: string
}

const props = defineProps<Props>()

const copying = ref(false)
const copied = ref(false)

async function copyToClipboard() {
  copying.value = true
  try {
    await navigator.clipboard.writeText(props.secret)
    copied.value = true
    showSuccessToast('密钥已复制到剪贴板')

    // Reset copied state after 2 seconds
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
    showFailToast('复制失败，请手动复制')
  } finally {
    copying.value = false
  }
}
</script>

<style scoped>
.secret-display {
  padding: 16px;
  background-color: var(--bg-primary);
  border-radius: 8px;
}

.secret-content {
  margin: 16px 0;
  padding: 16px;
  background-color: var(--bg-card);
  border-radius: 4px;
  border: 1px solid var(--border-color);
}

.secret-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.secret-value {
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  color: var(--text-primary);
  word-break: break-all;
  line-height: 1.6;
  padding: 8px;
  background-color: var(--bg-primary);
  border-radius: 4px;
}

.van-button {
  margin-top: 8px;
}

.van-icon {
  margin-right: 4px;
}
</style>
