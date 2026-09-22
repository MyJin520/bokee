<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const displayName = computed(() => {
  return userStore.userName || localStorage.getItem('userName') || ''
})

function goHome() {
  router.push('/')
}

function goToMyArticles() {
  router.push('/my-articles')
}

function goToCreateArticle() {
  router.push('/articles/create')
}

function goToLogin() {
  router.push('/login')
}

async function handleLogout() {
  await userStore.doLogout()
  router.push('/')
}
</script>

<template>
  <header class="header">
    <div class="header-inner">
      <div class="logo" @click="goHome">
        <span class="logo-icon">📝</span>
        <span class="logo-text">My Blog</span>
      </div>

      <nav class="nav">
        <button class="nav-link" @click="goHome">首页</button>

        <template v-if="userStore.isLoggedIn">
          <button class="nav-link" @click="goToMyArticles">我的文章</button>
          <button class="nav-link nav-link-primary" @click="goToCreateArticle">写文章</button>
          <div class="user-menu">
            <span class="user-name">{{ displayName }}</span>
            <button class="nav-link nav-link-logout" @click="handleLogout">退出</button>
          </div>
        </template>

        <template v-else>
          <button class="nav-link nav-link-login" @click="goToLogin">登录</button>
        </template>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.header {
  background: var(--color-white);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-sm);
}

.header-inner {
  max-width: 960px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.logo-icon {
  font-size: 1.5rem;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text);
}

.nav {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-link {
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-light);
  font-size: 0.9375rem;
  font-weight: 500;
  transition: all 0.2s;
}

.nav-link:hover {
  color: var(--color-primary);
  background: #f0f4f9;
}

.nav-link-primary {
  background: var(--color-primary);
  color: #fff;
}

.nav-link-primary:hover {
  background: var(--color-primary-hover);
  color: #fff;
}

.nav-link-login {
  color: var(--color-primary);
  font-weight: 600;
}

.nav-link-logout {
  color: var(--color-text-light);
  font-size: 0.875rem;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
  padding-left: 12px;
  border-left: 1px solid var(--color-border);
}

.user-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text);
  padding: 0 8px;
}
</style>
