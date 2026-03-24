import type { Supplier, ProcurementOrder, InventoryItem, SalesOrder, LogisticsOrder, DashboardStats } from '@/types'

const generateId = () => Math.random().toString(36).substr(2, 9)
const generateCode = (prefix: string) => `${prefix}${Date.now().toString(36).toUpperCase()}`

export const mockSuppliers: Supplier[] = [
  {
    id: '1',
    code: 'SUP001',
    name: '华东电子科技有限公司',
    contact: '张经理',
    phone: '13800138001',
    email: 'zhang@huadong.com',
    address: '上海市浦东新区张江高科技园区',
    level: 'A',
    category: '电子元器件',
    paymentTerms: '月结30天',
    bankName: '中国工商银行',
    bankAccount: '6222021234567890',
    taxNumber: '91310000MA1K5ABC12',
    rating: 4.8,
    status: 'active',
    remark: '主要供应商，合作多年',
    createdAt: '2024-01-15',
    updatedAt: '2024-03-01'
  },
  {
    id: '2',
    code: 'SUP002',
    name: '深圳智能制造有限公司',
    contact: '李总监',
    phone: '13900139002',
    email: 'li@smartmfg.com',
    address: '深圳市南山区科技园',
    level: 'A',
    category: '机械配件',
    paymentTerms: '月结45天',
    bankName: '招商银行',
    bankAccount: '6214831234567890',
    taxNumber: '91440300MA5EXABC34',
    rating: 4.5,
    status: 'active',
    remark: '',
    createdAt: '2024-02-01',
    updatedAt: '2024-03-05'
  },
  {
    id: '3',
    code: 'SUP003',
    name: '杭州材料科技有限公司',
    contact: '王经理',
    phone: '13700137003',
    email: 'wang@hzmat.com',
    address: '杭州市余杭区未来科技城',
    level: 'B',
    category: '原材料',
    paymentTerms: '月结15天',
    bankName: '中国建设银行',
    bankAccount: '6217001234567890',
    taxNumber: '91330100MA2BABC56',
    rating: 4.2,
    status: 'active',
    remark: '新供应商',
    createdAt: '2024-02-20',
    updatedAt: '2024-02-20'
  },
  {
    id: '4',
    code: 'SUP004',
    name: '广州包装材料有限公司',
    contact: '陈主管',
    phone: '13600136004',
    email: 'chen@gzpack.com',
    address: '广州市白云区工业园',
    level: 'B',
    category: '包装材料',
    paymentTerms: '月结30天',
    bankName: '中国农业银行',
    bankAccount: '6228481234567890',
    taxNumber: '91440100MA5DABC78',
    rating: 3.9,
    status: 'inactive',
    remark: '暂停合作',
    createdAt: '2023-10-10',
    updatedAt: '2024-01-15'
  }
]

export const mockInventory: InventoryItem[] = [
  {
    id: '1',
    productId: '1',
    productCode: 'PRD001',
    productName: '高精度传感器',
    sku: 'SKU-SENSOR-001',
    category: '电子元器件',
    warehouse: '上海仓库',
    quantity: 1500,
    availableQty: 1200,
    lockedQty: 300,
    minStock: 500,
    maxStock: 3000,
    unit: '个',
    costPrice: 25.50,
    salePrice: 45.00,
    location: 'A区-01-02',
    batchNo: 'B20240301',
    expiryDate: '',
    status: 'normal'
  },
  {
    id: '2',
    productId: '2',
    productCode: 'PRD002',
    productName: '工业控制器',
    sku: 'SKU-CTRL-001',
    category: '电子元器件',
    warehouse: '上海仓库',
    quantity: 320,
    availableQty: 280,
    lockedQty: 40,
    minStock: 200,
    maxStock: 1000,
    unit: '台',
    costPrice: 180.00,
    salePrice: 320.00,
    location: 'A区-02-01',
    batchNo: 'B20240215',
    expiryDate: '',
    status: 'normal'
  },
  {
    id: '3',
    productId: '3',
    productCode: 'PRD003',
    productName: '精密轴承',
    sku: 'SKU-BEAR-001',
    category: '机械配件',
    warehouse: '深圳仓库',
    quantity: 80,
    availableQty: 80,
    lockedQty: 0,
    minStock: 100,
    maxStock: 500,
    unit: '套',
    costPrice: 85.00,
    salePrice: 150.00,
    location: 'B区-01-03',
    batchNo: 'B20240305',
    expiryDate: '',
    status: 'low'
  },
  {
    id: '4',
    productId: '4',
    productCode: 'PRD004',
    productName: '不锈钢板材',
    sku: 'SKU-STEEL-001',
    category: '原材料',
    warehouse: '杭州仓库',
    quantity: 2500,
    availableQty: 2000,
    lockedQty: 500,
    minStock: 300,
    maxStock: 1500,
    unit: '张',
    costPrice: 120.00,
    salePrice: 180.00,
    location: 'C区-全仓',
    batchNo: 'B20240220',
    expiryDate: '',
    status: 'over'
  },
  {
    id: '5',
    productId: '5',
    productCode: 'PRD005',
    productName: '工业电机',
    sku: 'SKU-MOTOR-001',
    category: '机械配件',
    warehouse: '上海仓库',
    quantity: 45,
    availableQty: 30,
    lockedQty: 15,
    minStock: 50,
    maxStock: 200,
    unit: '台',
    costPrice: 450.00,
    salePrice: 680.00,
    location: 'A区-03-01',
    batchNo: 'B20240110',
    expiryDate: '',
    status: 'low'
  }
]

