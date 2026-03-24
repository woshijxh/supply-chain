<template>
  <div class="supplier-page">
    <header class="page-header">
      <h1 class="page-title">{{ t('supplier.title') }}</h1>
      <p class="page-subtitle">管理供应商信息和合作状态</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(SUPPLIER_PERMISSIONS.CREATE) || canCreate()" :label="t('common.add')" icon="ri-add-line" @click="openDialog()" />
      <div class="search-box">
        <InputText v-model="searchKey" :placeholder="t('common.search')" @keyup.enter="onSearch" />
        <Button icon="ri-search-line" severity="secondary" @click="onSearch" />
      </div>
    </div>

    <div class="card">
      <DataTable
        :value="suppliers"
        stripedRows
        :loading="loading"
        :lazy="true"
        :paginator="true"
        :rows="rows"
        v-model:first="first"
        :totalRecords="totalRecords"
        @page="onPage"
      >
        <Column field="code" :header="t('supplier.code')" style="width: 120px"></Column>
        <Column field="name" :header="t('supplier.name')" style="width: 200px"></Column>
        <Column field="contact" :header="t('supplier.contact')"></Column>
        <Column field="phone" :header="t('supplier.phone')"></Column>
        <Column field="level" :header="t('supplier.level')">
          <template #body="{ data }">
            <span :class="['tag', getLevelClass(data.level)]">{{ getLevelText(data.level) }}</span>
          </template>
        </Column>
        <Column field="category" :header="t('supplier.category')"></Column>
        <Column field="status" :header="t('supplier.status')">
          <template #body="{ data }">
            <span :class="['tag', data.status === 'active' ? 'success' : 'warning']">
              {{ data.status === 'active' ? t('common.enable') : t('common.disable') }}
            </span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width: 150px">
          <template #body="{ data }">
            <Button v-if="hasPermission(SUPPLIER_PERMISSIONS.UPDATE) || canEdit(data.creatorId, data.creatorDepartment)" text severity="info" icon="ri-edit-line" @click="openDialog(data)" />
            <Button v-if="hasPermission(SUPPLIER_PERMISSIONS.DELETE) || canDelete()" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="dialogVisible" :header="isEdit ? t('common.edit') : t('common.add')" :style="{ width: '600px' }">
      <div class="form-grid">
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">{{ t('supplier.name') }} *</label>
          <InputText v-model="form.name" :class="{ 'p-invalid': v$.name.$dirty && v$.name.$invalid }" @blur="v$.name.$touch()" />
          <small v-if="v$.name.$dirty && v$.name.$invalid" class="p-error">{{ v$.name.$errors[0]?.$message || '请输入供应商名称' }}</small>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.contact') }} *</label>
          <InputText v-model="form.contact" :class="{ 'p-invalid': v$.contact.$dirty && v$.contact.$invalid }" @blur="v$.contact.$touch()" />
          <small v-if="v$.contact.$dirty && v$.contact.$invalid" class="p-error">{{ v$.contact.$errors[0]?.$message || '请输入联系人' }}</small>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.phone') }} *</label>
          <InputText v-model="form.phone" :class="{ 'p-invalid': v$.phone.$dirty && v$.phone.$invalid }" @blur="v$.phone.$touch()" />
          <small v-if="v$.phone.$dirty && v$.phone.$invalid" class="p-error">{{ v$.phone.$errors[0]?.$message || '请输入联系电话' }}</small>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.email') }}</label>
          <InputText v-model="form.email" :class="{ 'p-invalid': v$.email.$dirty && v$.email.$invalid }" @blur="v$.email.$touch()" />
          <small v-if="v$.email.$dirty && v$.email.$invalid" class="p-error">{{ v$.email.$errors[0]?.$message || '请输入有效的邮箱地址' }}</small>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.level') }}</label>
          <Select v-model="form.level" :options="levelOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.category') }}</label>
          <InputText v-model="form.category" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.paymentTerms') }}</label>
          <InputText v-model="form.paymentTerms" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.status') }}</label>
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">{{ t('supplier.address') }}</label>
          <InputText v-model="form.address" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.bankName') }}</label>
          <InputText v-model="form.bankName" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.bankAccount') }}</label>
          <InputText v-model="form.bankAccount" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('supplier.taxNumber') }}</label>
          <InputText v-model="form.taxNumber" />
        </div>
        <div class="form-group" v-if="isEdit">
          <label class="form-label">供应商编码</label>
          <InputText v-model="form.code" disabled />
          <small class="p-text-secondary">系统自动生成</small>
        </div>
        <div class="form-group" v-else>
          <label class="form-label">供应商编码</label>
          <InputText value="保存后自动生成" disabled />
          <small class="p-text-secondary">格式: SUPYYYYMMDD####</small>
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
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useVuelidate } from '@vuelidate/core'
import { required, email as emailValidator, helpers } from '@vuelidate/validators'
import { supplierApi } from '@/api'
import type { Supplier } from '@/types'
import { usePermission } from '@/utils/usePermission'
import { SUPPLIER_PERMISSIONS } from '@/config/permissions'

const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const { canCreate, canEdit, canDelete, hasPermission, initUserPermissions } = usePermission()

const loading = ref(false)
const saving = ref(false)
const suppliers = ref<Supplier[]>([])
const searchKey = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref('')
const currentPage = ref(1)
const rows = ref(10)
const totalRecords = ref(0)
const first = ref(0)

const form = ref<Omit<Partial<Supplier>, 'status'> & { status: string }>({
  code: '',
  name: '',
  contact: '',
  phone: '',
  email: '',
  address: '',
  level: 'B',
  category: '',
  paymentTerms: '',
  bankName: '',
  bankAccount: '',
  taxNumber: '',
  status: 'active',
  remark: ''
})

