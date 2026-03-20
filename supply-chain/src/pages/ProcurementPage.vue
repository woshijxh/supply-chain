<template>
  <div class="procurement-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('procurement.title') }}</h1>
      <p class="page-subtitle">管理采购订单和入库流程</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(PROCUREMENT_PERMISSIONS.WRITE)" :label="t('common.add')" icon="ri-add-line" @click="openDialog()" />
      <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.status')" style="width: 150px" />
    </div>

    <div class="card">
      <DataTable :value="orders" stripedRows :paginator="true" :rows="10">
        <Column field="orderNo" :header="t('procurement.orderNo')" style="width: 150px"></Column>
        <Column field="supplier.name" :header="t('procurement.supplier')">
          <template #body="{ data }">
            {{ data.supplier?.name || data.supplierName || '-' }}
          </template>
        </Column>
        <Column field="orderDate" :header="t('procurement.orderDate')">
          <template #body="{ data }">
            {{ formatDate(data.orderDate) }}
          </template>
        </Column>
        <Column field="expectedDate" :header="t('procurement.expectedDate')">
          <template #body="{ data }">
            {{ formatDate(data.expectedDate) }}
          </template>
        </Column>
        <Column field="totalAmount" :header="t('procurement.totalAmount')">
          <template #body="{ data }">
            ¥{{ formatNumber(data.totalAmount) }}
          </template>
        </Column>
        <Column field="status" :header="t('procurement.status')">
          <template #body="{ data }">
            <span :class="['tag', getStatusClass(data.status)]">{{ getStatusText(data.status) }}</span>
          </template>
        </Column>
        <Column field="warehouse" :header="t('procurement.warehouse')"></Column>
        <Column :header="t('common.actions')" style="width: 180px">
          <template #body="{ data }">
            <Button text severity="info" icon="ri-eye-line" @click="viewDetail(data)" />
            <Button v-if="hasPermission(PROCUREMENT_PERMISSIONS.WRITE) && data.status === 'pending'" text severity="success" icon="ri-check-line" @click="approveOrder(data)" />
            <Button v-if="hasPermission(PROCUREMENT_PERMISSIONS.WRITE)" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="dialogVisible" header="新增采购订单" :style="{ width: '700px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">{{ t('procurement.supplier') }} *</label>
          <Select v-model="form.supplierId" :options="supplierOptions" optionLabel="label" optionValue="value" placeholder="请选择供应商" @change="onSupplierChange" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('procurement.warehouse') }}</label>
          <Select v-model="form.warehouse" :options="warehouseOptions" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('procurement.expectedDate') }}</label>
          <DatePicker v-model="form.expectedDateValue" dateFormat="yy-mm-dd" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">{{ t('common.remark') }}</label>
          <Textarea v-model="form.remark" rows="2" />
        </div>
      </div>

      <h4 style="margin: 20px 0 12px;">{{ t('procurement.items') }}</h4>
      <DataTable :value="form.items" stripedRows>
        <Column field="productName" :header="t('procurement.productName')">
          <template #body="{ index }">
            <InputText v-model="form.items[index].productName" />
          </template>
        </Column>
        <Column field="quantity" :header="t('procurement.quantity')" style="width: 100px">
          <template #body="{ index }">
            <InputNumber v-model="form.items[index].quantity" :min="1" @input="calcAmount(index)" />
          </template>
        </Column>
        <Column field="unit" :header="t('procurement.unit')" style="width: 80px">
          <template #body="{ index }">
            <InputText v-model="form.items[index].unit" />
          </template>
        </Column>
        <Column field="unitPrice" :header="t('procurement.unitPrice')" style="width: 120px">
          <template #body="{ index }">
            <InputNumber v-model="form.items[index].unitPrice" mode="currency" currency="CNY" @input="calcAmount(index)" />
          </template>
        </Column>
        <Column field="amount" :header="t('procurement.amount')" style="width: 100px">
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

    <Dialog v-model:visible="detailVisible" header="采购订单详情" :style="{ width: '700px' }">
      <div v-if="currentOrder" class="order-detail">
        <div class="detail-row">
          <span class="detail-label">{{ t('procurement.orderNo') }}:</span>
          <span>{{ currentOrder.orderNo }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('procurement.supplier') }}:</span>
          <span>{{ currentOrder.supplier?.name || currentOrder.supplierName || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('procurement.orderDate') }}:</span>
          <span>{{ formatDate(currentOrder.orderDate) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('procurement.expectedDate') }}:</span>
          <span>{{ formatDate(currentOrder.expectedDate) }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">{{ t('procurement.status') }}:</span>
          <span :class="['tag', getStatusClass(currentOrder.status)]">{{ getStatusText(currentOrder.status) }}</span>
        </div>
        <div class="detail-row" v-if="currentOrder.remark">
          <span class="detail-label">{{ t('common.remark') }}:</span>
          <span>{{ currentOrder.remark }}</span>
        </div>

        <h4 style="margin: 16px 0 8px;">{{ t('procurement.items') }}</h4>
        <DataTable :value="currentOrder.items" stripedRows>
          <Column field="productName" :header="t('procurement.productName')"></Column>
          <Column field="quantity" :header="t('procurement.quantity')"></Column>
          <Column field="unit" :header="t('procurement.unit')"></Column>
          <Column field="unitPrice" :header="t('procurement.unitPrice')">
            <template #body="{ data }">¥{{ data.unitPrice }}</template>
          </Column>
          <Column field="amount" :header="t('procurement.amount')">
            <template #body="{ data }">¥{{ formatNumber(data.amount) }}</template>
          </Column>
        </DataTable>
        
        <div class="order-total">
          <span>{{ t('procurement.totalAmount') }}: </span>
          <strong>¥{{ formatNumber(currentOrder.totalAmount) }}</strong>
        </div>
      </div>
    </Dialog>

    <Dialog v-model:visible="approveDialogVisible" header="审批确认" :style="{ width: '500px' }">
      <div class="approve-form">
        <p style="margin-bottom: 16px; color: var(--text-muted);">确认审批通过该采购订单？</p>
        <div class="form-group">
          <label class="form-label">附件链接（可选）</label>
          <InputText v-model="approveForm.attachmentUrl" placeholder="请输入附件链接地址" style="width: 100%" />
        </div>
        <div class="form-group">
          <label class="form-label">审批备注（可选）</label>
          <Textarea v-model="approveForm.remark" rows="3" placeholder="请输入审批备注" style="width: 100%" />
        </div>
      </div>
      <div class="dialog-footer">
        <Button label="取消" severity="secondary" @click="approveDialogVisible = false" />
        <Button label="确认审批" severity="success" @click="handleApproveConfirm" />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { procurementApi, supplierApi } from '@/api'
import type { ProcurementOrder, ProcurementItem, Supplier } from '@/types'
import { PROCUREMENT_PERMISSIONS } from '@/config/permissions'
import { usePermission } from '@/utils/usePermission'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const { hasPermission, initUserPermissions } = usePermission()

const loading = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const approveDialogVisible = ref(false)
const statusFilter = ref('all')
const currentOrder = ref<ProcurementOrder | null>(null)
const orders = ref<ProcurementOrder[]>([])
const suppliers = ref<Supplier[]>([])

const approveForm = ref({
  orderId: '',
  attachmentUrl: '',
  remark: ''
})

const form = ref({
  supplierId: '',
  supplierName: '',
  warehouse: '上海仓库',
  expectedDateValue: null as Date | null,
  remark: '',
  items: [] as ProcurementItem[]
})

const statusOptions = [
  { label: t('common.all'), value: 'all' },
  { label: t('procurement.statusPending'), value: 'pending' },
  { label: t('procurement.statusApproved'), value: 'approved' },
  { label: t('procurement.statusPurchasing'), value: 'purchasing' },
  { label: t('procurement.statusReceived'), value: 'received' }
]

const warehouseOptions = ['上海仓库', '深圳仓库', '杭州仓库', '广州仓库']

const supplierOptions = ref<{ label: string; value: string }[]>([])

const formatNumber = (num: number) => num.toLocaleString('zh-CN')

const formatDate = (date: string | Date | null) => {
  if (!date) return '-'
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审批',
    approved: '已审批',
    purchasing: '采购中',
    received: '已入库',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'info',
    purchasing: 'info',
    received: 'success',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res: any = await procurementApi.list(1, 100, statusFilter.value === 'all' ? '' : statusFilter.value)
    if (res.code === 0) {
      orders.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取采购订单失败', error)
  } finally {
    loading.value = false
  }
}

const fetchSuppliers = async () => {
  try {
    const res: any = await supplierApi.list(1, 100, '')
    if (res.code === 0) {
      suppliers.value = res.data.list || []
      supplierOptions.value = suppliers.value.map(s => ({ label: s.name, value: s.id }))
    }
  } catch (error) {
    console.error('获取供应商列表失败', error)
  }
}

const onSupplierChange = (event: any) => {
  const supplier = suppliers.value.find(s => s.id === event.value)
  form.value.supplierName = supplier?.name || ''
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
    receivedQty: 0
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
    supplierId: '',
    supplierName: '',
    warehouse: '上海仓库',
    expectedDateValue: null,
    remark: '',
    items: []
  }
  dialogVisible.value = true
}

