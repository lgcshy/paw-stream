// Theme Store

import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

// Theme mode types
export type ThemeMode = 'light' | 'dark' | 'auto'
export type EffectiveTheme = 'light' | 'dark'

const STORAGE_KEY = 'theme_mode'

/**
 * Detect system theme preference
 */
function getSystemTheme(): EffectiveTheme {
  if (typeof window === 'undefined') return 'light'
  
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)')
  return prefersDark.matches ? 'dark' : 'light'
}

/**
 * Load saved theme mode from localStorage
 */
function loadThemeMode(): ThemeMode {
  if (typeof window === 'undefined') return 'auto'
  
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'light' || saved === 'dark' || saved === 'auto') {
    return saved
  }
  return 'auto'
}

export const useThemeStore = defineStore('theme', () => {
  // State
  const mode = ref<ThemeMode>(loadThemeMode())
  const systemTheme = ref<EffectiveTheme>(getSystemTheme())

  // Getters
  const effectiveTheme = computed<EffectiveTheme>(() => {
    if (mode.value === 'auto') {
      return systemTheme.value
    }
    return mode.value
  })

  const isDark = computed(() => effectiveTheme.value === 'dark')

  // Actions
  function setMode(newMode: ThemeMode) {
    mode.value = newMode
    localStorage.setItem(STORAGE_KEY, newMode)
  }

  function setSystemTheme(theme: EffectiveTheme) {
    systemTheme.value = theme
  }

  function toggleTheme() {
    if (mode.value === 'light') {
      setMode('dark')
    } else if (mode.value === 'dark') {
      setMode('light')
    } else {
      // If auto, toggle to opposite of current effective theme
      setMode(effectiveTheme.value === 'light' ? 'dark' : 'light')
    }
  }

  // Initialize system theme listener
  function initSystemThemeListener() {
    if (typeof window === 'undefined') return

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    
    // Update system theme when it changes
    const listener = (e: MediaQueryListEvent) => {
      setSystemTheme(e.matches ? 'dark' : 'light')
    }

    // Modern browsers
    if (mediaQuery.addEventListener) {
      mediaQuery.addEventListener('change', listener)
    } else {
      // Fallback for older browsers
      mediaQuery.addListener(listener)
    }

    // Return cleanup function
    return () => {
      if (mediaQuery.removeEventListener) {
        mediaQuery.removeEventListener('change', listener)
      } else {
        mediaQuery.removeListener(listener)
      }
    }
  }

  // Watch effectiveTheme and apply to document
  watch(
    effectiveTheme,
    (theme) => {
      if (typeof document !== 'undefined') {
        if (theme === 'dark') {
          document.documentElement.classList.add('dark-theme')
        } else {
          document.documentElement.classList.remove('dark-theme')
        }
      }
    },
    { immediate: true }
  )

  return {
    // State
    mode,
    systemTheme,
    // Getters
    effectiveTheme,
    isDark,
    // Actions
    setMode,
    setSystemTheme,
    toggleTheme,
    initSystemThemeListener,
  }
})
