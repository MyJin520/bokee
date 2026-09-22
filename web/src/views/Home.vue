<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getArticleList } from '@/api/article'
import { useUiStore } from '@/stores/ui'
import { useUserStore } from '@/stores/user'
import type { ArticleListItem } from '@/types'
import ArticleRow from '@/components/ArticleRow.vue'
import { formatNumber, formatDate, initialOf } from '@/utils/format'

type FilterKey = 'all' | 'top' | 'covered'

const router = useRouter()
const ui = useUiStore()
const userStore = useUserStore()

const articles = ref<ArticleListItem[]>([])
const statsArticles = ref<ArticleListItem[]>([])
const loading = ref(true)
const errorMessage = ref('')
const searchTitle = ref('')
const activeFilter = ref<FilterKey>('all')
const searchInput = ref<HTMLInputElement | null>(null)
const newsletterEmail = ref('')
const newsletterStatus = ref('')
const featureFailed = ref(false)
let searchTimer: ReturnType<typeof setTimeout> | undefined

const filteredArticles = computed(() =>
  articles.value.filter((article) => {
    const keyword = searchTitle.value.trim().toLowerCase()
    const text = `${article.title} ${article.summary || ''}`.toLowerCase()
    const matchesQuery = !keyword || text.includes(keyword)
    const matchesFilter =
      activeFilter.value === 'all' ||
      (activeFilter.value === 'top' ? article.isTop : Boolean(article.cover))
    return matchesQuery && matchesFilter
  }),
)

const featuredArticle = computed(
  () => filteredArticles.value.find((article) => article.isTop) || filteredArticles.value[0],
)
const feedArticles = computed(() =>
  filteredArticles.value.filter((article) => article.id !== featuredArticle.value?.id),
)

const totalArticles = computed(() => statsArticles.value.length)
const totalViews = computed(() => statsArticles.value.reduce((sum, item) => sum + item.viewCount, 0))

interface ArchiveBucket {
  key: string
  label: string
  count: number
}

const archiveBuckets = computed<ArchiveBucket[]>(() => {
  const buckets = new Map<string, number>()
  statsArticles.value.forEach((article) => {
    const date = new Date(article.createdAt)
    if (Number.isNaN(date.getTime())) return
    const key = `${date.getFullYear()}-${date.getMonth()}`
    buckets.set(key, (buckets.get(key) || 0) + 1)
  })
  return Array.from(buckets.entries())
    .sort((a, b) => (a[0] < b[0] ? 1 : -1))
    .slice(0, 3)
    .map(([key, count]) => {
      const [yearStr, monthStr] = key.split('-')
      const date = new Date(Number(yearStr), Number(monthStr), 1)
      return {
        key,
        label: `${date.toLocaleDateString('zh-CN', { month: 'long' })} · ${date.toLocaleDateString('en-US', { month: 'long' })}`,
        count,
      }
    })
})

const todayLabel = new Date().toLocaleDateString('zh-CN', {
  year: 'numeric',
  month: '2-digit',
  day: '2-digit',
})

async function loadArticles() {
  loading.value = true
  errorMessage.value = ''
  try {
    const keyword = searchTitle.value.trim() || undefined
    const [feed, stats] = await Promise.all([
      getArticleList({ page: 1, pageSize: 20, title: keyword }),
      getArticleList({ page: 1, pageSize: 100 }),
    ])
    articles.value = feed.list
    statsArticles.value = stats.list
    featureFailed.value = false
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '文章加载失败，请稍后重试。'
    articles.value = []
    statsArticles.value = []
  } finally {
    loading.value = false
  }
}

/** 搜索输入防抖后请求后端标题模糊查询 */
watch(searchTitle, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(loadArticles, 350)
})

function selectFilter(filter: FilterKey) {
  activeFilter.value = filter
}

function selectTopic(filter: FilterKey) {
  activeFilter.value = filter
  document.getElementById('latest')?.scrollIntoView({ behavior: 'smooth' })
}

function openArticle(id: number) {
  router.push(`/article/${id}`)
}

