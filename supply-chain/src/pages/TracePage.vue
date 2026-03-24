<template>
  <div class="trace-page">
    <header class="page-header">
      <h1 class="page-title">商品追溯</h1>
      <p class="page-subtitle">通过产品编码、SKU或批次号追溯商品流转全链路</p>
    </header>

    <!-- 搜索区域 -->
    <div class="search-section">
      <div class="search-box">
        <InputText v-model="searchCode" placeholder="请输入产品编码、SKU或批次号" style="flex: 1" @keyup.enter="handleSearch" />
        <Button label="追溯" icon="ri-search-line" @click="handleSearch" :loading="loading" />
      </div>
    </div>

    <!-- 追溯结果 -->
    <template v-if="traceResult">
      <!-- 产品信息卡片 -->
      <div class="product-card" v-if="traceResult.product">
        <div class="product-main">
          <div class="product-icon">
            <i class="ri-box-3-line"></i>
          </div>
          <div class="product-info">
            <h2>{{ traceResult.product.name }}</h2>
            <div class="product-tags">
              <Tag>编码: {{ traceResult.product.code }}</Tag>
              <Tag>SKU: {{ traceResult.product.sku }}</Tag>
              <Tag severity="success">库存: {{ traceResult.product.totalStock }} {{ traceResult.product.unit }}</Tag>
            </div>
          </div>
        </div>
        <div class="product-stats">
          <div class="stat-item">
            <span class="stat-label">成本价</span>
            <span class="stat-value">¥{{ traceResult.product.costPrice }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">销售价</span>
            <span class="stat-value primary">¥{{ traceResult.product.salePrice }}</span>
          </div>
        </div>
      </div>

      <!-- 统计概览 -->
      <div class="overview-grid">
        <div class="overview-card" :class="{ active: activeTab === 'procurement' }" @click="activeTab = 'procurement'">
          <div class="overview-icon blue"><i class="ri-shopping-cart-line"></i></div>
          <div class="overview-content">
            <div class="overview-value">{{ traceResult.procurement.length }}</div>
            <div class="overview-label">采购记录</div>
          </div>
        </div>
        <div class="overview-card" :class="{ active: activeTab === 'inventory' }" @click="activeTab = 'inventory'">
          <div class="overview-icon orange"><i class="ri-exchange-line"></i></div>
          <div class="overview-content">
            <div class="overview-value">{{ traceResult.inventory.length }}</div>
            <div class="overview-label">库存变动</div>
          </div>
        </div>
        <div class="overview-card" :class="{ active: activeTab === 'sales' }" @click="activeTab = 'sales'">
          <div class="overview-icon green"><i class="ri-file-list-3-line"></i></div>
          <div class="overview-content">
            <div class="overview-value">{{ traceResult.sales.length }}</div>
            <div class="overview-label">销售记录</div>
          </div>
        </div>
        <div class="overview-card" :class="{ active: activeTab === 'logistics' }" @click="activeTab = 'logistics'">
          <div class="overview-icon purple"><i class="ri-truck-line"></i></div>
          <div class="overview-content">
            <div class="overview-value">{{ traceResult.logistics.length }}</div>
            <div class="overview-label">物流记录</div>
          </div>
        </div>
        <div class="overview-card" :class="{ active: activeTab === 'returns' }" @click="activeTab = 'returns'">
          <div class="overview-icon red"><i class="ri-arrow-go-back-line"></i></div>
          <div class="overview-content">
            <div class="overview-value">{{ traceResult.returns.length }}</div>
            <div class="overview-label">退货记录</div>
          </div>
        </div>
      </div>

      <!-- 筛选和操作栏 -->
      <div class="filter-bar">
        <div class="filter-left">
          <Select v-model="filterType" :options="filterOptions" optionLabel="label" optionValue="value" placeholder="筛选类型" style="width: 140px" />
          <DatePicker v-model="dateRange" selectionMode="range" :maxDate="new Date()" placeholder="选择时间范围" style="width: 220px" />
        </div>
        <div class="filter-right">
          <Button label="导出" icon="ri-download-line" severity="secondary" outlined size="small" />
        </div>
      </div>

      <!-- 详细数据表格 -->
      <div class="detail-section">
        <TabView v-model:activeIndex="activeTabIndex">
          <TabPanel header="采购记录" value="0">
            <DataTable :value="traceResult.procurement" stripedRows paginator :rows="10" :rowsPerPageOptions="[10, 20, 50]" tableStyle="min-width: 50rem">
              <Column field="orderNo" header="采购单号" sortable></Column>
              <Column field="supplierName" header="供应商" sortable></Column>
              <Column field="quantity" header="采购数量" sortable></Column>
              <Column field="receivedQty" header="已收货" sortable></Column>
              <Column field="unitPrice" header="单价" sortable>
                <template #body="{ data }">¥{{ data.unitPrice }}</template>
              </Column>
              <Column field="warehouse" header="仓库"></Column>
              <Column field="status" header="状态" sortable>
                <template #body="{ data }">
                  <Tag :value="getProcurementStatus(data.status)" :severity="getProcurementSeverity(data.status)" />
                </template>
              </Column>
            </DataTable>
          </TabPanel>

          <TabPanel header="库存变动" value="1">
            <DataTable :value="traceResult.inventory" stripedRows paginator :rows="10" :rowsPerPageOptions="[10, 20, 50]" tableStyle="min-width: 50rem">
              <Column field="createdAt" header="时间" sortable>
                <template #body="{ data }">{{ formatTime(data.createdAt) }}</template>
              </Column>
              <Column field="type" header="类型" sortable>
                <template #body="{ data }">
                  <Tag :value="getInventoryType(data.type)" :severity="getInventorySeverity(data.type)" />
                </template>
              </Column>
              <Column field="quantity" header="数量" sortable>
                <template #body="{ data }">
                  <span :class="data.type === 'in' ? 'text-success' : 'text-danger'">
                    {{ data.type === 'in' ? '+' : '-' }}{{ data.quantity }}
                  </span>
                </template>
              </Column>
              <Column field="warehouse" header="仓库"></Column>
              <Column field="refNo" header="关联单号"></Column>
              <Column field="operator" header="操作人"></Column>
            </DataTable>
          </TabPanel>

          <TabPanel header="销售记录" value="2">
            <DataTable :value="traceResult.sales" stripedRows paginator :rows="10" :rowsPerPageOptions="[10, 20, 50]" tableStyle="min-width: 50rem">
              <Column field="orderNo" header="销售单号" sortable></Column>
              <Column field="customerName" header="客户" sortable></Column>
              <Column field="quantity" header="数量" sortable></Column>
              <Column field="unitPrice" header="单价" sortable>
                <template #body="{ data }">¥{{ data.unitPrice }}</template>
              </Column>
              <Column field="orderDate" header="下单时间" sortable>
                <template #body="{ data }">{{ formatTime(data.orderDate) }}</template>
              </Column>
              <Column field="status" header="状态" sortable>
                <template #body="{ data }">
                  <Tag :value="getSalesStatus(data.status)" :severity="getSalesSeverity(data.status)" />
                </template>
              </Column>
              <Column field="trackingNo" header="物流单号"></Column>
            </DataTable>
          </TabPanel>

          <TabPanel header="物流记录" value="3">
            <DataTable :value="traceResult.logistics" stripedRows paginator :rows="10" :rowsPerPageOptions="[10, 20, 50]" tableStyle="min-width: 50rem">
              <Column field="trackingNo" header="物流单号" sortable></Column>
              <Column field="carrier" header="承运商"></Column>
              <Column field="salesOrderNo" header="销售单号"></Column>
              <Column field="receiverName" header="收货人"></Column>
              <Column field="receiverPhone" header="联系电话"></Column>
              <Column field="status" header="状态" sortable>
                <template #body="{ data }">
                  <Tag :value="getLogisticsStatus(data.status)" />
                </template>
              </Column>
              <Column field="shippingFee" header="运费">
                <template #body="{ data }">¥{{ data.shippingFee }}</template>
              </Column>
            </DataTable>
          </TabPanel>

          <TabPanel header="退货记录" value="4">
            <DataTable :value="traceResult.returns" stripedRows paginator :rows="10" :rowsPerPageOptions="[10, 20, 50]" tableStyle="min-width: 50rem">
              <Column field="returnNo" header="退货单号" sortable></Column>
              <Column field="type" header="类型" sortable>
                <template #body="{ data }">
                  <Tag :value="data.type === 'sales_return' ? '销售退货' : '采购退货'" />
                </template>
              </Column>
              <Column field="relatedName" header="关联方"></Column>
              <Column field="quantity" header="数量" sortable></Column>
              <Column field="reason" header="原因"></Column>
              <Column field="status" header="状态" sortable>
                <template #body="{ data }">
                  <Tag :value="getReturnStatus(data.status)" />
                </template>
              </Column>
            </DataTable>
          </TabPanel>
        </TabView>
      </div>

      <!-- 时间线 -->
      <div class="timeline-section">
        <div class="timeline-header">
          <h3><i class="ri-time-line"></i> 流转时间线</h3>
          <span class="record-count">共 {{ filteredTimeline.length }} 条记录</span>
        </div>

        <!-- 时间线列表 -->
        <div class="timeline-container" v-if="filteredTimeline.length > 0">
          <div class="timeline-list">
            <div class="timeline-item" v-for="(item, index) in filteredTimeline" :key="index">
              <div class="timeline-line">
                <div class="timeline-marker" :class="item.type">
                  <i :class="getIcon(item.type)"></i>
                </div>
                <div class="line-connector" v-if="index < filteredTimeline.length - 1"></div>
              </div>
              <div class="timeline-card" :class="item.type">
                <div class="timeline-top">
                  <span class="timeline-action">{{ item.action }}</span>
                  <span class="timeline-time">{{ formatTime(item.time) }}</span>
                </div>
                <div class="timeline-title">{{ item.title }}</div>
                <div class="timeline-desc">{{ item.description }}</div>
                <div class="timeline-meta" v-if="item.refNo || item.operator">
                  <span v-if="item.refNo"><i class="ri-file-text-line"></i> {{ item.refNo }}</span>
                  <span v-if="item.operator"><i class="ri-user-line"></i> {{ item.operator }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div class="empty-timeline" v-else>
          <i class="ri-inbox-line"></i>
          <p>暂无流转记录</p>
        </div>
      </div>
    </template>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <i class="ri-search-eye-line"></i>
      <p>输入编码开始追溯商品流转记录</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { traceApi } from '@/api'
import type { TraceResult } from '@/types'

const toast = useToast()
const loading = ref(false)
const searchCode = ref('')
const traceResult = ref<TraceResult | null>(null)

// 筛选相关
const filterType = ref('all')
const dateRange = ref<[Date, Date] | null>(null)
const activeTab = ref('procurement')

const filterOptions = [
  { label: '全部', value: 'all' },
  { label: '采购', value: 'procurement' },
  { label: '库存', value: 'inventory' },
  { label: '销售', value: 'sales' },
  { label: '物流', value: 'logistics' },
  { label: '退货', value: 'return' }
]

// Tab 索引映射
const activeTabIndex = computed({
  get: () => {
    const map: Record<string, number> = {
      procurement: 0,
      inventory: 1,
      sales: 2,
      logistics: 3,
      returns: 4
    }
    return map[activeTab.value] || 0
  },
  set: (val: number) => {
    const map: Record<number, string> = {
      0: 'procurement',
      1: 'inventory',
      2: 'sales',
      3: 'logistics',
      4: 'returns'
    }
    activeTab.value = map[val] || 'procurement'
  }
})

// 筛选后的时间线
const filteredTimeline = computed(() => {
  if (!traceResult.value?.timeline) return []

  let items = [...traceResult.value.timeline]

  // 按类型筛选
  if (filterType.value !== 'all') {
    items = items.filter(item => item.type === filterType.value)
  }

  // 按时间范围筛选
  if (dateRange.value && dateRange.value[0] && dateRange.value[1]) {
    const start = new Date(dateRange.value[0])
    start.setHours(0, 0, 0, 0)
    const end = new Date(dateRange.value[1])
    end.setHours(23, 59, 59, 999)

    items = items.filter(item => {
      const itemDate = new Date(item.time)
      return itemDate >= start && itemDate <= end
    })
  }

  return items
})

const handleSearch = async () => {
  if (!searchCode.value.trim()) {
    toast.add({ severity: 'warn', summary: '请输入追溯编码', life: 3000 })
    return
  }

  loading.value = true
  try {
    const res: any = await traceApi.trace(searchCode.value.trim())
    if (res.code === 0) {
      const data = res.data || {}
      traceResult.value = {
        product: data.product || null,
        procurement: data.procurement || [],
        inventory: data.inventory || [],
        sales: data.sales || [],
        logistics: data.logistics || [],
        returns: data.returns || [],
        timeline: data.timeline || []
      }
    } else {
      toast.add({ severity: 'error', summary: res.message || '查询失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.message || '查询失败', life: 5000 })
  } finally {
    loading.value = false
  }
}

// 工具函数
const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const getIcon = (type: string) => {
  const map: Record<string, string> = {
    procurement: 'ri-shopping-cart-line',
    inventory: 'ri-exchange-line',
    sales: 'ri-file-list-3-line',
    logistics: 'ri-truck-line',
    return: 'ri-arrow-go-back-line'
  }
  return map[type] || 'ri-information-line'
}

const getProcurementStatus = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审批', approved: '已审批', purchasing: '采购中',
    received: '已收货', cancelled: '已取消'
  }
  return map[status] || status
}

const getProcurementSeverity = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning', approved: 'info', purchasing: 'info',
    received: 'success', cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getInventoryType = (type: string) => {
  const map: Record<string, string> = {
    in: '入库', out: '出库', lock: '锁定', unlock: '解锁'
  }
  return map[type] || type
}

const getInventorySeverity = (type: string) => {
  const map: Record<string, string> = {
    in: 'success', out: 'danger', lock: 'warning', unlock: 'info'
  }
  return map[type] || 'info'
}

const getSalesStatus = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认', confirmed: '已确认', shipping: '配送中',
    completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

const getSalesSeverity = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning', confirmed: 'info', shipping: 'info',
    completed: 'success', cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getReturnStatus = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理', approved: '已批准', rejected: '已拒绝', completed: '已完成'
  }
  return map[status] || status
}

const getLogisticsStatus = (status: string) => {
  const map: Record<string, string> = {
    pending: '待发货', shipping: '运输中', delivered: '已签收', returned: '已退回'
  }
  return map[status] || status
}
</script>

<style scoped>
.trace-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px;
}

