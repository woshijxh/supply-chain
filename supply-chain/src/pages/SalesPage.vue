<template>
  <div class="sales-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('sales.title') }}</h1>
      <p class="page-subtitle">管理销售订单和客户信息</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(SALES_PERMISSIONS.CREATE)" :label="t('common.add')" icon="ri-add-line" @click="openDialog()" />
      <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.status')" style="width: 150px" />
    </div>

    <div class="card">
      <DataTable :value="orders" stripedRows :paginator="true" :rows="10">
        <Column field="orderNo" :header="t('sales.orderNo')" style="width: 150px"></Column>
        <Column field="customerName" :header="t('sales.customer')"></Column>
        <Column field="orderDate" :header="t('sales.orderDate')"></Column>
        <Column field="totalAmount" :header="t('sales.totalAmount')">
          <template #body="{ data }">
            ¥{{ formatNumber(data.totalAmount) }}
          </template>
        </Column>
        <Column field="status" :header="t('sales.status')">
          <template #body="{ data }">
            <span :class="['tag', getStatusClass(data.status)]">{{ getStatusText(data.status) }}</span>
          </template>
        </Column>
        <Column field="paymentStatus" :header="t('sales.paymentStatus')">
          <template #body="{ data }">
            <span :class="['tag', getPaymentClass(data.paymentStatus)]">{{ getPaymentText(data.paymentStatus) }}</span>
          </template>
        </Column>
        <Column field="paymentMethod" :header="t('sales.paymentMethod')"></Column>
        <Column :header="t('common.actions')" style="width: 220px">
          <template #body="{ data }">
            <Button text severity="info" icon="ri-eye-line" @click="viewDetail(data)" />
            <Button v-if="hasPermission(SALES_PERMISSIONS.UPDATE) && data.status === 'pending'" text severity="success" icon="ri-check-line" @click="confirmOrder(data)" />
            <Button v-if="hasPermission(SALES_PERMISSIONS.UPDATE) && (data.status === 'confirmed' || data.status === 'completed')" text severity="warning" icon="ri-reply-line" @click="openReturnDialog(data)" />
            <Button v-if="hasPermission(SALES_PERMISSIONS.DELETE)" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="dialogVisible" header="新增销售订单" :style="{ width: '700px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">客户名称 *</label>
          <InputText v-model="form.customerName" />
        </div>
        <div class="form-group">
          <label class="form-label">联系电话</label>
          <InputText v-model="form.customerPhone" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">客户地址</label>
          <InputText v-model="form.customerAddress" />
        </div>
        <div class="form-group">
          <label class="form-label">支付方式</label>
          <Select v-model="form.paymentMethod" :options="paymentOptions" />
        </div>
        <div class="form-group">
          <label class="form-label">发货日期</label>
          <DatePicker v-model="form.deliveryDateValue" dateFormat="yy-mm-dd" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">{{ t('common.remark') }}</label>
          <Textarea v-model="form.remark" rows="2" />
        </div>
      </div>

      <h4 style="margin: 20px 0 12px;">{{ t('sales.items') }}</h4>
      <DataTable :value="form.items" stripedRows>
        <Column field="productName" :header="t('sales.productName')" style="width: 200px">
          <template #body="{ index }">
            <Select
              v-model="form.items[index].productId"
              :options="products"
              optionLabel="name"
              optionValue="id"
              placeholder="选择产品"
              @change="onProductSelect(index)"
              style="width: 100%"
            />
          </template>
        </Column>
        <Column header="库存" style="width: 80px">
          <template #body="{ data }">
            <span :class="{ 'text-danger': getStock(data.productId) < data.quantity }">
              {{ getStock(data.productId) }}
            </span>
          </template>
        </Column>
        <Column field="quantity" :header="t('sales.quantity')" style="width: 100px">
          <template #body="{ index }">
            <InputNumber v-model="form.items[index].quantity" :min="1" @input="calcAmount(index)" />
          </template>
        </Column>
        <Column field="unitPrice" :header="t('sales.unitPrice')" style="width: 120px">
          <template #body="{ index }">
            <InputNumber v-model="form.items[index].unitPrice" mode="currency" currency="CNY" @input="calcAmount(index)" />
          </template>
        </Column>
        <Column field="amount" :header="t('sales.amount')" style="width: 100px">
          <template #body="{ data }">
            ¥{{ formatNumber(data.amount) }}
          </template>
        </Column>
        <Column style="width: 60px">
          <template #body="{ index }">
            <Button text severity="danger" icon="ri-delete-bin-line" @click="removeItem(index)" />
          </template>
        </Column>
      </DataTable>
      <Button text icon="ri-add-line" :label="t('common.add')" @click="addItem" style="margin-top: 8px" />

      <div class="dialog-footer">
        <Button :label="t('common.cancel')" severity="secondary" @click="dialogVisible = false" />
        <Button :label="t('common.save')" @click="handleSave" />
      </div>
    </Dialog>

    <Dialog v-model:visible="detailVisible" header="销售订单详情" :style="{ width: '700px' }">
      <div v-if="currentOrder" class="order-detail">
        <div class="detail-section">
          <h4>订单信息</h4>
          <div class="detail-row">
            <span class="detail-label">{{ t('sales.orderNo') }}:</span>
            <span>{{ currentOrder.orderNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('sales.customer') }}:</span>
            <span>{{ currentOrder.customerName }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">联系电话:</span>
            <span>{{ currentOrder.customerPhone }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">地址:</span>
            <span>{{ currentOrder.customerAddress }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ t('sales.status') }}:</span>
            <span :class="['tag', getStatusClass(currentOrder.status)]">{{ getStatusText(currentOrder.status) }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4>{{ t('sales.items') }}</h4>
          <DataTable :value="currentOrder.items" stripedRows>
            <Column field="productName" :header="t('sales.productName')"></Column>
            <Column field="quantity" :header="t('sales.quantity')"></Column>
            <Column field="unitPrice" :header="t('sales.unitPrice')">
              <template #body="{ data }">¥{{ data.unitPrice }}</template>
            </Column>
            <Column field="amount" :header="t('sales.amount')">
              <template #body="{ data }">¥{{ formatNumber(data.amount) }}</template>
            </Column>
          </DataTable>
        </div>

        <div class="order-summary">
          <div class="summary-row">
            <span>商品金额:</span>
            <span>¥{{ formatNumber(currentOrder.totalAmount) }}</span>
          </div>
          <div class="summary-row">
            <span>{{ t('sales.discount') }}:</span>
            <span>-¥{{ formatNumber(currentOrder.discount) }}</span>
          </div>
          <div class="summary-row">
            <span>{{ t('sales.tax') }}:</span>
            <span>¥{{ formatNumber(currentOrder.tax) }}</span>
          </div>
          <div class="summary-row">
            <span>{{ t('sales.shippingFee') }}:</span>
            <span>¥{{ formatNumber(currentOrder.shippingFee) }}</span>
          </div>
          <div class="summary-row total">
            <span>订单总额:</span>
            <strong>¥{{ formatNumber(currentOrder.totalAmount - currentOrder.discount + currentOrder.tax + currentOrder.shippingFee) }}</strong>
          </div>
        </div>
      </div>
    </Dialog>

    <!-- 退货弹窗 -->
    <Dialog v-model:visible="returnDialogVisible" header="销售退货" :style="{ width: '600px' }">
      <div class="return-info" style="margin-bottom: 16px; padding: 12px; background: var(--gray-50); border-radius: 8px;">
        <div><strong>订单号:</strong> {{ returnForm.salesOrderNo }}</div>
        <div><strong>客户:</strong> {{ returnForm.customerName }}</div>
      </div>
      <div class="form-group" style="margin-bottom: 16px;">
        <label class="form-label">退货原因</label>
        <Textarea v-model="returnForm.reason" rows="2" placeholder="请输入退货原因" />
      </div>
      <h4 style="margin-bottom: 12px;">退货商品</h4>
      <DataTable :value="returnForm.items" stripedRows>
        <Column field="productName" header="产品"></Column>
        <Column field="originalQty" header="原数量" style="width: 80px"></Column>
        <Column header="退货数量" style="width: 120px">
          <template #body="{ data }">
            <InputNumber v-model="data.returnQty" :min="0" :max="data.originalQty" style="width: 100px" />
          </template>
        </Column>
        <Column field="unitPrice" header="单价">
          <template #body="{ data }">¥{{ data.unitPrice }}</template>
        </Column>
      </DataTable>
      <div class="dialog-footer">
        <Button :label="t('common.cancel')" severity="secondary" @click="returnDialogVisible = false" />
        <Button label="提交退货" @click="handleReturn" :loading="returnLoading" />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { salesApi, productApi, inventoryApi, salesReturnApi } from '@/api'
import type { SalesOrder, SalesOrderItem, Product } from '@/types'
import { SALES_PERMISSIONS } from '@/config/permissions'
import { usePermission } from '@/utils/usePermission'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const { hasPermission, initUserPermissions } = usePermission()

const loading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const products = ref<Product[]>([])
const stockMap = ref<Record<number, { quantity: number; availableQty: number }>>({})
const statusFilter = ref('all')
const currentOrder = ref<SalesOrder | null>(null)
const orders = ref<SalesOrder[]>([])

// 退货相关
const returnDialogVisible = ref(false)
const returnLoading = ref(false)
const returnForm = ref({
  salesOrderId: 0,
  salesOrderNo: '',
  customerName: '',
  reason: '',
  items: [] as { productId: number; productName: string; originalQty: number; returnQty: number; unitPrice: number }[]
})

const form = ref({
  customerName: '',
  customerPhone: '',
  customerAddress: '',
  paymentMethod: '银行转账',
  deliveryDateValue: null as Date | null,
  remark: '',
  items: [] as SalesOrderItem[]
})

// 加载产品列表
const fetchProducts = async () => {
  try {
    const res: any = await productApi.list(1, 1000)
    if (res.code === 0) {
      products.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

// 加载库存信息
const fetchInventory = async () => {
  try {
    const res: any = await inventoryApi.list(1, 1000)
    if (res.code === 0) {
      const list = res.data.list || []
      stockMap.value = {}
      list.forEach((item: any) => {
        stockMap.value[item.productId] = {
          quantity: item.quantity,
          availableQty: item.availableQty
        }
      })
    }
  } catch (error) {
    console.error('获取库存信息失败', error)
  }
}

// 获取产品库存
const getStock = (productId: number) => {
  return stockMap.value[productId]?.availableQty ?? 0
}

// 选择产品时自动填充价格
const onProductSelect = (index: number) => {
  const item = form.value.items[index]
  const product = products.value.find(p => p.id === item.productId)
  if (product) {
    item.productName = product.name
    item.unitPrice = product.salePrice || 0
    calcAmount(index)
  }
}

const statusOptions = [
  { label: t('common.all'), value: 'all' },
  { label: t('sales.statusPending'), value: 'pending' },
  { label: t('sales.statusConfirmed'), value: 'confirmed' },
  { label: t('sales.statusShipping'), value: 'shipping' },
  { label: t('sales.statusCompleted'), value: 'completed' }
]

const paymentOptions = ['银行转账', '支付宝', '微信支付', '现金']

const formatNumber = (num: number) => num.toLocaleString('zh-CN')

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    shipping: '发货中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || status
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    confirmed: 'info',
    shipping: 'info',
    completed: 'success',
    cancelled: 'danger',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

const getPaymentText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    refunded: '已退款'
  }
  return map[status] || status
}

const getPaymentClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    refunded: 'danger'
  }
  return map[status] || 'info'
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res: any = await salesApi.list(1, 100, statusFilter.value === 'all' ? '' : statusFilter.value)
    if (res.code === 0) {
      orders.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取销售订单失败', error)
  } finally {
    loading.value = false
  }
}

const addItem = () => {
  form.value.items.push({
    id: 0,
    orderId: 0,
    productId: 0,
    productCode: '',
    productName: '',
    quantity: 1,
    unit: '个',
    unitPrice: 0,
    amount: 0,
    discount: 0
  })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const calcAmount = (index: number) => {
  const item = form.value.items[index]
  item.amount = (item.quantity || 0) * (item.unitPrice || 0)
}

const openDialog = () => {
  form.value = {
    customerName: '',
    customerPhone: '',
    customerAddress: '',
    paymentMethod: '银行转账',
    deliveryDateValue: null,
    remark: '',
    items: []
  }
  dialogVisible.value = true
}

const viewDetail = async (order: SalesOrder) => {
  try {
    const res: any = await salesApi.get(String(order.id))
    if (res.code === 0) {
      currentOrder.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error('获取订单详情失败', error)
  }
}

const confirmOrder = async (order: SalesOrder) => {
  try {
    const res: any = await salesApi.updateStatus(String(order.id), 'confirmed')
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '确认成功', life: 5000 })
      fetchOrders()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '确认失败', life: 5000 })
  }
}