function subscribeNewsletter() {
  const email = newsletterEmail.value.trim()
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    newsletterStatus.value = '请输入有效的邮箱地址。'
    return
  }
  // 后端暂无订阅接口，与原型一致以前端反馈形式承接
  newsletterStatus.value = '订阅成功，下一封信见。'
  newsletterEmail.value = ''
  ui.toast('订阅成功，下一封信见')
}

onMounted(loadArticles)

// 头部搜索按钮：滚动到列表并聚焦搜索框
watch(
  () => ui.homeSearchSignal,
  () => {
    window.setTimeout(() => searchInput.value?.focus(), 350)
  },
)
</script>

<template>
  <div class="home-page">
    <section class="hero-section" aria-label="博客介绍">
      <div class="hero-copy">
        <span class="eyebrow">A quiet corner on the internet</span>
        <h1>
          把日常写成<br /><em>值得回看的事。</em>
        </h1>
        <p>这里记录产品、技术与生活里的小发现。慢一点思考，认真留下每一次有用的连接。</p>
        <div class="hero-note">
          <span class="avatar-stack" aria-hidden="true"><i>林</i><i>默</i><i>客</i></span>
          <span>{{ totalArticles }} 篇文章 · 共 {{ formatNumber(totalViews) }} 次阅读</span>
        </div>
      </div>
      <div
        class="hero-image"
        role="button"
        tabindex="0"
        aria-label="阅读本周精选文章"
        @click="featuredArticle && openArticle(featuredArticle.id)"
        @keydown.enter="featuredArticle && openArticle(featuredArticle.id)"
      >
        <img
          src="https://images.unsplash.com/photo-1499750310107-5fef28a66643?auto=format&fit=crop&w=1200&q=85"
          alt="桌面上的笔记本与咖啡"
        />
        <span class="hero-stamp">NO.<br />{{ String(totalArticles).padStart(3, '0') }}</span>
        <div class="hero-caption">
          <small>本周精选 · {{ todayLabel }}</small>
          <strong>{{ featuredArticle?.title || '让一个人项目，拥有自己的节奏' }}</strong>
        </div>
      </div>
    </section>

    <section class="content-grid" id="latest">
      <div class="main-column">
        <div class="section-head">
          <div>
            <p class="section-kicker">Latest notes</p>
            <h2>最近写了什么</h2>
          </div>
          <router-link class="section-link" to="/articles">查看全部文章 ↗</router-link>
        </div>

        <div class="filters" role="group" aria-label="文章筛选">
          <button :class="['chip', { active: activeFilter === 'all' }]" type="button" @click="selectFilter('all')">全部</button>
          <button :class="['chip', { active: activeFilter === 'top' }]" type="button" @click="selectFilter('top')">置顶</button>
          <button :class="['chip', { active: activeFilter === 'covered' }]" type="button" @click="selectFilter('covered')">有封面</button>
          <label class="search-box">
            <input
              ref="searchInput"
              v-model="searchTitle"
              type="search"
              placeholder="搜索标题或摘要"
              aria-label="搜索文章"
            />
            <span>⌕</span>
          </label>
        </div>

        <div v-if="loading" class="state-block">
          <div class="spinner"></div>
          <p>正在读取最新文章…</p>
        </div>
        <div v-else-if="errorMessage" class="state-block">
          <p class="state-title">加载遇到问题</p>
          <p>{{ errorMessage }}</p>
          <button class="text-button" type="button" style="margin-top: 12px" @click="loadArticles">重新加载 →</button>
        </div>
        <template v-else-if="filteredArticles.length">
          <article v-if="featuredArticle" class="feature-card" @click="openArticle(featuredArticle.id)">
            <div class="feature-cover">
              <img
                v-if="featuredArticle.cover && !featureFailed"
                :src="featuredArticle.cover"
                :alt="featuredArticle.title"
                @error="featureFailed = true"
              />
              <div v-else class="cover-fallback">{{ initialOf(featuredArticle.title) }}</div>
            </div>
            <div class="feature-body">
              <div>
                <span class="tag">{{ featuredArticle.isTop ? '置顶文章' : '最新文章' }}</span>
                <h3>{{ featuredArticle.title }}</h3>
                <p>{{ featuredArticle.summary || '这篇文章还没有摘要，点击阅读完整内容。' }}</p>
              </div>
              <div class="article-meta">
                <span>{{ formatDate(featuredArticle.createdAt) }}</span>
                <span class="meta-sep">·</span>
                <span>阅读 {{ formatNumber(featuredArticle.viewCount) }}</span>
                <span class="meta-sep">·</span>
                <span>点赞 {{ formatNumber(featuredArticle.likeCount) }}</span>
              </div>
            </div>
          </article>

          <div class="article-list">
            <ArticleRow v-for="article in feedArticles" :key="article.id" :article="article" />
          </div>
        </template>
        <div v-else class="state-block">
          <p class="state-title">没有找到匹配的文章</p>
          <p>换一个关键词，或清除筛选条件试试。</p>
          <button
            class="text-button"
            type="button"
            style="margin-top: 12px"
            @click="searchTitle = ''; selectFilter('all')"
          >
            清除筛选 →
          </button>
        </div>
      </div>

      <aside class="sidebar" id="about">
        <!-- 作者卡片：已登录展示当前用户，未登录展示博客介绍 -->
        <section v-if="userStore.isLoggedIn" class="profile-card">
          <p class="side-title profile-side-title">My account</p>
          <div class="profile-top">
            <span class="profile-avatar">
              <img v-if="userStore.userInfo?.avatar" :src="userStore.userInfo.avatar" alt="头像" />
              <template v-else>{{ userStore.userName.slice(0, 1) || '我' }}</template>
            </span>
            <div>
              <div class="profile-name">{{ userStore.userName }}</div>
              <div class="profile-role">{{ userStore.userInfo?.email || '已登录的写作者' }}</div>
            </div>
          </div>
          <p>在这里管理你的资料、文章，继续未完成的写作。</p>
          <div class="profile-links">
            <router-link class="profile-link" to="/profile">个人中心 →</router-link>
            <router-link class="profile-link" to="/my-articles">我的文章 →</router-link>
          </div>
        </section>
        <section v-else class="profile-card">
          <p class="side-title profile-side-title">About the author</p>
          <div class="profile-top">
            <span class="profile-avatar">林</span>
            <div>
              <div class="profile-name">林默 / BOKEE</div>
              <div class="profile-role">产品设计 · 独立开发</div>
            </div>
          </div>
          <p>白天做产品，晚上写代码。相信好的记录，能让复杂的事情重新变得清楚。</p>
          <router-link class="profile-link" to="/login">登录后开始写作 →</router-link>
        </section>

        <section class="side-block" id="topics">
          <p class="side-title">Popular topics</p>
          <div class="topic-list">
            <button class="topic" type="button" @click="selectTopic('all')">全部文章</button>
            <button class="topic" type="button" @click="selectTopic('top')">置顶文章</button>
            <button class="topic" type="button" @click="selectTopic('covered')">精选封面</button>
            <button class="topic" type="button" @click="router.push('/articles')">归档浏览</button>
          </div>
        </section>

        <section class="side-block">
          <p class="side-title">Archive</p>
          <div class="archive-list">
            <router-link
              v-for="bucket in archiveBuckets"
              :key="bucket.key"
              class="archive-row"
              to="/articles"
            >
              <span>{{ bucket.label }}</span>
              <strong>{{ String(bucket.count).padStart(2, '0') }}</strong>
            </router-link>
            <div v-if="!archiveBuckets.length" class="archive-row">
              <span>暂无归档</span>
            </div>
          </div>
        </section>

        <section class="newsletter">
          <h3>每月一封信</h3>
          <p>精选文章、最近在读，以及一些还没来得及发布的想法。</p>
          <form class="newsletter-form" @submit.prevent="subscribeNewsletter">
            <input v-model="newsletterEmail" type="email" placeholder="你的邮箱地址" aria-label="邮箱地址" />
            <button type="submit">订阅 →</button>
          </form>
          <div class="newsletter-status">{{ newsletterStatus }}</div>
        </section>
      </aside>
    </section>
  </div>
