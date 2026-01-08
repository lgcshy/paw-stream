<template>
  <div class="device-form-page">
    <van-nav-bar :title="pageTitle" left-arrow @click-left="router.back()" fixed placeholder />

    <div class="form-container">
      <!-- Success State - Show Secret (Create Only) -->
      <div v-if="createdSecret" class="success-section">
        <van-notice-bar type="success" :scrollable="false">
          <template #default> 设备创建成功！请立即保存密钥 </template>
        </van-notice-bar>

        <SecretDisplay :secret="createdSecret" />

        <van-button type="success" block round @click="handleDone">完成</van-button>
      </div>

      <!-- Form -->
      <van-form v-else @submit="handleSubmit">
        <!-- Device Name -->
        <van-cell-group inset>
          <van-field
            v-model="form.name"
            name="name"
            label="设备名称"
            placeholder="例如：客厅摄像头"
            clearable
            :rules="[
              { required: true, message: '请输入设备名称' },
              { pattern: /^.{1,100}$/, message: '设备名称不能超过100个字符' },
            ]"
          />
        </van-cell-group>

        <!-- Location -->
        <van-cell-group inset>
          <van-field
            v-model="form.location"
            name="location"
            label="位置"
            placeholder="例如：客厅、卧室（可选）"
            clearable
            :rules="[{ pattern: /^.{0,200}$/, message: '位置不能超过200个字符' }]"
          />
        </van-cell-group>

        <!-- Enable/Disable Toggle (Edit Only) -->
        <van-cell-group v-if="isEditMode" inset>
          <van-cell title="设备状态">
            <template #right-icon>
              <van-switch v-model="form.enabled" size="24" />
            </template>
          </van-cell>
        </van-cell-group>

        <!-- Info Notice -->
        <div class="info-notice">
          <van-notice-bar v-if="!isEditMode" left-icon="info-o" color="#1989fa" background="#ecf9ff" :scrollable="false">
            <template #default> 创建后将生成密钥，请妥善保存，之后无法再次查看 </template>
          </van-notice-bar>
        </div>

        <!-- Submit Button -->
        <div class="form-actions">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            {{ isEditMode ? '保存' : '创建' }}
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useDeviceStore } from '@/stores/device'
import { showSuccessToast, showFailToast } from 'vant'
import SecretDisplay from '@/components/SecretDisplay.vue'

const router = useRouter()
const route = useRoute()
const deviceStore = useDeviceStore()

const loading = ref(false)
const createdSecret = ref<string | null>(null)

const form = ref({
  name: '',
  location: '',
  enabled: true,
})

// Mode detection
const isEditMode = computed(() => !!route.params.id)
const pageTitle = computed(() => (isEditMode.value ? '编辑设备' : '创建设备'))

// Form submission
async function handleSubmit() {
  loading.value = true
  try {
    if (isEditMode.value) {
      // Update existing device
      const deviceId = route.params.id as string
      await deviceStore.updateDevice(deviceId, {
        name: form.value.name,
        location: form.value.location || undefined,
        disabled: !form.value.enabled,
      })
      showSuccessToast('设备更新成功')
      router.back()
    } else {
      // Create new device
      const response = await deviceStore.createDevice({
        name: form.value.name,
        location: form.value.location || undefined,
      })
      createdSecret.value = response.secret
      // Don't navigate yet - show secret first
    }
  } catch (error: any) {
    console.error('Submit failed:', error)
    showFailToast(error.message || (isEditMode.value ? '更新失败' : '创建失败'))
  } finally {
    loading.value = false
  }
}

function handleDone() {
  // Navigate to device list
  router.push('/devices')
}

// Load device data for edit mode
onMounted(async () => {
  if (isEditMode.value) {
    const deviceId = route.params.id as string
    try {
      const device = await deviceStore.getDevice(deviceId)
      form.value.name = device.name
      form.value.location = device.location || ''
      form.value.enabled = !device.disabled
    } catch (error: any) {
      console.error('Load device failed:', error)
      showFailToast(error.message || '加载设备信息失败')
      router.back()
    }
  }
})
</script>

<style scoped>
.device-form-page {
  min-height: 100vh;
  background-color: var(--bg-primary);
}

.form-container {
  padding: 16px;
}

.success-section {
  margin-top: 20px;
}

.success-section .van-notice-bar {
  margin-bottom: 16px;
}

.success-section .van-button {
  margin-top: 16px;
}

.van-cell-group {
  margin-bottom: 12px;
}

.info-notice {
  margin: 16px 0;
}

.form-actions {
  margin-top: 24px;
}
</style>
