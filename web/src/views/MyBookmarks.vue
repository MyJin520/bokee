<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { listUserActions } from '@/api/userAction'
import { useReactions } from '@/composables/useReactions'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import type { UserActionListItem } from '@/types'
import { formatNumber } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()
const ui = useUiStore()
const { toggleSaved } = useReactions()

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
      actionType: 'bookmark',
      targetType: 'article',
      page: page.value,
      pageSize,
    })
    items.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '收藏加载失败，请稍后重试。'
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

async function removeBookmark(item: UserActionListItem) {
  if (!item.article) return
  try {
    await toggleSaved(item.article.id)
    ui.toast('已取消收藏')
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
  <div class="saved-page">
    <header class="page-header">
      <span class="eyebrow">My bookmarks</span>
      <h1>我的收藏</h1>
      <p>共 {{ formatNumber(total) }} 篇收藏的文章，随时回来继续读。</p>
    </header>

    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取收藏…</p>
    </div>
    <div v-else-if="errorMessage" class="state-block">
      <p class="state-title">加载遇到问题</p>
      <p>{{ errorMessage }}</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="load">重新加载 →</button>
    </div>
    <template v-else-if="items.length">
      <ul class="saved-list">
        <li v-for="item in items" :key="item.actionId">
          <button class="saved-main" type="button" @click="router.push(`/article/${item.article?.id}`)">
            <span class="saved-title">{{ item.article?.title || '未命名文章' }}</span>
            <span class="saved-date">ID {{ item.article?.id }}</span>
          </button>
          <button
            class="saved-remove"
            type="button"
            @click="removeBookmark(item)"
            :aria-label="`取消收藏《${item.article?.title || ''}》`"
          >
            取消收藏
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
      <p class="state-title">还没有收藏</p>
      <p>看到想收藏的文章，点一下收藏按钮，就会出现在这里。</p>
      <router-link class="btn btn-coral" to="/articles" style="margin-top: 18px">去文章列表看看 →</router-link>
    </div>
  </div>
</template>

<style scoped>
.saved-page {
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

.saved-list {
  list-style: none;
  margin: 0;
  padding: 0;
  border-top: 1px solid var(--line);
}

.saved-list li {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 16px 0;
  border-bottom: 1px solid var(--line);
}

.saved-main {
  flex: 1;
  min-width: 0;
  padding: 0;
  text-align: left;
  background: transparent;
}

.saved-title {
  display: block;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 18px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.saved-main:hover .saved-title {
  color: var(--coral-dark);
}

.saved-date {
  display: block;
  margin-top: 5px;
  color: var(--muted);
  font-size: 10px;
}

.saved-remove {
  flex-shrink: 0;
  padding: 7px 12px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
}

.saved-remove:hover {
  border-color: var(--coral-dark);
  color: var(--coral-dark);
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 10px;
}

@media (max-width: 600px) {
  .saved-list li {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
