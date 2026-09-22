<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getUserArticleList, deleteArticle } from '@/api/article'
import type { ArticleListItem } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(true)

async function loadArticles() {
  loading.value = true
  try {
    const userId = userStore.userInfo?.id
    if (!userId) {
      articles.value = []
      return
    }
    const res = await getUserArticleList(userId, currentPage.value, pageSize.value)
    articles.value = res.list
    total.value = res.total
  } catch (e: unknown) {
    console.error('加载文章列表失败', e)
    articles.value = []
  } finally {
    loading.value = false
  }
}

const totalPages = () => Math.ceil(total.value / pageSize.value)

function goToPage(page: number) {
  if (page < 1 || page > totalPages()) return
  currentPage.value = page
  loadArticles()
}

function goToEdit(id: number) {
  router.push(`/articles/${id}/edit`)
}

async function handleDelete(id: number) {
  if (!confirm('确定要删除这篇文章吗？此操作不可撤销。')) return

  try {
    await deleteArticle(id)
    alert('删除成功')
    loadArticles()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<template>
  <div class="my-articles-page">
    <div class="page-heading">
      <div>
        <h1>我的文章</h1>
        <p class="page-subtitle">共 {{ total }} 篇文章</p>
      </div>
      <button class="create-btn" @click="router.push('/articles/create')">写新文章</button>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 文章列表 -->
    <div v-else-if="articles.length > 0" class="article-list">
      <div v-for="article in articles" :key="article.id" class="article-item">
        <div class="item-body" @click="router.push(`/article/${article.id}`)">
          <h3 class="item-title">
            <span v-if="article.isTop" class="top-badge">置顶</span>
            {{ article.title }}
          </h3>
          <p class="item-summary">{{ article.summary || '暂无摘要' }}</p>
          <div class="item-meta">
            <span>{{ new Date(article.createdAt).toLocaleDateString('zh-CN') }}</span>
            <span class="meta-dot">·</span>
            <span>👁 {{ article.viewCount }}</span>
            <span class="meta-dot">·</span>
            <span>👍 {{ article.likeCount }}</span>
          </div>
        </div>
        <div class="item-actions">
          <button class="action-edit" @click="goToEdit(article.id)">编辑</button>
          <button class="action-delete" @click="handleDelete(article.id)">删除</button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">📝</div>
      <p class="empty-title">还没有文章</p>
      <p class="empty-desc">开始写你的第一篇文章吧！</p>
      <button class="create-btn empty-create" @click="router.push('/articles/create')">
        写新文章
      </button>
    </div>

    <!-- 分页 -->
    <div v-if="totalPages() > 1" class="pagination">
      <button
        class="page-btn"
        :disabled="currentPage <= 1"
        @click="goToPage(currentPage - 1)"
      >
        上一页
      </button>
      <span class="page-info">{{ currentPage }} / {{ totalPages() }}</span>
      <button
        class="page-btn"
        :disabled="currentPage >= totalPages()"
        @click="goToPage(currentPage + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<style scoped>
.my-articles-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.page-heading h1 {
  font-size: 1.5rem;
  font-weight: 700;
}

.page-subtitle {
  color: var(--color-text-light);
  font-size: 0.875rem;
  margin-top: 4px;
}

.create-btn {
  padding: 10px 24px;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-sm);
  font-weight: 600;
  font-size: 0.9375rem;
  transition: background 0.2s;
  flex-shrink: 0;
}

.create-btn:hover {
  background: var(--color-primary-hover);
}

.loading-state {
  text-align: center;
  padding: 60px 0;
  color: var(--color-text-light);
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-white);
  border-radius: var(--radius-md);
  padding: 20px 24px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.2s;
  gap: 16px;
}

.article-item:hover {
  box-shadow: var(--shadow-md);
}

.item-body {
  flex: 1;
  min-width: 0;
  cursor: pointer;
}

.item-title {
  font-size: 1.0625rem;
  font-weight: 700;
  color: var(--color-text);
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.top-badge {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 4px;
  background: linear-gradient(135deg, #f6d365, #fda085);
  color: #fff;
  font-size: 0.6875rem;
  font-weight: 700;
  flex-shrink: 0;
}

.item-summary {
  font-size: 0.875rem;
  color: var(--color-text-light);
  margin-top: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  font-size: 0.8125rem;
  color: var(--color-text-light);
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-dot {
  color: #d0d5dd;
}

.item-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.action-edit,
.action-delete {
  padding: 6px 16px;
  border-radius: var(--radius-sm);
  font-size: 0.8125rem;
  font-weight: 600;
  transition: all 0.2s;
}

.action-edit {
  background: #f0f4f9;
  color: var(--color-primary);
}

.action-edit:hover {
  background: var(--color-primary);
  color: #fff;
}

.action-delete {
  background: #fef2f2;
  color: var(--color-danger);
}

.action-delete:hover {
  background: var(--color-danger);
  color: #fff;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 12px;
}

.empty-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text);
}

.empty-desc {
  margin-top: 6px;
  color: var(--color-text-light);
  margin-bottom: 20px;
}

.empty-create {
  display: inline-block;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
}

.page-btn {
  padding: 8px 20px;
  border-radius: var(--radius-sm);
  background: var(--color-white);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  font-weight: 500;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 0.875rem;
  color: var(--color-text-light);
  font-weight: 500;
}
</style>
