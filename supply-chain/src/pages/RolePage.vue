<template>
  <div class="role-page">
    <header class="page-header">
      <h1 class="page-title">角色管理</h1>
      <p class="page-subtitle">管理系统角色和权限分配</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(ROLE_PERMISSIONS.WRITE) || isAdmin()" label="新增角色" icon="ri-add-line" @click="openDialog()" />
      <div class="search-box">
        <InputText v-model="searchKey" placeholder="搜索角色名称" @keyup.enter="onSearch" />
        <Button icon="ri-search-line" severity="secondary" @click="onSearch" />
      </div>
    </div>

    <div class="card">
      <DataTable
        :value="roles"
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
        <Column field="name" header="角色名称" style="width: 150px"></Column>
        <Column field="code" header="角色编码"></Column>
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
        <Column header="操作" style="width: 200px">
          <template #body="{ data }">
            <Button v-if="hasPermission(ROLE_PERMISSIONS.WRITE) || isAdmin()" text severity="info" icon="ri-edit-line" @click="openDialog(data)" />
            <Button v-if="hasPermission(ROLE_PERMISSIONS.WRITE) || isAdmin()" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
            <Button text severity="secondary" icon="ri-shield-check-line" @click="openPermissionDialog(data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 角色表单对话框 -->
    <Dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑角色' : '新增角色'" :style="{ width: '500px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">角色名称 *</label>
          <InputText v-model="form.name" />
        </div>
        <div class="form-group">
          <label class="form-label">角色编码 *</label>
          <InputText v-model="form.code" />
        </div>
        <div class="form-group" style="grid-column: span 2">
          <label class="form-label">描述</label>
          <Textarea v-model="form.description" rows="3" />
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" severity="secondary" @click="dialogVisible = false" />
        <Button label="确定" @click="handleSave" />
      </template>
    </Dialog>

    <!-- 权限分配对话框 -->
    <Dialog v-model:visible="permissionDialogVisible" header="分配权限" :style="{ width: '500px' }">
      <div class="permission-tree">
        <div v-for="group in permissionGroups" :key="group.type" class="permission-group">
          <div class="group-title">{{ group.label }}</div>
          <div class="permission-list">
            <div v-for="perm in group.permissions" :key="perm.id" class="permission-item">
              <Checkbox v-model="selectedPermissionIds" :inputId="`perm-${perm.id}`" :value="perm.id" />
              <label :for="`perm-${perm.id}`">{{ perm.name }}</label>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="取消" severity="secondary" @click="permissionDialogVisible = false" />
        <Button label="确定" @click="handleSetPermissions" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { roleApi, permissionApi } from '@/api'
import type { Role, Permission } from '@/types'
import { usePermission } from '@/utils/usePermission'
import { ROLE_PERMISSIONS } from '@/config/permissions'

const toast = useToast()
const { isAdmin, hasPermission, initUserPermissions } = usePermission()

const roles = ref<Role[]>([])
const allPermissions = ref<Permission[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)
const first = ref(0)
const searchKey = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const currentRoleId = ref<number | null>(null)

const form = reactive({
  name: '',
  code: '',
  description: '',
  status: 1
})

const permissionDialogVisible = ref(false)
const selectedRoleId = ref<number | null>(null)
const selectedPermissionIds = ref<number[]>([])

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 按类型分组权限
const permissionGroups = computed(() => {
  const groups: Record<string, { label: string; permissions: Permission[] }> = {}

  allPermissions.value.forEach(perm => {
    const type = perm.type || 'other'
    if (!groups[type]) {
      const labels: Record<string, string> = {
        menu: '菜单权限',
        button: '按钮权限',
        data: '数据权限',
        api: 'API权限',
        other: '其他权限'
      }
      groups[type] = { label: labels[type] || type, permissions: [] }
    }
    groups[type].permissions.push(perm)
  })

  return Object.values(groups)
})

function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res = await roleApi.list(1, rows.value, searchKey.value)
    if (res.code === 0) {
      roles.value = res.data.list
      totalRecords.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to load roles:', error)
  } finally {
    loading.value = false
  }
}

async function loadPermissions() {
  try {
    const res = await permissionApi.getAll()
    if (res.code === 0) {
      allPermissions.value = res.data
    }
  } catch (error) {
    console.error('Failed to load permissions:', error)
  }
}

function openDialog(role?: Role) {
  if (role) {
    isEdit.value = true
    currentRoleId.value = role.id
    form.name = role.name
    form.code = role.code
    form.description = role.description || ''
    form.status = role.status
  } else {
    isEdit.value = false
    currentRoleId.value = null
    form.name = ''
    form.code = ''
    form.description = ''
    form.status = 1
  }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (isEdit.value && currentRoleId.value) {
      await roleApi.update(String(currentRoleId.value), {
        name: form.name,
        code: form.code,
        description: form.description,
        status: form.status
      })
      toast.add({ severity: 'success', summary: '成功', detail: '角色更新成功', life: 5000 })
    } else {
      await roleApi.create({
        name: form.name,
        code: form.code,
        description: form.description,
        status: form.status
      })
      toast.add({ severity: 'success', summary: '成功', detail: '角色创建成功', life: 5000 })
    }
    dialogVisible.value = false
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '操作失败', life: 5000 })
  }
}

async function handleDelete(role: Role) {
  if (!confirm(`确定要删除角色 ${role.name} 吗？`)) return

  try {
    await roleApi.delete(String(role.id))
    toast.add({ severity: 'success', summary: '成功', detail: '角色删除成功', life: 5000 })
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '删除失败', life: 5000 })
  }
}

async function openPermissionDialog(role: Role) {
  selectedRoleId.value = role.id
  try {
    const res = await roleApi.getPermissions(String(role.id))
    if (res.code === 0) {
      selectedPermissionIds.value = res.data.map((p: Permission) => p.id)
    }
  } catch (error) {
    selectedPermissionIds.value = []
  }
  permissionDialogVisible.value = true
}

async function handleSetPermissions() {
  try {
    await roleApi.setPermissions(String(selectedRoleId.value), selectedPermissionIds.value)
    toast.add({ severity: 'success', summary: '成功', detail: '权限分配成功', life: 5000 })
    permissionDialogVisible.value = false
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '操作失败', life: 5000 })
  }
}

function onSearch() {
  first.value = 0
  loadData()
}

function onPage(event: any) {
  first.value = event.first
  rows.value = event.rows
  loadData()
}

onMounted(() => {
  initUserPermissions()
  loadData()
  loadPermissions()
})
</script>

<style scoped lang="scss">
.role-page {
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

.permission-tree {
  max-height: 400px;
  overflow-y: auto;
}

.permission-group {
  margin-bottom: 16px;

  .group-title {
    font-weight: 600;
    color: #333;
    margin-bottom: 8px;
    padding-bottom: 8px;
    border-bottom: 1px solid #eee;
  }
}

.permission-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.permission-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #f5f5f5;
  border-radius: 4px;

  label {
    cursor: pointer;
    font-size: 14px;
  }
}
</style>