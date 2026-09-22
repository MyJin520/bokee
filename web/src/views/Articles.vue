<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleList } from '@/api/article'
import type { ArticleListItem } from '@/types'
import ArticleRow from '@/components/ArticleRow.vue'
import { formatNumber } from '@/utils/format'

type FilterKey = 'all' | 'top' | 'covered'

const route = useRoute()

const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 8
const loading = ref(true)
const errorMessage = ref('')
const keyword = ref('')
const activeFilter = ref<FilterKey>('all')
const authorName = ref('')
const searchInput = ref<HTMLInputElement | null>(null)
let searchTimer: ReturnType<typeof setTimeout> | undefined

const filteredArticles = computed(() =>
  articles.value.filter((article) => {
    const kw = keyword.value.trim().toLowerCase()
    const matchesKeyword =
      !kw || `${article.title} ${article.summary || ''}`.toLowerCase().includes(kw)
    const matchesFilter =
      activeFilter.value === 'all' ||
      (activeFilter.value === 'top' ? article.isTop : Boolean(article.cover))
    return matchesKeyword && matchesFilter
  }),
)

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

async function loadArticles() {
  loading.value = true
  errorMessage.value = ''
  try {
    // 后端列表支持标题/作者名检索；置顶、有封面筛选在当前页客户端过滤
    const res = await getArticleList({
      page: page.value,
      pageSize,
      title: keyword.value.trim() || undefined,
      userName: authorName.value || undefined,
    })
    articles.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '文章加载失败，请稍后重试。'
    articles.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

watch(keyword, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    loadArticles()
  }, 350)
})

watch(activeFilter, () => {
  page.value = 1
  loadArticles()
})

function changePage(next: number) {
  if (next < 1 || next > totalPages.value || next === page.value) return
  page.value = next
  loadArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function clearAuthor() {
  authorName.value = ''
  page.value = 1
  loadArticles()
}

onMounted(() => {
  if (typeof route.query.user === 'string') authorName.value = route.query.user
  loadArticles()
  if (route.query.focus === '1') {
    window.setTimeout(() => searchInput.value?.focus(), 400)
  }
})
</script>

<template>
  <div class="archive-page">
    <header class="archive-header">
      <span class="eyebrow">All articles</span>
      <h1>全部文章</h1>
      <p>所有公开文章按时间排列，共 {{ formatNumber(total) }} 篇。慢慢看，也可以搜索。</p>
    </header>

    <div class="toolbar">
      <div class="chips" role="group" aria-label="文章筛选">
        <button :class="['chip', { active: activeFilter === 'all' }]" type="button" @click="activeFilter = 'all'">全部</button>
        <button :class="['chip', { active: activeFilter === 'top' }]" type="button" @click="activeFilter = 'top'">置顶</button>
        <button :class="['chip', { active: activeFilter === 'covered' }]" type="button" @click="activeFilter = 'covered'">有封面</button>
      </div>
      <label class="search-box">
        <input
          ref="searchInput"
          v-model="keyword"
          type="search"
          placeholder="搜索标题或摘要"
          aria-label="搜索文章"
        />
        <span>⌕</span>
      </label>
    </div>

    <div v-if="authorName" class="author-filter">
      正在浏览 <strong>{{ authorName }}</strong> 的文章
      <button class="text-button" type="button" style="margin-left: 8px" @click="clearAuthor">清除 ×</button>
    </div>

    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取文章…</p>
    </div>
    <div v-else-if="errorMessage" class="state-block">
      <p class="state-title">加载遇到问题</p>
      <p>{{ errorMessage }}</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="loadArticles">重新加载 →</button>
    </div>
    <template v-else-if="filteredArticles.length">
      <div class="archive-list">
        <ArticleRow v-for="article in filteredArticles" :key="article.id" :article="article" />
      </div>
      <div class="pagination">
        <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">← 上一页</button>
        <span class="page-info">{{ page }} / {{ totalPages }}</span>
        <button class="page-btn" type="button" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页 →</button>
      </div>
    </template>
    <div v-else class="state-block">
      <p class="state-title">没有找到文章</p>
      <p>换一个关键词，或清除筛选条件试试。</p>
      <button
        class="text-button"
        type="button"
        style="margin-top: 12px"
        @click="keyword = ''; activeFilter = 'all'; authorName = ''"
      >
        清除筛选 →
      </button>
    </div>
  </div>
</template>

<style scoped>
.archive-page {
  padding: 54px 0 80px;
}

.archive-header {
  margin-bottom: 34px;
}

.archive-header h1 {
  margin: 18px 0 12px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(38px, 5vw, 62px);
  font-weight: 500;
  letter-spacing: -0.03em;
}

.archive-header p {
  max-width: 520px;
  margin: 0;
  color: var(--muted);
  font-size: 14px;
  line-height: 1.7;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  margin-bottom: 8px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--line);
}

.chips {
  display: flex;
  gap: 7px;
  overflow-x: auto;
}

.search-box {
  display: flex;
  align-items: center;
  width: 230px;
  flex-shrink: 0;
  border-bottom: 1px solid var(--line);
}

.search-box input {
  width: 100%;
  padding: 8px 0;
  border: 0;
  outline: 0;
  background: transparent;
  font-size: 12px;
}

.search-box span {
  color: var(--muted);
}

.author-filter {
  margin: 14px 0;
  padding: 10px 14px;
  background: var(--coral-tint);
  color: var(--coral-dark);
  font-size: 12px;
}

.archive-list {
  border-top: 1px solid var(--line);
}

@media (max-width: 600px) {
  .archive-page {
    padding-top: 34px;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-box {
    width: 100%;
  }
}
</style>
