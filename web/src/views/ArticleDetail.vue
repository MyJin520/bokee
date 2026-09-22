<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router'
import { getArticleInfo, deleteArticle, getUserArticleList } from '@/api/article'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import { useReactions } from '@/composables/useReactions'
import { renderMarkdown } from '@/utils/markdown'
import { formatDate, formatShortDate, formatNumber, initialOf } from '@/utils/format'
import type { ArticleInfo, ArticleListItem } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const ui = useUiStore()
const { isLiked, isSaved, toggleLiked, toggleSaved } = useReactions()

const article = ref<ArticleInfo | null>(null)
const renderedContent = ref('')
const loading = ref(true)
const notFound = ref(false)
const authorArticles = ref<ArticleListItem[]>([])
const authorArticlesLoading = ref(false)
const coverFailed = ref(false)

const isOwner = computed(
  () => Boolean(article.value) && userStore.userInfo?.id === article.value?.userId,
)

watch(article, async (val) => {
  renderedContent.value = val?.content ? await renderMarkdown(val.content) : ''
  if (val?.userId) fetchAuthorArticles(val.userId, val.id)
})

async function fetchAuthorArticles(userId: number, currentArticleId: number) {
  authorArticlesLoading.value = true
  try {
    const res = await getUserArticleList(userId, 1, 10)
    authorArticles.value = res.list.filter((item) => item.id !== currentArticleId).slice(0, 6)
  } catch {
    authorArticles.value = []
  } finally {
    authorArticlesLoading.value = false
  }
}

async function fetchArticle(id: number) {
  if (!id) {
    router.replace('/')
    return
  }
  loading.value = true
  notFound.value = false
  coverFailed.value = false
  try {
    article.value = await getArticleInfo(id)
  } catch (error: unknown) {
    article.value = null
    notFound.value = true
    ui.toastError(error instanceof Error ? error.message : '文章不存在')
  } finally {
    loading.value = false
  }
}

function onLike() {
  if (!article.value) return
  const nowLiked = toggleLiked(article.value.id)
  article.value.likeCount += nowLiked ? 1 : -1
}

function onSave() {
  if (!article.value) return
  const nowSaved = toggleSaved(article.value.id)
  ui.toast(nowSaved ? '已收藏这篇文章' : '已取消收藏')
}

function goToEdit() {
  if (article.value) router.push(`/articles/${article.value.id}/edit`)
}

async function handleDelete() {
  if (!article.value) return
  const ok = await ui.confirm({
    title: '删除文章',
    message: `确定要删除《${article.value.title}》吗？此操作不可撤销。`,
    confirmText: '确认删除',
    danger: true,
  })
  if (!ok) return
  try {
    await deleteArticle(article.value.id)
    ui.toast('文章已删除')
    router.push('/my-articles')
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '删除失败，请稍后重试')
  }
}

onMounted(() => fetchArticle(Number(route.params.id)))
onBeforeRouteUpdate((to) => {
  window.scrollTo({ top: 0 })
  fetchArticle(Number(to.params.id))
})
</script>

<template>
  <div class="detail-page">
    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取文章…</p>
    </div>

    <div v-else-if="notFound" class="state-block">
      <p class="state-title">文章不存在</p>
      <p>它可能已被作者删除，或链接有误。</p>
      <router-link class="text-button" to="/articles" style="display: inline-block; margin-top: 14px">
        去文章列表看看 →
      </router-link>
    </div>

    <template v-else-if="article">
      <div class="detail-layout">
        <article class="detail-main">
          <button class="back-link" type="button" @click="router.back()">← 返回</button>

          <header class="detail-header">
            <span class="tag" v-if="article.isTop">置顶文章</span>
            <h1>{{ article.title }}</h1>
            <div class="detail-meta">
              <span class="meta-author">
                <span class="meta-avatar">
                  <img v-if="article.author?.avatar" :src="article.author.avatar" :alt="article.author.name" />
                  <template v-else>{{ initialOf(article.author?.name || article.title) }}</template>
                </span>
                <span class="meta-author-name">{{ article.author?.name || '匿名作者' }}</span>
              </span>
              <span class="meta-dot">·</span>
              <span>{{ formatDate(article.createdAt) }}</span>
              <span class="meta-dot">·</span>
              <span>阅读 {{ formatNumber(article.viewCount) }}</span>
            </div>
          </header>

          <figure v-if="article.cover && !coverFailed" class="detail-cover">
            <img :src="article.cover" :alt="article.title" @error="coverFailed = true" />
          </figure>

          <p v-if="article.summary" class="detail-summary">{{ article.summary }}</p>

          <div class="markdown-body detail-content" v-html="renderedContent"></div>

          <div class="detail-reactions">
            <button
              :class="['reaction', { active: isLiked(article.id) }]"
              type="button"
              @click="onLike"
            >
              ♥ 点赞 · {{ formatNumber(article.likeCount) }}
            </button>
            <button
              :class="['reaction', { active: isSaved(article.id) }]"
              type="button"
              @click="onSave"
            >
              {{ isSaved(article.id) ? '♥ 已收藏' : '♡ 收藏' }}
            </button>
          </div>

          <div v-if="isOwner" class="detail-actions">
            <button class="btn btn-outline btn-sm" type="button" @click="goToEdit">编辑文章</button>
            <button class="btn btn-danger-outline btn-sm" type="button" @click="handleDelete">删除文章</button>
          </div>
        </article>

        <aside class="detail-sidebar">
          <section class="author-card">
            <p class="side-title author-side-title">About the author</p>
            <span class="author-avatar">
              <img v-if="article.author?.avatar" :src="article.author.avatar" :alt="article.author.name" />
              <template v-else>{{ initialOf(article.author?.name || '作') }}</template>
            </span>
            <div class="author-name">{{ article.author?.name || '匿名作者' }}</div>
            <p class="author-bio">认真记录，慢慢表达。这里是 TA 的公开文章与思考。</p>
            <router-link class="author-link" :to="`/articles?user=${encodeURIComponent(article.author?.name || '')}`">
              查看 TA 的文章 →
            </router-link>
          </section>

          <section class="more-card">
            <p class="side-title">More from author</p>
            <div v-if="authorArticlesLoading" class="more-loading">加载中…</div>
            <ul v-else-if="authorArticles.length" class="more-list">
              <li v-for="item in authorArticles" :key="item.id">
                <router-link :to="`/article/${item.id}`">
                  <span class="more-title">{{ item.title }}</span>
                  <span class="more-date">{{ formatShortDate(item.createdAt) }}</span>
                </router-link>
              </li>
            </ul>
            <p v-else class="more-empty">暂无其他文章</p>
          </section>
        </aside>
      </div>
    </template>
  </div>
