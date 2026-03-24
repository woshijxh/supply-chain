<template>
  <div class="login-page">
    <div class="login-left">
      <div class="bg-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>
      <div class="brand-content">
        <div class="brand-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h1>供应链管理系统</h1>
        <p>Supply Chain Management System</p>
      </div>
    </div>
    <div class="login-right">
      <div class="login-card">
        <div class="login-header">
          <h2>欢迎回来</h2>
          <span>请登录您的账号</span>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label>用户名</label>
            <InputText v-model="form.username" placeholder="请输入用户名" />
          </div>
          <div class="form-group">
            <label>密码</label>
            <InputText v-model="form.password" type="password" placeholder="请输入密码" />
          </div>
          <div class="form-group captcha-group">
            <label>验证码</label>
            <div class="captcha-row">
              <InputText v-model="form.captcha" placeholder="请输入验证码" @keyup.enter="handleLogin" />
              <div class="captcha-img" @click="refreshCaptcha">
                <img v-if="captchaUrl" :src="captchaUrl" alt="验证码" />
                <span v-else class="captcha-loading">加载中...</span>
              </div>
            </div>
          </div>
          <Button type="submit" label="登 录" :loading="loading" />
        </form>

        <p class="login-tip">默认账号: admin / 123456</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { authApi } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const toast = useToast()
const userStore = useUserStore()

const loading = ref(false)
const captchaId = ref('')
const captchaUrl = ref('')
const form = ref({
  username: 'admin',
  password: '123456',
  captcha: ''
})

async function refreshCaptcha() {
  try {
    const response = await fetch('http://localhost:8080/api/auth/captcha')
    const blob = await response.blob()
    // Try different header name cases
    let captchaIdHeader = response.headers.get('X-Captcha-ID')
    if (!captchaIdHeader) captchaIdHeader = response.headers.get('X-Captcha-Id')
    if (!captchaIdHeader) captchaIdHeader = response.headers.get('x-captcha-id')
    if (captchaIdHeader) {
      captchaId.value = captchaIdHeader
    }
    const url = URL.createObjectURL(blob)
    captchaUrl.value = url
  } catch (error) {
    console.error('Failed to load captcha:', error)
  }
}

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    toast.add({ severity: 'warn', summary: '请填写用户名和密码', life: 5000 })
    return
  }

  if (!form.value.captcha) {
    toast.add({ severity: 'warn', summary: '请输入验证码', life: 5000 })
    return
  }

  loading.value = true
  try {
    const res: any = await authApi.login(form.value.username, form.value.password, captchaId.value, form.value.captcha)
    if (res.code === 0) {
      userStore.setUser(res.data.user, res.data.token)
      toast.add({ severity: 'success', summary: '登录成功', life: 2000 })
      router.push('/')
    } else {
      toast.add({ severity: 'error', summary: res.message || '登录失败', life: 5000 })
      refreshCaptcha()
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.message || '登录失败', life: 5000 })
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
}

.login-left {
  flex: 1;
  background: #1e293b;
  position: relative;
  overflow: hidden;
}

.bg-shapes {
  position: absolute;
  inset: 0;
}

.shape {
  position: absolute;
}

.shape-1 {
  width: 70%;
  height: 80%;
  background: linear-gradient(135deg, #334155 0%, #1e293b 100%);
  border-radius: 0 0 100% 0;
  top: 0;
  left: 0;
}

.shape-2 {
  width: 50%;
  height: 60%;
  background: linear-gradient(135deg, #475569 0%, #334155 100%);
  border-radius: 0 0 80% 0;
  top: 20%;
  left: 25%;
  opacity: 0.8;
}

.shape-3 {
  width: 35%;
  height: 45%;
  background: #334155;
  border-radius: 50%;
  bottom: 5%;
  right: 15%;
  opacity: 0.5;
}

.shape-4 {
  width: 25%;
  height: 35%;
  background: #475569;
  border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
  top: 10%;
  right: 8%;
  opacity: 0.4;
}

.brand-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;

  .brand-icon {
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, #0ea5e9 0%, #0369a1 100%);
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 24px;
    color: #fff;
    box-shadow: 0 20px 40px -10px rgba(14, 165, 233, 0.3);
  }

  h1 {
    font-size: 36px;
    font-weight: 700;
    color: #f8fafc;
    margin-bottom: 12px;
    letter-spacing: -0.5px;
  }

  p {
    font-size: 16px;
    color: #94a3b8;
    letter-spacing: 2px;
    text-transform: uppercase;
  }
}

.login-right {
  flex: 0 0 500px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  padding: 40px;
}

.login-card {
  width: 100%;
  max-width: 360px;
}

.login-header {
  margin-bottom: 36px;

  h2 {
    font-size: 28px;
    font-weight: 700;
    color: #0f172a;
    margin-bottom: 8px;
  }

  span {
    font-size: 14px;
    color: #64748b;
  }
}

.form-group {
  margin-bottom: 20px;

  label {
    display: block;
    margin-bottom: 8px;
    font-size: 13px;
    font-weight: 500;
    color: #374151;
  }

  :deep(.p-inputtext) {
    width: 100%;
    height: 48px;
    padding: 0 16px;
    border-radius: 12px;
    border: 1.5px solid #e2e8f0;
    background: #fff;
    font-size: 15px;
    transition: all 0.2s ease;

    &:focus {
      border-color: #0369a1;
      box-shadow: 0 0 0 4px rgba(3, 105, 161, 0.1);
    }

    &::placeholder {
      color: #94a3b8;
    }
  }

  &.captcha-group {
    label {
      display: block;
      margin-bottom: 8px;
    }

    .captcha-row {
      display: flex;
      gap: 12px;

      :deep(.p-inputtext) {
        flex: 1;
        width: 100%;
      }
    }

    .captcha-img {
      width: 120px;
      height: 48px;
      border-radius: 12px;
      border: 1.5px solid #e2e8f0;
      background: #f8fafc;
      cursor: pointer;
      overflow: hidden;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.2s ease;
      flex-shrink: 0;

      &:hover {
        border-color: #0369a1;
      }

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .captcha-loading {
        font-size: 12px;
        color: #94a3b8;
      }
    }
  }
}

button[type="submit"] {
  width: 100%;
  margin-top: 24px;
  height: 48px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #0369a1 0%, #0284c7 100%);
  border: none;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 30px -8px rgba(3, 105, 161, 0.4);
  }
}

.login-tip {
  margin-top: 28px;
  color: #94a3b8;
  font-size: 13px;
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #e2e8f0;
}
</style>
