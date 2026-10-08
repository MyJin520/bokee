<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { listUserActions } from '@/api/userAction'
import { useReactions } from '@/composables/useReactions'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import type { UserActionListItem } from '@/types'
import { formatNumber } from '@/utils/format'

const userStore = useUserStore()
const ui = useUiStore()
const { toggleFollow } = useReactions()

const items = ref<UserActionListItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(true)
const errorMessage = ref('')

const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

async function load() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await listUserActions({
      actionType: 'follow',
      targetType: 'author',
      page: page.value,
      pageSize,
    })
    items.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '关注加载失败，请稍后重试。'
    items.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function changePage(next: number) {
  if (next < 1 || next > pages.value || next === page.value) return
  page.value = next
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function unfollow(item: UserActionListItem) {
  if (!item.author) return
  try {
    await toggleFollow(item.author.id)
    ui.toast(`已取消关注 ${item.author.name}`)
    if (items.value.length === 1 && page.value > 1) page.value -= 1
    load()
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '操作失败，请稍后重试')
  }
}

onMounted(async () => {
  // 页面直接刷新进入时，等待根组件恢复登录态后再加载
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  load()
})
</script>

<template>
  <div class="follows-page">
    <header class="page-header">
      <span class="eyebrow">My follows</span>
      <h1>我的关注</h1>
      <p>共 {{ formatNumber(total) }} 位关注的作者，关注后更容易找到 TA 的新文章。</p>
    </header>

    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取关注…</p>
    </div>
    <div v-else-if="errorMessage" class="state-block">
      <p class="state-title">加载遇到问题</p>
      <p>{{ errorMessage }}</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="load">重新加载 →</button>
    </div>
    <template v-else-if="items.length">
      <ul class="follows-list">
        <li v-for="item in items" :key="item.actionId">
          <router-link
            class="follows-main"
            :to="`/articles?user=${encodeURIComponent(item.author?.name || '')}`"
          >
            <span class="follows-avatar">{{ (item.author?.name || '?').slice(0, 1) }}</span>
            <span class="follows-name">{{ item.author?.name || '未命名作者' }}</span>
          </router-link>
          <button
            class="follows-remove"
            type="button"
            @click="unfollow(item)"
            :aria-label="`取消关注${item.author?.name || ''}`"
          >
            取消关注
          </button>
        </li>
      </ul>
      <div class="pagination">
        <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">← 上一页</button>
        <span class="page-info">{{ page }} / {{ pages }}</span>
        <button class="page-btn" type="button" :disabled="page >= pages" @click="changePage(page + 1)">下一页 →</button>
      </div>
    </template>
    <div v-else class="state-block empty-state">
      <p class="state-title">还没有关注</p>
      <p>读到喜欢的文章，关注它的作者，新文章会第一时间看到。</p>
      <router-link class="btn btn-coral" to="/articles" style="margin-top: 18px">去文章列表看看 →</router-link>
    </div>
  </div>
</template>

<style scoped>
.follows-page {
  padding: 48px 0 80px;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 14px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(34px, 4.6vw, 54px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.page-header p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.follows-list {
  list-style: none;
  margin: 0;
  padding: 0;
  border-top: 1px solid var(--line);
}

.follows-list li {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
}

.follows-main {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.follows-avatar {
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: var(--coral);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 16px;
}

.follows-name {
  font-size: 14px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.follows-main:hover .follows-name {
  color: var(--coral-dark);
}

.follows-remove {
  flex-shrink: 0;
  padding: 7px 12px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
}

.follows-remove:hover {
  border-color: var(--coral-dark);
  color: var(--coral-dark);
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 10px;
}

@media (max-width: 600px) {
  .follows-list li {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
