<template>
  <div class="return-page">
    <header class="page-header">
      <h1 class="page-title">退货管理</h1>
      <p class="page-subtitle">管理销售退货和采购退货</p>
    </header>

    <TabView v-model:activeIndex="activeTab">
      <TabPanel header="销售退货" :value="0">
        <div class="toolbar">
          <Select v-model="salesStatusFilter" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="状态筛选" style="width: 150px" />
        </div>
        <DataTable :value="salesReturns" stripedRows :loading="salesLoading" :paginator="true" :rows="10" scrollable scrollHeight="400px">
          <Column field="returnNo" header="退货单号" style="width: 150px"></Column>
          <Column field="salesOrderNo" header="销售单号" style="width: 150px"></Column>
          <Column field="customerName" header="客户"></Column>
          <Column field="totalAmount" header="退货金额">
            <template #body="{ data }">¥{{ formatNumber(data.totalAmount) }}</template>
          </Column>
          <Column field="refundAmount" header="退款金额">
            <template #body="{ data }">¥{{ formatNumber(data.refundAmount || 0) }}</template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <span :class="['tag', getStatusClass(data.status)]">{{ getStatusText(data.status) }}</span>
            </template>
          </Column>
          <Column field="refundStatus" header="退款状态">
            <template #body="{ data }">
              <span :class="['tag', getRefundClass(data.refundStatus)]">{{ getRefundText(data.refundStatus) }}</span>
            </template>
          </Column>
          <Column field="createdAt" header="创建时间">
            <template #body="{ data }">{{ formatDate(data.createdAt) }}</template>
          </Column>
          <Column header="操作" style="width: 180px">
            <template #body="{ data }">
              <Button text severity="info" icon="ri-eye-line" @click="viewSalesReturn(data)" />
              <Button v-if="data.status === 'pending'" text severity="success" icon="ri-check-line" @click="approveSalesReturn(data)" />
              <Button v-if="data.status === 'pending'" text severity="warning" icon="ri-close-line" @click="rejectSalesReturn(data)" />
              <Button v-if="data.status === 'approved'" text severity="primary" icon="ri-money-cny-circle-line" @click="completeSalesReturn(data)" />
            </template>
          </Column>
        </DataTable>
      </TabPanel>

      <TabPanel header="采购退货" :value="1">
        <div class="toolbar">
          <Select v-model="procurementStatusFilter" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="状态筛选" style="width: 150px" />
        </div>
        <DataTable :value="procurementReturns" stripedRows :loading="procurementLoading" :paginator="true" :rows="10" scrollable scrollHeight="400px">
          <Column field="returnNo" header="退货单号" style="width: 150px"></Column>
          <Column field="procurementOrderNo" header="采购单号" style="width: 150px"></Column>
          <Column field="supplierName" header="供应商"></Column>
          <Column field="totalAmount" header="退货金额">
            <template #body="{ data }">¥{{ formatNumber(data.totalAmount) }}</template>
          </Column>
          <Column field="refundAmount" header="退款金额">
            <template #body="{ data }">¥{{ formatNumber(data.refundAmount || 0) }}</template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <span :class="['tag', getStatusClass(data.status)]">{{ getStatusText(data.status) }}</span>
            </template>
          </Column>
          <Column field="refundStatus" header="退款状态">
            <template #body="{ data }">
              <span :class="['tag', getRefundClass(data.refundStatus)]">{{ getRefundText(data.refundStatus) }}</span>
            </template>
          </Column>
          <Column field="createdAt" header="创建时间">
            <template #body="{ data }">{{ formatDate(data.createdAt) }}</template>
          </Column>
          <Column header="操作" style="width: 180px">
            <template #body="{ data }">
              <Button text severity="info" icon="ri-eye-line" @click="viewProcurementReturn(data)" />
              <Button v-if="data.status === 'pending'" text severity="success" icon="ri-check-line" @click="approveProcurementReturn(data)" />
              <Button v-if="data.status === 'pending'" text severity="warning" icon="ri-close-line" @click="rejectProcurementReturn(data)" />
              <Button v-if="data.status === 'approved'" text severity="primary" icon="ri-money-cny-circle-line" @click="completeProcurementReturn(data)" />
            </template>
          </Column>
        </DataTable>
      </TabPanel>
    </TabView>

    <!-- 销售退货详情弹窗 -->
    <Dialog v-model:visible="salesDetailVisible" header="销售退货详情" :style="{ width: '600px' }">
      <div v-if="currentSalesReturn" class="return-detail">
        <div class="detail-section">
          <h4>退货信息</h4>
          <div class="detail-row">
            <span class="detail-label">退货单号:</span>
            <span>{{ currentSalesReturn.returnNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">销售单号:</span>
            <span>{{ currentSalesReturn.salesOrderNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">客户:</span>
            <span>{{ currentSalesReturn.customerName }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">状态:</span>
            <span :class="['tag', getStatusClass(currentSalesReturn.status)]">{{ getStatusText(currentSalesReturn.status) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">退货原因:</span>
            <span>{{ currentSalesReturn.reason }}</span>
          </div>
        </div>
        <div class="detail-section">
          <h4>退货商品</h4>
          <DataTable :value="currentSalesReturn.items" stripedRows>
            <Column field="productName" header="产品"></Column>
            <Column field="quantity" header="数量"></Column>
            <Column field="unitPrice" header="单价">
              <template #body="{ data }">¥{{ data.unitPrice }}</template>
            </Column>
            <Column field="amount" header="金额">
              <template #body="{ data }">¥{{ formatNumber(data.amount) }}</template>
            </Column>
          </DataTable>
        </div>
        <div class="detail-section">
          <div class="summary-row">
            <span>退货总额:</span>
            <strong>¥{{ formatNumber(currentSalesReturn.totalAmount) }}</strong>
          </div>
        </div>
      </div>
    </Dialog>

    <!-- 采购退货详情弹窗 -->
    <Dialog v-model:visible="procurementDetailVisible" header="采购退货详情" :style="{ width: '600px' }">
      <div v-if="currentProcurementReturn" class="return-detail">
        <div class="detail-section">
          <h4>退货信息</h4>
          <div class="detail-row">
            <span class="detail-label">退货单号:</span>
            <span>{{ currentProcurementReturn.returnNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">采购单号:</span>
            <span>{{ currentProcurementReturn.procurementOrderNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">供应商:</span>
            <span>{{ currentProcurementReturn.supplierName }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">状态:</span>
            <span :class="['tag', getStatusClass(currentProcurementReturn.status)]">{{ getStatusText(currentProcurementReturn.status) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">退货原因:</span>
            <span>{{ currentProcurementReturn.reason }}</span>
          </div>
        </div>
        <div class="detail-section">
          <h4>退货商品</h4>
          <DataTable :value="currentProcurementReturn.items" stripedRows>
            <Column field="productName" header="产品"></Column>
            <Column field="quantity" header="数量"></Column>
            <Column field="unitPrice" header="单价">
              <template #body="{ data }">¥{{ data.unitPrice }}</template>
            </Column>
            <Column field="amount" header="金额">
              <template #body="{ data }">¥{{ formatNumber(data.amount) }}</template>
            </Column>
          </DataTable>
        </div>
        <div class="detail-section">
          <div class="summary-row">
            <span>退货总额:</span>
            <strong>¥{{ formatNumber(currentProcurementReturn.totalAmount) }}</strong>
          </div>
        </div>
      </div>
    </Dialog>

    <!-- 退款金额输入弹窗 -->
    <Dialog v-model:visible="refundDialogVisible" header="确认退款" :style="{ width: '400px' }">
      <div class="form-group">
        <label class="form-label">退款金额</label>
        <InputNumber v-model="refundAmount" mode="currency" currency="CNY" style="width: 100%" />
      </div>
      <div class="dialog-footer">
        <Button label="取消" severity="secondary" @click="refundDialogVisible = false" />
        <Button label="确认" @click="confirmRefund" :loading="refundLoading" />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { salesReturnApi, procurementReturnApi } from '@/api'

const toast = useToast()
const activeTab = ref(0)

// 销售退货
const salesReturns = ref<any[]>([])
const salesLoading = ref(false)
const salesStatusFilter = ref('all')
const salesDetailVisible = ref(false)
const currentSalesReturn = ref<any>(null)

// 采购退货
const procurementReturns = ref<any[]>([])
const procurementLoading = ref(false)
const procurementStatusFilter = ref('all')
const procurementDetailVisible = ref(false)
const currentProcurementReturn = ref<any>(null)

// 退款相关
const refundDialogVisible = ref(false)
const refundAmount = ref(0)
const refundLoading = ref(false)
const currentRefundType = ref('')
const currentRefundId = ref(0)

const statusOptions = [
  { label: '全部', value: 'all' },
  { label: '待处理', value: 'pending' },
  { label: '已批准', value: 'approved' },
  { label: '已拒绝', value: 'rejected' },
  { label: '已完成', value: 'completed' }
]

const formatNumber = (num: number) => num?.toLocaleString('zh-CN') || '0'

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    approved: '已批准',
    rejected: '已拒绝',
    completed: '已完成'
  }
  return map[status] || status
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'info',
    rejected: 'danger',
    completed: 'success'
  }
  return map[status] || 'info'
}

const getRefundText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待退款',
    refunded: '已退款'
  }
  return map[status] || status
}

const getRefundClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    refunded: 'success'
  }
  return map[status] || 'info'
}

// 加载销售退货列表
const fetchSalesReturns = async () => {
  salesLoading.value = true
  try {
    const res: any = await salesReturnApi.list(1, 100, salesStatusFilter.value === 'all' ? '' : salesStatusFilter.value)
    if (res.code === 0) {
      salesReturns.value = res.data?.list || []
    }
  } catch (error) {
    console.error('获取销售退货列表失败', error)
  } finally {
    salesLoading.value = false
  }
}

// 加载采购退货列表
const fetchProcurementReturns = async () => {
  procurementLoading.value = true
  try {
    const res: any = await procurementReturnApi.list(1, 100, procurementStatusFilter.value === 'all' ? '' : procurementStatusFilter.value)
    if (res.code === 0) {
      procurementReturns.value = res.data?.list || []
    }
  } catch (error) {
    console.error('获取采购退货列表失败', error)
  } finally {
    procurementLoading.value = false
  }
}

// 查看销售退货详情
const viewSalesReturn = async (item: any) => {
  try {
    const res: any = await salesReturnApi.get(String(item.id))
    if (res.code === 0) {
      currentSalesReturn.value = res.data
      salesDetailVisible.value = true
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '获取详情失败', life: 3000 })
  }
}

// 查看采购退货详情
const viewProcurementReturn = async (item: any) => {
  try {
    const res: any = await procurementReturnApi.get(String(item.id))
    if (res.code === 0) {
      currentProcurementReturn.value = res.data
      procurementDetailVisible.value = true
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '获取详情失败', life: 3000 })
  }
}

// 批准销售退货
const approveSalesReturn = async (item: any) => {
  try {
    const res: any = await salesReturnApi.approve(String(item.id))
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '已批准', life: 3000 })
      fetchSalesReturns()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '操作失败', life: 3000 })
  }
}

// 拒绝销售退货
const rejectSalesReturn = async (item: any) => {
  const reason = prompt('请输入拒绝原因')
  if (!reason) return
  try {
    const res: any = await salesReturnApi.reject(String(item.id), reason)
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '已拒绝', life: 3000 })
      fetchSalesReturns()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '操作失败', life: 3000 })
  }
}