.page-subtitle {
  color: var(--text-muted);
  margin: 0;
}

/* 搜索区域 */
.search-section {
  background: var(--surface-card);
  padding: 20px;
  border-radius: var(--radius-lg);
  margin-bottom: 24px;
}

.search-box {
  display: flex;
  gap: 12px;
}

/* 产品卡片 */
.product-card {
  background: linear-gradient(135deg, var(--primary-50) 0%, var(--surface-card) 100%);
  border: 1px solid var(--primary-100);
  padding: 24px;
  border-radius: var(--radius-lg);
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-main {
  display: flex;
  align-items: center;
  gap: 16px;
}

.product-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.product-info h2 {
  margin: 0 0 8px;
  font-size: 20px;
}

.product-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.product-stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
}

.stat-value.primary {
  color: var(--primary);
}

/* 统计概览 */
.overview-grid {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.overview-card {
  background: var(--surface-card);
  padding: 12px 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid #d1d5db;
  flex: 1;
  min-width: 140px;
  max-width: 180px;
}

.overview-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.overview-card.active {
  border-color: var(--primary);
  background: var(--primary-50);
}

.overview-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: white;
  flex-shrink: 0;
}

.overview-icon.blue { background: var(--blue-500); }
.overview-icon.orange { background: var(--orange-500); }
.overview-icon.green { background: var(--green-500); }
.overview-icon.purple { background: var(--purple-500); }
.overview-icon.red { background: var(--red-500); }