export const mockProcurementOrders: ProcurementOrder[] = [
  {
    id: '1',
    orderNo: 'PO20240301001',
    supplierId: '1',
    supplierName: '华东电子科技有限公司',
    totalAmount: 52500.00,
    orderDate: '2024-03-01',
    expectedDate: '2024-03-15',
    actualDate: '',
    status: 'purchasing',
    warehouse: '上海仓库',
    remark: '急用',
    attachmentUrl: '',
    items: [
      { id: 1, orderId: 1, productId: 1, productCode: 'PRD001', productName: '高精度传感器', quantity: 1000, unit: '个', unitPrice: 35.00, amount: 35000, receivedQty: 0 },
      { id: 2, orderId: 1, productId: 2, productCode: 'PRD002', productName: '工业控制器', quantity: 100, unit: '台', unitPrice: 175.00, amount: 17500, receivedQty: 0 }
    ],
    createdAt: '2024-03-01',
    updatedAt: '2024-03-01'
  },
  {
    id: '2',
    orderNo: 'PO20240302001',
    supplierId: '3',
    supplierName: '杭州材料科技有限公司',
    totalAmount: 36000.00,
    orderDate: '2024-03-02',
    expectedDate: '2024-03-20',
    actualDate: '',
    status: 'approved',
    warehouse: '杭州仓库',
    remark: '',
    attachmentUrl: '',
    items: [
      { id: 1, orderId: 2, productId: 4, productCode: 'PRD004', productName: '不锈钢板材', quantity: 300, unit: '张', unitPrice: 120.00, amount: 36000, receivedQty: 0 }
    ],
    createdAt: '2024-03-02',
    updatedAt: '2024-03-02'
  },
  {
    id: '3',
    orderNo: 'PO20240225001',
    supplierId: '2',
    supplierName: '深圳智能制造有限公司',
    totalAmount: 8500.00,
    orderDate: '2024-02-25',
    expectedDate: '2024-03-05',
    actualDate: '2024-03-04',
    status: 'received',
    warehouse: '深圳仓库',
    remark: '',
    attachmentUrl: '',
    items: [
      { id: 1, orderId: 3, productId: 3, productCode: 'PRD003', productName: '精密轴承', quantity: 100, unit: '套', unitPrice: 85.00, amount: 8500, receivedQty: 100 }
    ],
    createdAt: '2024-02-25',
    updatedAt: '2024-03-04'
  }
]

export const mockSalesOrders: SalesOrder[] = [
  {
    id: '1',
    orderNo: 'SO20240305001',
    customerId: 'C001',
    customerName: '北京智能装备有限公司',
    customerPhone: '010-12345678',
    customerAddress: '北京市海淀区中关村科技园',
    totalAmount: 45800.00,
    discount: 500.00,
    tax: 4580.00,
    shippingFee: 0,
    orderDate: '2024-03-05',
    deliveryDate: '',
    status: 'confirmed',
    paymentMethod: '银行转账',
    paymentStatus: 'paid',
    remark: '',
    items: [
      { id: '1', orderId: '1', productId: 1, productCode: 'PRD001', productName: '高精度传感器', quantity: 200, unit: '个', unitPrice: 45.00, amount: 9000, discount: 0 },
      { id: '2', orderId: '1', productId: 2, productCode: 'PRD002', productName: '工业控制器', quantity: 50, unit: '台', unitPrice: 320.00, amount: 16000, discount: 0 },
      { id: '3', orderId: '1', productId: 5, productCode: 'PRD005', productName: '工业电机', quantity: 30, unit: '台', unitPrice: 680.00, amount: 20400, discount: 100 }
    ],
    createdAt: '2024-03-05',
    updatedAt: '2024-03-05'
  },
  {
    id: '2',
    orderNo: 'SO20240306001',
    customerId: 'C002',
    customerName: '苏州精密仪器厂',
    customerPhone: '0512-87654321',
    customerAddress: '苏州市工业园区',
    totalAmount: 15000.00,
    discount: 0,
    tax: 1500.00,
    shippingFee: 200.00,
    orderDate: '2024-03-06',
    deliveryDate: '',
    status: 'pending',
    paymentMethod: '支付宝',
    paymentStatus: 'pending',
    remark: '需要发票',
    items: [
      { id: '1', orderId: '2', productId: 3, productCode: 'PRD003', productName: '精密轴承', quantity: 100, unit: '套', unitPrice: 150.00, amount: 15000, discount: 0 }
    ],
    createdAt: '2024-03-06',
    updatedAt: '2024-03-06'
  }
]

