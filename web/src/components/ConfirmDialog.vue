<template>
  <van-dialog
    v-model:show="visible"
    :title="title"
    :message="message"
    :confirm-button-text="confirmText"
    :cancel-button-text="cancelText"
    show-cancel-button
    @confirm="handleConfirm"
    @cancel="handleCancel"
  >
    <template v-if="$slots.default" #default>
      <slot />
    </template>
  </van-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  show: boolean
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  danger?: boolean
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '确认操作',
  message: '确定要执行此操作吗？',
  confirmText: '确认',
  cancelText: '取消',
  danger: false,
})

const emit = defineEmits<Emits>()

const visible = ref(props.show)

// Sync with parent
watch(
  () => props.show,
  (newValue) => {
    visible.value = newValue
  }
)

watch(visible, (newValue) => {
  emit('update:show', newValue)
})

function handleConfirm() {
  emit('confirm')
  visible.value = false
}

function handleCancel() {
  emit('cancel')
  visible.value = false
}
</script>

<style scoped>
/* Custom styles can be added here if needed */
</style>