</template>

<style scoped>
.home-page {
  padding-bottom: 80px;
}

.hero-section {
  padding: 54px 0 46px;
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(280px, 0.75fr);
  gap: 30px;
  align-items: stretch;
}

.hero-copy {
  padding: 10px 0 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

h1 {
  max-width: 680px;
  margin: 22px 0 19px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(42px, 6vw, 78px);
  font-weight: 500;
  line-height: 0.98;
  letter-spacing: -0.035em;
}

h1 em {
  color: var(--coral-dark);
  font-style: normal;
}

.hero-copy p {
  max-width: 500px;
  margin: 0;
  color: var(--muted);
  font-size: 16px;
  line-height: 1.75;
}

.hero-note {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 35px;
  color: var(--muted);
  font-size: 12px;
}

.avatar-stack {
  display: flex;
  padding-left: 8px;
}

.avatar-stack i {
  width: 28px;
  height: 28px;
  margin-left: -8px;
  border: 2px solid var(--paper);
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--ink);
  color: #fff;
  font-size: 10px;
  font-style: normal;
}

.avatar-stack i:nth-child(2) {
  background: var(--coral);
}

.avatar-stack i:nth-child(3) {
  background: var(--lime);
  color: var(--ink);
}

.hero-image {
  position: relative;
  min-height: 410px;
  overflow: hidden;
  background: var(--blue);
  box-shadow: var(--shadow);
  cursor: pointer;
}

