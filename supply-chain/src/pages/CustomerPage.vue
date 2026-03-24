<template>
  <div class="customer-page">
    <header class="page-header">
      <h1 class="page-title">客户管理</h1>
      <p class="page-subtitle">管理客户信息和等级</p>
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
        :value="customers"
        stripedRows
        :loading="loading"
        :paginator="true"
        :rows="10"
      >
        <Column field="code" header="客户编码" style="width: 120px"></Column>
        <Column field="name" header="客户名称" style="width: 150px"></Column>
        <Column field="contact" header="联系人"></Column>
        <Column field="phone" header="联系电话"></Column>
        <Column field="level" header="等级">
          <template #body="{ data }">
            <span :class="['tag', getLevelClass(data.level)]">{{ getLevelText(data.level) }}</span>
          </template>
        </Column>
        <Column field="source" header="来源"></Column>
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

    <Dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑客户' : '新增客户'" :style="{ width: '600px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">客户名称 *</label>
          <InputText v-model="form.name" />
        </div>
        <div class="form-group">
          <label class="form-label">联系人</label>
          <InputText v-model="form.contact" />
        </div>
        <div class="form-group">
          <label class="form-label">联系电话</label>
          <InputText v-model="form.phone" />
        </div>
        <div class="form-group">
          <label class="form-label">邮箱</label>
          <InputText v-model="form.email" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">地址</label>
          <InputText v-model="form.address" />
        </div>
        <div class="form-group">
          <label class="form-label">客户等级</label>
          <Select v-model="form.level" :options="levelOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group">
          <label class="form-label">客户来源</label>
          <InputText v-model="form.source" />
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group" v-if="isEdit">
          <label class="form-label">客户编码</label>
          <InputText v-model="form.code" disabled />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">{{ t('common.remark') }}</label>
          <Textarea v-model="form.remark" rows="3" />
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
import { customerApi } from '@/api'
import type { Customer } from '@/types'
import { useUserStore } from '@/stores/user'
import { CUSTOMER_PERMISSIONS } from '@/config/permissions'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const userStore = useUserStore()

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const customers = ref<Customer[]>([])
const searchKey = ref('')
const currentId = ref<number | null>(null)

const form = ref({
  name: '',
  contact: '',
  phone: '',
  email: '',
  address: '',
  level: 'C',
  source: '',
  status: 1,
  remark: '',
  code: ''
})

const levelOptions = [
  { label: 'A类客户', value: 'A' },
  { label: 'B类客户', value: 'B' },
  { label: 'C类客户', value: 'C' }
]

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

const hasCreatePermission = computed(() => userStore.hasPermission(CUSTOMER_PERMISSIONS.CREATE))
const hasUpdatePermission = computed(() => userStore.hasPermission(CUSTOMER_PERMISSIONS.UPDATE))
const hasDeletePermission = computed(() => userStore.hasPermission(CUSTOMER_PERMISSIONS.DELETE))

const getLevelText = (level: string) => {
  const map: Record<string, string> = { A: 'A类', B: 'B类', C: 'C类' }
  return map[level] || level
}

const getLevelClass = (level: string) => {
  const map: Record<string, string> = { A: 'success', B: 'info', C: 'secondary' }
  return map[level] || 'secondary'
}

const fetchCustomers = async () => {
  loading.value = true
  try {
    const res: any = await customerApi.list(1, 100, searchKey.value)
    if (res.code === 0) {
      customers.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取客户列表失败', error)
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  fetchCustomers()
}

const openDialog = (customer?: Customer) => {
  if (customer) {
    isEdit.value = true
    currentId.value = customer.id
    form.value = {
      name: customer.name,
      contact: customer.contact || '',
      phone: customer.phone || '',
      email: customer.email || '',
      address: customer.address || '',
      level: customer.level || 'C',
      source: customer.source || '',
      status: customer.status,
      remark: customer.remark || '',
      code: customer.code || ''
    }
  } else {
    isEdit.value = false
    currentId.value = null
    form.value = {
      name: '',
      contact: '',
      phone: '',
      email: '',
      address: '',
      level: 'C',
      source: '',
      status: 1,
      remark: '',
      code: ''
    }
  }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!form.value.name) {
    toast.add({ severity: 'warn', summary: '请填写客户名称', life: 5000 })
    return
  }

  saving.value = true
  try {
    const data = {
      name: form.value.name,
      contact: form.value.contact,
      phone: form.value.phone,
      email: form.value.email,
      address: form.value.address,
      level: form.value.level,
      source: form.value.source,
      status: form.value.status,
      remark: form.value.remark
    }

    let res: any
    if (isEdit.value && currentId.value) {
      res = await customerApi.update(String(currentId.value), data)
    } else {
      res = await customerApi.create(data)
    }

    if (res.code === 0) {
      toast.add({ severity: 'success', summary: t('common.saveSuccess'), life: 5000 })
      dialogVisible.value = false
      fetchCustomers()
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

const handleDelete = (customer: Customer) => {
  confirm.require({
    message: t('common.confirmDelete'),
    header: t('common.confirm'),
    accept: async () => {
      try {
        const res: any = await customerApi.delete(String(customer.id))
        if (res.code === 0) {
          toast.add({ severity: 'success', summary: t('common.deleteSuccess'), life: 5000 })
          fetchCustomers()
        }
      } catch (error) {
        toast.add({ severity: 'error', summary: '删除失败', life: 5000 })
      }
    }
  })
}

onMounted(() => {
  fetchCustomers()
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

.card {
  background: var(--card-bg);
  border-radius: var(--radius-lg);
  padding: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-weight: 500;
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
}

.tag {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.tag.success {
  background: var(--green-100);
  color: var(--green-700);
}

.tag.warning {
  background: var(--yellow-100);
  color: var(--yellow-700);
}

.tag.info {
  background: var(--blue-100);
  color: var(--blue-700);
}

.tag.secondary {
  background: var(--gray-100);
  color: var(--gray-700);
}
</style>