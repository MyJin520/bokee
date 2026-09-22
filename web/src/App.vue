<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import AppHeader from '@/components/AppHeader.vue'

const userStore = useUserStore()

onMounted(() => {
  // 页面刷新时尝试恢复登录态
  userStore.restoreSession()
})
</script>

<template>
  <div class="app-container">
    <AppHeader />
    <main class="main-content">
      <router-view />
    </main>
    <footer class="app-footer">
      <p>© 2026 My Blog. All rights reserved.</p>
    </footer>
  </div>
</template>

<style>
/* ===== 全局重置 ===== */
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --color-primary: #4a90d9;
  --color-primary-hover: #357abd;
  --color-bg: #f5f7fa;
  --color-white: #ffffff;
  --color-text: #2c3e50;
  --color-text-light: #6b7a8d;
  --color-border: #e8ecf1;
  --color-danger: #e74c3c;
  --color-success: #27ae60;
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.08);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.12);
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;

  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-bg);
  -webkit-font-smoothing: antialiased;
}

a {
  color: var(--color-primary);
  text-decoration: none;
  transition: color 0.2s;
}

a:hover {
  color: var(--color-primary-hover);
}

img {
  max-width: 100%;
  height: auto;
}

button {
  cursor: pointer;
  border: none;
  font-family: inherit;
  font-size: inherit;
}
</style>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
  padding: 28px 20px;
}

.app-footer {
  text-align: center;
  padding: 24px 20px;
  color: var(--color-text-light);
  font-size: 0.875rem;
  border-top: 1px solid var(--color-border);
  margin-top: 40px;
}
</style>
