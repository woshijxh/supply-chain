import api from './request'
import type { Supplier } from '@/types'

export const authApi = {
  login: (username: string, password: string, captchaId?: string, captcha?: string) =>
    api.post('/auth/login', { username, password, captchaId, captcha }),

  register: (data: { username: string; password: string; email?: string }) =>
    api.post('/auth/register', data),

  profile: () =>
    api.get('/auth/profile'),

  changePassword: (currentPassword: string, newPassword: string) =>
    api.put('/auth/password', { currentPassword, newPassword }),

  updateAvatar: (avatar: string) =>
    api.put('/auth/avatar', { avatar }),

  getCaptcha: () =>
    api.get('/auth/captcha', { responseType: 'blob' })
}

export const productApi = {
  list: (page = 1, pageSize = 100, keyword = '') =>
    api.get('/products', { params: { page, pageSize, keyword } }),

  get: (id: string | number) =>
    api.get(`/products/${id}`),

  create: (data: any) =>
    api.post('/products', data),

  update: (id: string | number, data: any) =>
    api.put(`/products/${id}`, data),

  delete: (id: string | number) =>
    api.delete(`/products/${id}`)
}

export const supplierApi = {
  list: (page = 1, pageSize = 10, keyword = '') => 
    api.get('/suppliers', { params: { page, pageSize, keyword } }),
  
  get: (id: string) => 
    api.get(`/suppliers/${id}`),
  
  create: (data: Partial<Supplier>) => 
    api.post('/suppliers', data),
  
  update: (id: string, data: Partial<Supplier>) => 
    api.put(`/suppliers/${id}`, data),
  
  delete: (id: string) => 
    api.delete(`/suppliers/${id}`)
}

export const inventoryApi = {
  list: (page = 1, pageSize = 10, status = '', warehouse = '') =>
    api.get('/inventory', { params: { page, pageSize, status, warehouse } }),

  get: (id: string) =>
    api.get(`/inventory/${id}`),

  stockIn: (productId: string, quantity: number, warehouse = '默认仓库', procurementId?: string, procurementItemId?: string) =>
    api.post('/inventory/stock-in', { productId, quantity, warehouse, procurementId, procurementItemId }),

  stockOut: (productId: string, quantity: number, warehouse = '默认仓库') =>
    api.post('/inventory/stock-out', { productId, quantity, warehouse }),

  stats: () =>
    api.get('/inventory/stats')
}

export const procurementApi = {
  list: (page = 1, pageSize = 10, status = '') => 
    api.get('/procurement', { params: { page, pageSize, status } }),
  
  get: (id: string) => 
    api.get(`/procurement/${id}`),
  
  create: (data: any) => 
    api.post('/procurement', data),
  
  updateStatus: (id: string, status: string, extra?: { attachmentUrl?: string; remark?: string }) =>
    api.put(`/procurement/${id}/status`, { status, ...extra }),
  
  delete: (id: string) => 
    api.delete(`/procurement/${id}`)
}

export const salesApi = {
  list: (page = 1, pageSize = 10, status = '') => 
    api.get('/sales', { params: { page, pageSize, status } }),
  
  get: (id: string) => 
    api.get(`/sales/${id}`),
  
  create: (data: any) => 
    api.post('/sales', data),
  
  updateStatus: (id: string, status: string) => 
    api.put(`/sales/${id}/status`, { status }),
  
  delete: (id: string) => 
    api.delete(`/sales/${id}`)
}

export const logisticsApi = {
  list: (page = 1, pageSize = 10, status = '') => 
    api.get('/logistics', { params: { page, pageSize, status } }),
  
  get: (id: string) => 
    api.get(`/logistics/${id}`),
  
  create: (data: any) => 
    api.post('/logistics', data),
  
  updateStatus: (id: string, status: string) => 
    api.put(`/logistics/${id}/status`, { status }),
  
  delete: (id: string) => 
    api.delete(`/logistics/${id}`)
}

export const dashboardApi = {
  stats: () =>
    api.get('/dashboard/stats'),
  topProducts: () =>
    api.get('/dashboard/top-products'),
  recentOrders: () =>
    api.get('/dashboard/recent-orders'),
  lowStockItems: () =>
    api.get('/dashboard/low-stock'),
  salesTrend: (days = 7) =>
    api.get('/dashboard/sales-trend', { params: { days } }),
  inventoryDistribution: () =>
    api.get('/dashboard/inventory-distribution')
}

