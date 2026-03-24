import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export interface Tab {
  path: string
  name: string
  title: string
  icon?: string
  closable: boolean
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<Tab[]>([])
  const activeTab = ref<string>('')

  // 固定的首页标签
  const homeTab: Tab = {
    path: '/',
    name: 'dashboard',
    title: '仪表盘',
    icon: 'ri-dashboard-3-line',
    closable: false
  }

  // 确保已初始化
  function ensureInit() {
    if (tabs.value.length === 0) {
      tabs.value = [homeTab]
      activeTab.value = '/'
    }
  }

  // 初始化
  function init() {
    ensureInit()
  }

  // 添加标签
  function addTab(route: RouteLocationNormalized) {
    ensureInit()

    const path = route.path
    const name = route.name as string
    const title = (route.meta?.title as string) || name || path
    const icon = route.meta?.icon as string | undefined

    // 不添加登录页
    if (path === '/login') return

    // 检查是否已存在
    const exists = tabs.value.find(t => t.path === path)
    if (!exists) {
      tabs.value.push({
        path,
        name,
        title,
        icon,
        closable: path !== '/'
      })
    }
    activeTab.value = path
  }

  // 关闭标签
  function closeTab(path: string) {
    const index = tabs.value.findIndex(t => t.path === path)
    if (index === -1) return

    const tab = tabs.value[index]
    if (!tab.closable) return

    tabs.value.splice(index, 1)

    // 如果关闭的是当前激活的标签，切换到前一个或后一个
    if (activeTab.value === path) {
      const newIndex = Math.min(index, tabs.value.length - 1)
      activeTab.value = tabs.value[newIndex]?.path || '/'
    }

    return activeTab.value
  }

  // 关闭其他标签
  function closeOtherTabs(path: string) {
    tabs.value = tabs.value.filter(t => t.path === path || !t.closable)
    activeTab.value = path
  }

  // 关闭左侧标签
  function closeLeftTabs(path: string) {
    const index = tabs.value.findIndex(t => t.path === path)
    if (index === -1) return

    tabs.value = tabs.value.filter((t, i) => i >= index || !t.closable)
    activeTab.value = path
  }

  // 关闭右侧标签
  function closeRightTabs(path: string) {
    const index = tabs.value.findIndex(t => t.path === path)
    if (index === -1) return

    tabs.value = tabs.value.filter((t, i) => i <= index || !t.closable)
    activeTab.value = path
  }

  // 关闭所有可关闭的标签
  function closeAllTabs() {
    tabs.value = [homeTab]
    activeTab.value = '/'
    return '/'
  }

  // 设置激活标签
  function setActiveTab(path: string) {
    activeTab.value = path
  }

  // 获取当前激活标签
  const currentTab = computed(() => tabs.value.find(t => t.path === activeTab.value))

  return {
    tabs,
    activeTab,
    currentTab,
    init,
    addTab,
    closeTab,
    closeOtherTabs,
    closeLeftTabs,
    closeRightTabs,
    closeAllTabs,
    setActiveTab
  }
})