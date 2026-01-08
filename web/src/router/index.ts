import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/streams',
    name: 'StreamList',
    component: () => import('@/views/StreamListView.vue'),
    meta: { requiresAuth: true, showBottomNav: true },
  },
  {
    path: '/stream/:id',
    name: 'StreamPlayer',
    component: () => import('@/views/StreamPlayerView.vue'),
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/devices',
    name: 'DeviceList',
    component: () => import('@/views/DeviceListView.vue'),
    meta: { requiresAuth: true, showBottomNav: true },
  },
  {
    path: '/devices/new',
    name: 'DeviceCreate',
    component: () => import('@/views/DeviceFormView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/devices/:id',
    name: 'DeviceDetail',
    component: () => import('@/views/DeviceDetailView.vue'),
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/devices/:id/edit',
    name: 'DeviceEdit',
    component: () => import('@/views/DeviceFormView.vue'),
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true, showBottomNav: true },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation guards
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false // Default to true

  // Try to load auth state if not already loaded
  if (!authStore.isAuthenticated && localStorage.getItem('auth_token')) {
    try {
      await authStore.loadToken()
    } catch (error) {
      // Token invalid, will redirect to login
    }
  }

  // Check authentication
  if (requiresAuth && !authStore.isAuthenticated) {
    // Redirect to login, save the original destination
    next({
      name: 'Login',
      query: { redirect: to.fullPath },
    })
  } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
    // Already logged in, redirect to streams
    next({ name: 'StreamList' })
  } else {
    // Proceed normally
    next()
  }
})

export default router
