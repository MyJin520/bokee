<script lang="ts">
import { marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.min.css'

// 配置 marked：注册代码高亮扩展
// 放在 <script> 中仅在模块加载时执行一次，避免每次组件挂载重复注册导致代码被多次高亮
marked.use(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code: string, lang: string) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    },
  }),
)
</script>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router'
import { getArticleInfo, deleteArticle } from '@/api/article'
import { useUserStore } from '@/stores/user'
import type { ArticleInfo } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref<ArticleInfo | null>(null)
const renderedContent = ref('')
const loading = ref(true)

// 当 article 变化时异步解析 Markdown
watch(article, async (val) => {
  if (!val?.content) {
    renderedContent.value = ''
    return
  }
  renderedContent.value = await marked.parse(val.content)
})

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

/** 根据 route param 加载文章 */
async function fetchArticle(id: number) {
  if (!id) {
    router.push('/')
    return
  }
  loading.value = true
  try {
    article.value = await getArticleInfo(id)
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '文章不存在')
    router.push('/')
  } finally {
    loading.value = false
  }
}

// 页面刷新 / 首次进入
onMounted(() => {
  fetchArticle(Number(route.params.id))
})

// 同一组件内切换文章（如从 /article/1 → /article/2）
onBeforeRouteUpdate((to) => {
  fetchArticle(Number(to.params.id))
})

// 跳转到编辑页
function goToEdit() {
  if (article.value) {
    router.push(`/articles/${article.value.id}/edit`)
  }
}

// 删除当前文章
async function handleDelete() {
  if (!article.value) return
  if (!confirm('确定要删除这篇文章吗？此操作不可撤销。')) return
  try {
    await deleteArticle(article.value.id)
    alert('删除成功')
    router.push('/')
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}
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

      <div class="detail-content markdown-body" v-html="renderedContent"></div>

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
.detail-content :deep(h3),
.detail-content :deep(h4) {
  margin-top: 1.5em;
  margin-bottom: 0.6em;
  font-weight: 700;
  color: var(--color-text);
}

.detail-content :deep(h1) { font-size: 1.75rem; }
.detail-content :deep(h2) { font-size: 1.45rem; border-bottom: 1px solid var(--color-border); padding-bottom: 0.3em; }
.detail-content :deep(h3) { font-size: 1.2rem; }
.detail-content :deep(h4) { font-size: 1.05rem; }

.detail-content :deep(p) {
  margin-bottom: 1em;
}

.detail-content :deep(ul),
.detail-content :deep(ol) {
  margin: 0.5em 0 1em 1.5em;
  line-height: 1.8;
}

.detail-content :deep(li) {
  margin-bottom: 0.3em;
}

/* 行内代码（不在 pre 内）*/
.detail-content :deep(p code),
.detail-content :deep(li code),
.detail-content :deep(h1 code),
.detail-content :deep(h2 code),
.detail-content :deep(h3 code),
.detail-content :deep(h4 code),
.detail-content :deep(blockquote code) {
  padding: 2px 7px;
  background: #f0f2f5;
  border-radius: 4px;
  font-size: 0.875em;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  color: #d63384;
}

/* highlight.js 生成的代码块 - 完全交由 highlight.js CSS 控制 */
.detail-content :deep(pre code) {
  all: revert;
  background: none;
  padding: 0;
  color: inherit;
  font-size: inherit;
}

.detail-content :deep(pre) {
  border-radius: var(--radius-sm);
  overflow-x: auto;
  margin: 1.2em 0;
  line-height: 1.55;
  font-size: 0.875rem;
}

.detail-content :deep(img) {
  border-radius: var(--radius-sm);
  margin: 1em 0;
  max-width: 100%;
}

.detail-content :deep(blockquote) {
  padding: 12px 20px;
  margin: 1em 0;
  border-left: 4px solid var(--color-primary);
  background: #f8f9fc;
  color: var(--color-text-light);
  font-style: italic;
}

.detail-content :deep(blockquote p) {
  margin-bottom: 0;
}

.detail-content :deep(a) {
  color: var(--color-primary);
  text-decoration: underline;
}

.detail-content :deep(hr) {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 1.5em 0;
}

.detail-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1em 0;
  font-size: 0.9375rem;
}

.detail-content :deep(th),
.detail-content :deep(td) {
  padding: 8px 14px;
  border: 1px solid var(--color-border);
  text-align: left;
}

.detail-content :deep(th) {
  background: #f8f9fc;
  font-weight: 700;
}

.detail-content :deep(strong) {
  font-weight: 700;
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