const viewDetail = async (order: ProcurementOrder) => {
  try {
    const res: any = await procurementApi.get(String(order.id))
    if (res.code === 0) {
      currentOrder.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error('获取订单详情失败', error)
  }
}

const approveOrder = (order: ProcurementOrder) => {
  approveForm.value = {
    orderId: String(order.id),
    attachmentUrl: '',
    remark: ''
  }
  approveDialogVisible.value = true
}

const handleApproveConfirm = async () => {
  try {
    const res: any = await procurementApi.updateStatus(approveForm.value.orderId, 'approved', {
      attachmentUrl: approveForm.value.attachmentUrl,
      remark: approveForm.value.remark
    })
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '审批成功', life: 5000 })
      approveDialogVisible.value = false
      fetchOrders()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '审批失败', life: 5000 })
  }
}

const handleSave = async () => {
  if (!form.value.supplierId) {
    toast.add({ severity: 'warn', summary: '请选择供应商', life: 5000 })
    return
  }
  
  const totalAmount = form.value.items.reduce((sum, item) => sum + item.amount, 0)
  const expectedDate = form.value.expectedDateValue
    ? `${form.value.expectedDateValue.getFullYear()}-${String(form.value.expectedDateValue.getMonth() + 1).padStart(2, '0')}-${String(form.value.expectedDateValue.getDate()).padStart(2, '0')}`
    : ''
  
  try {
    const res: any = await procurementApi.create({
      supplierId: form.value.supplierId,
      supplierName: form.value.supplierName,
      warehouse: form.value.warehouse,
      expectedDate,
      remark: form.value.remark,
      items: form.value.items,
      totalAmount,
      status: 'pending'
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

const handleDelete = (order: ProcurementOrder) => {
  confirm.require({
    message: t('common.confirmDelete'),
    header: t('common.confirm'),
    accept: async () => {
      try {
        const res: any = await procurementApi.delete(String(order.id))
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
  fetchSuppliers()
})
</script>

<style scoped>
.order-detail {
  .detail-row {
    display: flex;
    gap: 12px;
    margin-bottom: 8px;
    .detail-label {
      width: 100px;
      color: var(--text-muted);
    }
  }
  .order-total {
    text-align: right;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--border);
    font-size: 16px;
    strong {
      font-size: 20px;
      color: var(--primary);
    }
  }
}
</style>