<template>
  <div class="permission-page">
    <header class="page-header">
      <h1 class="page-title">权限管理</h1>
      <p class="page-subtitle">管理系统权限点和权限配置</p>
    </header>

    <div class="toolbar">
      <Button v-if="isAdmin()" label="新增权限" icon="ri-add-line" @click="openDialog()" />
      <div class="search-box">
        <InputText v-model="searchKey" placeholder="搜索权限名称" @keyup.enter="onSearch" />
        <Button icon="ri-search-line" severity="secondary" @click="onSearch" />
      </div>
    </div>

    <div class="card">
      <DataTable
        :value="permissions"
        stripedRows
        :loading="loading"
        :lazy="true"
        :paginator="true"
        :rows="rows"
        :rowsPerPageOptions="[10, 20, 50, 100]"
        v-model:first="first"
        :totalRecords="totalRecords"
        @page="onPage"
      >
        <Column field="id" header="ID" style="width: 80px"></Column>
        <Column field="name" header="权限名称" style="width: 150px"></Column>
        <Column field="code" header="权限编码"></Column>
        <Column field="type" header="类型">
          <template #body="{ data }">
            <span :class="['tag', getTypeClass(data.type)]">{{ getTypeText(data.type) }}</span>
          </template>
        </Column>
        <Column field="description" header="描述"></Column>
        <Column field="status" header="状态">
          <template #body="{ data }">
            <span :class="['tag', data.status === 1 ? 'success' : 'warning']">
              {{ data.status === 1 ? '启用' : '禁用' }}
            </span>
          </template>
        </Column>
        <Column field="createdAt" header="创建时间">
          <template #body="{ data }">
            {{ formatDate(data.createdAt) }}
          </template>
        </Column>
        <Column header="操作" style="width: 120px">
          <template #body="{ data }">
            <Button v-if="isAdmin()" text icon="ri-edit-line" @click="openEditDialog(data)" />
            <Button v-if="isAdmin()" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 权限表单对话框 -->
    <Dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑权限' : '新增权限'" :style="{ width: '500px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">权限名称 *</label>
          <InputText v-model="form.name" />
        </div>
        <div class="form-group">
          <label class="form-label">权限编码 *</label>
          <InputText v-model="form.code" />
        </div>
        <div class="form-group">
          <label class="form-label">权限类型</label>
          <Select v-model="form.type" :options="typeOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">描述</label>
          <Textarea v-model="form.description" rows="3" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" severity="secondary" @click="dialogVisible = false" />
        <Button :label="isEdit ? '更新' : '确定'" @click="handleSave" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { permissionApi } from '@/api'
import type { Permission } from '@/types'
import { usePermission } from '@/utils/usePermission'

const toast = useToast()
const { isAdmin, initUserPermissions } = usePermission()

const permissions = ref<Permission[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)
const first = ref(0)
const currentPage = ref(1)
const searchKey = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)

const form = reactive({
  name: '',
  code: '',
  type: 'button',
  description: '',
  status: 1
})

const typeOptions = [
  { label: '菜单权限', value: 'menu' },
  { label: '按钮权限', value: 'button' },
  { label: '数据权限', value: 'data' },
  { label: 'API权限', value: 'api' },
  { label: '其他', value: 'other' }
]

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

function getTypeClass(type: string) {
  const map: Record<string, string> = {
    menu: 'info',
    button: 'success',
    data: 'warning',
    api: 'danger',
    other: 'info'
  }
  return map[type] || 'info'
}

function getTypeText(type: string) {
  const map: Record<string, string> = {
    menu: '菜单',
    button: '按钮',
    data: '数据',
    api: 'API',
    other: '其他'
  }
  return map[type] || type
}

function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res = await permissionApi.list(currentPage.value, rows.value, searchKey.value)
    if (res.code === 0) {
      permissions.value = res.data.list
      totalRecords.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to load permissions:', error)
  } finally {
    loading.value = false
  }
}

function openDialog() {
  isEdit.value = false
  editId.value = null
  form.name = ''
  form.code = ''
  form.type = 'button'
  form.description = ''
  form.status = 1
  dialogVisible.value = true
}

function openEditDialog(permission: Permission) {
  isEdit.value = true
  editId.value = permission.id
  form.name = permission.name
  form.code = (permission as any).code || permission.name
  form.type = (permission as any).type || 'button'
  form.description = permission.description || ''
  form.status = (permission as any).status || 1
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (isEdit.value && editId.value) {
      await permissionApi.update(String(editId.value), {
        name: form.name,
        code: form.code,
        type: form.type,
        description: form.description,
        status: form.status
      })
      toast.add({ severity: 'success', summary: '成功', detail: '权限更新成功', life: 5000 })
    } else {
      await permissionApi.create({
        name: form.name,
        code: form.code,
        type: form.type,
        description: form.description,
        status: form.status
      })
      toast.add({ severity: 'success', summary: '成功', detail: '权限创建成功', life: 5000 })
    }
    dialogVisible.value = false
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '操作失败', life: 5000 })
  }
}

async function handleDelete(permission: Permission) {
  if (!confirm(`确定要删除权限 ${permission.name} 吗？`)) return

  try {
    await permissionApi.delete(String(permission.id))
    toast.add({ severity: 'success', summary: '成功', detail: '权限删除成功', life: 5000 })
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '删除失败', life: 5000 })
  }
}

function onSearch() {
  first.value = 0
  currentPage.value = 1
  loadData()
}

function onPage(event: any) {
  first.value = event.first
  rows.value = event.rows
  currentPage.value = Math.floor(event.first / event.rows) + 1
  loadData()
}

onMounted(() => {
  initUserPermissions()
  loadData()
})
</script>

<style scoped lang="scss">
.permission-page {
  padding: 24px;
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
  color: #666;
  margin: 0;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;

  .search-box {
    display: flex;
    gap: 8px;

    :deep(.p-inputtext) {
      width: 200px;
    }
  }
}

.card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
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

  &.full-width {
    grid-column: span 2;
  }
}

.form-label {
  font-weight: 500;
  color: #333;
}

.tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;

  &.success {
    background: #dcfce7;
    color: #166534;
  }

  &.warning {
    background: #fef3c7;
    color: #92400e;
  }

  &.info {
    background: #dbeafe;
    color: #1e40af;
  }

  &.danger {
    background: #fee2e2;
    color: #991b1b;
  }
}
</style>