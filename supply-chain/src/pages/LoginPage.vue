<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="login-logo">SC</div>
        <h1>供应链管理系统</h1>
        <p>Supply Chain Management System</p>
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  background: var(--panel);
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
  
  .login-logo {
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 16px;
    color: #fff;
    font-size: 24px;
    font-weight: 700;
  }
  
  h1 {
    font-size: 24px;
    margin-bottom: 4px;
  }
  
  p {
    color: var(--text-muted);
    font-size: 14px;
  }
}

.form-group {
  margin-bottom: 20px;
  
  label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: var(--text-muted);
  }
  
  input {
    width: 100%;
  }
}

button[type="submit"] {
  width: 100%;
  margin-top: 10px;
}

.login-tip {
  text-align: center;
  margin-top: 20px;
  color: var(--text-muted);
  font-size: 12px;
}
</style>