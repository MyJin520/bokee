<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const ui = useUiStore()
const menuOpen = ref(false)

const displayName = computed(
  () => userStore.userInfo?.name || localStorage.getItem('userName') || '我的账户',
)
const avatarUrl = computed(() => userStore.userInfo?.avatar || '')

const navItems = [
  { label: '首页', to: '/', match: (path: string) => path === '/' },
  { label: '文章', to: '/articles', match: (path: string) => path.startsWith('/article') },
  { label: '专题', to: '/#topics', match: () => false },
  { label: '关于', to: '/#about', match: () => false },
]

function isActive(item: (typeof navItems)[number]) {
  return item.match(route.path)
}

function go(path: string) {
  menuOpen.value = false
  if (path.includes('#')) {
    const [routePath, hash] = path.split('#')
    if (route.path === routePath || (routePath === '' && route.path === '/')) {
      router.push({ path: routePath || '/', hash: `#${hash}` })
      return
    }
    router.push({ path: routePath || '/', hash: `#${hash}` })
    return
  }
  router.push(path)
}

/** 头部搜索：首页内滚动聚焦搜索框，其他页面进入文章归档页并聚焦 */
function focusSearch() {
  menuOpen.value = false
  if (route.path === '/') {
    document.getElementById('latest')?.scrollIntoView({ behavior: 'smooth' })
    ui.focusHomeSearch()
  } else {
    router.push({ path: '/articles', query: { focus: '1' } })
  }
}

async function logout() {
  menuOpen.value = false
  const ok = await ui.confirm({
    title: '退出登录',
    message: '确定要退出当前账号吗？',
    confirmText: '退出',
    danger: true,
  })
  if (!ok) return
  await userStore.doLogout()
  ui.toast('已退出登录')
  go('/')
}
</script>

<template>
  <header class="site-header">
    <div class="header-inner">
      <button class="brand" type="button" aria-label="BOKEE 首页" @click="go('/')">
        <span class="brand-mark">B</span><span>BOKEE</span>
      </button>

      <nav class="desktop-nav" aria-label="主导航">
        <button
          v-for="item in navItems"
          :key="item.label"
          :class="['nav-link', { active: isActive(item) }]"
          type="button"
          @click="go(item.to)"
        >
          {{ item.label }}
        </button>
      </nav>

      <div class="header-actions">
        <button class="header-icon" type="button" aria-label="搜索文章" title="搜索" @click="focusSearch">
          ⌕
        </button>
        <template v-if="userStore.isLoggedIn">
          <button class="write-button" type="button" @click="go('/articles/create')">＋ 写文章</button>
          <button class="account" type="button" @click="go('/profile')">
            <span class="account-avatar">
              <img v-if="avatarUrl" :src="avatarUrl" :alt="displayName" />
              <template v-else>{{ displayName.slice(0, 1) }}</template>
            </span>
            <span class="account-name">{{ displayName }}</span>
          </button>
          <button class="logout-inline" type="button" title="退出登录" @click="logout">退出</button>
        </template>
        <button v-else class="login-button" type="button" @click="go('/login')">登录 / 注册</button>
        <button
          class="menu-button"
          type="button"
          aria-label="打开菜单"
          @click="menuOpen = !menuOpen"
        >
          {{ menuOpen ? '×' : '☰' }}
        </button>
      </div>
    </div>

    <nav v-if="menuOpen" class="mobile-nav" aria-label="移动端导航">
      <button
        v-for="item in navItems"
        :key="item.label"
        :class="{ active: isActive(item) }"
        type="button"
        @click="go(item.to)"
      >
        {{ item.label }}
      </button>
      <template v-if="userStore.isLoggedIn">
        <button type="button" @click="go('/articles/create')">写文章</button>
        <button type="button" @click="go('/my-articles')">我的文章</button>
        <button type="button" @click="go('/profile')">个人中心</button>
        <button type="button" @click="logout">退出登录</button>
      </template>
      <template v-else>
        <button type="button" @click="go('/login')">登录</button>
        <button type="button" @click="go('/register')">注册</button>
      </template>
    </nav>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 20;
  border-bottom: 1px solid var(--line);
  background: rgba(251, 250, 247, 0.94);
  backdrop-filter: blur(12px);
}

.header-inner {
  max-width: var(--max);
  min-height: 76px;
  margin: 0 auto;
  padding: 0 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0;
  color: var(--ink);
  background: transparent;
  font-size: 19px;
  font-weight: 800;
  letter-spacing: 0.04em;
}

.brand-mark {
  width: 31px;
  height: 31px;
  display: grid;
  place-items: center;
  background: var(--coral);
  color: #fff;
  font-size: 17px;
  transform: rotate(-7deg);
}

.desktop-nav {
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-link {
  position: relative;
  padding: 10px 13px;
  color: var(--muted);
  background: transparent;
  font-size: 13px;
  font-weight: 650;
}

.nav-link:hover,
.nav-link.active {
  color: var(--ink);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  left: 13px;
  right: 13px;
  bottom: 4px;
  height: 2px;
  background: var(--coral);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon,
.menu-button {
  width: 37px;
  height: 37px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--ink);
  font-size: 16px;
}

.header-icon:hover,
.menu-button:hover {
  border-color: var(--ink);
}

.write-button,
.login-button {
  padding: 10px 15px;
  background: var(--ink);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
}

.write-button {
  background: var(--coral);
}

.write-button:hover {
  background: var(--coral-dark);
}

.login-button:hover {
  background: var(--coral);
}

.account {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 4px 4px 4px;
  background: transparent;
}

.account-avatar {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--ink);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
}

.account-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.account-name {
  max-width: 110px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  font-weight: 700;
}

.account:hover .account-name {
  color: var(--coral-dark);
}

.logout-inline {
  padding: 5px 8px;
  background: transparent;
  color: var(--muted);
  font-size: 11px;
}

.logout-inline:hover {
  color: var(--coral-dark);
}

.menu-button {
  display: none;
}

.mobile-nav {
  display: none;
  border-top: 1px solid var(--line);
  background: var(--surface);
  padding: 8px 20px 13px;
  grid-template-columns: repeat(4, 1fr);
  gap: 4px;
}

.mobile-nav button {
  padding: 10px 4px;
  background: transparent;
  color: var(--muted);
  font-size: 12px;
}

.mobile-nav button.active,
.mobile-nav button:hover {
  color: var(--ink);
}

@media (max-width: 880px) {
  .header-inner {
    padding: 0 20px;
  }

  .desktop-nav {
    display: none;
  }

  .menu-button {
    display: grid;
  }

  .mobile-nav {
    display: grid;
  }
}

@media (max-width: 600px) {
  .header-inner {
    min-height: 67px;
    padding: 0 16px;
  }

  .write-button,
  .logout-inline,
  .account-name {
    display: none;
  }

  .login-button {
    padding: 9px 11px;
  }

  .mobile-nav {
    grid-template-columns: repeat(3, 1fr);
    padding-left: 16px;
    padding-right: 16px;
  }
}
</style>