const handleSave = async () => {
  if (!form.value.customerName) {
    toast.add({ severity: 'warn', summary: '请填写客户名称', life: 5000 })
    return
  }

  if (form.value.items.length === 0) {
    toast.add({ severity: 'warn', summary: '请添加订单明细', life: 5000 })
    return
  }

  const totalAmount = form.value.items.reduce((sum, item) => sum + item.amount, 0)
  const deliveryDate = form.value.deliveryDateValue
    ? `${form.value.deliveryDateValue.getFullYear()}-${String(form.value.deliveryDateValue.getMonth() + 1).padStart(2, '0')}-${String(form.value.deliveryDateValue.getDate()).padStart(2, '0')}`
    : ''

  try {
    const res: any = await salesApi.create({
      customerName: form.value.customerName,
      customerPhone: form.value.customerPhone,
      customerAddress: form.value.customerAddress,
      paymentMethod: form.value.paymentMethod,
      deliveryDate,
      remark: form.value.remark,
      items: form.value.items,
      totalAmount,
      status: 'pending',
      paymentStatus: 'pending'
    })
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: t('common.saveSuccess'), life: 5000 })
      dialogVisible.value = false
      fetchOrders()
    } else {
      toast.add({ severity: 'error', summary: res.message || '创建失败', life: 5000 })
    }
  } catch (error: any) {
    const message = error.response?.data?.message || error.message || '创建失败'
    toast.add({ severity: 'error', summary: message, life: 5000 })
  }
}