.hero-image img {
  width: 100%;
  height: 100%;
  min-height: 410px;
  display: block;
  object-fit: cover;
  filter: saturate(0.82);
}

.hero-image::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(140deg, rgba(12, 34, 35, 0.04), rgba(12, 34, 35, 0.32));
  pointer-events: none;
}

.hero-stamp {
  position: absolute;
  z-index: 1;
  top: 20px;
  right: 20px;
  width: 92px;
  height: 92px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.68);
  border-radius: 50%;
  color: #fff;
  font-family: Georgia, serif;
  font-size: 12px;
  line-height: 1.3;
  text-align: center;
  transform: rotate(12deg);
}

.hero-caption {
  position: absolute;
  z-index: 1;
  left: 22px;
  right: 22px;
  bottom: 19px;
  color: #fff;
}

.hero-caption small {
  font-size: 11px;
  opacity: 0.8;
}

.hero-caption strong {
  display: block;
  margin-top: 4px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 25px;
  line-height: 1.12;
  font-weight: 500;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 52px;
}

.section-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin: 0 0 21px;
}

.section-head h2 {
  margin: 0;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 31px;
  font-weight: 500;
  letter-spacing: -0.02em;
}

.section-link {
  color: var(--muted);
  font-size: 12px;
  border-bottom: 1px solid var(--muted);
  padding-bottom: 4px;
}

.section-link:hover {
  color: var(--ink);
  border-color: var(--ink);
}

.filters {
  display: flex;
  gap: 7px;
  align-items: center;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 3px;
}

.search-box {
  margin-left: auto;
  display: flex;
  align-items: center;
  width: 190px;
  border-bottom: 1px solid var(--line);
}

.search-box input {
  width: 100%;
  padding: 8px 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--ink);
  font-size: 12px;
}

.search-box span {
  color: var(--muted);
  font-size: 14px;
}

.feature-card {
  display: grid;
  grid-template-columns: 44% 56%;
  min-height: 262px;
  background: var(--surface);
  border: 1px solid var(--line);
  cursor: pointer;
}

.feature-card:hover {
  box-shadow: var(--shadow);
}

.feature-cover {
  overflow: hidden;
  background: #d6dfd6;
}

.feature-cover img,
.cover-fallback {
  width: 100%;
  height: 100%;
  min-height: 262px;
  display: block;
  object-fit: cover;
  transition: transform 0.45s ease;
}

.feature-card:hover .feature-cover img {
  transform: scale(1.04);
}

