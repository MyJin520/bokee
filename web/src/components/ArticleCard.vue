<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { ArticleListItem } from '@/types'

const props = defineProps<{
  article: ArticleListItem
}>()

const router = useRouter()

function goToDetail() {
  router.push(`/article/${props.article.id}`)
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>

<template>
  <article class="article-card" @click="goToDetail">
    <div class="card-body">
      <div class="card-header">
        <span v-if="article.isTop" class="top-badge">置顶</span>
        <h2 class="card-title">{{ article.title }}</h2>
      </div>

      <p v-if="article.summary" class="card-summary">{{ article.summary }}</p>
      <p v-else class="card-summary card-summary-empty">暂无摘要</p>

      <div class="card-meta">
        <span class="meta-item meta-date">{{ formatDate(article.createdAt) }}</span>
        <span class="meta-divider">·</span>
        <span class="meta-item meta-views">👁 {{ article.viewCount }}</span>
        <span class="meta-divider">·</span>
        <span class="meta-item meta-likes">👍 {{ article.likeCount }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.article-card {
  background: var(--color-white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  padding: 24px 28px;
  cursor: pointer;
  transition: all 0.25s ease;
  border: 1px solid var(--color-border);
}

.article-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--color-primary);
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.top-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  background: linear-gradient(135deg, #f6d365, #fda085);
  color: #fff;
  font-size: 0.75rem;
  font-weight: 700;
  flex-shrink: 0;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-summary {
  font-size: 0.9375rem;
  color: var(--color-text-light);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-summary-empty {
  color: #ccc;
  font-style: italic;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  color: var(--color-text-light);
  margin-top: 4px;
}

.meta-divider {
  color: #d0d5dd;
}
</style>
