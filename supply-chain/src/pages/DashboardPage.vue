<template>
  <div class="dashboard-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('dashboard.title') }}</h1>
      <p class="page-subtitle">{{ t('dashboard.overview') }}</p>
    </header>

    <section class="stat-grid">
      <div class="stat-card">
        <div class="stat-icon blue">
          <i class="ri-money-cny-circle-line"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('dashboard.todaySales') }}</div>
          <div class="stat-value">¥{{ formatNumber(stats.todaySales) }}</div>
          <div class="stat-change up">
            <i class="ri-arrow-up-line"></i>
            {{ stats.salesGrowth }}%
          </div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <i class="ri-file-list-3-line"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('dashboard.todayOrders') }}</div>
          <div class="stat-value">{{ stats.todayOrders }}</div>
          <div class="stat-change up">
            <i class="ri-arrow-up-line"></i>
            {{ stats.ordersGrowth }}%
          </div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <i class="ri-alert-line"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('dashboard.inventoryAlert') }}</div>
          <div class="stat-value">{{ stats.inventoryAlert }}</div>
          <div class="stat-change down">
            <i class="ri-arrow-down-line"></i>
            需要补货
          </div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon purple">
          <i class="ri-shopping-cart-line"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('dashboard.pendingProcurement') }}</div>
          <div class="stat-value">{{ stats.pendingProcurement }}</div>
          <div class="stat-change">待处理订单</div>
        </div>
      </div>
    </section>

    <section class="grid-2">
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <span class="card-icon blue"><i class="ri-line-chart-line"></i></span>
            {{ t('dashboard.salesTrend') }}
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="salesChartRef"></canvas>
        </div>
      </div>
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <span class="card-icon green"><i class="ri-pie-chart-line"></i></span>
            {{ t('dashboard.inventoryDistribution') }}
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="inventoryChartRef"></canvas>
        </div>
      </div>
    </section>

    <section class="grid-2" style="margin-top: 20px;">
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <span class="card-icon orange"><i class="ri-fire-line"></i></span>
            {{ t('dashboard.topProducts') }}
          </div>
        </div>
        <DataTable :value="topProducts" stripedRows>
          <Column field="rank" header="#" style="width: 50px"></Column>
          <Column field="name" header="产品名称"></Column>
          <Column field="sales" header="销售额">
            <template #body="{ data }">
              ¥{{ formatNumber(data.sales) }}
            </template>
          </Column>
          <Column field="quantity" header="销量"></Column>
        </DataTable>
      </div>
      <div class="card">
        <div class="card-header">
          <div class="card-title">
            <span class="card-icon purple"><i class="ri-time-line"></i></span>
            {{ t('dashboard.recentOrders') }}
          </div>
        </div>
        <DataTable :value="recentOrders" stripedRows>
          <Column field="orderNo" header="订单号"></Column>
          <Column field="customer" header="客户"></Column>
          <Column field="amount" header="金额">
            <template #body="{ data }">
              ¥{{ formatNumber(data.amount) }}
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <span :class="['tag', getStatusClass(data.status)]">{{ data.status }}</span>
            </template>
          </Column>
        </DataTable>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Chart, registerables } from 'chart.js'
import { dashboardApi } from '@/api'

Chart.register(...registerables)

const { t } = useI18n()

const stats = ref({
  todaySales: 0,
  todayOrders: 0,
  inventoryAlert: 0,
  pendingProcurement: 0,
  salesGrowth: 0,
  ordersGrowth: 0
})
const salesChartRef = ref<HTMLCanvasElement | null>(null)
const inventoryChartRef = ref<HTMLCanvasElement | null>(null)
let salesChart: Chart | null = null
let inventoryChart: Chart | null = null

const topProducts = ref([
  { rank: 1, name: '工业电机', sales: 156800, quantity: 230 },
  { rank: 2, name: '高精度传感器', sales: 98500, quantity: 2180 },
  { rank: 3, name: '工业控制器', sales: 76200, quantity: 238 },
  { rank: 4, name: '精密轴承', sales: 45600, quantity: 304 },
  { rank: 5, name: '不锈钢板材', sales: 32800, quantity: 182 }
])

const recentOrders = ref([
  { orderNo: 'SO20240306001', customer: '苏州精密仪器厂', amount: 15000, status: '待确认' },
  { orderNo: 'SO20240305001', customer: '北京智能装备有限公司', amount: 45800, status: '已确认' },
  { orderNo: 'SO20240304001', customer: '广州机械制造', amount: 28900, status: '发货中' },
  { orderNo: 'SO20240303001', customer: '深圳科技股份', amount: 67200, status: '已完成' }
])

const formatNumber = (num: number) => num.toLocaleString('zh-CN')

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    '待确认': 'warning',
    '已确认': 'info',
    '发货中': 'info',
    '已完成': 'success'
  }
  return map[status] || 'info'
}

const fetchStats = async () => {
  try {
    const res: any = await dashboardApi.stats()
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

const initCharts = () => {
  if (salesChartRef.value) {
    salesChart = new Chart(salesChartRef.value, {
      type: 'line',
      data: {
        labels: ['1月', '2月', '3月', '4月', '5月', '6月'],
        datasets: [{
          label: '销售额',
          data: [65000, 78000, 90000, 85000, 120000, 156800],
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.4,
          fill: true
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false }
        },
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  if (inventoryChartRef.value) {
    inventoryChart = new Chart(inventoryChartRef.value, {
      type: 'doughnut',
      data: {
        labels: ['电子元器件', '机械配件', '原材料', '包装材料'],
        datasets: [{
          data: [45, 25, 20, 10],
          backgroundColor: ['#3b82f6', '#22c55e', '#f59e0b', '#a855f7']
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { position: 'bottom' }
        }
      }
    })
  }
}

onMounted(() => {
  fetchStats()
  initCharts()
})

onBeforeUnmount(() => {
  salesChart?.destroy()
  inventoryChart?.destroy()
})
</script>