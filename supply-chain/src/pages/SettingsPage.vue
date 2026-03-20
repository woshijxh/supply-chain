<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('settings.title') }}</h1>
    </div>

    <div class="settings-container">
      <!-- 个人信息 -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="ri-user-line"></i>
            {{ t('settings.profile') }}
          </h3>
          <Button
            :label="t('settings.changePassword')"
            severity="secondary"
            size="small"
            @click="showPasswordDialog = true"
          />
        </div>

        <!-- 头像上传 -->
        <div class="avatar-section">
          <div class="avatar-preview" @click="triggerAvatarUpload">
            <img v-if="avatarPreview" :src="avatarPreview" alt="avatar" />
            <i v-else class="ri-user-line"></i>
          </div>
          <input
            ref="avatarInput"
            type="file"
            accept="image/*"
            @change="handleAvatarChange"
            style="display: none"
          />
          <Button
            :label="t('settings.uploadAvatar')"
            severity="secondary"
            size="small"
            @click="triggerAvatarUpload"
          />
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label class="form-label">{{ t('settings.username') }}</label>
            <InputText v-model="profileForm.username" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.email') }}</label>
            <InputText v-model="profileForm.email" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.phone') }}</label>
            <InputText v-model="profileForm.phone" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.role') }}</label>
            <InputText v-model="profileForm.role" disabled />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.department') }}</label>
            <InputText v-model="profileForm.department" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.position') }}</label>
            <InputText v-model="profileForm.position" />
          </div>
        </div>
        <div class="dialog-footer">
          <Button :label="t('common.save')" @click="saveProfile" :loading="savingProfile" />
        </div>
      </div>

      <!-- 偏好设置 -->
      <div class="card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="ri-settings-line"></i>
            {{ t('settings.preferences') }}
          </h3>
        </div>
        <div class="preference-item">
          <div class="preference-label">
            <i class="ri-moon-line"></i>
            {{ t('settings.theme') }}
          </div>
          <div class="preference-control">
            <SelectButton
              v-model="themeValue"
              :options="themeOptions"
              optionLabel="label"
              optionValue="value"
              @change="handleThemeChange"
            />
          </div>
        </div>
        <div class="preference-item">
          <div class="preference-label">
            <i class="ri-global-line"></i>
            {{ t('settings.language') }}
          </div>
          <div class="preference-control">
            <SelectButton
              v-model="languageValue"
              :options="languageOptions"
              optionLabel="label"
              optionValue="value"
              @change="handleLanguageChange"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <Dialog
      v-model:visible="showPasswordDialog"
      :header="t('settings.changePassword')"
      :modal="true"
      :style="{ width: '400px' }"
      @hide="clearPasswordForm"
    >
      <div class="password-form">
        <div class="form-group">
          <label class="form-label">{{ t('settings.currentPassword') }}</label>
          <InputText v-model="passwordForm.currentPassword" type="password" class="w-full" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('settings.newPassword') }}</label>
          <InputText v-model="passwordForm.newPassword" type="password" class="w-full" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('settings.confirmPassword') }}</label>
          <InputText v-model="passwordForm.confirmPassword" type="password" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" @click="showPasswordDialog = false" />
        <Button :label="t('common.confirm')" @click="changePassword" :loading="changingPassword" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import SelectButton from 'primevue/selectbutton'
import Dialog from 'primevue/dialog'
import { authApi } from '@/api'

const { t, locale } = useI18n()
const toast = useToast()

// 用户信息
const user = computed(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch {
      return null
    }
  }
  return null
})

// 头像相关
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarPreview = ref('')
const uploadingAvatar = ref(false)

// 个人信息表单
const savingProfile = ref(false)
const profileForm = reactive({
  username: '',
  email: '',
  phone: '',
  role: '',
  department: '',
  position: ''
})