// vuelidate 验证规则对象
const formRules = computed(() => ({
  name: form.value.name,
  contact: form.value.contact,
  phone: form.value.phone,
  email: form.value.email,
  level: form.value.level,
  address: form.value.address,
  category: form.value.category,
  bankAccount: form.value.bankAccount,
  taxNumber: form.value.taxNumber
}))

// 手机号验证规则
const phoneValidator = (value: string) => {
  if (!value) return true
  // 支持手机号和座机号
  const phoneRegex = /^1[3-9]\d{9}$|^0\d{2,3}-?\d{7,8}$/
  return phoneRegex.test(value)
}

const rules = computed(() => ({
  name: { required: helpers.withMessage('供应商名称不能为空', required) },
  contact: { required: helpers.withMessage('联系人不能为空', required) },
  phone: {
    required: helpers.withMessage('联系电话不能为空', required),
    phone: helpers.withMessage('请输入有效的电话号码', phoneValidator)
  },
  email: { email: helpers.withMessage('请输入有效的邮箱地址', emailValidator) },
  level: {},
  address: { maxLength: helpers.withMessage('地址最多255个字符', (value: string) => !value || value.length <= 255) },
  category: { maxLength: helpers.withMessage('分类最多50个字符', (value: string) => !value || value.length <= 50) },
  bankAccount: { maxLength: helpers.withMessage('银行账号最多50个字符', (value: string) => !value || value.length <= 50) },
  taxNumber: { maxLength: helpers.withMessage('税号最多50个字符', (value: string) => !value || value.length <= 50) }
}))

const v$ = useVuelidate(rules, formRules)

const levelOptions = [
  { label: t('supplier.levelA'), value: 'A' },
  { label: t('supplier.levelB'), value: 'B' },
  { label: t('supplier.levelC'), value: 'C' }
]

const statusOptions = [
  { label: t('common.enable'), value: 'active' },
  { label: t('common.disable'), value: 'inactive' }
]

const getLevelText = (level: string) => {
  const map: Record<string, string> = { A: 'A级', B: 'B级', C: 'C级' }
  return map[level] || level
}

const getLevelClass = (level: string) => {
  const map: Record<string, string> = { A: 'success', B: 'info', C: 'warning' }
  return map[level] || 'info'
}

// 后端status是int8: 1=active, 0=inactive，前端转换为字符串
const mapStatusFromBackend = (status: number | string): string => {
  if (typeof status === 'string') return status
  return status === 1 ? 'active' : 'inactive'
}

const mapStatusToBackend = (status: string): number => {
  return status === 'active' ? 1 : 0
}

const fetchSuppliers = async () => {
  loading.value = true
  try {
    const res: any = await supplierApi.list(currentPage.value, rows.value, searchKey.value)
    if (res.code === 0) {
      suppliers.value = (res.data.list || []).map((s: any) => ({
        ...s,
        status: mapStatusFromBackend(s.status)
      }))
      totalRecords.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取供应商列表失败', error)
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  currentPage.value = 1
  first.value = 0
  fetchSuppliers()
}

const onPage = (event: any) => {
  currentPage.value = (event.first / event.rows) + 1
  rows.value = event.rows
  fetchSuppliers()
}

const openDialog = (supplier?: Supplier) => {
  if (supplier) {
    isEdit.value = true
    editingId.value = String(supplier.id)
    form.value = {
      ...supplier,
      status: mapStatusFromBackend(supplier.status)
    }
  } else {
    isEdit.value = false
    editingId.value = ''
    form.value = {
      code: '',
      name: '',
      contact: '',
      phone: '',
      email: '',
      address: '',
      level: 'B',
      category: '',
      paymentTerms: '',
      bankName: '',
      bankAccount: '',
      taxNumber: '',
      status: 'active',
      remark: ''
    }
  }
  v$.value.$reset()
  dialogVisible.value = true
}

const handleSave = async () => {
  const isValid = await v$.value.$validate()
  if (!isValid) {
    toast.add({ severity: 'warn', summary: '请检查表单填写是否正确', life: 5000 })
    return
  }

  saving.value = true
  try {
    let res: any
    // 提交时将status转换为后端格式
    const submitData: any = {
      ...form.value,
      status: mapStatusToBackend(form.value.status as string)
    }
    // 删除编码字段，让后端自动生成
    if (!isEdit.value) {
      delete submitData.code
    }
    if (isEdit.value) {
      res = await supplierApi.update(editingId.value, submitData)
    } else {
      res = await supplierApi.create(submitData)
    }

    if (res.code === 0) {
      toast.add({ severity: 'success', summary: t('common.saveSuccess'), life: 5000 })
      dialogVisible.value = false
      fetchSuppliers()
    } else {
      toast.add({ severity: 'error', summary: res.message || '保存失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.message || '保存失败', life: 5000 })
  } finally {
    saving.value = false
  }
}

const handleDelete = (supplier: Supplier) => {
  confirm.require({
    message: t('common.confirmDelete'),
    header: t('common.confirm'),
    icon: 'ri-error-warning-line',
    accept: async () => {
      try {
        const res: any = await supplierApi.delete(String(supplier.id))
        if (res.code === 0) {
          toast.add({ severity: 'success', summary: t('common.deleteSuccess'), life: 5000 })
          fetchSuppliers()
        } else {
          toast.add({ severity: 'error', summary: res.message || '删除失败', life: 5000 })
        }
      } catch (error) {
        toast.add({ severity: 'error', summary: '删除失败', life: 5000 })
      }
    }
  })
}

onMounted(() => {
  initUserPermissions()
  fetchSuppliers()
})
</script>