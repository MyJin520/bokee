<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { findUsers } from '@/api/user'
import { useReactions } from '@/composables/useReactions'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import type { UserBrief } from '@/types'
import { formatNumber, initialOf } from '@/utils/format'

const userStore = useUserStore()
const ui = useUiStore()
const { isFollowed, toggleFollow, ensureReactions } = useReactions()

const pageSize = 12

const authors = ref<UserBrief[]>([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const reactingId = ref(0)

const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

async function load() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await findUsers({ name: keyword.value.trim(), page: page.value, pageSize })
    authors.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    errorMessage.value = error instanceof Error ? error.message : '作者加载失败，请稍后重试。'
    authors.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function clearKeyword() {
  keyword.value = ''
  search()
}

function changePage(next: number) {
  if (next < 1 || next > pages.value || next === page.value) return
  page.value = next
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

/** 关注 / 取消关注作者；本人不可关注自己 */
async function onFollow(author: UserBrief) {
  if (author.id === userStore.userInfo?.id || reactingId.value) return
  reactingId.value = author.id
  try {
    const nowFollowing = await toggleFollow(author.id)
    // 本地同步粉丝数：关注 +1，取消关注 -1
    author.followersCount += nowFollowing ? 1 : -1
    ui.toast(nowFollowing ? '已关注 ' + author.name : '已取消关注 ' + author.name)
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '操作失败，请稍后重试')
  } finally {
    reactingId.value = 0
  }
}

onMounted(async () => {
  // 页面直接刷新进入时，等待根组件恢复登录态后再加载
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  // 先同步关注状态，避免关注映射未加载时误判为「未关注」
  await ensureReactions().catch(() => {})
  load()
})
</script>

<template>
  <div class="authors-page">
    <header class="page-header">
      <span class="eyebrow">Discover · authors</span>
      <h1>发现作者</h1>
      <p>共 {{ formatNumber(total) }} 位作者，搜索名字即可进入 TA 的主页，关注后第一时间读到新文章。</p>
    </header>

    <form class="search-bar" @submit.prevent="search">
      <label class="search-field">
        <span class="field-label">作者名字</span>
        <input v-model="keyword" class="field-input" type="search" placeholder="模糊搜索作者名字" />
      </label>
      <div class="search-actions">
        <button class="btn btn-coral btn-sm" type="submit" :disabled="loading">搜索</button>
        <button class="btn btn-outline btn-sm" type="button" :disabled="loading" @click="clearKeyword">
          重置
        </button>
      </div>
    </form>

    <div v-if="loading" class="state-block">
      <div class="spinner"></div>
      <p>正在查找作者…</p>
    </div>
    <div v-else-if="errorMessage" class="state-block">
      <p class="state-title">加载遇到问题</p>
      <p>{{ errorMessage }}</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="load">重新加载 →</button>
    </div>
    <template v-else-if="authors.length">
      <ul class="author-grid">
        <li v-for="author in authors" :key="author.id" class="author-card surface-card">
          <router-link class="author-main" :to="{ name: 'user-profile', params: { id: author.id } }">
            <span class="author-avatar">
              <img v-if="author.avatar" :src="author.avatar" :alt="author.name" />
              <template v-else>{{ initialOf(author.name) }}</template>
            </span>
            <span class="author-name">{{ author.name }}</span>
            <span class="author-fans">{{ formatNumber(author.followersCount) }} 位粉丝</span>
            <span class="author-hint">查看主页与文章 →</span>
          </router-link>
          <button
            v-if="author.id !== userStore.userInfo?.id"
            :class="['follow-btn', { active: isFollowed(author.id) }]"
            type="button"
            :disabled="reactingId === author.id"
            @click="onFollow(author)"
          >
            {{ isFollowed(author.id) ? '✓ 已关注' : '＋ 关注' }}
          </button>
          <span v-else class="self-hint">这是你自己</span>
        </li>
      </ul>
      <div class="pagination">
        <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">← 上一页</button>
        <span class="page-info">{{ page }} / {{ pages }}</span>
        <button class="page-btn" type="button" :disabled="page >= pages" @click="changePage(page + 1)">
          下一页 →
        </button>
      </div>
    </template>
    <div v-else class="state-block empty-state">
      <p class="state-title">没有找到作者</p>
      <p>换一个名字试试，或清空搜索看看全部作者。</p>
      <button class="text-button" type="button" style="margin-top: 12px" @click="clearKeyword">查看全部 →</button>
    </div>
  </div>
</template>

<style scoped>
.authors-page {
  padding: 48px 0 80px;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 14px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(34px, 4.6vw, 54px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.page-header p {
  max-width: 560px;
  margin: 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.7;
}

.search-bar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
  gap: 18px;
  padding: 22px 24px;
  border: 1px solid var(--line);
  background: var(--surface);
}

.search-actions {
  display: flex;
  gap: 8px;
}

.author-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
  gap: 18px;
  list-style: none;
  margin: 22px 0 0;
  padding: 0;
}

.author-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 26px 20px 22px;
}

.author-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  text-align: center;
}

.author-avatar {
  width: 62px;
  height: 62px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--ink);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 26px;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-name {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 20px;
  font-weight: 500;
}

.author-fans {
  color: var(--muted);
  font-size: 12px;
}

.author-main:hover .author-name {
  color: var(--coral-dark);
}

.author-hint {
  color: var(--muted);
  font-size: 11px;
}

.follow-btn {
  width: 100%;
  padding: 9px 12px;
  border: 1px solid var(--line);
  background: transparent;
  color: var(--ink);
  font-size: 11px;
  font-weight: 700;
}

.follow-btn:hover:not(:disabled) {
  border-color: var(--coral-dark);
  color: var(--coral-dark);
}

.follow-btn.active {
  border-color: var(--coral);
  background: var(--coral-tint);
  color: var(--coral-dark);
}

.follow-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.self-hint {
  color: var(--muted);
  font-size: 11px;
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 22px;
}

@media (max-width: 600px) {
  .search-bar {
    grid-template-columns: minmax(0, 1fr);
  }

  .author-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 12px;
  }
}
</style>
