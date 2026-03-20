<template>
  <div class="user-page">
    <header class="page-header">
      <h1 class="page-title">用户管理</h1>
      <p class="page-subtitle">管理系统用户和权限分配</p>
    </header>

    <div class="toolbar">
      <Button v-if="hasPermission(USER_PERMISSIONS.WRITE) || canCreate()" label="新增用户" icon="ri-add-line" @click="openDialog()" />
      <div class="search-box">
        <InputText v-model="searchKey" placeholder="搜索用户名/邮箱" @keyup.enter="onSearch" />
        <Button icon="ri-search-line" severity="secondary" @click="onSearch" />
      </div>
    </div>

    <div class="card">
      <DataTable
        :value="users"
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
        <Column field="username" header="用户名" style="width: 150px"></Column>
        <Column field="email" header="邮箱"></Column>
        <Column field="phone" header="手机号"></Column>
        <Column header="角色" style="width: 150px">
          <template #body="{ data }">
            <div class="role-tags">
              <span v-for="role in getUserRoles(data)" :key="role.id" :class="['tag', getRoleClass(role.name)]">
                {{ role.name }}
              </span>
            </div>
          </template>
        </Column>
        <Column header="权限" style="width: 220px">
          <template #body="{ data }">
            <div class="permission-tags">
              <template v-if="getUserPermissions(data).length <= 3">
                <span v-for="perm in getUserPermissions(data)" :key="perm.id" class="perm-tag">
                  {{ perm.name }}
                </span>
              </template>
              <template v-else>
                <span v-for="perm in getUserPermissions(data).slice(0, 3)" :key="perm.id" class="perm-tag">
                  {{ perm.name }}
                </span>
                <span class="perm-tag perm-more" v-tooltip.top="getUserPermissions(data).map(p => p.name).join('\n')">
                  +{{ getUserPermissions(data).length - 3 }}
                </span>
              </template>
            </div>
          </template>
        </Column>
        <Column field="department" header="部门"></Column>
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
        <Column header="操作" style="width: 250px">
          <template #body="{ data }">
            <Button v-if="hasPermission(USER_PERMISSIONS.WRITE) || canEdit()" text severity="info" icon="ri-edit-line" @click="openDialog(data)" />
            <Button v-if="hasPermission(USER_PERMISSIONS.WRITE) || canDelete()" text severity="danger" icon="ri-delete-bin-line" @click="handleDelete(data)" />
            <Button v-if="isAdmin()" text severity="secondary" icon="ri-shield-line" @click="openRoleDialog(data)" title="分配角色" />
            <Button v-if="isAdmin()" text severity="warning" icon="ri-key-line" @click="openPermDialog(data)" title="分配权限" />
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 用户表单对话框 -->
    <Dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑用户' : '新增用户'" :style="{ width: '500px' }">
      <div class="form-grid">
        <div class="form-group">
          <label class="form-label">用户名 *</label>
          <InputText v-model="form.username" :disabled="isEdit" />
        </div>
        <div class="form-group">
          <label class="form-label">邮箱</label>
          <InputText v-model="form.email" />
        </div>
        <div class="form-group">
          <label class="form-label">手机号</label>
          <InputText v-model="form.phone" />
        </div>
        <div class="form-group">
          <label class="form-label">密码 {{ isEdit ? '' : '*' }}</label>
          <InputText v-model="form.password" type="password" />
        </div>
        <div class="form-group">
          <label class="form-label">部门</label>
          <InputText v-model="form.department" />
        </div>
        <div class="form-group">
          <label class="form-label">职位</label>
          <InputText v-model="form.position" />
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

    <!-- 角色分配对话框 -->
    <Dialog v-model:visible="roleDialogVisible" header="分配角色" :style="{ width: '400px' }">
      <div class="form-group">
        <label class="form-label">选择角色</label>
        <div class="checkbox-list">
          <div v-for="role in allRoles" :key="role.id" class="checkbox-item">
            <Checkbox v-model="selectedRoleIds" :inputId="`role-${role.id}`" :value="role.id" />
            <label :for="`role-${role.id}`">{{ role.name }}</label>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="取消" severity="secondary" @click="roleDialogVisible = false" />
        <Button label="确定" @click="handleSetRoles" />
      </template>
    </Dialog>

    <!-- 权限分配对话框 -->
    <Dialog v-model:visible="permDialogVisible" header="分配权限" :style="{ width: '500px' }">
      <div class="form-group">
        <label class="form-label">选择权限</label>
        <div class="checkbox-list permission-list">
          <div v-for="perm in allPermissions" :key="perm.id" class="checkbox-item">
            <Checkbox v-model="selectedPermIds" :inputId="`perm-${perm.id}`" :value="perm.id" />
            <label :for="`perm-${perm.id}`">
              <span class="perm-name">{{ perm.name }}</span>
              <span class="perm-desc">{{ perm.description }}</span>
            </label>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="取消" severity="secondary" @click="permDialogVisible = false" />
        <Button label="确定" @click="handleSetPermissions" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { userApi, roleApi, permissionApi } from '@/api'
import type { User, Role } from '@/types'
import { usePermission } from '@/utils/usePermission'
import { USER_PERMISSIONS } from '@/config/permissions'

const toast = useToast()
const { canCreate, canEdit, canDelete, isAdmin, hasPermission, initUserPermissions } = usePermission()

const users = ref<User[]>([])
const loading = ref(false)
const totalRecords = ref(0)
const rows = ref(10)
const first = ref(0)
const searchKey = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const currentUserId = ref<number | null>(null)

