<template>
  <div class="product-page">
    <header class="page-header">
      <h1 class="page-title">产品管理</h1>
      <p class="page-subtitle">管理产品信息和价格</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasCreatePermission" :label="t('common.add')" icon="ri-add-line" @click="openDialog()" />
      <div class="search-box">
        <InputText v-model="searchKey" :placeholder="t('common.search')" @keyup.enter="onSearch" />
        <Button icon="ri-search-line" severity="secondary" @click="onSearch" />
      </div>
    </div>

    <div class="card">
      <DataTable
        :value="products"
        stripedRows
        :loading="loading"
        :paginator="true"
        :rows="10"
        scrollable
        scrollHeight="600px"
      >
        <Column field="code" header="产品编码" style="width: 120px"></Column>
        <Column field="name" header="产品名称" style="width: 200px"></Column>
        <Column field="sku" header="SKU" style="width: 120px"></Column>
        <Column field="category" header="分类"></Column>
        <Column field="unit" header="单位" style="width: 80px"></Column>
        <Column field="costPrice" header="成本价">
          <template #body="{ data }">
            ¥{{ formatNumber(data.costPrice) }}
          </template>
        </Column>
        <Column field="salePrice" header="销售价">
          <template #body="{ data }">
            ¥{{ formatNumber(data.salePrice) }}
          </template>
        </Column>
        <Column field="minStock" header="最低库存"></Column>
        <Column field="maxStock" header="最高库存"></Column>
        <Column field="status" header="状态">
          <template #body="{ data }">
            <span :class="['tag', data.status === 1 ? 'success' : 'warning']">
              {{ data.status === 1 ? '启用' : '禁用' }}
            </span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width: 120px">
          <template #body="{ data }">
            <Button v-if="hasUpdatePermission" text severity="info" icon="ri-edit-line" @click="openDialog(data)" />
            <Button v-if="hasDeletePermission" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑产品' : '新增产品'" :style="{ width: '600px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">产品名称 *</label>
          <InputText v-model="form.name" />
        </div>
        <div class="form-group">
          <label class="form-label">SKU</label>
          <InputText v-model="form.sku" />
        </div>
        <div class="form-group">
          <label class="form-label">分类</label>
          <InputText v-model="form.category" />
        </div>
        <div class="form-group">
          <label class="form-label">单位</label>
          <InputText v-model="form.unit" placeholder="个/件/箱" />
        </div>
        <div class="form-group">
          <label class="form-label">成本价</label>
          <InputNumber v-model="form.costPrice" mode="currency" currency="CNY" :minFractionDigits="2" />
        </div>
        <div class="form-group">
          <label class="form-label">销售价</label>
          <InputNumber v-model="form.salePrice" mode="currency" currency="CNY" :minFractionDigits="2" />
        </div>
        <div class="form-group">
          <label class="form-label">最低库存</label>
          <InputNumber v-model="form.minStock" :min="0" />
        </div>
        <div class="form-group">
          <label class="form-label">最高库存</label>
          <InputNumber v-model="form.maxStock" :min="0" />
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group" v-if="isEdit">
          <label class="form-label">产品编码</label>
          <InputText v-model="form.code" disabled />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">描述</label>
          <Textarea v-model="form.description" rows="3" />
        </div>
      </div>
      <div class="dialog-footer">
        <Button :label="t('common.cancel')" severity="secondary" @click="dialogVisible = false" />
        <Button :label="t('common.save')" @click="handleSave" :loading="saving" />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { productApi } from '@/api'
import type { Product } from '@/types'
import { usePermission } from '@/utils/usePermission'
import { PRODUCT_PERMISSIONS } from '@/config/permissions'
import { useUserStore } from '@/stores/user'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const { initUserPermissions } = usePermission()
const userStore = useUserStore()

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const products = ref<Product[]>([])
const searchKey = ref('')
const currentId = ref<number | null>(null)

const form = ref({
  name: '',
  sku: '',
  category: '',
  unit: '个',
  costPrice: 0,
  salePrice: 0,
  minStock: 0,
  maxStock: 0,
  status: 1,
  description: '',
  code: ''
})

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const formatNumber = (num: number) => {
  if (num === null || num === undefined) return '0.00'
  return num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const hasCreatePermission = computed(() => userStore.hasPermission(PRODUCT_PERMISSIONS.CREATE))
const hasUpdatePermission = computed(() => userStore.hasPermission(PRODUCT_PERMISSIONS.UPDATE))
const hasDeletePermission = computed(() => userStore.hasPermission(PRODUCT_PERMISSIONS.DELETE))

const fetchProducts = async () => {
  loading.value = true
  try {
    const res: any = await productApi.list(1, 100, searchKey.value)
    if (res.code === 0) {
      products.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取产品列表失败', error)
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  fetchProducts()
}

const openDialog = (product?: Product) => {
  if (product) {
    isEdit.value = true
    currentId.value = product.id
    form.value = {
      name: product.name,
      sku: product.sku || '',
      category: product.category || '',
      unit: product.unit || '个',
      costPrice: product.costPrice || 0,
      salePrice: product.salePrice || 0,
      minStock: product.minStock || 0,
      maxStock: product.maxStock || 0,
      status: product.status ?? 1,
      description: product.description || '',
      code: product.code || ''
    }
  } else {
    isEdit.value = false
    currentId.value = null
    form.value = {
      name: '',
      sku: '',
      category: '',
      unit: '个',
      costPrice: 0,
      salePrice: 0,
      minStock: 0,
      maxStock: 0,
      status: 1,
      description: '',
      code: ''
    }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.value.name) {
    toast.add({ severity: 'warn', summary: '请填写产品名称', life: 5000 })
    return
  }

  saving.value = true
  try {
    const data = {
      name: form.value.name,
      sku: form.value.sku,
      category: form.value.category,
      unit: form.value.unit,
      costPrice: form.value.costPrice,
      salePrice: form.value.salePrice,
      minStock: form.value.minStock,
      maxStock: form.value.maxStock,
      status: form.value.status,
      description: form.value.description
    }

    let res: any
    if (isEdit.value && currentId.value) {
      res = await productApi.update(currentId.value, data)
    } else {
      res = await productApi.create(data)
    }

    if (res.code === 0) {
      toast.add({ severity: 'success', summary: t('common.saveSuccess'), life: 5000 })
      dialogVisible.value = false
      fetchProducts()
    } else {
      toast.add({ severity: 'error', summary: res.message || '保存失败', life: 5000 })
    }
  } catch (error: any) {
    const message = error.response?.data?.message || error.message || '保存失败'
    toast.add({ severity: 'error', summary: message, life: 5000 })
  } finally {
    saving.value = false
  }
}

const handleDelete = (_product: Product) => {
  confirm.require({
    message: t('common.confirmDelete'),
    header: t('common.confirm'),
    accept: async () => {
      try {
        // Note: productApi doesn't have delete method, need to add it
        toast.add({ severity: 'info', summary: '产品暂不支持删除', life: 5000 })
      } catch (error) {
        toast.add({ severity: 'error', summary: '删除失败', life: 5000 })
      }
    }
  })
}

onMounted(() => {
  initUserPermissions()
  fetchProducts()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.search-box {
  display: flex;
  gap: 8px;
}
</style>