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
        <RouterLink to="/" class="menu-item" :class="{ active: isActive('/') }">
          <i class="ri-dashboard-3-line"></i>
          <span>{{ t('menu.dashboard') }}</span>
        </RouterLink>
        <RouterLink v-if="canViewSuppliers" to="/suppliers" class="menu-item" :class="{ active: isActive('/suppliers') }">
          <i class="ri-building-line"></i>
          <span>{{ t('menu.suppliers') }}</span>
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
        <RouterLink v-if="canViewLogistics" to="/logistics" class="menu-item" :class="{ active: isActive('/logistics') }">
          <i class="ri-truck-line"></i>
          <span>{{ t('menu.logistics') }}</span>
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
        <RouterLink to="/settings" class="user-avatar" :title="user?.username">
          <img v-if="user?.avatar" :src="user.avatar" alt="avatar" />
          <i v-else-if="!user?.username" class="ri-user-line"></i>
          <span v-else class="avatar-text">{{ userInitial }}</span>
        </RouterLink>
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
import { ref, watch, computed } from 'vue'
import { useRoute, useRouter, RouterLink, RouterView } from 'vue-router'
import { useI18n } from 'vue-i18n'
import type { User } from '@/types'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const isDark = ref(localStorage.getItem('darkMode') === 'true')

// 从 localStorage 获取用户信息
const user = computed<User | null>(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch {
      return null
    }
  }
  return null
})

// 检查用户是否有权限访问特定模块
const hasPermission = (module: string, action: 'read' | 'write' = 'read') => {
  if (!user.value) return false

  // admin 拥有所有权限
  if (user.value.role === 'admin') return true

  // 检查 permissions 数组中是否有对应权限
  const permKey = `${module}:${action}`
  if (user.value.permissions?.includes(permKey)) return true

  // 兼容旧的 role 字段权限检查
  const rolePerms: Record<string, string[]> = {
    admin: ['supplier:read', 'supplier:write', 'inventory:read', 'inventory:write', 'procurement:read', 'procurement:write', 'sales:read', 'sales:write', 'logistics:read', 'logistics:write', 'user:read', 'user:write', 'role:read', 'role:write', 'permission:read', 'permission:write'],
    manager: ['supplier:read', 'inventory:read', 'inventory:write', 'procurement:read', 'procurement:write', 'sales:read', 'sales:write', 'logistics:read', 'logistics:write', 'user:read'],
    operator: ['dashboard:read', 'inventory:read', 'procurement:read', 'sales:read']
  }

  const perms = rolePerms[user.value.role] || []
  return perms.includes(permKey) || perms.includes(`${module}:read`)
}

// 菜单项权限控制
const canViewSuppliers = computed(() => hasPermission('supplier', 'read'))
const canViewProcurement = computed(() => hasPermission('procurement', 'read'))
const canViewInventory = computed(() => hasPermission('inventory', 'read'))
const canViewSales = computed(() => hasPermission('sales', 'read'))
const canViewLogistics = computed(() => hasPermission('logistics', 'read'))
const canViewUsers = computed(() => hasPermission('user', 'read'))
const canViewRoles = computed(() => hasPermission('role', 'read'))
const canViewPermissions = computed(() => hasPermission('permission', 'read'))

const isLoginPage = computed(() => route.path === '/login')
const isActive = (path: string) => route.path === path

// 获取用户名的首字母
const userInitial = computed(() => {
  if (!user.value?.username) return '?'
  return user.value.username.charAt(0).toUpperCase()
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}

watch(isDark, (val) => {
  document.documentElement.classList.toggle('dark-mode', val)
}, { immediate: true })
</script>