import { ref } from 'vue'

// 用户权限信息
const currentUserPermissions = ref<string[]>([])
const currentUserRoles = ref<string[]>([])
const currentUserRole = ref<string>('')
const currentUserDepartment = ref<string>('')

// 从 localStorage 恢复用户信息
function initUserPermissions() {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      currentUserRole.value = user.role || ''
      currentUserDepartment.value = user.department || ''

      const permissions: string[] = []

      // 从 user.permissions 提取 - 支持字符串数组格式（后端新格式）
      if (user.permissions) {
        const permList = Array.isArray(user.permissions) ? user.permissions : []
        permList.forEach((p: any) => {
          // 后端返回的权限可能是字符串直接是权限名
          if (typeof p === 'string') {
            permissions.push(p)
          } else if (p && typeof p === 'object') {
            // 旧格式，权限对象有 name 字段
            permissions.push(p.name || '')
          }
        })
      }

      // 从 user.roles 提取 - 支持字符串数组格式（后端新格式）
      if (user.roles) {
        const rolesList = Array.isArray(user.roles) ? user.roles : []
        rolesList.forEach((role: any) => {
          // 新格式：role 是字符串
          if (typeof role === 'string') {
            if (role && !currentUserRoles.value.includes(role)) {
              currentUserRoles.value.push(role)
            }
          } else if (role && typeof role === 'object') {
            // 旧格式：role 是对象，有 name 字段
            const roleName = role.name || role.code || ''
            if (roleName && !currentUserRoles.value.includes(roleName)) {
              currentUserRoles.value.push(roleName)
            }
            // 从角色对象中提取权限（旧格式）
            if (role.permissions) {
              const perms = Array.isArray(role.permissions) ? role.permissions : []
              perms.forEach((p: any) => {
                if (p && typeof p === 'object') {
                  const permName = p.name || ''
                  if (permName && !permissions.includes(permName)) {
                    permissions.push(permName)
                  }
                } else if (typeof p === 'string' && !permissions.includes(p)) {
                  permissions.push(p)
                }
              })
            }
          }
        })
      }

      currentUserPermissions.value = permissions.filter(p => p)
    } catch (e) {
      console.error('Failed to parse user info:', e)
    }
  }
}

// 设置用户权限信息
function setUserPermissions(permissions: string[], roles: string[], role: string, department: string) {
  currentUserPermissions.value = permissions
  currentUserRoles.value = roles
  currentUserRole.value = role
  currentUserDepartment.value = department
}

// 检查用户是否有特定权限
function hasPermission(permission: string): boolean {
  // admin 拥有所有权限
  if (currentUserRole.value === 'admin') {
    return true
  }
  return currentUserPermissions.value.includes(permission)
}

// 检查用户角色
function hasRole(role: 'admin' | 'manager' | 'operator'): boolean {
  return currentUserRole.value === role
}

// 检查是否有创建权限
function canCreate(): boolean {
  // admin 和 manager 可以创建
  return currentUserRole.value === 'admin' || currentUserRole.value === 'manager'
}

// 检查是否有编辑权限
function canEdit(_creatorId?: number, creatorDepartment?: string): boolean {
  // admin 可以编辑所有
  if (currentUserRole.value === 'admin') {
    return true
  }
  // manager 可以编辑本部门数据
  if (currentUserRole.value === 'manager') {
    return creatorDepartment === currentUserDepartment.value
  }
  return false
}

// 检查是否有删除权限
function canDelete(): boolean {
  // 只有 admin 可以删除
  return currentUserRole.value === 'admin'
}

// 检查是否是超级管理员
function isAdmin(): boolean {
  return currentUserRole.value === 'admin'
}

// 检查是否是管理员或经理
function isAdminOrManager(): boolean {
  return currentUserRole.value === 'admin' || currentUserRole.value === 'manager'
}

// 清除用户权限信息（登出时使用）
function clearUserPermissions() {
  currentUserPermissions.value = []
  currentUserRoles.value = []
  currentUserRole.value = ''
  currentUserDepartment.value = ''
}

export function usePermission() {
  return {
    currentUserPermissions,
    currentUserRoles,
    currentUserRole,
    currentUserDepartment,
    initUserPermissions,
    setUserPermissions,
    hasPermission,
    hasRole,
    canCreate,
    canEdit,
    canDelete,
    isAdmin,
    isAdminOrManager,
    clearUserPermissions
  }
}