// 完成销售退货（退款）
const completeSalesReturn = (item: any) => {
  currentRefundType.value = 'sales'
  currentRefundId.value = item.id
  refundAmount.value = item.totalAmount
  refundDialogVisible.value = true
}

// 批准采购退货
const approveProcurementReturn = async (item: any) => {
  try {
    const res: any = await procurementReturnApi.approve(String(item.id))
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '已批准', life: 3000 })
      fetchProcurementReturns()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '操作失败', life: 3000 })
  }
}

// 拒绝采购退货
const rejectProcurementReturn = async (item: any) => {
  const reason = prompt('请输入拒绝原因')
  if (!reason) return
  try {
    const res: any = await procurementReturnApi.reject(String(item.id), reason)
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '已拒绝', life: 3000 })
      fetchProcurementReturns()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '操作失败', life: 3000 })
  }
}

// 完成采购退货（退款）
const completeProcurementReturn = (item: any) => {
  currentRefundType.value = 'procurement'
  currentRefundId.value = item.id
  refundAmount.value = item.totalAmount
  refundDialogVisible.value = true
}

// 确认退款
const confirmRefund = async () => {
  refundLoading.value = true
  try {
    let res: any
    if (currentRefundType.value === 'sales') {
      res = await salesReturnApi.complete(String(currentRefundId.value), refundAmount.value)
    } else {
      res = await procurementReturnApi.complete(String(currentRefundId.value), refundAmount.value)
    }
    if (res.code === 0) {
      toast.add({ severity: 'success', summary: '退款完成', life: 3000 })
      refundDialogVisible.value = false
      fetchSalesReturns()
      fetchProcurementReturns()
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: '操作失败', life: 3000 })
  } finally {
    refundLoading.value = false
  }
}

// 监听筛选变化
watch(salesStatusFilter, fetchSalesReturns)
watch(procurementStatusFilter, fetchProcurementReturns)

onMounted(() => {
  fetchSalesReturns()
  fetchProcurementReturns()
})
</script>

<style scoped>
.return-detail {
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
  .summary-row {
    display: flex;
    justify-content: space-between;
    padding: 12px;
    background: var(--bg);
    border-radius: var(--radius-md);
    strong {
      color: var(--primary);
      font-size: 18px;
    }
  }
}

.tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}
.tag.warning { background: #fef3cd; color: #856404; }
.tag.info { background: #cce5ff; color: #004085; }
.tag.success { background: #d4edda; color: #155724; }
.tag.danger { background: #f8d7da; color: #721c24; }
</style>