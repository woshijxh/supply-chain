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
        <!-- 仪表盘 -->
        <RouterLink v-if="canViewDashboard" to="/" class="menu-item" :class="{ active: isActive('/') }">
          <i class="ri-dashboard-3-line"></i>
          <span>{{ t('menu.dashboard') }}</span>
        </RouterLink>

        <!-- 基础数据 -->
        <div v-if="canViewSuppliers || canViewCustomers || canViewProducts" class="menu-group">
          <div class="menu-group-title" @click="toggleGroup('base')">
            <i class="ri-database-2-line"></i>
            <span>基础数据</span>
            <i class="ri-arrow-down-s-line menu-arrow" :class="{ expanded: expandedGroups.base }"></i>
          </div>
          <div class="menu-group-items" :class="{ expanded: expandedGroups.base }">
            <RouterLink v-if="canViewSuppliers" to="/suppliers" class="menu-sub-item" :class="{ active: isActive('/suppliers') }">
              <i class="ri-building-line"></i>
              <span>{{ t('menu.suppliers') }}</span>
            </RouterLink>
            <RouterLink v-if="canViewCustomers" to="/customers" class="menu-sub-item" :class="{ active: isActive('/customers') }">
              <i class="ri-user-heart-line"></i>
              <span>客户管理</span>
            </RouterLink>
            <RouterLink v-if="canViewProducts" to="/products" class="menu-sub-item" :class="{ active: isActive('/products') }">
              <i class="ri-box-3-line"></i>
              <span>产品管理</span>
            </RouterLink>
          </div>
        </div>

        <!-- 采购管理 -->
        <div v-if="canViewProcurement" class="menu-group">
          <div class="menu-group-title" @click="toggleGroup('procurement')">
            <i class="ri-shopping-cart-line"></i>
            <span>采购管理</span>
            <i class="ri-arrow-down-s-line menu-arrow" :class="{ expanded: expandedGroups.procurement }"></i>
          </div>
          <div class="menu-group-items" :class="{ expanded: expandedGroups.procurement }">
            <RouterLink to="/procurement" class="menu-sub-item" :class="{ active: isActive('/procurement') }">
              <i class="ri-file-list-3-line"></i>
              <span>采购订单</span>
            </RouterLink>
          </div>
        </div>

        <!-- 库存管理 -->
        <div v-if="canViewInventory" class="menu-group">
          <div class="menu-group-title" @click="toggleGroup('inventory')">
            <i class="ri-store-2-line"></i>
            <span>库存管理</span>
            <i class="ri-arrow-down-s-line menu-arrow" :class="{ expanded: expandedGroups.inventory }"></i>
          </div>
          <div class="menu-group-items" :class="{ expanded: expandedGroups.inventory }">
            <RouterLink to="/inventory" class="menu-sub-item" :class="{ active: isActive('/inventory') }">
              <i class="ri-stack-line"></i>
              <span>库存查询</span>
            </RouterLink>
            <RouterLink to="/inventory/logs" class="menu-sub-item" :class="{ active: isActive('/inventory/logs') }">
              <i class="ri-file-list-3-line"></i>
              <span>库存流水</span>
            </RouterLink>
          </div>
        </div>

        <!-- 销售管理 -->
        <div v-if="canViewSales || canViewLogistics" class="menu-group">
          <div class="menu-group-title" @click="toggleGroup('sales')">
            <i class="ri-handbag-line"></i>
            <span>销售管理</span>
            <i class="ri-arrow-down-s-line menu-arrow" :class="{ expanded: expandedGroups.sales }"></i>
          </div>
          <div class="menu-group-items" :class="{ expanded: expandedGroups.sales }">
            <RouterLink v-if="canViewSales" to="/sales" class="menu-sub-item" :class="{ active: isActive('/sales') }">
              <i class="ri-file-list-3-line"></i>
              <span>销售订单</span>
            </RouterLink>
            <RouterLink v-if="canViewSales" to="/returns" class="menu-sub-item" :class="{ active: isActive('/returns') }">
              <i class="ri-arrow-go-back-line"></i>
              <span>退货管理</span>
            </RouterLink>
            <RouterLink v-if="canViewLogistics" to="/logistics" class="menu-sub-item" :class="{ active: isActive('/logistics') }">
              <i class="ri-truck-line"></i>
              <span>物流跟踪</span>
            </RouterLink>
          </div>
        </div>

        <!-- 商品追溯 -->
        <RouterLink v-if="canViewTrace" to="/trace" class="menu-item" :class="{ active: isActive('/trace') }">
          <i class="ri-search-eye-line"></i>
          <span>商品追溯</span>
        </RouterLink>

        <!-- 系统管理 -->
        <div v-if="canViewUsers || canViewRoles || canViewPermissions" class="menu-group">
          <div class="menu-group-title" @click="toggleGroup('system')">
            <i class="ri-settings-4-line"></i>
            <span>系统管理</span>
            <i class="ri-arrow-down-s-line menu-arrow" :class="{ expanded: expandedGroups.system }"></i>
          </div>
          <div class="menu-group-items" :class="{ expanded: expandedGroups.system }">
            <RouterLink v-if="canViewUsers" to="/users" class="menu-sub-item" :class="{ active: isActive('/users') }">
              <i class="ri-user-settings-line"></i>
              <span>用户管理</span>
            </RouterLink>
            <RouterLink v-if="canViewRoles" to="/roles" class="menu-sub-item" :class="{ active: isActive('/roles') }">
              <i class="ri-shield-line"></i>
              <span>角色管理</span>
            </RouterLink>
            <RouterLink v-if="canViewPermissions" to="/permissions" class="menu-sub-item" :class="{ active: isActive('/permissions') }">
              <i class="ri-lock-line"></i>
              <span>权限管理</span>
            </RouterLink>
          </div>
        </div>
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
      <!-- 标签栏 -->
      <div v-if="!isLoginPage" class="tabs-bar">
        <div class="tabs-wrapper">
          <div
            v-for="tab in tabsStore.tabs"
            :key="tab.path"
            class="tab-item"
            :class="{ active: tabsStore.activeTab === tab.path }"
            @click="handleTabClick(tab.path)"
            @contextmenu.prevent="showContextMenu($event, tab)"
          >
            <i v-if="tab.icon" :class="tab.icon" class="tab-icon"></i>
            <span class="tab-title">{{ tab.title }}</span>
            <i
              v-if="tab.closable"
              class="ri-close-line tab-close"
              @click.stop="handleTabClose(tab.path)"
            ></i>
          </div>
        </div>
      </div>
      <!-- 右键菜单 -->
      <div
        v-if="contextMenuVisible"
        class="context-menu"
        :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }"
      >
        <div class="context-menu-item" @click="handleCloseOther">关闭其他</div>
        <div class="context-menu-item" @click="handleCloseRight">关闭右侧</div>
        <div class="context-menu-item" @click="handleCloseAll">关闭所有</div>
      </div>
      <!-- 内容区 -->
      <div class="page-content" :class="{ 'no-tabs': isLoginPage }">
        <RouterView v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" :key="route.path" />
          </keep-alive>
        </RouterView>
      </div>
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
import { useTabsStore, type Tab } from '@/stores/tabs'
import { setWatermark, removeWatermark } from '@/utils/watermark'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const userStore = useUserStore()
const tabsStore = useTabsStore()