const form = reactive({
  username: '',
  email: '',
  phone: '',
  password: '',
  role: 'operator',
  department: '',
  position: '',
  status: 1
})

const roleDialogVisible = ref(false)
const selectedUserId = ref<number | null>(null)
const selectedRoleIds = ref<number[]>([])
const allRoles = ref<any[]>([])

const permDialogVisible = ref(false)
const selectedPermIds = ref<number[]>([])
const allPermissions = ref<any[]>([])

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

function getRoleClass(role: string) {
  const map: Record<string, string> = {
    admin: 'danger',
    manager: 'info',
    operator: 'success'
  }
  return map[role] || 'info'
}

function getUserRoles(user: User) {
  if (user.roles && user.roles.length > 0) {
    return user.roles
  }
  // 如果没有 roles 字段，从 role 字段兼容
  return user.role ? [{ id: 0, name: user.role, code: user.role, description: '', status: 1, permissions: [], createdAt: '' }] : []
}

function getUserPermissions(user: User) {
  // 从 roles 中收集所有权限
  const perms: any[] = []
  if (user.roles) {
    for (const role of user.roles) {
      if (role.permissions) {
        for (const perm of role.permissions) {
          // 去重
          if (!perms.find(p => p.id === perm.id)) {
            perms.push(perm)
          }
        }
      }
    }
  }
  return perms
}

function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const res = await userApi.list(1, rows.value, searchKey.value)
    if (res.code === 0) {
      users.value = res.data.list
      totalRecords.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  try {
    const res = await roleApi.getAll()
    if (res.code === 0) {
      allRoles.value = res.data
    }
  } catch (error) {
    console.error('Failed to load roles:', error)
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

async function openPermDialog(user: User) {
  selectedUserId.value = Number(user.id)
  try {
    const res = await userApi.getPermissions(String(user.id))
    if (res.code === 0) {
      selectedPermIds.value = res.data.map((p: any) => p.id)
    }
  } catch (error) {
    selectedPermIds.value = []
  }
  permDialogVisible.value = true
}

function openDialog(user?: User) {
  if (user) {
    isEdit.value = true
    currentUserId.value = Number(user.id)
    form.username = user.username
    form.email = user.email || ''
    form.phone = user.phone || ''
    form.password = ''
    form.role = user.role
    form.department = user.department || ''
    form.position = user.profile?.position || ''
    form.status = user.status
  } else {
    isEdit.value = false
    currentUserId.value = null
    form.username = ''
    form.email = ''
    form.phone = ''
    form.password = ''
    form.role = 'operator'
    form.department = ''
    form.position = ''
    form.status = 1
  }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (isEdit.value && currentUserId.value) {
      const updateData: any = {
        email: form.email,
        phone: form.phone,
        role: form.role,
        department: form.department,
        position: form.position,
        status: form.status
      }
      if (form.password) {
        updateData.password = form.password
      }
      await userApi.update(String(currentUserId.value), updateData)
      toast.add({ severity: 'success', summary: '成功', detail: '用户更新成功', life: 5000 })
    } else {
      await userApi.create({
        username: form.username,
        email: form.email,
        phone: form.phone,
        password: form.password,
        role: form.role,
        department: form.department,
        position: form.position,
        status: form.status
      })
      toast.add({ severity: 'success', summary: '成功', detail: '用户创建成功', life: 5000 })
    }
    dialogVisible.value = false
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '操作失败', life: 5000 })
  }
}

async function handleDelete(user: User) {
  if (!confirm(`确定要删除用户 ${user.username} 吗？`)) return

  try {
    await userApi.delete(String(user.id))
    toast.add({ severity: 'success', summary: '成功', detail: '用户删除成功', life: 5000 })
    loadData()
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '删除失败', life: 5000 })
  }
}

async function openRoleDialog(user: User) {
  selectedUserId.value = Number(user.id)
  try {
    const res = await userApi.getRoles(String(user.id))
    if (res.code === 0) {
      selectedRoleIds.value = res.data.map((r: Role) => r.id)
    }
  } catch (error) {
    selectedRoleIds.value = []
  }
  roleDialogVisible.value = true
}

async function handleSetRoles() {
  try {
    await userApi.setRoles(String(selectedUserId.value), selectedRoleIds.value)
    toast.add({ severity: 'success', summary: '成功', detail: '角色分配成功', life: 5000 })
    roleDialogVisible.value = false
    loadData() // 刷新列表
  } catch (error: any) {
    toast.add({ severity: 'error', summary: '错误', detail: error.response?.data?.message || '操作失败', life: 5000 })
  }
}

async function handleSetPermissions() {
  try {
    await userApi.setPermissions(String(selectedUserId.value), selectedPermIds.value)
    toast.add({ severity: 'success', summary: '成功', detail: '权限分配成功', life: 5000 })
    permDialogVisible.value = false
    loadData() // 刷新列表
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
  loadRoles()
  loadPermissions()
})
</script>

<style scoped lang="scss">
.user-page {
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

.checkbox-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;

  label {
    display: flex;
    flex-direction: column;
    cursor: pointer;
    flex: 1;

    .perm-name {
      font-weight: 500;
    }

    .perm-desc {
      font-size: 12px;
      color: #666;
    }
  }
}

.role-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.permission-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;

  .perm-tag {
    background: #e0e7ff;
    color: #3730a3;
    padding: 2px 6px;
    border-radius: 3px;
    font-size: 11px;
  }

  .perm-more {
    background: #fef3c7;
    color: #92400e;
    padding: 2px 6px;
    border-radius: 3px;
    font-size: 11px;
    cursor: pointer;
    font-weight: 500;
  }
}
</style>