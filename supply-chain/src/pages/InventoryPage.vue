<template>
  <div class="inventory-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('inventory.title') }}</h1>
      <p class="page-subtitle">实时监控库存状态和库存预警</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(INVENTORY_PERMISSIONS.WRITE)" :label="t('inventory.stockIn')" icon="ri-login-box-line" severity="success" @click="openStockDialog('in')" />
      <Button v-if="hasPermission(INVENTORY_PERMISSIONS.WRITE)" :label="t('inventory.stockOut')" icon="ri-logout-box-line" severity="warning" @click="openStockDialog('out')" />
      <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.status')" style="width: 150px" />
      <Select v-model="warehouseFilter" :options="warehouseOptions" optionLabel="label" optionValue="value" placeholder="仓库" style="width: 150px" />
    </div>

    <div class="stat-grid" style="margin-bottom: 20px;">
      <div class="stat-card">
        <div class="stat-icon blue"><i class="ri-box-3-line"></i></div>
        <div class="stat-content">
          <div class="stat-label">总SKU数</div>
          <div class="stat-value">{{ totalSku }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green"><i class="ri-checkbox-circle-line"></i></div>
        <div class="stat-content">
          <div class="stat-label">正常库存</div>
          <div class="stat-value">{{ normalCount }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange"><i class="ri-error-warning-line"></i></div>
        <div class="stat-content">
          <div class="stat-label">库存不足</div>
          <div class="stat-value">{{ lowCount }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon purple"><i class="ri-archive-line"></i></div>
        <div class="stat-content">
          <div class="stat-label">库存过量</div>
          <div class="stat-value">{{ overCount }}</div>
        </div>
      </div>
    </div>

    <div class="card">
      <DataTable :value="filteredItems" stripedRows :paginator="true" :rows="10">
        <Column field="productCode" :header="t('inventory.productCode')" style="width: 120px"></Column>
        <Column field="productName" :header="t('inventory.productName')"></Column>
        <Column field="sku" :header="t('inventory.sku')" style="width: 140px"></Column>
        <Column field="category" :header="t('inventory.category')"></Column>
        <Column field="warehouse" :header="t('inventory.warehouse')"></Column>
        <Column field="quantity" :header="t('inventory.quantity')" style="width: 100px">
          <template #body="{ data }">
            <span :class="{ 'low-stock': data.status === 'low' }">{{ data.quantity }}</span>
          </template>
        </Column>
        <Column field="availableQty" :header="t('inventory.availableQty')" style="width: 100px"></Column>
        <Column field="lockedQty" :header="t('inventory.lockedQty')" style="width: 100px"></Column>
        <Column field="minStock" :header="t('inventory.minStock')" style="width: 80px"></Column>
        <Column field="maxStock" :header="t('inventory.maxStock')" style="width: 80px"></Column>
        <Column field="costPrice" :header="t('inventory.costPrice')">
          <template #body="{ data }">¥{{ data.costPrice }}</template>
        </Column>
        <Column field="salePrice" :header="t('inventory.salePrice')">
          <template #body="{ data }">¥{{ data.salePrice }}</template>
        </Column>
        <Column field="status" :header="t('inventory.status')">
          <template #body="{ data }">
            <span :class="['tag', getStatusClass(data.status)]">{{ getStatusText(data.status) }}</span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width: 100px">
          <template #body="{ data }">
            <Button text severity="info" icon="ri-eye-line" @click="viewDetail(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="detailVisible" header="库存详情" :style="{ width: '500px' }">
      <div v-if="currentItem" class="inventory-detail">
        <div class="detail-grid">
          <div class="detail-item">
            <label>产品编码</label>
            <span>{{ currentItem.productCode }}</span>
          </div>
          <div class="detail-item">
            <label>产品名称</label>
            <span>{{ currentItem.productName }}</span>
          </div>
          <div class="detail-item">
            <label>SKU</label>
            <span>{{ currentItem.sku }}</span>
          </div>
          <div class="detail-item">
            <label>分类</label>
            <span>{{ currentItem.category }}</span>
          </div>
          <div class="detail-item">
            <label>仓库</label>
            <span>{{ currentItem.warehouse }}</span>
          </div>
          <div class="detail-item">
            <label>库位</label>
            <span>{{ currentItem.location }}</span>
          </div>
          <div class="detail-item">
            <label>库存数量</label>
            <span>{{ currentItem.quantity }} {{ currentItem.unit }}</span>
          </div>
          <div class="detail-item">
            <label>可用数量</label>
            <span>{{ currentItem.availableQty }} {{ currentItem.unit }}</span>
          </div>
          <div class="detail-item">
            <label>锁定数量</label>
            <span>{{ currentItem.lockedQty }} {{ currentItem.unit }}</span>
          </div>
          <div class="detail-item">
            <label>库存范围</label>
            <span>{{ currentItem.minStock }} - {{ currentItem.maxStock }}</span>
          </div>
          <div class="detail-item">
            <label>成本价</label>
            <span>¥{{ currentItem.costPrice }}</span>
          </div>
          <div class="detail-item">
            <label>销售价</label>
            <span>¥{{ currentItem.salePrice }}</span>
          </div>
          <div class="detail-item">
            <label>批次号</label>
            <span>{{ currentItem.batchNo }}</span>
          </div>
          <div class="detail-item">
            <label>状态</label>
            <span :class="['tag', getStatusClass(currentItem.status)]">{{ getStatusText(currentItem.status) }}</span>
          </div>
        </div>
      </div>
    </Dialog>

    <Dialog v-model:visible="stockDialogVisible" :header="stockType === 'in' ? '入库' : '出库'" :style="{ width: '450px' }">
      <div class="form-grid">
        <div class="form-group" style="grid-column: span 2" v-if="stockType === 'in'">
          <label class="form-label">关联采购订单（可选）</label>
          <Select v-model="stockForm.procurementId" :options="procurementOptions" optionLabel="label" optionValue="value" placeholder="选择采购订单自动带入产品和仓库" style="width: 100%" @change="onProcurementChange" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">选择产品 *</label>
          <Select v-model="stockForm.productId" :options="productOptions" optionLabel="label" optionValue="value" placeholder="请选择产品" style="width: 100%" :disabled="!!stockForm.procurementId" />
        </div>
        <div class="form-group">
          <label class="form-label">数量 *</label>
          <InputNumber v-model="stockForm.quantity" :min="1" style="width: 100%" />
        </div>
        <div class="form-group">
          <label class="form-label">仓库</label>
          <Select v-model="stockForm.warehouse" :options="warehouseList" placeholder="选择仓库" style="width: 100%" :disabled="!!stockForm.procurementId" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">备注</label>
          <Textarea v-model="stockForm.remark" rows="2" />
        </div>
      </div>
      <div class="dialog-footer">
        <Button :label="t('common.cancel')" severity="secondary" @click="stockDialogVisible = false" />
        <Button :label="t('common.confirm')" @click="handleStock" />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { inventoryApi, productApi, procurementApi } from '@/api'
import type { InventoryItem } from '@/types'
import { INVENTORY_PERMISSIONS } from '@/config/permissions'
import { usePermission } from '@/utils/usePermission'

const { t } = useI18n()
const toast = useToast()
const { hasPermission, initUserPermissions } = usePermission()

const loading = ref(false)
const detailVisible = ref(false)
const stockDialogVisible = ref(false)
const stockType = ref<'in' | 'out'>('in')
const statusFilter = ref('all')
const warehouseFilter = ref('all')
const currentItem = ref<InventoryItem | null>(null)
const inventoryItems = ref<InventoryItem[]>([])
const inventoryStats = ref({ total: 0, normal: 0, low: 0, over: 0 })
const products = ref<any[]>([])
const procurements = ref<any[]>([])

const stockForm = ref({
  productId: '',
  quantity: 1,
  warehouse: '上海仓库',
  remark: '',
  procurementId: '',
  procurementItemId: ''
})

const statusOptions = [
  { label: t('common.all'), value: 'all' },
  { label: t('inventory.statusNormal'), value: 'normal' },
  { label: t('inventory.statusLow'), value: 'low' },
  { label: t('inventory.statusOver'), value: 'over' },
  { label: t('inventory.statusLocked'), value: 'locked' }
]

const warehouseOptions = [
  { label: '全部', value: 'all' },
  { label: '上海仓库', value: '上海仓库' },
  { label: '深圳仓库', value: '深圳仓库' },
  { label: '杭州仓库', value: '杭州仓库' },
  { label: '广州仓库', value: '广州仓库' }
]

const warehouseList = ['上海仓库', '深圳仓库', '杭州仓库', '广州仓库']

const productOptions = computed(() => {
  return products.value.map((product: any) => ({
    label: `${product.name} (${product.code})`,
    value: product.id
  }))
})

const procurementOptions = computed(() => {
  return procurements.value.map((proc: any) => ({
    label: `采购#${proc.id} - ${proc.supplierName} (${proc.warehouse})`,
    value: proc.id,
    items: proc.items || []
  }))
})

const filteredItems = computed(() => {
  let items = inventoryItems.value
  if (statusFilter.value !== 'all') {
    items = items.filter(i => i.status === statusFilter.value)
  }
  if (warehouseFilter.value !== 'all') {
    items = items.filter(i => i.warehouse === warehouseFilter.value)
  }
  return items
})

const totalSku = computed(() => inventoryStats.value.total)
const normalCount = computed(() => inventoryStats.value.normal)
const lowCount = computed(() => inventoryStats.value.low)
const overCount = computed(() => inventoryStats.value.over)

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    normal: '正常',
    low: '库存不足',
    over: '库存过量',
    locked: '已锁定'
  }
  return map[status] || status
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    normal: 'success',
    low: 'warning',
    over: 'info',
    locked: 'danger'
  }
  return map[status] || 'info'
}

