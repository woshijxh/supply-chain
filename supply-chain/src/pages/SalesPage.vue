<template>
  <div class="sales-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('sales.title') }}</h1>
      <p class="page-subtitle">管理销售订单和客户信息</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(SALES_PERMISSIONS.WRITE)" :label="t('common.add')" icon="ri-add-line" @click="openDialog()" />
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
        <Column :header="t('common.actions')" style="width: 180px">
          <template #body="{ data }">
            <Button text severity="info" icon="ri-eye-line" @click="viewDetail(data)" />
            <Button v-if="hasPermission(SALES_PERMISSIONS.WRITE) && data.status === 'pending'" text severity="success" icon="ri-check-line" @click="confirmOrder(data)" />
            <Button v-if="hasPermission(SALES_PERMISSIONS.WRITE)" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
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
        <Column field="productName" :header="t('sales.productName')">
          <template #body="{ index }">
            <InputText v-model="form.items[index].productName" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { salesApi } from '@/api'
import type { SalesOrder, SalesOrderItem } from '@/types'
import { SALES_PERMISSIONS } from '@/config/permissions'
import { usePermission } from '@/utils/usePermission'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const { hasPermission, initUserPermissions } = usePermission()

const loading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const statusFilter = ref('all')
const currentOrder = ref<SalesOrder | null>(null)
const orders = ref<SalesOrder[]>([])

const form = ref({
  customerName: '',
  customerPhone: '',
  customerAddress: '',
  paymentMethod: '银行转账',
  deliveryDateValue: null as Date | null,
  remark: '',
  items: [] as SalesOrderItem[]
})

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
    id: '',
    orderId: '',
    productId: '',
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
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '创建失败', life: 5000 })
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

onMounted(() => {
  initUserPermissions()
  fetchOrders()
})
</script>

<style scoped>
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