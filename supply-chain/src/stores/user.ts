import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  // 初始化：从 localStorage 读取
  function init() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')

    if (storedToken) {
      token.value = storedToken
    }

    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        user.value = null
      }
    }
  }

  // 设置用户信息（登录成功后调用）
  function setUser(userData: User, tokenValue: string) {
    user.value = userData
    token.value = tokenValue
    localStorage.setItem('token', tokenValue)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  // 清除用户信息（退出登录时调用）
  function clearUser() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 是否已登录
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // 检查是否有权限
  function hasPermission(permCode: string): boolean {
    if (!user.value) return false

    // admin 拥有所有权限
    if (user.value.role === 'admin') return true

    // 检查 permissions 数组中是否有对应权限
    // 支持两种格式：p.code 或 p.name（字符串时直接比较）
    if (user.value.permissions?.some(p => {
      if (typeof p === 'string') return p === permCode
      return p.code === permCode || p.name === permCode
    })) return true

    return false
  }

  // 获取用户名首字母
  const userInitial = computed(() => {
    if (!user.value?.username) return '?'
    return user.value.username.charAt(0).toUpperCase()
  })

  return {
    user,
    token,
    isLoggedIn,
    userInitial,
    init,
    setUser,
    clearUser,
    hasPermission
  }
})