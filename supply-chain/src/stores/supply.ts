import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Supplier, ProcurementOrder, InventoryItem, SalesOrder, LogisticsOrder, DashboardStats } from '@/types'

export const useSupplyStore = defineStore('supply', () => {
  const suppliers = ref<Supplier[]>([])
  const procurementOrders = ref<ProcurementOrder[]>([])
  const inventoryItems = ref<InventoryItem[]>([])
  const salesOrders = ref<SalesOrder[]>([])
  const logisticsOrders = ref<LogisticsOrder[]>([])
  const dashboardStats = ref<DashboardStats | null>(null)
  const loading = ref(false)

  const activeSuppliers = computed(() => suppliers.value.filter(s => s.status === 1))
  const pendingProcurement = computed(() => procurementOrders.value.filter(p => p.status === 'pending'))
  const lowStockItems = computed(() => inventoryItems.value.filter(i => i.status === 'low'))
  const pendingSales = computed(() => salesOrders.value.filter(s => s.status === 'pending'))

  function setSuppliers(data: Supplier[]) {
    suppliers.value = data
  }

  function addSupplier(supplier: Supplier) {
    suppliers.value.push(supplier)
  }

  function updateSupplier(id: number, data: Partial<Supplier>) {
    const index = suppliers.value.findIndex(s => s.id === id)
    if (index !== -1) {
      suppliers.value[index] = { ...suppliers.value[index], ...data }
    }
  }

  function deleteSupplier(id: number) {
    suppliers.value = suppliers.value.filter(s => s.id !== id)
  }

  function setProcurementOrders(data: ProcurementOrder[]) {
    procurementOrders.value = data
  }

  function addProcurementOrder(order: ProcurementOrder) {
    procurementOrders.value.push(order)
  }

  function updateProcurementOrder(id: number, data: Partial<ProcurementOrder>) {
    const index = procurementOrders.value.findIndex(o => o.id === id)
    if (index !== -1) {
      procurementOrders.value[index] = { ...procurementOrders.value[index], ...data }
    }
  }

  function setInventoryItems(data: InventoryItem[]) {
    inventoryItems.value = data
  }

  function setSalesOrders(data: SalesOrder[]) {
    salesOrders.value = data
  }

  function addSalesOrder(order: SalesOrder) {
    salesOrders.value.push(order)
  }

  function setLogisticsOrders(data: LogisticsOrder[]) {
    logisticsOrders.value = data
  }

  function setDashboardStats(data: DashboardStats) {
    dashboardStats.value = data
  }

  return {
    suppliers,
    procurementOrders,
    inventoryItems,
    salesOrders,
    logisticsOrders,
    dashboardStats,
    loading,
    activeSuppliers,
    pendingProcurement,
    lowStockItems,
    pendingSales,
    setSuppliers,
    addSupplier,
    updateSupplier,
    deleteSupplier,
    setProcurementOrders,
    addProcurementOrder,
    updateProcurementOrder,
    setInventoryItems,
    setSalesOrders,
    addSalesOrder,
    setLogisticsOrders,
    setDashboardStats
  }
})