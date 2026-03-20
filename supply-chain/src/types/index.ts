export interface Supplier {
  id: string
  code: string
  name: string
  contact: string
  phone: string
  email: string
  address: string
  level: 'A' | 'B' | 'C'
  category: string
  paymentTerms: string
  bankName: string
  bankAccount: string
  taxNumber: string
  rating: number
  status: 'active' | 'inactive'
  remark: string
  createdAt: string
  updatedAt: string
}

export interface Product {
  id: string
  code: string
  name: string
  sku: string
  category: string
  unit: string
  costPrice: number
  salePrice: number
  minStock: number
  maxStock: number
  description: string
  status: 'active' | 'inactive'
  createdAt: string
}

export interface InventoryItem {
  id: string
  productId: string
  productCode: string
  productName: string
  sku: string
  category: string
  warehouse: string
  quantity: number
  availableQty: number
  lockedQty: number
  minStock: number
  maxStock: number
  unit: string
  costPrice: number
  salePrice: number
  location: string
  batchNo: string
  expiryDate: string
  status: 'normal' | 'low' | 'over' | 'locked'
}

export interface ProcurementOrder {
  id: string
  orderNo: string
  supplierId: string
  supplierName: string
  supplier?: Supplier
  totalAmount: number
  orderDate: string
  expectedDate: string
  actualDate: string
  status: 'pending' | 'approved' | 'purchasing' | 'received' | 'cancelled'
  warehouse: string
  remark: string
  attachmentUrl: string
  items: ProcurementItem[]
  createdAt: string
  updatedAt: string
}

export interface ProcurementItem {
  id: string
  orderId: string
  productId: string
  productCode: string
  productName: string
  quantity: number
  unit: string
  unitPrice: number
  amount: number
  receivedQty: number
}

export interface SalesOrder {
  id: string
  orderNo: string
  customerId: string
  customerName: string
  customerPhone: string
  customerAddress: string
  totalAmount: number
  discount: number
  tax: number
  shippingFee: number
  orderDate: string
  deliveryDate: string
  status: 'pending' | 'confirmed' | 'shipping' | 'completed' | 'cancelled' | 'refunded'
  paymentMethod: string
  paymentStatus: 'pending' | 'paid' | 'refunded'
  remark: string
  items: SalesOrderItem[]
  createdAt: string
  updatedAt: string
}

export interface SalesOrderItem {
  id: string
  orderId: string
  productId: string
  productCode: string
  productName: string
  quantity: number
  unit: string
  unitPrice: number
  amount: number
  discount: number
}

export interface LogisticsOrder {
  id: string
  trackingNo: string
  carrier: string
  salesOrderId: string
  salesOrderNo: string
  status: 'pending' | 'picked' | 'in_transit' | 'delivering' | 'delivered' | 'returned'
  senderName: string
  senderPhone: string
  senderAddress: string
  receiverName: string
  receiverPhone: string
  receiverAddress: string
  weight: number
  shippingFee: number
  estimatedDelivery: string
  actualDelivery: string
  timeline: LogisticsTimelineItem[]
  createdAt: string
  updatedAt: string
}

export interface LogisticsTimelineItem {
  time: string
  status: string
  location: string
  description: string
}

export interface DashboardStats {
  todaySales: number
  todayOrders: number
  inventoryAlert: number
  pendingProcurement: number
  salesGrowth: number
  ordersGrowth: number
}

export interface ChartData {
  labels: string[]
  datasets: {
    label: string
    data: number[]
    backgroundColor?: string | string[]
    borderColor?: string
    tension?: number
    fill?: boolean
  }[]
}

// 权限
export interface Permission {
  id: number
  name: string
  code: string
  description: string
  type: string
  parentId: number
  status: number
  createdAt: string
}

// 角色
export interface Role {
  id: number
  name: string
  code: string
  description: string
  status: number
  permissions?: Permission[]
  createdAt: string
}

// 用户扩展信息
export interface UserProfile {
  userId: number
  department: string
  position: string
}

// 用户
export interface User {
  id: number
  username: string
  email: string
  phone: string
  role: 'admin' | 'manager' | 'operator'
  department: string
  position: string
  avatar: string
  status: number
  profile?: UserProfile
  roles?: Role[]
  permissions?: Permission[]
  createdAt: string
}