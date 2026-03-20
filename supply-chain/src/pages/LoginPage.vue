<template>
  <div class="login-page">
    <div class="login-bg">
      <div class="grid-lines"></div>
      <div class="glow glow-1"></div>
      <div class="glow glow-2"></div>
      <div class="glow glow-3"></div>
    </div>
    <div class="login-card">
      <div class="login-header">
        <div class="login-logo">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h1>供应链管理系统</h1>
        <p>Supply Chain Management</p>
      </div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <div class="input-wrapper">
            <InputText v-model="form.username" placeholder="用户名" />
          </div>
        </div>
        <div class="form-group">
          <div class="input-wrapper">
            <InputText v-model="form.password" type="password" placeholder="密码" />
          </div>
        </div>
        <Button type="submit" label="登 录" :loading="loading" />
      </form>

      <p class="login-tip">默认账号: admin / admin123</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { authApi } from '@/api'

const router = useRouter()
const toast = useToast()

const loading = ref(false)
const form = ref({
  username: 'admin',
  password: 'admin123'
})

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    toast.add({ severity: 'warn', summary: '请填写用户名和密码', life: 5000 })
    return
  }

  loading.value = true
  try {
    const res: any = await authApi.login(form.value.username, form.value.password)
    if (res.code === 0) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      toast.add({ severity: 'success', summary: '登录成功', life: 2000 })
      router.push('/')
    } else {
      toast.add({ severity: 'error', summary: res.message || '登录失败', life: 5000 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: error.message || '登录失败', life: 5000 })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #09090b;
  position: relative;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.grid-lines {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(to right, rgba(255,255,255,0.02) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255,255,255,0.02) 1px, transparent 1px);
  background-size: 60px 60px;
  mask-image: radial-gradient(ellipse 80% 60% at 50% 50%, black 40%, transparent 100%);
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  animation: pulse 8s ease-in-out infinite;
}

.glow-1 {
  width: 400px;
  height: 400px;
  background: #22d3ee;
  top: -100px;
  left: 20%;
  opacity: 0.3;
  animation-delay: 0s;
}

.glow-2 {
  width: 500px;
  height: 500px;
  background: #a855f7;
  bottom: -150px;
  right: 10%;
  opacity: 0.25;
  animation-delay: -3s;
}

.glow-3 {
  width: 300px;
  height: 300px;
  background: #f472b6;
  top: 40%;
  right: 30%;
  opacity: 0.2;
  animation-delay: -5s;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.3;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.4;
  }
}

.login-card {
  width: 380px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(40px);
  border-radius: 20px;
  padding: 44px 36px;
  box-shadow:
    0 0 0 1px rgba(255,255,255,0.1),
    0 20px 50px -15px rgba(0, 0, 0, 0.5),
    0 0 100px -20px rgba(34, 211, 238, 0.15);
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 36px;

  .login-logo {
    width: 56px;
    height: 56px;
    background: linear-gradient(135deg, #18181b 0%, #27272a 100%);
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 18px;
    color: #fafafa;
    box-shadow: 0 8px 20px -8px rgba(0,0,0,0.3);
  }

  h1 {
    font-size: 22px;
    font-weight: 600;
    color: #18181b;
    margin-bottom: 6px;
    letter-spacing: -0.3px;
  }

  p {
    color: #71717a;
    font-size: 12px;
    letter-spacing: 2px;
    text-transform: uppercase;
    font-weight: 500;
  }
}

.form-group {
  margin-bottom: 18px;

  .input-wrapper {
    :deep(.p-inputtext) {
      width: 100%;
      height: 48px;
      padding: 0 16px;
      border-radius: 12px;
      border: 1.5px solid #e4e4e7;
      background: #fafafa;
      font-size: 14px;
      transition: all 0.25s ease;

      &:focus {
        border-color: #18181b;
        background: #fff;
        box-shadow: 0 0 0 4px rgba(24, 24, 27, 0.08);
      }

      &::placeholder {
        color: #a1a1aa;
      }
    }
  }
}

button[type="submit"] {
  width: 100%;
  margin-top: 20px;
  height: 48px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 3px;
  background: linear-gradient(135deg, #18181b 0%, #27272a 100%);
  border: none;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.35);
  }

  &:active {
    transform: translateY(0);
  }
}

.login-tip {
  text-align: center;
  margin-top: 24px;
  color: #a1a1aa;
  font-size: 12px;
  padding-top: 20px;
  border-top: 1px solid #f4f4f5;
}
</style>