const handleDelete = (order: SalesOrder) => {
  confirm.require({
    message: t('common.confirmDelete'),
    header: t('common.confirm'),
    accept: async () => {
      try {
        const res: any = await salesApi.delete(String(order.id))
        if (res.code === 0) {
          toast.add({ severity: 'success', summary: t('common.deleteSuccess'), life: 5000 })
          fetchOrders()
        }
      } catch (error) {
        toast.add({ severity: 'error', summary: '删除失败', life: 5000 })
      }
    }
  })
}

// 退货相关方法
const openReturnDialog = async (order: SalesOrder) => {
  try {
    const res: any = await salesApi.get(String(order.id))
    if (res.code === 0) {
      const detail = res.data
      returnForm.value = {
        salesOrderId: detail.id,
        salesOrderNo: detail.orderNo,
        customerName: detail.customerName,
        reason: '',
        items: detail.items.map((item: any) => ({
          productId: item.productId,
          productName: item.productName,
          originalQty: item.quantity,
          returnQty: item.quantity,
          unitPrice: item.unitPrice
        }))
      }
      returnDialogVisible.value = true
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '获取订单详情失败', life: 5000 })
  }
}

const handleReturn = async () => {
  const items = returnForm.value.items.filter(item => item.returnQty > 0)
  if (items.length === 0) {
    toast.add({ severity: 'warn', summary: '请选择退货商品', life: 5000 })
    return
  }

  returnLoading.value = true
  try {
    const res: any = await salesReturnApi.create({
      salesOrderId: returnForm.value.salesOrderId,
      salesOrderNo: returnForm.value.salesOrderNo,
      customerName: returnForm.value.customerName,
      reason: returnForm.value.reason,
      items: items.map(item => ({
        productId: item.productId,
        productName: item.productName,
        quantity: item.returnQty,
        unitPrice: item.unitPrice
      }))
    })
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '退货申请已提交', life: 5000 })
      returnDialogVisible.value = false
    } else {
      toast.add({ severity: 'error', summary: res.message || '提交失败', life: 5000 })
    }
  } catch (error: any) {
    const message = error.response?.data?.message || error.message || '提交失败'
    toast.add({ severity: 'error', summary: message, life: 5000 })
  } finally {
    returnLoading.value = false
  }
}

onMounted(() => {
  initUserPermissions()
  fetchOrders()
  fetchProducts()
  fetchInventory()
})
</script>

<style scoped>
.text-danger {
  color: #ef4444;
  font-weight: 500;
}

.order-detail {
  .detail-section {
    margin-bottom: 20px;
    h4 {
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid var(--border);
    }
  }
  .detail-row {
    display: flex;
    gap: 12px;
    margin-bottom: 8px;
    .detail-label {
      width: 80px;
      color: var(--text-muted);
    }
  }
  .order-summary {
    margin-top: 16px;
    padding: 16px;
    background: var(--bg);
    border-radius: var(--radius-md);
    .summary-row {
      display: flex;
      justify-content: space-between;
      margin-bottom: 8px;
      &.total {
        margin-top: 12px;
        padding-top: 12px;
        border-top: 1px solid var(--border);
        font-size: 18px;
        strong {
          color: var(--primary);
        }
      }
    }
  }
}
</style>