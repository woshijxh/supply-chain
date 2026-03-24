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
  id: number
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
  id: number
  orderId: number
  productId: number
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
  productId: number
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

// 客户
export interface Customer {
  id: number
  code: string
  name: string
  contact: string
  phone: string
  email: string
  address: string
  level: string
  source: string
  status: number
  remark: string
  createdAt: string
}

// 库存流水
export interface InventoryLog {
  id: number
  productId: number
  productName: string
  type: 'in' | 'out' | 'lock' | 'unlock'
  quantity: number
  beforeQty: number
  afterQty: number
  warehouse: string
  refType: string
  refId: number
  refNo: string
  operator: string
  remark: string
  createdAt: string
}

// 销售退货单
export interface SalesReturn {
  id: number
  returnNo: string
  salesOrderId: number
  salesOrderNo: string
  customerName: string
  totalAmount: number
  refundAmount: number
  status: string
  refundStatus: string
  reason: string
  remark: string
  items: SalesReturnItem[]
  createdAt: string
}

// 销售退货明细
export interface SalesReturnItem {
  id: number
  returnId: number
  productId: number
  productName: string
  quantity: number
  unitPrice: number
  amount: number
  reason: string
}

// 采购退货单
export interface ProcurementReturn {
  id: number
  returnNo: string
  procurementOrderId: number
  procurementOrderNo: string
  supplierName: string
  totalAmount: number
  refundAmount: number
  status: string
  refundStatus: string
  reason: string
  remark: string
  items: ProcurementReturnItem[]
  createdAt: string
}

// 采购退货明细
export interface ProcurementReturnItem {
  id: number
  returnId: number
  productId: number
  productName: string
  quantity: number
  unitPrice: number
  amount: number
  reason: string
}

// ========== 追溯相关类型 ==========

// 追溯结果
export interface TraceResult {
  product: ProductInfo
  procurement: ProcurementTrace[]
  inventory: InventoryTrace[]
  sales: SalesTrace[]
  logistics: LogisticsTrace[]
  returns: ReturnTrace[]
  timeline: TimelineEvent[]
}

// 产品追溯信息
export interface ProductInfo {
  id: number
  code: string
  name: string
  sku: string
  category: string
  unit: string
  costPrice: number
  salePrice: number
  totalStock: number
}

// 采购追溯节点
export interface ProcurementTrace {
  orderId: number
  orderNo: string
  supplierId: number
  supplierName: string
  quantity: number
  receivedQty: number
  unitPrice: number
  orderDate: string
  actualDate: string
  warehouse: string
  status: string
}

// 库存追溯节点
export interface InventoryTrace {
  id: number
  type: 'in' | 'out' | 'lock' | 'unlock'
  quantity: number
  beforeQty: number
  afterQty: number
  warehouse: string
  batchNo: string
  refType: string
  refNo: string
  operator: string
  remark: string
  createdAt: string
}

// 销售追溯节点
export interface SalesTrace {
  orderId: number
  orderNo: string
  customerName: string
  customerPhone: string
  quantity: number
  unitPrice: number
  orderDate: string
  deliveryDate: string
  status: string
  paymentStatus: string
  trackingNo: string
  logisticsStatus: string
}

// 物流追溯节点
export interface LogisticsTrace {
  logisticsId: number
  trackingNo: string
  carrier: string
  salesOrderNo: string
  receiverName: string
  receiverPhone: string
  status: string
  shippingFee: number
  estimatedDelivery: string
  actualDelivery: string
  timeline: LogisticsTimelineItem[]
}

// 物流轨迹项
export interface LogisticsTimelineItem {
  time: string
  status: string
  location: string
  description: string
}

// 退货追溯节点
export interface ReturnTrace {
  returnId: number
  returnNo: string
  type: 'sales_return' | 'procurement_return'
  orderNo: string
  relatedName: string
  quantity: number
  unitPrice: number
  reason: string
  status: string
  createdAt: string
}

// 时间线事件
export interface TimelineEvent {
  time: string
  type: 'procurement' | 'inventory' | 'sales' | 'logistics' | 'return'
  action: string
  title: string
  description: string
  refNo: string
  refType: string
  operator: string
  quantity: number
  amount: number
}