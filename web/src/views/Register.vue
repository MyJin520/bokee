<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/user'

const router = useRouter()

const form = ref({
  userName: '',
  password: '',
  confirmPassword: '',
  phone: '',
  email: '',
})
const loading = ref(false)
const errorMsg = ref('')

async function handleRegister() {
  errorMsg.value = ''

  if (!form.value.userName || !form.value.password) {
    errorMsg.value = '用户名和密码为必填项'
    return
  }

  if (form.value.password.length < 6) {
    errorMsg.value = '密码至少6位'
    return
  }

  if (form.value.password !== form.value.confirmPassword) {
    errorMsg.value = '两次密码输入不一致'
    return
  }

  loading.value = true
  try {
    await register({
      userName: form.value.userName,
      password: form.value.password,
      phone: form.value.phone || undefined,
      email: form.value.email || undefined,
    })
    alert('注册成功，请登录')
    router.push('/login')
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="auth-card">
      <h1 class="auth-title">注册</h1>
      <p class="auth-subtitle">创建你的账号，开始写博客</p>

      <form class="auth-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">用户名 <span class="required">*</span></label>
          <input
            v-model="form.userName"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label class="form-label">密码 <span class="required">*</span></label>
          <input
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="至少6位密码"
            autocomplete="new-password"
          />
        </div>

        <div class="form-group">
          <label class="form-label">确认密码 <span class="required">*</span></label>
          <input
            v-model="form.confirmPassword"
            type="password"
            class="form-input"
            placeholder="再次输入密码"
            autocomplete="new-password"
          />
        </div>

        <div class="form-group">
          <label class="form-label">手机号</label>
          <input
            v-model="form.phone"
            type="text"
            class="form-input"
            placeholder="选填"
            autocomplete="tel"
          />
        </div>

        <div class="form-group">
          <label class="form-label">邮箱</label>
          <input
            v-model="form.email"
            type="email"
            class="form-input"
            placeholder="选填"
            autocomplete="email"
          />
        </div>

        <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <p class="auth-switch">
        已有账号？
        <router-link to="/login">立即登录</router-link>
      </p>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  display: flex;
  justify-content: center;
  padding-top: 40px;
}

.auth-card {
  width: 100%;
  max-width: 420px;
  background: var(--color-white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: 40px 36px;
  border: 1px solid var(--color-border);
}

.auth-title {
  font-size: 1.5rem;
  font-weight: 700;
  text-align: center;
  color: var(--color-text);
}

.auth-subtitle {
  text-align: center;
  color: var(--color-text-light);
  font-size: 0.9375rem;
  margin-top: 8px;
  margin-bottom: 28px;
}

.required {
  color: var(--color-danger);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text);
}

.form-input {
  padding: 10px 14px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: 0.9375rem;
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: var(--color-primary);
}

.form-error {
  color: var(--color-danger);
  font-size: 0.875rem;
  text-align: center;
}

.submit-btn {
  padding: 12px;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-sm);
  font-size: 1rem;
  font-weight: 600;
  transition: background 0.2s;
  margin-top: 4px;
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-switch {
  text-align: center;
  margin-top: 24px;
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.auth-switch a {
  font-weight: 600;
}
</style>
