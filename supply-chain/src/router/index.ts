import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { User } from '@/types'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/LoginPage.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { title: 'dashboard', requiresAuth: true, permission: 'dashboard:read' }
  },
  {
    path: '/suppliers',
    name: 'suppliers',
    component: () => import('@/pages/SupplierPage.vue'),
    meta: { title: 'suppliers', requiresAuth: true, permission: 'supplier:read' }
  },
  {
    path: '/procurement',
    name: 'procurement',
    component: () => import('@/pages/ProcurementPage.vue'),
    meta: { title: 'procurement', requiresAuth: true, permission: 'procurement:read' }
  },
  {
    path: '/inventory',
    name: 'inventory',
    component: () => import('@/pages/InventoryPage.vue'),
    meta: { title: 'inventory', requiresAuth: true, permission: 'inventory:read' }
  },
  {
    path: '/sales',
    name: 'sales',
    component: () => import('@/pages/SalesPage.vue'),
    meta: { title: 'sales', requiresAuth: true, permission: 'sales:read' }
  },
  {
    path: '/logistics',
    name: 'logistics',
    component: () => import('@/pages/LogisticsPage.vue'),
    meta: { title: 'logistics', requiresAuth: true, permission: 'logistics:read' }
  },
  {
    path: '/users',
    name: 'users',
    component: () => import('@/pages/UserPage.vue'),
    meta: { title: '用户管理', requiresAuth: true, permission: 'user:read' }
  },
  {
    path: '/roles',
    name: 'roles',
    component: () => import('@/pages/RolePage.vue'),
    meta: { title: '角色管理', requiresAuth: true, permission: 'role:read' }
  },
  {
    path: '/permissions',
    name: 'permissions',
    component: () => import('@/pages/PermissionPage.vue'),
    meta: { title: '权限管理', requiresAuth: true, permission: 'permission:read' }
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/SettingsPage.vue'),
    meta: { title: '用户设置', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 检查用户是否有权限访问特定模块
const hasPermission = (user: User | null, permission?: string): boolean => {
  if (!user) return false

  // admin 拥有所有权限
  if (user.role === 'admin') return true

  // 如果没有设置权限要求，则允许访问
  if (!permission) return true

  // 检查 permissions 数组中是否有对应权限
  if (user.permissions?.includes(permission)) return true

  // 兼容旧的 role 字段权限检查
  const rolePerms: Record<string, string[]> = {
    admin: ['supplier:read', 'supplier:write', 'inventory:read', 'inventory:write', 'procurement:read', 'procurement:write', 'sales:read', 'sales:write', 'logistics:read', 'logistics:write', 'user:read', 'user:write', 'role:read', 'role:write', 'permission:read', 'permission:write', 'dashboard:read'],
    manager: ['supplier:read', 'inventory:read', 'inventory:write', 'procurement:read', 'procurement:write', 'sales:read', 'sales:write', 'logistics:read', 'logistics:write', 'user:read', 'dashboard:read'],
    operator: ['dashboard:read', 'inventory:read', 'procurement:read', 'sales:read']
  }

  const perms = rolePerms[user.role] || []
  return perms.includes(permission)
}

// 获取用户信息
const getUser = (): User | null => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch {
      return null
    }
  }
  return null
}

router.beforeEach((to, _from, next) => {
  const title = to.meta.title as string
  document.title = title ? `${title} - 供应链管理系统` : '供应链管理系统'

  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else if (to.meta.requiresAuth && token) {
    // 检查权限
    const user = getUser()
    const permission = to.meta.permission as string | undefined
    if (hasPermission(user, permission)) {
      next()
    } else {
      // 没有权限，跳转到首页或显示提示
      next('/')
    }
  } else {
    next()
  }
})

export default router