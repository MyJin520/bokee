<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { deleteArticle, getUserArticleList } from '@/api/article'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import type { ArticleListItem } from '@/types'
import ArticleRow from '@/components/ArticleRow.vue'
import { formatNumber } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()
const ui = useUiStore()

const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 8
const loading = ref(true)
const errorMessage = ref('')

const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

async function loadArticles() {
  if (!userStore.userInfo?.id) return
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await getUserArticleList(userStore.userInfo.id, page.value, pageSize)
    articles.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '文章加载失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

function changePage(next: number) {
  if (next < 1 || next > pages.value || next === page.value) return
  page.value = next
  loadArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function editArticle(id: number) {
  router.push(`/articles/${id}/edit`)
}

async function removeArticle(id: number, title: string) {
  const ok = await ui.confirm({
    title: '删除文章',
    message: `确定要删除《${title}》吗？此操作不可撤销。`,
    confirmText: '确认删除',
    danger: true,
  })
  if (!ok) return
  try {
    await deleteArticle(id)
    ui.toast('文章已删除')
    if (articles.value.length === 1 && page.value > 1) page.value -= 1
    loadArticles()
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '删除失败，请稍后重试')
  }
}

onMounted(async () => {
  // 页面直接刷新进入时，等待根组件恢复登录态后再加载
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  loadArticles()
})
</script>

<template>
  <div class="my-page">
    <header class="my-header">
      <div>
        <span class="eyebrow">My articles</span>
        <h1>我的文章</h1>
        <p>共 {{ formatNumber(total) }} 篇，在这里继续写作或整理旧文。</p>
      </div>
      <router-link class="btn btn-coral" to="/articles/create">＋ 写新文章</router-link>
    </header>

    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取你的文章…</p>
    </div>
    <div v-else-if="errorMessage" class="state-block">
      <p class="state-title">加载遇到问题</p>
      <p>{{ errorMessage }}</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="loadArticles">重新加载 →</button>
    </div>
    <template v-else-if="articles.length">
      <div class="my-list">
        <ArticleRow v-for="article in articles" :key="article.id" :article="article">
          <template #actions>
            <button class="row-manage" type="button" @click.stop="editArticle(article.id)">编辑</button>
            <button class="row-manage danger" type="button" @click.stop="removeArticle(article.id, article.title)">
              删除
            </button>
          </template>
        </ArticleRow>
      </div>
      <div class="pagination">
        <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">← 上一页</button>
        <span class="page-info">{{ page }} / {{ pages }}</span>
        <button class="page-btn" type="button" :disabled="page >= pages" @click="changePage(page + 1)">下一页 →</button>
      </div>
    </template>
    <div v-else class="state-block empty-state">
      <p class="state-title">还没有文章</p>
      <p>写下第一篇，把今天值得记录的事留下来。</p>
      <router-link class="btn btn-coral" to="/articles/create" style="margin-top: 18px">开始写作 →</router-link>
    </div>
  </div>
</template>

<style scoped>
.my-page {
  padding: 48px 0 80px;
}

.my-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 30px;
}

.my-header h1 {
  margin: 14px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(34px, 4.6vw, 54px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.my-header p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.my-list {
  border-top: 1px solid var(--line);
}

.row-manage {
  padding: 7px 12px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--ink);
  font-size: 11px;
  font-weight: 700;
}

.row-manage:hover {
  border-color: var(--ink);
}

.row-manage.danger {
  color: var(--color-danger);
}

.row-manage.danger:hover {
  border-color: var(--color-danger);
  background: var(--color-danger);
  color: #fff;
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 10px;
}

@media (max-width: 600px) {
  .my-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