.cover-fallback {
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, var(--blue), var(--lime));
  color: rgba(24, 34, 31, 0.72);
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 90px;
}

.feature-body {
  padding: 28px 30px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.feature-body h3 {
  max-width: 460px;
  margin: 18px 0 10px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 31px;
  line-height: 1.13;
  font-weight: 500;
}

.feature-body p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.65;
}

.article-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
  color: var(--muted);
  font-size: 11px;
}

.meta-sep {
  color: #b2bbb5;
}

.article-list {
  margin-top: 18px;
  border-top: 1px solid var(--line);
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.profile-card {
  padding: 21px;
  background: var(--ink);
  color: #fff;
}

.profile-side-title {
  color: #a9b7ae !important;
}

.profile-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.profile-avatar {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  border-radius: 50%;
  border: 2px solid var(--lime);
  overflow: hidden;
  display: grid;
  place-items: center;
  background: var(--coral);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 20px;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-name {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 20px;
}

.profile-role {
  margin-top: 3px;
  color: #b7c3bb;
  font-size: 11px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 190px;
}

.profile-card > p:not(.side-title) {
  margin: 17px 0 18px;
  color: #c4cec7;
  font-size: 12px;
  line-height: 1.65;
}

.profile-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profile-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--lime);
  font-size: 11px;
  font-weight: 700;
}

.topic-list {
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
}

.topic {
  padding: 7px 10px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--muted);
  font-size: 11px;
}

.topic:hover {
  border-color: var(--ink);
  color: var(--ink);
}

.archive-list {
  border-top: 1px solid var(--line);
}

.archive-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 11px 0;
  border-bottom: 1px solid var(--line);
  color: var(--muted);
  font-size: 12px;
}

.archive-row strong {
  color: var(--ink);
  font-size: 11px;
}

a.archive-row:hover span {
  color: var(--ink);
}

.newsletter {
  padding: 20px;
  background: var(--mint);
}

.newsletter h3 {
  margin: 0;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 24px;
  font-weight: 500;
}

.newsletter p {
  margin: 9px 0 16px;
  color: #55635e;
  font-size: 11px;
  line-height: 1.6;
}

.newsletter-form {
  display: flex;
  border-bottom: 1px solid #849992;
}

.newsletter-form input {
  flex: 1;
  min-width: 0;
  padding: 9px 0;
  border: 0;
  outline: 0;
  background: transparent;
  font-size: 11px;
}

.newsletter-form button {
  padding: 8px 0 8px 9px;
  color: var(--ink);
  background: transparent;
  font-size: 11px;
  font-weight: 800;
}

.newsletter-status {
  min-height: 18px;
  margin-top: 8px;
  color: var(--coral-dark);
  font-size: 10px;
}

@media (max-width: 880px) {
  .hero-section {
    grid-template-columns: 1fr;
    padding-top: 37px;
  }

  .hero-image,
  .hero-image img {
    min-height: 340px;
  }

  .content-grid {
    grid-template-columns: 1fr;
    gap: 56px;
  }

  .sidebar {
    display: grid;
    grid-template-columns: 1fr 1fr;
    align-items: start;
  }

  .sidebar .profile-card,
  .sidebar .newsletter {
    grid-column: 1 / -1;
  }
}

@media (max-width: 600px) {
  .hero-section {
    padding-bottom: 36px;
  }

  .hero-image,
  .hero-image img {
    min-height: 275px;
  }

  .filters {
    flex-wrap: wrap;
  }

  .search-box {
    order: -1;
    width: 100%;
    margin-bottom: 3px;
  }

  .feature-card {
    grid-template-columns: 1fr;
  }

  .feature-cover,
  .feature-cover img,
  .feature-cover .cover-fallback {
    height: 190px;
    min-height: 190px;
  }

  .feature-body {
    padding: 22px;
  }

  .feature-body h3 {
    font-size: 26px;
  }

  .section-head {
    align-items: flex-start;
    gap: 12px;
  }

  .section-head h2 {
    font-size: 27px;
  }

  .sidebar {
    display: flex;
  }
}
</style>