.overview-value {
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
}

.overview-label {
  font-size: 12px;
  color: var(--text-muted);
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.filter-left {
  display: flex;
  gap: 12px;
}

/* 时间线 */
.timeline-section {
  background: var(--surface-card);
  border-radius: var(--radius-lg);
  margin-bottom: 24px;
  overflow: hidden;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--surface-border);
}

.timeline-header h3 {
  margin: 0;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.record-count {
  font-size: 13px;
  color: var(--text-muted);
}

.timeline-container {
  padding: 20px;
  overflow-x: hidden;
}

.timeline-list {
  max-height: 280px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-left: 8px;
}

.timeline-item {
  display: flex;
  gap: 16px;
  padding: 8px 0;
  position: relative;
  min-width: 0;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-line {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 32px;
  flex-shrink: 0;
}

.timeline-marker {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 1;
}

.line-connector {
  flex: 1;
  width: 3px;
  min-height: 20px;
  background: var(--surface-border);
  border-radius: 2px;
  margin-top: 4px;
}

.timeline-marker.procurement { background: linear-gradient(135deg, #3b82f6, #1d4ed8); }
.timeline-marker.inventory { background: linear-gradient(135deg, #f97316, #c2410c); }
.timeline-marker.sales { background: linear-gradient(135deg, #22c55e, #15803d); }
.timeline-marker.logistics { background: linear-gradient(135deg, #a855f7, #7c3aed); }
.timeline-marker.return { background: linear-gradient(135deg, #ef4444, #b91c1c); }

.timeline-card {
  flex: 1;
  min-width: 0;
  border-radius: 12px;
  padding: 14px 18px;
  border-left: 4px solid;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.timeline-card.procurement {
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.08), rgba(59, 130, 246, 0.02));
  border-left-color: #3b82f6;
}

.timeline-card.inventory {
  background: linear-gradient(90deg, rgba(249, 115, 22, 0.08), rgba(249, 115, 22, 0.02));
  border-left-color: #f97316;
}

.timeline-card.sales {
  background: linear-gradient(90deg, rgba(34, 197, 94, 0.08), rgba(34, 197, 94, 0.02));
  border-left-color: #22c55e;
}

.timeline-card.logistics {
  background: linear-gradient(90deg, rgba(168, 85, 247, 0.08), rgba(168, 85, 247, 0.02));
  border-left-color: #a855f7;
}

.timeline-card.return {
  background: linear-gradient(90deg, rgba(239, 68, 68, 0.08), rgba(239, 68, 68, 0.02));
  border-left-color: #ef4444;
}

.timeline-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.timeline-action {
  font-weight: 600;
  font-size: 14px;
}

.timeline-card.procurement .timeline-action { color: #1d4ed8; }
.timeline-card.inventory .timeline-action { color: #c2410c; }
.timeline-card.sales .timeline-action { color: #15803d; }
.timeline-card.logistics .timeline-action { color: #7c3aed; }
.timeline-card.return .timeline-action { color: #b91c1c; }

.timeline-time {
  font-size: 12px;
  color: var(--text-muted);
}

.timeline-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.timeline-desc {
  font-size: 13px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.timeline-card:hover {
  transform: translateX(4px);
  transition: transform 0.2s ease;
}

.timeline-meta {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-muted);
  display: flex;
  gap: 16px;
}

.timeline-meta i {
  margin-right: 4px;
}

.empty-timeline {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
}

.empty-timeline i {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

/* 详情表格 */
.detail-section {
  background: var(--surface-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.text-success { color: var(--green-500); }
.text-danger { color: var(--red-500); }

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-muted);
}

.empty-state i {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 16px;
}

/* 响应式 */
@media (max-width: 1200px) {
  .overview-grid {
    flex-wrap: wrap;
  }

  .overview-card {
    min-width: 130px;
  }
}

@media (max-width: 768px) {
  .overview-card {
    min-width: 120px;
  }

  .product-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .filter-bar {
    flex-direction: column;
    gap: 12px;
  }

  .filter-left {
    width: 100%;
  }
}
</style>