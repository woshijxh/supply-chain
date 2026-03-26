<template>
  <div class="dashboard-page">
    <header class="page-header">
      <h1 class="page-title">仪表盘</h1>
      <p class="page-subtitle">供应链管理系统数据概览</p>
    </header>

    <!-- 核心指标 -->
    <section class="metric-grid">
      <div class="metric-card sales">
        <div class="metric-icon"><i class="ri-money-cny-circle-line"></i></div>
        <div class="metric-body">
          <div class="metric-label">今日销售额</div>
          <div class="metric-value">¥{{ formatNumber(stats.todaySales) }}</div>
          <div class="metric-trend" :class="stats.salesGrowth >= 0 ? 'up' : 'down'">
            <i :class="stats.salesGrowth >= 0 ? 'ri-arrow-up-line' : 'ri-arrow-down-line'"></i>
            {{ Math.abs(stats.salesGrowth) }}%
          </div>
        </div>
      </div>
      <div class="metric-card orders">
        <div class="metric-icon"><i class="ri-file-list-3-line"></i></div>
        <div class="metric-body">
          <div class="metric-label">今日订单</div>
          <div class="metric-value">{{ stats.todayOrders }}</div>
          <div class="metric-trend" :class="stats.ordersGrowth >= 0 ? 'up' : 'down'">
            <i :class="stats.ordersGrowth >= 0 ? 'ri-arrow-up-line' : 'ri-arrow-down-line'"></i>
            {{ Math.abs(stats.ordersGrowth) }}%
          </div>
        </div>
      </div>
      <div class="metric-card inventory">
        <div class="metric-icon"><i class="ri-store-2-line"></i></div>
        <div class="metric-body">
          <div class="metric-label">库存总量</div>
          <div class="metric-value">{{ stats.inventoryStats?.total || 0 }}</div>
          <div class="metric-trend warning" v-if="stats.inventoryAlert > 0">
            <i class="ri-alert-line"></i>
            {{ stats.inventoryAlert }} 项预警
          </div>
          <div class="metric-trend ok" v-else>
            <i class="ri-checkbox-circle-line"></i> 库存正常
          </div>
        </div>
      </div>
      <div class="metric-card procurement">
        <div class="metric-icon"><i class="ri-shopping-cart-line"></i></div>
        <div class="metric-body">
          <div class="metric-label">待处理采购</div>
          <div class="metric-value">{{ stats.pendingProcurement }}</div>
          <div class="metric-trend neutral">采购订单</div>
        </div>
      </div>
    </section>

    <!-- 业务概览 -->
    <section class="overview-section">
      <div class="overview-card">
        <div class="overview-header">
          <h3><i class="ri-building-line"></i> 业务伙伴</h3>
        </div>
        <div class="overview-stats">
          <div class="stat-row" @click="navigateTo('/suppliers')">
            <span class="stat-name">供应商</span>
            <span class="stat-num blue">{{ stats.activeSuppliers || 0 }}</span>
          </div>
          <div class="stat-row" @click="navigateTo('/customers')">
            <span class="stat-name">客户</span>
            <span class="stat-num green">{{ stats.customerCount || 0 }}</span>
          </div>
          <div class="stat-row" @click="navigateTo('/products')">
            <span class="stat-name">产品</span>
            <span class="stat-num purple">{{ stats.productCount || 0 }}</span>
          </div>
        </div>
      </div>

      <div class="overview-card">
        <div class="overview-header">
          <h3><i class="ri-exchange-line"></i> 库存状态</h3>
        </div>
        <div class="overview-stats">
          <div class="stat-row">
            <span class="stat-name">正常库存</span>
            <span class="stat-num green">{{ stats.inventoryStats?.normal || 0 }}</span>
          </div>
          <div class="stat-row">
            <span class="stat-name">低库存预警</span>
            <span class="stat-num orange">{{ stats.inventoryStats?.low || 0 }}</span>
          </div>
          <div class="stat-row">
            <span class="stat-name">总SKU数</span>
            <span class="stat-num">{{ stats.inventoryStats?.totalSku || 0 }}</span>
          </div>
        </div>
      </div>

      <div class="overview-card">
        <div class="overview-header">
          <h3><i class="ri-truck-line"></i> 订单状态</h3>
        </div>
        <div class="overview-stats">
          <div class="stat-row">
            <span class="stat-name">待确认销售</span>
            <span class="stat-num orange">{{ stats.pendingSales || 0 }}</span>
          </div>
          <div class="stat-row">
            <span class="stat-name">配送中物流</span>
            <span class="stat-num blue">{{ stats.shippingLogistics || 0 }}</span>
          </div>
          <div class="stat-row">
            <span class="stat-name">待处理退货</span>
            <span class="stat-num red">{{ stats.pendingReturns || 0 }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 图表区域 -->
    <section class="chart-section">
      <div class="chart-card">
        <div class="card-header">
          <h3><i class="ri-line-chart-line"></i> 销售趋势</h3>
          <div class="chart-tabs">
            <span :class="{ active: salesPeriod === 'week' }" @click="changeSalesPeriod('week')">近7天</span>
            <span :class="{ active: salesPeriod === 'month' }" @click="changeSalesPeriod('month')">近30天</span>
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="salesChartRef"></canvas>
        </div>
      </div>
      <div class="chart-card">
        <div class="card-header">
          <h3><i class="ri-pie-chart-line"></i> 库存分布</h3>
        </div>
        <div class="chart-container">
          <canvas ref="inventoryChartRef"></canvas>
        </div>
      </div>
    </section>

    <!-- 列表区域 -->
    <section class="list-section">
      <div class="list-card">
        <div class="card-header">
          <h3><i class="ri-fire-line"></i> 热销产品</h3>
          <a class="view-all" @click="navigateTo('/products')">查看全部</a>
        </div>
        <div class="product-list">
          <div class="product-item" v-for="(item, index) in topProducts" :key="index">
            <span class="rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</span>
            <div class="product-info">
              <span class="product-name">{{ item.name }}</span>
              <span class="product-sales">销量 {{ item.quantity }}</span>
            </div>
            <span class="product-amount">¥{{ formatNumber(item.sales) }}</span>
          </div>
        </div>
      </div>
      <div class="list-card">
        <div class="card-header">
          <h3><i class="ri-time-line"></i> 最近订单</h3>
          <a class="view-all" @click="navigateTo('/sales')">查看全部</a>
        </div>
        <div class="order-list">
          <div class="order-item" v-for="(order, index) in recentOrders" :key="index">
            <div class="order-main">
              <span class="order-no">{{ order.orderNo }}</span>
              <span class="order-customer">{{ order.customer }}</span>
            </div>
            <div class="order-meta">
              <span class="order-amount">¥{{ formatNumber(order.amount) }}</span>
              <Tag :value="order.status" :severity="getOrderSeverity(order.status)" size="small" />
            </div>
          </div>
        </div>
      </div>
      <div class="list-card">
        <div class="card-header">
          <h3><i class="ri-alert-line"></i> 库存预警</h3>
          <a class="view-all" @click="navigateTo('/inventory')">查看全部</a>
        </div>
        <div class="alert-list">
          <div class="alert-item" v-for="(item, index) in lowStockItems" :key="index">
            <div class="alert-info">
              <span class="alert-name">{{ item.name }}</span>
              <span class="alert-sku">{{ item.sku }}</span>
            </div>
            <div class="alert-stock">
              <span class="current">{{ item.quantity }}</span>
              <span class="divider">/</span>
              <span class="threshold">{{ item.minStock }}</span>
            </div>
          </div>
          <div class="empty-alert" v-if="lowStockItems.length === 0">
            <i class="ri-checkbox-circle-line"></i>
            <span>库存充足，暂无预警</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Chart, registerables } from 'chart.js'