// 修改密码弹窗
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 清空密码表单
const clearPasswordForm = () => {
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

// 主题设置
const isDark = ref(localStorage.getItem('darkMode') === 'true')
const themeValue = ref(isDark.value ? 'dark' : 'light')
const themeOptions = computed(() => [
  { label: t('settings.themeLight'), value: 'light' },
  { label: t('settings.themeDark'), value: 'dark' }
])

// 语言设置
const languageValue = ref(locale.value)
const languageOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'EN', value: 'en-US' }
]

onMounted(() => {
  // 加载用户信息到表单
  if (user.value) {
    profileForm.username = user.value.username || ''
    profileForm.email = user.value.email || ''
    profileForm.phone = user.value.phone || ''
    profileForm.role = user.value.role || ''
    profileForm.department = user.value.department || ''
    profileForm.position = user.value.position || ''
    avatarPreview.value = user.value.avatar || ''
  }
})

// 触发头像上传
const triggerAvatarUpload = () => {
  avatarInput.value?.click()
}

// 处理头像选择
const handleAvatarChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 预览头像
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)

  // 上传到服务器
  try {
    uploadingAvatar.value = true
    const res: any = await authApi.updateAvatar(avatarPreview.value)
    if (res.code === 0) {
      // 更新本地存储的用户信息
      const updatedUser = { ...user.value, avatar: avatarPreview.value }
      localStorage.setItem('user', JSON.stringify(updatedUser))
      toast.add({ severity: 'success', summary: t('settings.avatarUpdated'), life: 5000 })
    } else {
      toast.add({ severity: 'error', summary: res.message || '头像上传失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.response?.data?.message || '头像上传失败', life: 5000 })
  } finally {
    uploadingAvatar.value = false
    target.value = '' // 清空input，允许重复选择同一文件
  }
}

// 保存个人信息
const saveProfile = () => {
  if (!user.value) return

  // 模拟保存成功
  const updatedUser = { ...user.value, ...profileForm }
  localStorage.setItem('user', JSON.stringify(updatedUser))
  toast.add({ severity: 'success', summary: t('settings.profileUpdated'), life: 5000 })
}

// 修改密码
const changePassword = async () => {
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    toast.add({ severity: 'warn', summary: t('settings.passwordMismatch'), life: 5000 })
    return
  }

  if (!passwordForm.currentPassword || !passwordForm.newPassword) {
    toast.add({ severity: 'warn', summary: '请填写完整密码信息', life: 5000 })
    return
  }

  try {
    changingPassword.value = true
    const res: any = await authApi.changePassword(passwordForm.currentPassword, passwordForm.newPassword)
    if (res.code === 0) {
      passwordForm.currentPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
      showPasswordDialog.value = false
      toast.add({ severity: 'success', summary: t('settings.passwordChanged'), life: 5000 })
    } else {
      toast.add({ severity: 'error', summary: res.message || '修改密码失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.response?.data?.message || '修改密码失败', life: 5000 })
  } finally {
    changingPassword.value = false
  }
}

// 切换主题
const handleThemeChange = () => {
  const dark = themeValue.value === 'dark'
  isDark.value = dark
  localStorage.setItem('darkMode', String(dark))
  document.documentElement.classList.toggle('dark-mode', dark)
}

// 切换语言
const handleLanguageChange = () => {
  locale.value = languageValue.value
  localStorage.setItem('locale', languageValue.value)
}
</script>

<style scoped lang="scss">
.settings-page {
  max-width: 900px;
}

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border);
}

.avatar-preview {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-strong) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.2s ease;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  i {
    font-size: 32px;
    color: #fff;
  }

  &:hover {
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  }
}

.preference-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid var(--border);

  &:last-child {
    border-bottom: none;
  }
}

.preference-label {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 500;

  i {
    font-size: 18px;
    color: var(--primary);
  }
}

.preference-control {
  :deep(.p-selectbutton) {
    .p-button {
      padding: 8px 16px;
    }
  }
}

.form-group {
  margin-bottom: 16px;

  .p-inputtext {
    width: 100%;
  }
}

.w-full {
  width: 100%;
}

.password-form {
  padding: 8px 0;
}
</style>
