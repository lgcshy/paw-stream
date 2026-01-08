<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'
import { RouterView } from 'vue-router'
import { ConfigProvider } from 'vant'
import { useThemeStore } from '@/stores/theme'

// Initialize theme store
const themeStore = useThemeStore()

// System theme listener cleanup function
let cleanupThemeListener: (() => void) | undefined

onMounted(() => {
  // Initialize system theme listener
  cleanupThemeListener = themeStore.initSystemThemeListener()
})

onUnmounted(() => {
  // Cleanup system theme listener
  if (cleanupThemeListener) {
    cleanupThemeListener()
  }
})

// Vant ConfigProvider theme configuration
const themeVars = computed(() => {
  if (themeStore.isDark) {
    // Dark theme configuration
    return {
      // Background Colors
      background: '#1a1a1a',
      background2: '#2c2c2c',
      
      // Text Colors
      textColor: '#e5e5e5',
      textColor2: '#b0b0b0',
      textColor3: '#888888',
      
      // Border Colors
      borderColor: '#404040',
      
      // Component Colors
      cellBackground: '#2c2c2c',
      cellGroupBackground: '#1a1a1a',
      
      // Cell specific
      cellTextColor: '#e5e5e5',
      cellLabelColor: '#a0a0a0',
      cellValueColor: '#b0b0b0',
      
      // NavBar
      navBarBackground: '#2c2c2c',
      navBarTextColor: '#e5e5e5',
      navBarTitleTextColor: '#e5e5e5',
      
      // Tabbar
      tabbarBackground: '#2c2c2c',
      tabbarItemTextColor: '#888888',
      tabbarItemActiveColor: '#64B5F6',
      tabbarItemActiveBackground: 'rgba(100, 181, 246, 0.1)',
      
      // Button
      buttonPrimaryBackground: '#4da6ff',
      buttonPrimaryBorderColor: '#4da6ff',
      buttonPrimaryColor: '#ffffff',
      
      // Field
      fieldLabelColor: '#b0b0b0',
      fieldInputTextColor: '#e5e5e5',
      fieldPlaceholderTextColor: '#666666',
      
      // Dialog & Popup
      dialogBackground: '#2c2c2c',
      popupBackground: '#2c2c2c',
      overlayBackgroundColor: 'rgba(0, 0, 0, 0.7)',
      
      // Tag
      tagTextColor: '#e5e5e5',
      
      // Empty
      emptyTextColor: '#b0b0b0',
      
      // Loading
      loadingTextColor: '#b0b0b0',
      
      // ActionSheet
      actionSheetItemBackground: '#2c2c2c',
      actionSheetItemTextColor: '#e5e5e5',
      actionSheetCancelTextColor: '#b0b0b0',
      actionSheetSubnameColor: '#64B5F6',
      
      // Other
      activeColor: '#3a3a3a',
      cardBackground: '#2c2c2c',
    }
  } else {
    // Light theme - use Vant defaults
    return undefined
  }
})
</script>

<template>
  <ConfigProvider :theme-vars="themeVars">
    <div id="app">
      <RouterView />
    </div>
  </ConfigProvider>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

body {
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 只针对 Vue 应用的根容器，不影响 Toast/Dialog 等挂载到 body 的元素 */
#app {
  width: 100%;
  height: 100%;
}

/* Ensure ConfigProvider doesn't break height chain */
/* ConfigProvider wraps content in a div, we need to ensure it inherits full height */
.van-config-provider {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.van-config-provider > * {
  flex: 1;
  min-height: 0;
}
</style>