import { dashboardApi } from '@/api'

Chart.register(...registerables)

const router = useRouter()

const stats = ref<any>({
  todaySales: 0,
  todayOrders: 0,
  inventoryAlert: 0,
  pendingProcurement: 0,
  salesGrowth: 0,
  ordersGrowth: 0,
  inventoryStats: { total: 0, normal: 0, low: 0, totalSku: 0 },
  activeSuppliers: 0,
  customerCount: 0,
  productCount: 0,
  pendingSales: 0,
  shippingLogistics: 0,
  pendingReturns: 0
})

const salesPeriod = ref('week')
const salesChartRef = ref<HTMLCanvasElement | null>(null)
const inventoryChartRef = ref<HTMLCanvasElement | null>(null)
let salesChart: Chart | null = null
let inventoryChart: Chart | null = null

const topProducts = ref<any[]>([])
const recentOrders = ref<any[]>([])
const lowStockItems = ref<any[]>([])

const formatNumber = (num: number) => num?.toLocaleString('zh-CN') || '0'

const getOrderSeverity = (status: string) => {
  const map: Record<string, any> = {
    '待确认': 'warning',
    '已确认': 'info',
    '配送中': 'info',
    '已完成': 'success',
    '已取消': 'danger'
  }
  return map[status] || 'secondary'
}

