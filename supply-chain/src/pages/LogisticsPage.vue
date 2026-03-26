<template>
  <div class="logistics-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('logistics.title') }}</h1>
      <p class="page-subtitle">实时追踪物流状态和配送进度</p>
    </header>

    <div class="toolbar">
      <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.status')" style="width: 150px" />
      <InputText v-model="searchKey" :placeholder="'搜索运单号'" style="width: 200px" />
    </div>

    <div class="card">
      <DataTable :value="filteredOrders" stripedRows :paginator="true" :rows="10" scrollable scrollHeight="600px">
        <Column field="trackingNo" :header="t('logistics.trackingNo')" style="width: 150px"></Column>
        <Column field="carrier" :header="t('logistics.carrier')"></Column>
        <Column field="salesOrderNo" :header="t('logistics.orderNo')"></Column>
        <Column field="receiverName" :header="t('logistics.receiver')"></Column>
        <Column field="receiverPhone" :header="t('logistics.receiverPhone')"></Column>
        <Column field="weight" :header="t('logistics.weight')">
          <template #body="{ data }">{{ data.weight }}kg</template>
        </Column>
        <Column field="shippingFee" :header="t('logistics.shippingFee')">
          <template #body="{ data }">¥{{ data.shippingFee }}</template>
        </Column>
        <Column field="status" :header="t('logistics.status')">
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

    <Dialog v-model:visible="detailVisible" header="物流详情" :style="{ width: '600px' }">
      <div v-if="currentOrder" class="logistics-detail">
        <div class="detail-grid">
          <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.trackingNo') }}:</span>
              <span>{{ currentOrder.trackingNo }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.carrier') }}:</span>
              <span>{{ currentOrder.carrier }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.status') }}:</span>
              <span :class="['tag', getStatusClass(currentOrder.status)]">{{ getStatusText(currentOrder.status) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.weight') }}:</span>
              <span>{{ currentOrder.weight }}kg</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.shippingFee') }}:</span>
              <span>¥{{ currentOrder.shippingFee }}</span>
            </div>
          </div>

          <div class="detail-section">
            <h4>发件人信息</h4>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.sender') }}:</span>
              <span>{{ currentOrder.senderName }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.senderPhone') }}:</span>
              <span>{{ currentOrder.senderPhone }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.senderAddress') }}:</span>
              <span>{{ currentOrder.senderAddress }}</span>
            </div>
          </div>

          <div class="detail-section">
            <h4>收件人信息</h4>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.receiver') }}:</span>
              <span>{{ currentOrder.receiverName }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.receiverPhone') }}:</span>
              <span>{{ currentOrder.receiverPhone }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">{{ t('logistics.receiverAddress') }}:</span>
              <span>{{ currentOrder.receiverAddress }}</span>
            </div>
          </div>
        </div>

        <div class="timeline-section">
          <h4>{{ t('logistics.timeline') }}</h4>
          <div class="timeline">
            <div v-for="(item, index) in currentOrder.timeline" :key="index" class="timeline-item">
              <div class="timeline-dot" :class="{ active: index === 0 }"></div>
              <div class="timeline-content">
                <div class="timeline-time">{{ item.time }}</div>
                <div class="timeline-status">{{ item.status }}</div>
                <div class="timeline-location">{{ item.location }}</div>
                <div class="timeline-desc">{{ item.description }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { logisticsApi } from '@/api'
import type { LogisticsOrder } from '@/types'
import { usePermission } from '@/utils/usePermission'

const { t } = useI18n()
const { initUserPermissions } = usePermission()

const loading = ref(false)
const detailVisible = ref(false)
const statusFilter = ref('all')
const searchKey = ref('')
const currentOrder = ref<LogisticsOrder | null>(null)
const orders = ref<LogisticsOrder[]>([])

const statusOptions = [
  { label: t('common.all'), value: 'all' },
  { label: t('logistics.statusPending'), value: 'pending' },
  { label: t('logistics.statusPicked'), value: 'picked' },
  { label: t('logistics.statusInTransit'), value: 'in_transit' },
  { label: t('logistics.statusDelivering'), value: 'delivering' },
  { label: t('logistics.statusDelivered'), value: 'delivered' }
]

const filteredOrders = computed(() => {
  let result = orders.value
  if (statusFilter.value !== 'all') {
    result = result.filter(o => o.status === statusFilter.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(o => o.trackingNo.toLowerCase().includes(key))
  }
  return result
})

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待发货',
    picked: '已揽收',
    in_transit: '运输中',
    delivering: '派送中',
    delivered: '已签收',
    returned: '已退回'
  }
  return map[status] || status
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    picked: 'info',
    in_transit: 'info',
    delivering: 'info',
    delivered: 'success',
    returned: 'danger'
  }
  return map[status] || 'info'
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const res: any = await logisticsApi.list(1, 100, statusFilter.value === 'all' ? '' : statusFilter.value)
    if (res.code === 0) {
      orders.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取物流订单失败', error)
  } finally {
    loading.value = false
  }
}

const viewDetail = async (order: LogisticsOrder) => {
  try {
    const res: any = await logisticsApi.get(String(order.id))
    if (res.code === 0) {
      currentOrder.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error('获取物流详情失败', error)
  }
}

onMounted(() => {
  initUserPermissions()
  fetchOrders()
})
</script>

<style scoped>
.logistics-detail {
  .detail-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
    margin-bottom: 24px;
    
    @media (max-width: 900px) {
      grid-template-columns: 1fr;
    }
  }
  
  .detail-section {
    h4 {
      margin-bottom: 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid var(--border);
      font-size: 14px;
      color: var(--text-muted);
    }
  }
  
  .detail-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
    font-size: 13px;
    
    .detail-label {
      color: var(--text-muted);
      min-width: 70px;
    }
  }
}

.timeline-section {
  h4 {
    margin-bottom: 16px;
    font-size: 14px;
    color: var(--text-muted);
  }
}

.timeline {
  position: relative;
  padding-left: 24px;
  
  .timeline-item {
    position: relative;
    padding-bottom: 20px;
    
    &:last-child {
      padding-bottom: 0;
    }
  }
  
  .timeline-dot {
    position: absolute;
    left: -24px;
    top: 4px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: var(--border);
    
    &.active {
      background: var(--primary);
      box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.2);
    }
  }
  
  .timeline-content {
    .timeline-time {
      font-size: 12px;
      color: var(--text-muted);
      margin-bottom: 4px;
    }
    
    .timeline-status {
      font-weight: 600;
      margin-bottom: 2px;
    }
    
    .timeline-location {
      font-size: 13px;
      color: var(--text-muted);
      margin-bottom: 2px;
    }
    
    .timeline-desc {
      font-size: 13px;
    }
  }
}
</style>