const fetchInventory = async () => {
  loading.value = true
  try {
    const res: any = await inventoryApi.list(1, 100, statusFilter.value === 'all' ? '' : statusFilter.value, warehouseFilter.value === 'all' ? '' : warehouseFilter.value)
    if (res.code === 0) {
      inventoryItems.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取库存列表失败', error)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res: any = await inventoryApi.stats()
    if (res.code === 0) {
      inventoryStats.value = res.data
    }
  } catch (error) {
    console.error('获取库存统计失败', error)
  }
}

const fetchProducts = async () => {
  try {
    const res: any = await productApi.list(1, 100, '')
    if (res.code === 0) {
      products.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

const fetchPendingProcurements = async () => {
  try {
    const res: any = await procurementApi.list(1, 50, 'purchasing')
    if (res.code === 0) {
      procurements.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取采购订单失败', error)
  }
}

const onProcurementChange = (event: any) => {
  const procurementId = event?.value || ''
  if (!procurementId) {
    stockForm.value.procurementItemId = ''
    return
  }
  const proc = procurements.value.find((p: any) => p.id === procurementId)
  if (proc && proc.items && proc.items.length > 0) {
    const item = proc.items[0]
    stockForm.value.productId = item.productId
    stockForm.value.warehouse = proc.warehouse
    stockForm.value.procurementItemId = item.id
    const pendingQty = item.quantity - item.receivedQty
    if (pendingQty > 0) {
      stockForm.value.quantity = pendingQty
    }
  }
}

const viewDetail = (item: InventoryItem) => {
  currentItem.value = item
  detailVisible.value = true
}

const openStockDialog = (type: 'in' | 'out') => {
  stockType.value = type
  stockForm.value = { productId: '', quantity: 1, warehouse: '上海仓库', remark: '', procurementId: '', procurementItemId: '' }
  stockDialogVisible.value = true
  if (type === 'in') {
    fetchPendingProcurements()
  }
}

const handleStock = async () => {
  if (!stockForm.value.productId || stockForm.value.quantity < 1) {
    toast.add({ severity: 'warn', summary: '请选择产品并填写数量', life: 5000 })
    return
  }

  try {
    const apiCall = stockType.value === 'in' ? inventoryApi.stockIn : inventoryApi.stockOut
    const res: any = await apiCall(
      stockForm.value.productId,
      stockForm.value.quantity,
      stockForm.value.warehouse,
      stockForm.value.procurementId || undefined,
      stockForm.value.procurementItemId || undefined
    )
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: `${stockType.value === 'in' ? '入库' : '出库'}成功`, life: 5000 })
      stockDialogVisible.value = false
      fetchInventory()
      fetchStats()
    } else {
      toast.add({ severity: 'error', summary: res.message || '操作失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.message || '操作失败', life: 5000 })
  }
}

onMounted(() => {
  initUserPermissions()
  fetchInventory()
  fetchStats()
  fetchProducts()
})
</script>

<style scoped>
.low-stock {
  color: var(--warning);
  font-weight: 600;
}

.inventory-detail {
  .detail-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
  .detail-item {
    label {
      display: block;
      font-size: 12px;
      color: var(--text-muted);
      margin-bottom: 4px;
    }
    span {
      font-weight: 500;
    }
  }
}
</style>