export const mockLogisticsOrders: LogisticsOrder[] = [
  {
    id: '1',
    trackingNo: 'SF1234567890',
    carrier: '顺丰速运',
    salesOrderId: '1',
    salesOrderNo: 'SO20240305001',
    status: 'in_transit',
    senderName: '发货中心',
    senderPhone: '400-123-4567',
    senderAddress: '上海市浦东新区仓库',
    receiverName: '北京智能装备有限公司',
    receiverPhone: '010-12345678',
    receiverAddress: '北京市海淀区中关村科技园',
    weight: 15.5,
    shippingFee: 45.00,
    estimatedDelivery: '2024-03-08',
    actualDelivery: '',
    timeline: [
      { time: '2024-03-06 14:30', status: '已发货', location: '上海', description: '快件已发出' },
      { time: '2024-03-06 18:00', status: '运输中', location: '上海转运中心', description: '快件已到达上海转运中心' },
      { time: '2024-03-07 08:00', status: '运输中', location: '北京转运中心', description: '快件已到达北京转运中心' }
    ],
    createdAt: '2024-03-06',
    updatedAt: '2024-03-07'
  }
]

export const mockDashboardStats: DashboardStats = {
  todaySales: 156800.00,
  todayOrders: 28,
  inventoryAlert: 3,
  pendingProcurement: 2,
  salesGrowth: 12.5,
  ordersGrowth: 8.3
}

export function generateSupplier(data: Partial<Supplier>): Supplier {
  return {
    id: generateId(),
    code: data.code || generateCode('SUP'),
    name: data.name || '',
    contact: data.contact || '',
    phone: data.phone || '',
    email: data.email || '',
    address: data.address || '',
    level: data.level || 'B',
    category: data.category || '',
    paymentTerms: data.paymentTerms || '',
    bankName: data.bankName || '',
    bankAccount: data.bankAccount || '',
    taxNumber: data.taxNumber || '',
    rating: data.rating || 4.0,
    status: data.status || 'active',
    remark: data.remark || '',
    createdAt: new Date().toISOString().split('T')[0],
    updatedAt: new Date().toISOString().split('T')[0]
  }
}

export function generateProcurementOrder(data: Partial<ProcurementOrder>): ProcurementOrder {
  return {
    id: generateId(),
    orderNo: data.orderNo || generateCode('PO'),
    supplierId: data.supplierId || '',
    supplierName: data.supplierName || '',
    totalAmount: data.totalAmount || 0,
    orderDate: data.orderDate || new Date().toISOString().split('T')[0],
    expectedDate: data.expectedDate || '',
    actualDate: data.actualDate || '',
    status: data.status || 'pending',
    warehouse: data.warehouse || '',
    remark: data.remark || '',
    attachmentUrl: data.attachmentUrl || '',
    items: data.items || [],
    createdAt: new Date().toISOString().split('T')[0],
    updatedAt: new Date().toISOString().split('T')[0]
  }
}

export function generateSalesOrder(data: Partial<SalesOrder>): SalesOrder {
  return {
    id: generateId(),
    orderNo: data.orderNo || generateCode('SO'),
    customerId: data.customerId || '',
    customerName: data.customerName || '',
    customerPhone: data.customerPhone || '',
    customerAddress: data.customerAddress || '',
    totalAmount: data.totalAmount || 0,
    discount: data.discount || 0,
    tax: data.tax || 0,
    shippingFee: data.shippingFee || 0,
    orderDate: data.orderDate || new Date().toISOString().split('T')[0],
    deliveryDate: data.deliveryDate || '',
    status: data.status || 'pending',
    paymentMethod: data.paymentMethod || '',
    paymentStatus: data.paymentStatus || 'pending',
    remark: data.remark || '',
    items: data.items || [],
    createdAt: new Date().toISOString().split('T')[0],
    updatedAt: new Date().toISOString().split('T')[0]
  }
}