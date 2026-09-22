<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getArticleList } from '@/api/article'
import ArticleCard from '@/components/ArticleCard.vue'
import type { ArticleListItem } from '@/types'

const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const searchTitle = ref('')

async function loadArticles() {
  loading.value = true
  try {
    const res = await getArticleList({
      title: searchTitle.value || undefined,
      page: currentPage.value,
      pageSize: pageSize.value,
    })
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
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function onSearch() {
  currentPage.value = 1
  loadArticles()
}

onMounted(() => {
  loadArticles()
})
</script>

<template>
  <div class="home-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">我的博客</h1>
      <p class="page-subtitle">记录技术点滴，分享学习心得</p>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <input
        v-model="searchTitle"
        type="text"
        class="search-input"
        placeholder="搜索文章标题..."
        @keyup.enter="onSearch"
      />
      <button class="search-btn" @click="onSearch">搜索</button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 文章列表 -->
    <div v-else-if="articles.length > 0" class="article-list">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">📭</div>
      <p class="empty-title">暂无文章</p>
      <p class="empty-desc">还没有人发表文章，快来写下第一篇吧！</p>
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
.home-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  text-align: center;
  padding: 32px 0 16px;
}

.page-title {
  font-size: 2rem;
  font-weight: 800;
  background: linear-gradient(135deg, #4a90d9, #7c5cbf);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.page-subtitle {
  margin-top: 8px;
  color: var(--color-text-light);
  font-size: 1rem;
}

.search-bar {
  display: flex;
  gap: 8px;
  max-width: 480px;
  margin: 0 auto;
  width: 100%;
}

.search-input {
  flex: 1;
  padding: 10px 16px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: 0.9375rem;
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus {
  border-color: var(--color-primary);
}

.search-btn {
  padding: 10px 24px;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-sm);
  font-weight: 600;
  transition: background 0.2s;
}

.search-btn:hover {
  background: var(--color-primary-hover);
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