const navigateTo = (path: string) => {
  router.push(path)
}

const fetchStats = async () => {
  try {
    const res: any = await dashboardApi.stats()
    if (res.code === 0) {
      stats.value = { ...stats.value, ...res.data }
    }
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

const fetchTopProducts = async () => {
  try {
    const res: any = await dashboardApi.topProducts()
    if (res.code === 0) {
      topProducts.value = res.data || []
    }
  } catch (error) {
    console.error('获取热销产品失败', error)
  }
}

const fetchRecentOrders = async () => {
  try {
    const res: any = await dashboardApi.recentOrders()
    if (res.code === 0) {
      recentOrders.value = (res.data || []).map((order: any) => ({
        ...order,
        status: getStatusText(order.status),
        customer: order.customer_name || order.customer || '-',
        amount: order.total_amount || order.amount || 0
      }))
    }
  } catch (error) {
    console.error('获取最近订单失败', error)
  }
}

const fetchLowStockItems = async () => {
  try {
    const res: any = await dashboardApi.lowStockItems()
    if (res.code === 0) {
      lowStockItems.value = (res.data || []).map((item: any) => ({
        name: item.product?.name || item.productName || '未知产品',
        sku: item.product?.sku || '',
        quantity: item.quantity,
        minStock: item.product?.minStock || 10
      }))
    }
  } catch (error) {
    console.error('获取库存预警失败', error)
  }
}

const fetchSalesTrend = async () => {
  try {
    const res: any = await dashboardApi.salesTrend(salesPeriod.value === 'week' ? 7 : 30)
    if (res.code === 0 && salesChart) {
      const data = res.data || []
      const labels = data.map((item: any) => {
        const date = new Date(item.date)
        return salesPeriod.value === 'week'
          ? ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][date.getDay()]
          : `${date.getDate()}日`
      })
      const values = data.map((item: any) => item.amount || 0)
      salesChart.data.labels = labels
      salesChart.data.datasets[0].data = values
      salesChart.update()
    }
  } catch (error) {
    console.error('获取销售趋势失败', error)
  }
}

const fetchInventoryDistribution = async () => {
  try {
    const res: any = await dashboardApi.inventoryDistribution()
    if (res.code === 0 && inventoryChart) {
      const data = res.data || []
      const labels = data.map((item: any) => item.category || '未分类')
      const values = data.map((item: any) => item.total || 0)
      inventoryChart.data.labels = labels
      inventoryChart.data.datasets[0].data = values
      inventoryChart.update()
    }
  } catch (error) {
    console.error('获取库存分布失败', error)
  }
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    'pending': '待确认',
    'confirmed': '已确认',
    'shipping': '配送中',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return map[status] || status
}

const changeSalesPeriod = (period: string) => {
  salesPeriod.value = period
  fetchSalesTrend()
}

const initCharts = () => {
  if (salesChartRef.value) {
    salesChart = new Chart(salesChartRef.value, {
      type: 'line',
      data: {
        labels: [],
        datasets: [{
          label: '销售额',
          data: [],
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.4,
          fill: true,
          pointRadius: 4,
          pointHoverRadius: 6
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: (value) => '¥' + (Number(value) / 10000).toFixed(1) + '万'
            }
          }
        }
      }
    })
  }

  if (inventoryChartRef.value) {
    inventoryChart = new Chart(inventoryChartRef.value, {
      type: 'doughnut',
      data: {
        labels: [],
        datasets: [{
          data: [],
          backgroundColor: ['#3b82f6', '#22c55e', '#f59e0b', '#a855f7', '#6b7280'],
          borderWidth: 0
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        cutout: '65%',
        plugins: {
          legend: {
            position: 'right',
            labels: {
              boxWidth: 12,
              padding: 16
            }
          }
        }
      }
    })
  }
}

onMounted(() => {
  fetchStats()
  fetchTopProducts()
  fetchRecentOrders()
  fetchLowStockItems()
  initCharts()
  fetchSalesTrend()
  fetchInventoryDistribution()
})

onBeforeUnmount(() => {
  salesChart?.destroy()
  inventoryChart?.destroy()
})
</script>

<style scoped>
.dashboard-page {
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

/* 核心指标 */
.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card {
  background: var(--surface-card);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-left: 4px solid;
  transition: transform 0.2s, box-shadow 0.2s;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.metric-card.sales { border-left-color: #3b82f6; }
.metric-card.orders { border-left-color: #22c55e; }
.metric-card.inventory { border-left-color: #f59e0b; }
.metric-card.procurement { border-left-color: #a855f7; }

.metric-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}

.metric-card.sales .metric-icon { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.metric-card.orders .metric-icon { background: rgba(34, 197, 94, 0.1); color: #22c55e; }
.metric-card.inventory .metric-icon { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
.metric-card.procurement .metric-icon { background: rgba(168, 85, 247, 0.1); color: #a855f7; }

.metric-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 26px;
  font-weight: 600;
}

.metric-trend {
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
}

.metric-trend.up { color: #22c55e; }
.metric-trend.down { color: #ef4444; }
.metric-trend.warning { color: #f59e0b; }
.metric-trend.ok { color: #22c55e; }
.metric-trend.neutral { color: var(--text-muted); }

/* 业务概览 */
.overview-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.overview-card {
  background: var(--surface-card);
  border-radius: 12px;
  padding: 16px 20px;
}

.overview-header h3 {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
}

.overview-stats {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--surface-border);
  cursor: pointer;
  transition: background 0.2s;
}

.stat-row:last-child {
  border-bottom: none;
}

.stat-row:hover {
  background: var(--surface-ground);
  margin: 0 -12px;
  padding: 8px 12px;
  border-radius: 6px;
}

.stat-name {
  font-size: 14px;
  color: var(--text-secondary);
}

.stat-num {
  font-size: 16px;
  font-weight: 600;
}

.stat-num.blue { color: #3b82f6; }
.stat-num.green { color: #22c55e; }
.stat-num.purple { color: #a855f7; }
.stat-num.orange { color: #f59e0b; }
.stat-num.red { color: #ef4444; }

/* 图表区域 */
.chart-section {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.chart-card {
  background: var(--surface-card);
  border-radius: 12px;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-header h3 {
  font-size: 15px;
  font-weight: 600;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-tabs {
  display: flex;
  gap: 8px;
}

.chart-tabs span {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  color: var(--text-muted);
  transition: all 0.2s;
}

.chart-tabs span:hover {
  background: var(--surface-ground);
}

.chart-tabs span.active {
  background: var(--primary);
  color: white;
}

.chart-container {
  height: 250px;
}

/* 列表区域 */
.list-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.list-card {
  background: var(--surface-card);
  border-radius: 12px;
  padding: 20px;
}

.view-all {
  font-size: 13px;
  color: var(--primary);
  cursor: pointer;
}

.view-all:hover {
  text-decoration: underline;
}

/* 产品列表 */
.product-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--surface-border);
}

.product-item:last-child {
  border-bottom: none;
}

.rank {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  background: var(--surface-ground);
  color: var(--text-muted);
}

.rank.rank-1 { background: #fef3c7; color: #d97706; }
.rank.rank-2 { background: #e5e7eb; color: #6b7280; }
.rank.rank-3 { background: #fed7aa; color: #c2410c; }

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  display: block;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-sales {
  font-size: 12px;
  color: var(--text-muted);
}

.product-amount {
  font-size: 14px;
  font-weight: 600;
  color: var(--primary);
}

/* 订单列表 */
.order-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--surface-border);
}

.order-item:last-child {
  border-bottom: none;
}

.order-no {
  display: block;
  font-size: 14px;
  font-weight: 500;
}

.order-customer {
  font-size: 12px;
  color: var(--text-muted);
}

.order-amount {
  font-size: 14px;
  font-weight: 500;
  margin-right: 8px;
}

/* 预警列表 */
.alert-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: #fef2f2;
  border-radius: 8px;
  border-left: 3px solid #ef4444;
}

.alert-name {
  display: block;
  font-size: 14px;
  font-weight: 500;
}

.alert-sku {
  font-size: 12px;
  color: var(--text-muted);
}

.alert-stock {
  font-size: 14px;
}

.alert-stock .current {
  color: #ef4444;
  font-weight: 600;
}

.alert-stock .divider {
  color: var(--text-muted);
  margin: 0 4px;
}

.alert-stock .threshold {
  color: var(--text-muted);
}

.empty-alert {
  text-align: center;
  padding: 24px;
  color: var(--text-muted);
}

.empty-alert i {
  display: block;
  font-size: 32px;
  color: #22c55e;
  margin-bottom: 8px;
}

/* 响应式 */
@media (max-width: 1200px) {
  .metric-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .chart-section {
    grid-template-columns: 1fr;
  }

  .list-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .metric-grid {
    grid-template-columns: 1fr;
  }

  .overview-section {
    grid-template-columns: 1fr;
  }

  .list-section {
    grid-template-columns: 1fr;
  }
}
</style>