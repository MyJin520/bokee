<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import AppHeader from '@/components/AppHeader.vue'
import AppToast from '@/components/AppToast.vue'
import AppConfirm from '@/components/AppConfirm.vue'

const userStore = useUserStore()
onMounted(() => {
  userStore.restoreSession()
})
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="main-content">
      <router-view />
    </main>
    <footer class="app-footer">
      <div class="footer-inner">
        <span>© 2026 BOKEE · 写给愿意慢下来的人</span>
        <div class="footer-links">
          <a href="#rss">RSS</a>
          <a href="https://github.com" target="_blank" rel="noreferrer">GitHub</a>
          <a href="mailto:hello@bokee.local">联系我</a>
        </div>
      </div>
    </footer>
    <AppToast />
    <AppConfirm />
  </div>
</template>

<style>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  width: 100%;
  max-width: var(--max);
  flex: 1;
  margin: 0 auto;
  padding: 0 28px;
}

.app-footer {
  width: 100%;
  max-width: var(--max);
  margin: auto auto 0;
  padding: 26px 28px 34px;
  color: var(--muted);
  font-size: 11px;
}

.footer-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding-top: 23px;
  border-top: 1px solid var(--line);
}

.footer-links {
  display: flex;
  gap: 17px;
}

.footer-links a:hover {
  color: var(--ink);
}

@media (max-width: 880px) {
  .main-content {
    padding: 0 20px;
  }

  .app-footer {
    padding-left: 20px;
    padding-right: 20px;
  }
}

@media (max-width: 600px) {
  .main-content {
    padding: 0 16px;
  }

  .app-footer {
    padding-left: 16px;
    padding-right: 16px;
  }

  .footer-inner {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
