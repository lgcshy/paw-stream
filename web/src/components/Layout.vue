<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NavBar, Tabbar, TabbarItem } from 'vant'

const router = useRouter()
const route = useRoute()

const onClickLeft = () => {
  router.back()
}

// Show bottom nav based on route meta
const showBottomNav = computed(() => route.meta.showBottomNav === true)

// Show back button only on sub-pages (not on main pages with bottom nav)
const showBackButton = computed(() => !showBottomNav.value)

// Active tab - use ref for two-way binding
const activeTab = ref('streams')

// Update active tab when route changes
watch(
  () => route.path,
  (newPath) => {
    if (newPath.startsWith('/streams')) {
      activeTab.value = 'streams'
    } else if (newPath.startsWith('/devices')) {
      activeTab.value = 'devices'
    } else if (newPath.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

function setActiveTab(name: string) {
  if (name === 'streams') {
    router.push('/streams')
  } else if (name === 'devices') {
    router.push('/devices')
  } else if (name === 'profile') {
    router.push('/profile')
  }
}
</script>

<template>
  <div class="layout">
    <NavBar title="PawStream" :left-arrow="showBackButton" @click-left="onClickLeft" class="navbar" />
    <div class="content" :class="{ 'with-bottom-nav': showBottomNav }">
      <slot />
    </div>

    <!-- Bottom Navigation -->
    <Tabbar v-if="showBottomNav" v-model="activeTab" @change="setActiveTab" fixed placeholder>
      <TabbarItem name="streams" icon="video">直播</TabbarItem>
      <TabbarItem name="devices" icon="apps-o">设备</TabbarItem>
      <TabbarItem name="profile" icon="user-o">我的</TabbarItem>
    </Tabbar>
  </div>
</template>

<style scoped>
.layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.navbar {
  flex-shrink: 0;
}

.content {
  flex: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  background-color: var(--bg-primary); /* 添加背景色，这样页面不需要设置 min-height */
}

.content.with-bottom-nav {
  padding-bottom: 70px; /* 增加到70px，确保底部导航不遮挡内容 */
}
</style>
