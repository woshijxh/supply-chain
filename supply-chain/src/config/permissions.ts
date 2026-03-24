/**
 * 权限码定义 - 与后端 permission 表的 code 字段对应
 */

// 仪表盘权限
export const DASHBOARD_PERMISSIONS = {
  READ: 'dashboard:read',
} as const

// 供应商权限
export const SUPPLIER_PERMISSIONS = {
  READ: 'supplier:read',
  CREATE: 'supplier:create',
  UPDATE: 'supplier:update',
  DELETE: 'supplier:delete',
} as const

// 产品权限
export const PRODUCT_PERMISSIONS = {
  READ: 'product:read',
  CREATE: 'product:create',
  UPDATE: 'product:update',
  DELETE: 'product:delete',
} as const

// 库存权限
export const INVENTORY_PERMISSIONS = {
  READ: 'inventory:read',
  CREATE: 'inventory:create',
  UPDATE: 'inventory:update',
  DELETE: 'inventory:delete',
} as const

// 客户权限
export const CUSTOMER_PERMISSIONS = {
  READ: 'customer:read',
  CREATE: 'customer:create',
  UPDATE: 'customer:update',
  DELETE: 'customer:delete',
} as const

// 采购权限
export const PROCUREMENT_PERMISSIONS = {
  READ: 'procurement:read',
  CREATE: 'procurement:create',
  UPDATE: 'procurement:update',
  DELETE: 'procurement:delete',
} as const

// 销售权限
export const SALES_PERMISSIONS = {
  READ: 'sales:read',
  CREATE: 'sales:create',
  UPDATE: 'sales:update',
  DELETE: 'sales:delete',
} as const

// 物流权限
export const LOGISTICS_PERMISSIONS = {
  READ: 'logistics:read',
  CREATE: 'logistics:create',
  UPDATE: 'logistics:update',
  DELETE: 'logistics:delete',
} as const

// 用户管理权限
export const USER_PERMISSIONS = {
  READ: 'user:read',
  CREATE: 'user:create',
  UPDATE: 'user:update',
  DELETE: 'user:delete',
} as const

// 角色管理权限
export const ROLE_PERMISSIONS = {
  READ: 'role:read',
  CREATE: 'role:create',
  UPDATE: 'role:update',
  DELETE: 'role:delete',
} as const

// 权限管理权限
export const PERMISSION_PERMISSIONS = {
  READ: 'permission:read',
  CREATE: 'permission:create',
  UPDATE: 'permission:update',
  DELETE: 'permission:delete',
} as const

// 追溯权限
export const TRACE_PERMISSIONS = {
  READ: 'trace:read',
} as const

// 所有权限映射
export const PERMISSIONS = {
  ...DASHBOARD_PERMISSIONS,
  ...SUPPLIER_PERMISSIONS,
  ...PRODUCT_PERMISSIONS,
  ...INVENTORY_PERMISSIONS,
  ...CUSTOMER_PERMISSIONS,
  ...PROCUREMENT_PERMISSIONS,
  ...SALES_PERMISSIONS,
  ...LOGISTICS_PERMISSIONS,
  ...USER_PERMISSIONS,
  ...ROLE_PERMISSIONS,
  ...PERMISSION_PERMISSIONS,
  ...TRACE_PERMISSIONS,
} as const

// 权限类型
export type PermissionCode = (typeof PERMISSIONS)[keyof typeof PERMISSIONS]