</template>

<style scoped>
.detail-page {
  padding: 40px 0 80px;
}

.detail-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 52px;
  align-items: start;
}

.back-link {
  padding: 0;
  margin-bottom: 26px;
  background: transparent;
  color: var(--muted);
  font-size: 12px;
}

.back-link:hover {
  color: var(--ink);
}

.detail-header h1 {
  margin: 16px 0 18px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(32px, 4.4vw, 52px);
  font-weight: 500;
  line-height: 1.08;
  letter-spacing: -0.02em;
}

.detail-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 9px;
  color: var(--muted);
  font-size: 11px;
}

.meta-author {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.meta-avatar,
.author-avatar {
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--coral);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
}

.meta-avatar {
  width: 26px;
  height: 26px;
  font-size: 11px;
}

.meta-avatar img,
.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.meta-author-name {
  color: var(--ink);
  font-weight: 700;
}

.meta-dot {
  color: #b2bbb5;
}

.detail-cover {
  margin: 28px 0 0;
}

.detail-cover img {
  width: 100%;
  max-height: 480px;
  display: block;
  object-fit: cover;
  border: 1px solid var(--line);
}

.detail-summary {
  margin: 26px 0 0;
  padding: 16px 20px;
  border-left: 3px solid var(--coral);
  background: var(--surface);
  color: var(--muted);
  font-size: 14px;
  line-height: 1.8;
  font-style: italic;
}

.detail-content {
  margin-top: 30px;
}

.detail-reactions {
  display: flex;
  gap: 10px;
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
}

.reaction {
  padding: 9px 16px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--muted);
  font-size: 12px;
  font-weight: 700;
}

.reaction:hover,
.reaction.active {
  border-color: var(--coral);
  color: var(--coral-dark);
  background: var(--coral-tint);
}

.detail-actions {
  display: flex;
  gap: 10px;
  margin-top: 18px;
}

.detail-sidebar {
  position: sticky;
  top: 100px;
  display: flex;
  flex-direction: column;
  gap: 26px;
}

.author-card {
  padding: 24px 21px;
  background: var(--ink);
  color: #fff;
  text-align: center;
}

.author-side-title {
  color: #a9b7ae !important;
  text-align: left;
}

.author-avatar {
  width: 68px;
  height: 68px;
  margin: 12px auto 12px;
  border: 2px solid var(--lime);
  font-size: 28px;
}

.author-name {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 21px;
}

.author-bio {
  margin: 12px 0 16px;
  color: #c4cec7;
  font-size: 12px;
  line-height: 1.7;
}

.author-link {
  color: var(--lime);
  font-size: 11px;
  font-weight: 700;
}

.more-card {
  padding: 21px;
  background: var(--surface);
  border: 1px solid var(--line);
}

.more-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.more-list li {
  border-bottom: 1px solid var(--line);
}

.more-list li:last-child {
  border-bottom: 0;
}

.more-list a {
  display: block;
  padding: 11px 0;
}

.more-title {
  display: block;
  font-size: 12.5px;
  line-height: 1.5;
}

.more-list a:hover .more-title {
  color: var(--coral-dark);
}

.more-date {
  display: block;
  margin-top: 3px;
  color: var(--muted);
  font-size: 10px;
}

.more-loading,
.more-empty {
  color: var(--muted);
  font-size: 12px;
}

@media (max-width: 880px) {
  .detail-layout {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .detail-sidebar {
    position: static;
  }
}
</style>
