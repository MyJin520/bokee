<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/user'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const ui = useUiStore()

const form = reactive({
  userName: '',
  password: '',
  confirmPassword: '',
  email: '',
  phone: '',
})
const loading = ref(false)
const errors = reactive<Record<string, string>>({})

function validate(): boolean {
  Object.keys(errors).forEach((key) => delete errors[key])

  if (!form.userName.trim()) {
    errors.userName = '请输入用户名'
  } else if (form.userName.trim().length < 2 || form.userName.trim().length > 10) {
    errors.userName = '用户名长度需为 2-10 个字符'
  }
  if (!form.password) {
    errors.password = '请输入密码'
  } else if (form.password.length < 6) {
    errors.password = '密码长度至少 6 位'
  }
  if (form.confirmPassword !== form.password) {
    errors.confirmPassword = '两次输入的密码不一致'
  }
  if (!form.email.trim() && !form.phone.trim()) {
    errors.email = '邮箱与手机号至少填写一项'
  }
  if (form.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) {
    errors.email = '邮箱格式不正确'
  }
  if (form.phone.trim() && !/^1\d{10}$/.test(form.phone.trim())) {
    errors.phone = '请输入正确的 11 位手机号'
  }
  return Object.keys(errors).length === 0
}

async function handleRegister() {
  if (!validate()) return
  loading.value = true
  try {
    await register({
      name: form.userName.trim(),
      password: form.password,
      email: form.email.trim() || undefined,
      phone: form.phone.trim() || undefined,
    })
    ui.toast('注册成功，请登录')
    router.push('/login')
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-layout">
      <section class="auth-brand">
        <router-link class="auth-brand-mark" to="/">BOKEE</router-link>
        <div class="auth-quote">
          <h2>从今天起，<br />开始记录。</h2>
          <p>注册后即可发表文章、上传封面，并拥有自己的写作主页。</p>
        </div>
        <span class="auth-brand-foot">JOIN THE READERS · 2026</span>
      </section>

      <section class="auth-panel">
        <div class="auth-box">
          <span class="eyebrow">Create account</span>
          <h1>注册账户</h1>
          <p class="auth-sub">填写以下信息，几秒钟完成注册。</p>

          <form @submit.prevent="handleRegister">
            <label class="field-label"><span class="required">*</span> 用户名</label>
            <input v-model="form.userName" class="field-input" type="text" placeholder="2-10 个字符" autocomplete="username" />
            <p v-if="errors.userName" class="form-error">{{ errors.userName }}</p>

            <label class="field-label"><span class="required">*</span> 密码</label>
            <input v-model="form.password" class="field-input" type="password" placeholder="至少 6 位" autocomplete="new-password" />
            <p v-if="errors.password" class="form-error">{{ errors.password }}</p>

            <label class="field-label"><span class="required">*</span> 确认密码</label>
            <input v-model="form.confirmPassword" class="field-input" type="password" placeholder="再次输入密码" autocomplete="new-password" />
            <p v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</p>

            <label class="field-label">邮箱<span class="required">（与手机号至少填一项）</span></label>
            <input v-model="form.email" class="field-input" type="email" placeholder="you@example.com" autocomplete="email" />
            <p v-if="errors.email" class="form-error">{{ errors.email }}</p>

            <label class="field-label">手机号</label>
            <input v-model="form.phone" class="field-input" type="tel" placeholder="11 位手机号" autocomplete="tel" />
            <p v-if="errors.phone" class="form-error">{{ errors.phone }}</p>

            <button class="btn btn-coral btn-block" type="submit" style="margin-top: 26px" :disabled="loading">
              {{ loading ? '注册中…' : '注 册' }}
            </button>
          </form>

          <p class="auth-switch">
            已有账户？
            <router-link to="/login">返回登录 →</router-link>
          </p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: calc(100vh - 76px);
  padding: 44px 0;
  display: grid;
  place-items: center;
}

.auth-layout {
  width: 100%;
  max-width: 980px;
  display: grid;
  grid-template-columns: 0.9fr 1.1fr;
  border: 1px solid var(--line);
  background: var(--surface);
  box-shadow: var(--shadow);
}

.auth-brand {
  position: relative;
  padding: 34px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  color: #fff;
  background: var(--ink);
}

.auth-brand-mark {
  display: inline-flex;
  align-items: center;
  gap: 9px;
  width: max-content;
  color: #fff;
  font-size: 15px;
  font-weight: 800;
  letter-spacing: 0.06em;
}

.auth-brand-mark::before {
  content: 'B';
  width: 27px;
  height: 27px;
  display: grid;
  place-items: center;
  background: var(--coral);
  font-size: 14px;
  transform: rotate(-7deg);
}

.auth-quote h2 {
  margin: 0 0 14px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 38px;
  font-weight: 500;
  line-height: 1.1;
}

.auth-quote p {
  max-width: 280px;
  margin: 0;
  color: #b7c3bb;
  font-size: 12px;
  line-height: 1.8;
}

.auth-brand-foot {
  color: #81908a;
  font-size: 10px;
  letter-spacing: 0.14em;
}

.auth-panel {
  display: grid;
  place-items: center;
  padding: 34px 40px;
}

.auth-box {
  width: 100%;
  max-width: 360px;
}

.auth-box h1 {
  margin: 16px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 34px;
  font-weight: 500;
}

.auth-sub {
  margin: 0 0 18px;
  color: var(--muted);
  font-size: 12px;
}

.auth-switch {
  margin: 20px 0 0;
  text-align: center;
  color: var(--muted);
  font-size: 12px;
}

.auth-switch a {
  color: var(--coral-dark);
  font-weight: 700;
}

@media (max-width: 820px) {
  .auth-page {
    padding: 24px 0;
  }

  .auth-layout {
    grid-template-columns: 1fr;
  }

  .auth-brand {
    min-height: 200px;
  }
}
</style>
