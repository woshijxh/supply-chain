<template>
  <div class="inventory-log-page">
    <header class="page-header">
      <h1 class="page-title">库存流水</h1>
      <p class="page-subtitle">查看所有库存变动记录</p>
    </header>

    <div class="toolbar">
      <Select v-model="logTypeFilter" :options="logTypeOptions" optionLabel="label" optionValue="value" placeholder="全部类型" style="width: 150px" @change="fetchLogs" />
      <Select v-model="productFilter" :options="productOptions" optionLabel="label" optionValue="value" placeholder="全部产品" style="width: 200px" @change="fetchLogs" filter />
    </div>

    <div class="card">
      <DataTable :value="inventoryLogs" stripedRows :paginator="true" :rows="15" :loading="loading" scrollable scrollHeight="600px">
        <Column field="createdAt" header="时间" style="width: 160px">
          <template #body="{ data }">
            {{ formatTime(data.createdAt) }}
          </template>
        </Column>
        <Column field="productName" header="产品" style="width: 150px"></Column>
        <Column field="type" header="类型" style="width: 100px">
          <template #body="{ data }">
            <span :class="['log-tag', getLogTypeClass(data.type)]">{{ getLogTypeText(data.type) }}</span>
          </template>
        </Column>
        <Column field="quantity" header="数量" style="width: 100px">
          <template #body="{ data }">
            <span :class="{ 'text-success': data.type === 'in' || data.type === 'unlock', 'text-danger': data.type === 'out' || data.type === 'lock' }">
              {{ data.type === 'in' || data.type === 'unlock' ? '+' : '-' }}{{ data.quantity }}
            </span>
          </template>
        </Column>
        <Column field="beforeQty" header="变更前" style="width: 100px"></Column>
        <Column field="afterQty" header="变更后" style="width: 100px"></Column>
        <Column field="warehouse" header="仓库" style="width: 120px"></Column>
        <Column field="refNo" header="关联单号" style="width: 150px">
          <template #body="{ data }">
            <span v-if="data.refNo" class="ref-no">{{ data.refNo }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </Column>
        <Column field="operator" header="操作人" style="width: 100px">
          <template #body="{ data }">
            {{ data.operator || '-' }}
          </template>
        </Column>
        <Column field="remark" header="备注"></Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { inventoryLogApi, productApi } from '@/api'
import type { InventoryLog } from '@/types'

const loading = ref(false)
const inventoryLogs = ref<InventoryLog[]>([])
const logTypeFilter = ref('')
const productFilter = ref('')
const products = ref<any[]>([])

const logTypeOptions = [
  { label: '全部类型', value: '' },
  { label: '入库', value: 'in' },
  { label: '出库', value: 'out' },
  { label: '锁定', value: 'lock' },
  { label: '解锁', value: 'unlock' }
]

const productOptions = computed(() => {
  const options = [{ label: '全部产品', value: '' }]
  return options.concat(products.value.map((p: any) => ({
    label: p.name,
    value: String(p.id)
  })))
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const res: any = await inventoryLogApi.list(1, 200, productFilter.value, logTypeFilter.value)
    if (res.code === 0) {
      inventoryLogs.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取库存流水失败', error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res: any = await productApi.list(1, 200, '')
    if (res.code === 0) {
      products.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

const formatTime = (time: string) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

const getLogTypeText = (type: string) => {
  const map: Record<string, string> = {
    in: '入库',
    out: '出库',
    lock: '锁定',
    unlock: '解锁'
  }
  return map[type] || type
}

const getLogTypeClass = (type: string) => {
  const map: Record<string, string> = {
    in: 'success',
    out: 'danger',
    lock: 'warning',
    unlock: 'info'
  }
  return map[type] || 'secondary'
}

onMounted(() => {
  fetchProducts()
  fetchLogs()
})
</script>

<style scoped>
.log-tag {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.log-tag.success {
  background: rgba(34, 197, 94, 0.1);
  color: #16a34a;
}

.log-tag.danger {
  background: rgba(239, 68, 68, 0.1);
  color: #dc2626;
}

.log-tag.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #d97706;
}

.log-tag.info {
  background: rgba(59, 130, 246, 0.1);
  color: #2563eb;
}

.log-tag.secondary {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.text-success {
  color: #16a34a;
  font-weight: 600;
}

.text-danger {
  color: #dc2626;
  font-weight: 600;
}

.text-muted {
  color: var(--text-muted);
}

.ref-no {
  font-family: monospace;
  font-size: 12px;
  background: var(--surface-ground);
  padding: 2px 6px;
  border-radius: 4px;
}
</style>