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
import { getArticleInfo, deleteArticle, getUserArticleList } from '@/api/article'
import { useUserStore } from '@/stores/user'
import type { ArticleInfo, ArticleListItem } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref<ArticleInfo | null>(null)
const renderedContent = ref('')
const loading = ref(true)
const authorArticles = ref<ArticleListItem[]>([])
const authorArticlesLoading = ref(false)

// 当 article 变化时异步解析 Markdown 并加载作者其他文章
watch(article, async (val) => {
  if (!val?.content) {
    renderedContent.value = ''
    return
  }
  renderedContent.value = await marked.parse(val.content)
  if (val.userId) {
    fetchAuthorArticles(val.userId, val.id)
  }
})

/** 加载作者的其他文章（排除当前文章） */
async function fetchAuthorArticles(userId: number, currentArticleId: number) {
  authorArticlesLoading.value = true
  try {
    const res = await getUserArticleList(userId, 1, 10)
    authorArticles.value = res.list.filter((a) => a.id !== currentArticleId)
  } catch (e) {
    console.error('加载作者其他文章失败', e)
    authorArticles.value = []
  } finally {
    authorArticlesLoading.value = false
  }
}

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

function formatShortDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
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

function goToArticle(id: number) {
  router.push(`/article/${id}`)
}
</script>

<template>
  <div v-if="loading" class="loading-state">
    <div class="spinner"></div>
    <p>加载中...</p>
  </div>

  <div v-else-if="article" class="article-detail">
    <div class="detail-layout">
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

      <!-- 右侧边栏：作者信息 + 作者其他文章 -->
      <aside class="detail-sidebar">
        <!-- 作者信息卡片 -->
        <div class="author-card">
          <div class="author-avatar">
            <img v-if="article.author.avatar" :src="article.author.avatar" :alt="article.author.name" />
            <span v-else class="avatar-placeholder">{{ article.author.name?.charAt(0) || 'U' }}</span>
          </div>
          <div class="author-name">{{ article.author.name || '匿名用户' }}</div>
        </div>

        <!-- 作者其他文章 -->
        <div class="author-articles">
          <h3 class="sidebar-title">作者其他文章</h3>
          <div v-if="authorArticlesLoading" class="sidebar-loading">加载中...</div>
          <ul v-else-if="authorArticles.length" class="author-article-list">
            <li
              v-for="item in authorArticles"
              :key="item.id"
              class="author-article-item"
              @click="goToArticle(item.id)"
            >
              <span class="article-item-title">{{ item.title }}</span>
              <span class="article-item-meta">{{ formatShortDate(item.createdAt) }}</span>
            </li>
          </ul>
          <p v-else class="sidebar-empty">暂无其他文章</p>
        </div>
      </aside>
    </div>

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

.detail-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 24px;
  align-items: start;
}

.detail-card {
  background: var(--color-white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  padding: 40px;
  border: 1px solid var(--color-border);
  min-width: 0;
}

/* 右侧边栏 */
.detail-sidebar {
  position: sticky;
  top: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 作者信息卡片 */
.author-card {
  background: var(--color-white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  padding: 28px 20px;
  text-align: center;
}

.author-avatar {
  width: 72px;
  height: 72px;
  margin: 0 auto 14px;
  border-radius: 50%;
  overflow: hidden;
  background: linear-gradient(135deg, var(--color-primary), #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 1.75rem;
  font-weight: 700;
  color: #fff;
}

.author-name {
  font-size: 1.0625rem;
  font-weight: 700;
  color: var(--color-text);
}

/* 作者其他文章 */
.author-articles {
  background: var(--color-white);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  padding: 20px;
}

.sidebar-title {
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-border);
}

.author-article-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.author-article-item {
  cursor: pointer;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  transition: background 0.2s;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-article-item:hover {
  background: #f8f9fc;
}

.article-item-title {
  font-size: 0.875rem;
  color: var(--color-text);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-item-meta {
  font-size: 0.75rem;
  color: var(--color-text-light);
}

.sidebar-loading,
.sidebar-empty {
  font-size: 0.875rem;
  color: var(--color-text-light);
  text-align: center;
  padding: 16px 0;
}

/* 响应式：窄屏时边栏下移 */
@media (max-width: 960px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }

  .detail-sidebar {
    position: static;
  }
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
