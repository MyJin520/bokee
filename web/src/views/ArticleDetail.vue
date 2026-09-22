<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticleInfo, deleteArticle } from '@/api/article'
import { useUserStore } from '@/stores/user'
import type { ArticleInfo } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref<ArticleInfo | null>(null)
const loading = ref(true)

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function goToEdit() {
  if (article.value) {
    router.push(`/articles/${article.value.id}/edit`)
  }
}

async function handleDelete() {
  if (!article.value) return
  if (!confirm('确定要删除这篇文章吗？')) return

  try {
    await deleteArticle(article.value.id)
    alert('删除成功')
    router.push('/')
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    router.push('/')
    return
  }
  try {
    article.value = await getArticleInfo(id)
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '文章不存在')
    router.push('/')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="loading-state">
    <div class="spinner"></div>
    <p>加载中...</p>
  </div>

  <div v-else-if="article" class="article-detail">
    <article class="detail-card">
      <h1 class="detail-title">{{ article.title }}</h1>

      <div class="detail-meta">
        <span class="meta-item">{{ formatDate(article.createdAt) }}</span>
        <span class="meta-divider">·</span>
        <span class="meta-item">👁 {{ article.viewCount }} 阅读</span>
        <span class="meta-divider">·</span>
        <span class="meta-item">👍 {{ article.likeCount }} 点赞</span>
        <span v-if="article.isTop" class="top-badge">置顶</span>
      </div>

      <div v-if="article.summary" class="detail-summary">
        <p>{{ article.summary }}</p>
      </div>

      <div class="detail-content markdown-body" v-html="article.content"></div>

      <!-- 操作按钮：仅文章作者可见 -->
      <div
        v-if="userStore.userInfo?.id === article.userId"
        class="detail-actions"
      >
        <button class="action-btn action-edit" @click="goToEdit">编辑</button>
        <button class="action-btn action-delete" @click="handleDelete">删除</button>
      </div>
    </article>

    <div class="back-link">
      <router-link to="/">← 返回首页</router-link>
    </div>
  </div>

  <div v-else class="not-found">
    <p>文章不存在</p>
    <router-link to="/">返回首页</router-link>
  </div>
</template>

<style scoped>
.loading-state {
  text-align: center;
  padding: 80px 0;
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

.article-detail {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-card {
  background: var(--color-white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  padding: 40px;
  border: 1px solid var(--color-border);
}

.detail-title {
  font-size: 1.875rem;
  font-weight: 800;
  line-height: 1.3;
  color: var(--color-text);
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  font-size: 0.875rem;
  color: var(--color-text-light);
  flex-wrap: wrap;
}

.meta-divider {
  color: #d0d5dd;
}

.top-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  background: linear-gradient(135deg, #f6d365, #fda085);
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
}

.detail-summary {
  margin-top: 24px;
  padding: 16px 20px;
  background: #f8f9fc;
  border-left: 4px solid var(--color-primary);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  color: var(--color-text-light);
  font-size: 0.9375rem;
  line-height: 1.7;
}

.detail-content {
  margin-top: 28px;
  font-size: 1rem;
  line-height: 1.8;
  color: var(--color-text);
}

.detail-content :deep(h1),
.detail-content :deep(h2),
.detail-content :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.6em;
  font-weight: 700;
}

.detail-content :deep(p) {
  margin-bottom: 1em;
}

.detail-content :deep(code) {
  padding: 2px 6px;
  background: #f0f2f5;
  border-radius: 4px;
  font-size: 0.875em;
}

.detail-content :deep(pre) {
  padding: 16px 20px;
  background: #282c34;
  border-radius: var(--radius-sm);
  overflow-x: auto;
  margin: 1em 0;
}

.detail-content :deep(pre code) {
  background: none;
  padding: 0;
  color: #abb2bf;
}

.detail-content :deep(img) {
  border-radius: var(--radius-sm);
  margin: 1em 0;
}

.detail-content :deep(blockquote) {
  padding: 12px 20px;
  margin: 1em 0;
  border-left: 4px solid var(--color-primary);
  background: #f8f9fc;
  color: var(--color-text-light);
}

.detail-content :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}

.detail-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid var(--color-border);
}

.action-btn {
  padding: 10px 24px;
  border-radius: var(--radius-sm);
  font-weight: 600;
  font-size: 0.9375rem;
  transition: all 0.2s;
}

.action-edit {
  background: var(--color-primary);
  color: #fff;
}

.action-edit:hover {
  background: var(--color-primary-hover);
}

.action-delete {
  background: var(--color-white);
  color: var(--color-danger);
  border: 1px solid var(--color-danger);
}

.action-delete:hover {
  background: #fef2f2;
}

.back-link {
  text-align: center;
}

.back-link a {
  font-size: 0.9375rem;
  font-weight: 500;
}

.not-found {
  text-align: center;
  padding: 80px 0;
  font-size: 1.125rem;
}

.not-found a {
  display: inline-block;
  margin-top: 12px;
  font-weight: 600;
}
</style>