// 客户管理 API
export const customerApi = {
  list: (page = 1, pageSize = 10, keyword = '') =>
    api.get('/customers', { params: { page, pageSize, keyword } }),

  getAll: () =>
    api.get('/customers/all'),

  get: (id: string) =>
    api.get(`/customers/${id}`),

  create: (data: any) =>
    api.post('/customers', data),

  update: (id: string, data: any) =>
    api.put(`/customers/${id}`, data),

  delete: (id: string) =>
    api.delete(`/customers/${id}`)
}

// 库存流水 API
export const inventoryLogApi = {
  list: (page = 1, pageSize = 20, productId = '', type = '') =>
    api.get('/inventory/logs', { params: { page, pageSize, productId, type } })
}

// 销售退货 API
export const salesReturnApi = {
  list: (page = 1, pageSize = 10, status = '') =>
    api.get('/sales-returns', { params: { page, pageSize, status } }),

  get: (id: string) =>
    api.get(`/sales-returns/${id}`),

  create: (data: any) =>
    api.post('/sales-returns', data),

  approve: (id: string) =>
    api.put(`/sales-returns/${id}/approve`),

  reject: (id: string, reason: string) =>
    api.put(`/sales-returns/${id}/reject`, { reason }),

  complete: (id: string, refundAmount: number) =>
    api.put(`/sales-returns/${id}/complete`, { refundAmount }),

  delete: (id: string) =>
    api.delete(`/sales-returns/${id}`)
}

// 采购退货 API
export const procurementReturnApi = {
  list: (page = 1, pageSize = 10, status = '') =>
    api.get('/procurement-returns', { params: { page, pageSize, status } }),

  get: (id: string) =>
    api.get(`/procurement-returns/${id}`),

  create: (data: any) =>
    api.post('/procurement-returns', data),

  approve: (id: string) =>
    api.put(`/procurement-returns/${id}/approve`),

  reject: (id: string, reason: string) =>
    api.put(`/procurement-returns/${id}/reject`, { reason }),

  complete: (id: string, refundAmount: number) =>
    api.put(`/procurement-returns/${id}/complete`, { refundAmount }),

  delete: (id: string) =>
    api.delete(`/procurement-returns/${id}`)
}

// 用户管理 API
export const userApi = {
  list: (page = 1, pageSize = 10, keyword = '') =>
    api.get('/users', { params: { page, pageSize, keyword } }),

  get: (id: string) =>
    api.get(`/users/${id}`),

  create: (data: any) =>
    api.post('/users', data),

  update: (id: string, data: any) =>
    api.put(`/users/${id}`, data),

  delete: (id: string) =>
    api.delete(`/users/${id}`),

  getRoles: (id: string) =>
    api.get(`/users/${id}/roles`),

  setRoles: (id: string, roleIds: number[]) =>
    api.put(`/users/${id}/roles`, { roleIds }),

  getPermissions: (id: string) =>
    api.get(`/users/${id}/permissions`),

  setPermissions: (id: string, permissionIds: number[]) =>
    api.put(`/users/${id}/permissions`, { permissionIds })
}

// 角色管理 API
export const roleApi = {
  list: (page = 1, pageSize = 10, keyword = '') =>
    api.get('/roles', { params: { page, pageSize, keyword } }),

  getAll: () =>
    api.get('/roles/all'),

  get: (id: string) =>
    api.get(`/roles/${id}`),

  create: (data: any) =>
    api.post('/roles', data),

  update: (id: string, data: any) =>
    api.put(`/roles/${id}`, data),

  delete: (id: string) =>
    api.delete(`/roles/${id}`),

  getPermissions: (id: string) =>
    api.get(`/roles/${id}/permissions`),

  setPermissions: (id: string, permissionIds: number[]) =>
    api.post(`/roles/${id}/permissions`, { permissionIds })
}

// 权限管理 API
export const permissionApi = {
  list: (page = 1, pageSize = 10, keyword = '') =>
    api.get('/permissions', { params: { page, pageSize, keyword } }),

  getAll: () =>
    api.get('/permissions/all'),

  get: (id: string) =>
    api.get(`/permissions/${id}`),

  create: (data: any) =>
    api.post('/permissions', data),

  update: (id: string, data: any) =>
    api.put(`/permissions/${id}`, data),

  delete: (id: string) =>
    api.delete(`/permissions/${id}`)
}

// 追溯管理 API
export const traceApi = {
  trace: (code: string) =>
    api.get('/trace', { params: { code } })
}