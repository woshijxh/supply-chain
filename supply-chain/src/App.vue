<template>
  <div class="app-shell" :class="{ 'dark-mode': isDark, 'login-page': isLoginPage }">
    <aside v-if="!isLoginPage" class="sidebar">
      <div class="brand">
        <div class="brand-mark">
          <span class="brand-initials">SC</span>
          <svg class="brand-pulse" viewBox="0 0 32 16">
            <path d="M1 8H7L10 3L14 13L18 8H31" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="brand-text">
          <div class="brand-title">{{ t('app.title') }}</div>
          <div class="brand-sub">{{ t('app.subtitle') }}</div>
        </div>
      </div>
      <nav class="menu">
        <RouterLink v-if="canViewDashboard" to="/" class="menu-item" :class="{ active: isActive('/') }">
          <i class="ri-dashboard-3-line"></i>
          <span>{{ t('menu.dashboard') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewSuppliers" to="/suppliers" class="menu-item" :class="{ active: isActive('/suppliers') }">
          <i class="ri-building-line"></i>
          <span>{{ t('menu.suppliers') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewCustomers" to="/customers" class="menu-item" :class="{ active: isActive('/customers') }">
          <i class="ri-user-heart-line"></i>
          <span>客户管理</span>
        </RouterLink>
        <RouterLink v-if="canViewProducts" to="/products" class="menu-item" :class="{ active: isActive('/products') }">
          <i class="ri-box-3-line"></i>
          <span>产品管理</span>
        </RouterLink>
        <RouterLink v-if="canViewProcurement" to="/procurement" class="menu-item" :class="{ active: isActive('/procurement') }">
          <i class="ri-shopping-cart-line"></i>
          <span>{{ t('menu.procurement') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewInventory" to="/inventory" class="menu-item" :class="{ active: isActive('/inventory') }">
          <i class="ri-store-2-line"></i>
          <span>{{ t('menu.inventory') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewSales" to="/sales" class="menu-item" :class="{ active: isActive('/sales') }">
          <i class="ri-file-list-3-line"></i>
          <span>{{ t('menu.sales') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewSales" to="/returns" class="menu-item" :class="{ active: isActive('/returns') }">
          <i class="ri-arrow-go-back-line"></i>
          <span>退货管理</span>
        </RouterLink>
        <RouterLink v-if="canViewLogistics" to="/logistics" class="menu-item" :class="{ active: isActive('/logistics') }">
          <i class="ri-truck-line"></i>
          <span>{{ t('menu.logistics') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewTrace" to="/trace" class="menu-item" :class="{ active: isActive('/trace') }">
          <i class="ri-search-eye-line"></i>
          <span>商品追溯</span>
        </RouterLink>
        <div v-if="canViewUsers" class="menu-divider"></div>
        <RouterLink v-if="canViewUsers" to="/users" class="menu-item" :class="{ active: isActive('/users') }">
          <i class="ri-user-settings-line"></i>
          <span>用户管理</span>
        </RouterLink>
        <RouterLink v-if="canViewRoles" to="/roles" class="menu-item" :class="{ active: isActive('/roles') }">
          <i class="ri-shield-line"></i>
          <span>角色管理</span>
        </RouterLink>
        <RouterLink v-if="canViewPermissions" to="/permissions" class="menu-item" :class="{ active: isActive('/permissions') }">
          <i class="ri-lock-line"></i>
          <span>权限管理</span>
        </RouterLink>
      </nav>
      <div class="sidebar-footer">
        <div class="user-info">
          <RouterLink to="/settings" class="user-avatar" :title="userStore.user?.username">
            <img v-if="userStore.user?.avatar" :src="userStore.user.avatar" alt="avatar" />
            <i v-else-if="!userStore.user?.username" class="ri-user-line"></i>
            <span v-else class="avatar-text">{{ userStore.userInitial }}</span>
          </RouterLink>
          <span class="user-name">{{ userStore.user?.username || '用户' }}</span>
        </div>
        <button class="logout-btn" @click="handleLogout">
          <span>{{ t('settings.logout') }}</span>
        </button>
      </div>
    </aside>
    <main class="main-content" :class="{ 'full-screen': isLoginPage }">
      <RouterView />
    </main>
    <Toast />
    <ConfirmDialog />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter, RouterLink, RouterView } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { setWatermark, removeWatermark } from '@/utils/watermark'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const userStore = useUserStore()

const isDark = ref(localStorage.getItem('darkMode') === 'true')

// 初始化用户状态
onMounted(() => {
  userStore.init()
})

// 菜单项权限控制
const canViewDashboard = computed(() => userStore.hasPermission('dashboard:read'))
const canViewSuppliers = computed(() => userStore.hasPermission('supplier:read'))
const canViewCustomers = computed(() => userStore.hasPermission('customer:read'))
const canViewProducts = computed(() => userStore.hasPermission('product:read'))
const canViewProcurement = computed(() => userStore.hasPermission('procurement:read'))
const canViewInventory = computed(() => userStore.hasPermission('inventory:read'))
const canViewSales = computed(() => userStore.hasPermission('sales:read'))
const canViewLogistics = computed(() => userStore.hasPermission('logistics:read'))
const canViewTrace = computed(() => userStore.hasPermission('trace:read'))
const canViewUsers = computed(() => userStore.hasPermission('user:read'))
const canViewRoles = computed(() => userStore.hasPermission('role:read'))
const canViewPermissions = computed(() => userStore.hasPermission('permission:read'))

const isLoginPage = computed(() => route.path === '/login')
const isActive = (path: string) => route.path === path

const handleLogout = () => {
  userStore.clearUser()
  removeWatermark()
  router.push('/login')
}

watch(isDark, (val) => {
  document.documentElement.classList.toggle('dark-mode', val)
}, { immediate: true })

// 水印功能：登录后显示用户名水印
watch(() => userStore.user, (newUser) => {
  if (newUser?.username) {
    setWatermark({
      text: newUser.username,
      fontSize: 14,
      color: '#666',
      gap: 120,
      opacity: 0.12,
      rotate: -20
    })
  } else {
    removeWatermark()
  }
}, { immediate: true })

onUnmounted(() => {
  removeWatermark()
})
</script>