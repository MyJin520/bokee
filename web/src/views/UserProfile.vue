<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router'
import { getUserArticleList } from '@/api/article'
import { findUsers } from '@/api/user'
import { useReactions } from '@/composables/useReactions'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import type { ArticleListItem, UserBrief } from '@/types'
import ArticleRow from '@/components/ArticleRow.vue'
import { formatDate, formatNumber, initialOf } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const ui = useUiStore()
const { isFollowed, toggleFollow, ensureReactions } = useReactions()

const pageSize = 8

const user = ref<UserBrief | null>(null)
const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(true)
const articlesLoading = ref(false)
const notFound = ref(false)
const errorMessage = ref('')
const reacting = ref(false)

const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const isSelf = computed(() => Boolean(user.value) && user.value?.id === userStore.userInfo?.id)
const followed = computed(() => (user.value ? isFollowed(user.value.id) : false))

async function loadArticles(userId: number) {
  articlesLoading.value = true
  try {
    const res = await getUserArticleList(userId, page.value, pageSize)
    articles.value = res.list
    total.value = res.total
  } catch {
    articles.value = []
    total.value = 0
  } finally {
    articlesLoading.value = false
  }
}

async function load(id: number) {
  if (!id) {
    router.replace({ name: 'authors' })
    return
  }
  loading.value = true
  notFound.value = false
  errorMessage.value = ''
  page.value = 1
  try {
    // 只用 find_users 的公开字段（id / 昵称 / 头像），无需读取他人资料接口
    const res = await findUsers({ id, page: 1, pageSize: 1 })
    const author = res.list[0]
    if (!author) throw new Error('用户不存在')
    user.value = author
    loadArticles(id)
  } catch (error: unknown) {
    user.value = null
    articles.value = []
    total.value = 0
    notFound.value = true
    errorMessage.value = error instanceof Error ? error.message : '用户不存在'
  } finally {
    loading.value = false
  }
}

async function onFollow() {
  if (!user.value || isSelf.value || reacting.value) return
  reacting.value = true
  try {
    const nowFollowing = await toggleFollow(user.value.id)
    // 本地同步粉丝数：关注 +1，取消关注 -1（数据库已落库，下次进入页面会重新拉取最新值）
    user.value.followersCount += nowFollowing ? 1 : -1
    ui.toast(nowFollowing ? '已关注 ' + user.value.name : '已取消关注 ' + user.value.name)
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '操作失败，请稍后重试')
  } finally {
    reacting.value = false
  }
}

function changePage(next: number) {
  if (next < 1 || next > pages.value || next === page.value) return
  page.value = next
  if (user.value) loadArticles(user.value.id)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  // 页面直接刷新进入时，等待根组件恢复登录态后再加载
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  ensureReactions().catch(() => {})
  load(Number(route.params.id))
})

onBeforeRouteUpdate((to) => {
  window.scrollTo({ top: 0 })
  load(Number(to.params.id))
})
</script>

<template>
  <div class="user-page">
    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在读取用户信息…</p>
    </div>

    <div v-else-if="notFound" class="state-block">
      <p class="state-title">用户不存在</p>
      <p>{{ errorMessage }}</p>
      <router-link class="text-button" :to="{ name: 'authors' }" style="display: inline-block; margin-top: 14px">
        去发现作者 →
      </router-link>
    </div>

    <template v-else-if="user">
      <header class="user-header">
        <span class="user-avatar">
          <img v-if="user.avatar" :src="user.avatar" :alt="user.name" />
          <template v-else>{{ initialOf(user.name) }}</template>
        </span>
        <div class="user-intro">
          <span class="eyebrow">Author profile</span>
          <h1>{{ user.name }}</h1>
          <p class="user-meta">
            {{ formatDate(user.createdAt) }} 加入 · 共 {{ formatNumber(total) }} 篇公开文章 ·
            {{ formatNumber(user.followersCount) }} 位粉丝
          </p>
          <div class="user-actions">
            <router-link v-if="isSelf" class="btn btn-outline btn-sm" :to="{ name: 'profile' }">
              编辑个人资料
            </router-link>
            <button
              v-else
              :class="['btn', 'btn-sm', followed ? 'btn-outline' : 'btn-coral']"
              type="button"
              :disabled="reacting"
              @click="onFollow"
            >
              {{ followed ? '✓ 已关注' : '＋ 关注 TA' }}
            </button>
            <router-link class="btn btn-outline btn-sm" :to="{ name: 'articles' }">浏览全部文章</router-link>
          </div>
        </div>
      </header>

      <section class="user-articles">
        <p class="side-title">Articles · TA 的文章</p>
        <div v-if="articlesLoading" class="state-block">
          <div class="spinner"></div>
          <p>正在读取文章…</p>
        </div>
        <template v-else-if="articles.length">
          <div class="user-list">
            <ArticleRow v-for="article in articles" :key="article.id" :article="article" />
          </div>
          <div class="pagination">
            <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">
              ← 上一页
            </button>
            <span class="page-info">{{ page }} / {{ pages }}</span>
            <button class="page-btn" type="button" :disabled="page >= pages" @click="changePage(page + 1)">
              下一页 →
            </button>
          </div>
        </template>
        <div v-else class="state-block empty-state">
          <p class="state-title">TA 还没有公开文章</p>
          <p>过段时间再来看看，或者先去发现其他作者。</p>
          <router-link class="btn btn-coral" :to="{ name: 'authors' }" style="margin-top: 18px">
            去发现作者 →
          </router-link>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.user-page {
  padding: 48px 0 80px;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 26px;
  padding-bottom: 30px;
  border-bottom: 1px solid var(--line);
}

.user-avatar {
  width: 96px;
  height: 96px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--ink);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 40px;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-intro {
  min-width: 0;
}

.user-intro h1 {
  margin: 12px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(30px, 4vw, 46px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.user-meta {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.user-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 18px;
}

.user-articles {
  margin-top: 34px;
}

.user-list {
  border-top: 1px solid var(--line);
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 10px;
}

@media (max-width: 600px) {
  .user-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }

  .user-avatar {
    width: 76px;
    height: 76px;
    font-size: 32px;
  }
}
</style>
