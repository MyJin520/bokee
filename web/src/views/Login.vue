<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const ui = useUiStore()

const userName = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  errorMsg.value = ''
  if (!userName.value.trim() || !password.value) {
    errorMsg.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  try {
    await userStore.login(userName.value.trim(), password.value)
    ui.toast('登录成功，欢迎回来')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error: unknown) {
    errorMsg.value = error instanceof Error ? error.message : '登录失败，请稍后重试'
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
          <h2>慢一点，<br />写得久一点。</h2>
          <p>登录后可以收藏文章、管理自己的文章，并继续你的写作。</p>
        </div>
        <span class="auth-brand-foot">SINCE 2026 · A QUIET CORNER</span>
      </section>

      <section class="auth-panel">
        <div class="auth-box">
          <span class="eyebrow">Welcome back</span>
          <h1>登录账户</h1>
          <p class="auth-sub">使用用户名和密码进入你的写作后台。</p>

          <form @submit.prevent="handleLogin">
            <label class="field-label"><span class="required">*</span> 用户名</label>
            <input
              v-model="userName"
              class="field-input"
              type="text"
              placeholder="请输入用户名"
              autocomplete="username"
            />

            <label class="field-label"><span class="required">*</span> 密码</label>
            <input
              v-model="password"
              class="field-input"
              type="password"
              placeholder="请输入密码"
              autocomplete="current-password"
            />

            <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>

            <button class="btn btn-block" type="submit" style="margin-top: 26px" :disabled="loading">
              {{ loading ? '登录中…' : '登 录' }}
            </button>
          </form>

          <p class="auth-switch">
            还没有账户？
            <router-link to="/register">立即注册 →</router-link>
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
  min-height: 560px;
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
  padding: 40px;
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
  margin: 0 0 22px;
  color: var(--muted);
  font-size: 12px;
}

.auth-switch {
  margin: 22px 0 0;
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
    min-height: 210px;
  }
}
</style>
