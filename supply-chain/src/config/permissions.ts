/**
 * 权限码定义 - 与后端 permission 表的 name 字段对应
 */

// 供应商权限
export const SUPPLIER_PERMISSIONS = {
  READ: 'supplier:read',
  WRITE: 'supplier:write',
} as const

// 库存权限
export const INVENTORY_PERMISSIONS = {
  READ: 'inventory:read',
  WRITE: 'inventory:write',
} as const

// 采购权限
export const PROCUREMENT_PERMISSIONS = {
  READ: 'procurement:read',
  WRITE: 'procurement:write',
} as const

// 销售权限
export const SALES_PERMISSIONS = {
  READ: 'sales:read',
  WRITE: 'sales:write',
} as const

// 物流权限
export const LOGISTICS_PERMISSIONS = {
  READ: 'logistics:read',
  WRITE: 'logistics:write',
} as const

// 用户管理权限
export const USER_PERMISSIONS = {
  READ: 'user:read',
  WRITE: 'user:write',
} as const

// 角色管理权限
export const ROLE_PERMISSIONS = {
  READ: 'role:read',
  WRITE: 'role:write',
} as const

// 所有权限映射
export const PERMISSIONS = {
  ...SUPPLIER_PERMISSIONS,
  ...INVENTORY_PERMISSIONS,
  ...PROCUREMENT_PERMISSIONS,
  ...SALES_PERMISSIONS,
  ...LOGISTICS_PERMISSIONS,
  ...USER_PERMISSIONS,
  ...ROLE_PERMISSIONS,
} as const

// 权限类型
export type PermissionCode = (typeof PERMISSIONS)[keyof typeof PERMISSIONS]