// 确保 tabs 初始化
tabsStore.init()

const isDark = ref(localStorage.getItem('darkMode') === 'true')

// 右键菜单状态
const contextMenuVisible = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const selectedTab = ref<Tab | null>(null)

// 初始化用户状态
onMounted(() => {
  userStore.init()
  // 添加初始路由
  if (route.meta.requiresAuth) {
    tabsStore.addTab(route)
  }
})

// 使用路由守卫监听路由变化
router.afterEach((to) => {
  if (to.meta.requiresAuth) {
    tabsStore.addTab(to)
  }
})

// 点击其他地方关闭右键菜单
const handleClickOutside = () => {
  contextMenuVisible.value = false
}

watch(contextMenuVisible, (visible) => {
  if (visible) {
    document.addEventListener('click', handleClickOutside)
  } else {
    document.removeEventListener('click', handleClickOutside)
  }
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

// 菜单分组展开状态
const expandedGroups = ref<Record<string, boolean>>({
  base: true,
  procurement: true,
  inventory: true,
  sales: true,
  system: true
})

// 切换分组展开状态
const toggleGroup = (group: string) => {
  expandedGroups.value[group] = !expandedGroups.value[group]
}

const isLoginPage = computed(() => route.path === '/login')
const isActive = (path: string) => {
  if (path === '/') {
    return route.path === '/'
  }
  return route.path === path || route.path.startsWith(path + '/')
}

// 标签点击
const handleTabClick = (path: string) => {
  tabsStore.setActiveTab(path)
  router.push(path)
}

// 标签关闭
const handleTabClose = (path: string) => {
  const newPath = tabsStore.closeTab(path)
  if (newPath && newPath !== route.path) {
    router.push(newPath)
  }
}

// 显示右键菜单
const showContextMenu = (event: MouseEvent, tab: Tab) => {
  event.preventDefault()
  selectedTab.value = tab
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  contextMenuVisible.value = true
}

// 关闭其他
const handleCloseOther = () => {
  if (selectedTab.value) {
    tabsStore.closeOtherTabs(selectedTab.value.path)
    if (route.path !== selectedTab.value.path) {
      router.push(selectedTab.value.path)
    }
  }
  contextMenuVisible.value = false
}

// 关闭右侧
const handleCloseRight = () => {
  if (selectedTab.value) {
    tabsStore.closeRightTabs(selectedTab.value.path)
  }
  contextMenuVisible.value = false
}

// 关闭所有
const handleCloseAll = () => {
  const path = tabsStore.closeAllTabs()
  if (route.path !== path) {
    router.push(path)
  }
  contextMenuVisible.value = false
}

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

<style scoped>
.tabs-bar {
  display: flex;
  align-items: center;
  background: var(--panel);
  padding: 8px 16px 0;
  height: 48px;
  border-bottom: 1px solid var(--border);
}

.tabs-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  flex: 1;
  overflow-x: auto;
  scrollbar-width: none;
}

.tabs-wrapper::-webkit-scrollbar {
  display: none;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  height: 40px;
  background: var(--surface-ground);
  border: 1px solid var(--border);
  border-bottom: none;
  border-radius: 8px 8px 0 0;
  cursor: pointer;
  white-space: nowrap;
  font-size: 13px;
  color: var(--text-muted);
  transition: all 0.2s ease;
  user-select: none;
  position: relative;
  margin-bottom: -1px;
}

.tab-item:hover {
  background: var(--bg);
  color: var(--text);
}

.tab-item.active {
  background: var(--bg);
  color: var(--primary);
  font-weight: 500;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--bg);
}

.tab-title {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-icon {
  font-size: 15px;
  opacity: 0.7;
}

.tab-item.active .tab-icon {
  opacity: 1;
}

.tab-close {
  font-size: 14px;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  opacity: 0.6;
}

.tab-close:hover {
  opacity: 1;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.context-menu {
  position: fixed;
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  padding: 4px 0;
  min-width: 120px;
  z-index: 1000;
}

.context-menu-item {
  padding: 8px 16px;
  font-size: 13px;
  color: var(--text);
  cursor: pointer;
  transition: background 0.2s ease;
}

.context-menu-item:hover {
  background: var(--primary-soft);
  color: var(--primary);
}

.page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.page-content.no-tabs {
  padding: 0;
}
</style>