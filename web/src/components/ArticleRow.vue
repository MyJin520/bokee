<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { ArticleListItem } from '@/types'
import { useReactions } from '@/composables/useReactions'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import { formatShortDate, formatNumber, initialOf } from '@/utils/format'

const props = defineProps<{
  article: ArticleListItem
}>()

const router = useRouter()
const ui = useUiStore()
const userStore = useUserStore()
const { isLiked, isSaved, toggleLiked, toggleSaved, ensureReactions } = useReactions()

const liked = computed(() => isLiked(props.article.id))
const saved = computed(() => isSaved(props.article.id))
const reacting = ref(false)
const coverFailed = ref(false)

onMounted(() => {
  // 同步当前用户的点赞/收藏状态（已加载则跳过）
  ensureReactions().catch(() => {})
})

function openArticle() {
  router.push(`/article/${props.article.id}`)
}

/** 未登录时提示并跳转登录页 */
function requireLogin(): boolean {
  if (userStore.isLoggedIn) return true
  ui.toast('请先登录后再操作')
  router.push({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
  return false
}

async function onLike(event: MouseEvent) {
  event.stopPropagation()
  if (!requireLogin() || reacting.value) return
  reacting.value = true
  try {
    const nowLiked = await toggleLiked(props.article.id)
    ui.toast(nowLiked ? '已点赞这篇文章' : '已取消点赞')
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '操作失败，请稍后重试')
  } finally {
    reacting.value = false
  }
}

async function onSave(event: MouseEvent) {
  event.stopPropagation()
  if (!requireLogin() || reacting.value) return
  reacting.value = true
  try {
    const nowSaved = await toggleSaved(props.article.id)
    ui.toast(nowSaved ? '已收藏这篇文章' : '已取消收藏')
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '操作失败，请稍后重试')
  } finally {
    reacting.value = false
  }
}
</script>

<template>
  <article class="article-row" @click="openArticle">
    <div class="row-cover">
      <img v-if="article.cover && !coverFailed" :src="article.cover" :alt="article.title" loading="lazy" @error="coverFailed = true" />
      <div v-else class="cover-fallback small">{{ initialOf(article.title) }}</div>
    </div>
    <div class="row-info">
      <h3>
        <span v-if="article.isTop" class="top-mark">置顶</span>
        {{ article.title }}
      </h3>
      <p>{{ article.summary || '暂无摘要，点击查看文章详情。' }}</p>
      <div class="row-meta">
        <span>{{ formatShortDate(article.createdAt) }}</span>
        <span>阅读 {{ formatNumber(article.viewCount) }}</span>
        <span>点赞 {{ formatNumber(article.likeCount) }}</span>
      </div>
    </div>
    <div class="row-actions" @click.stop>
      <slot name="actions">
        <button
          :class="['row-action', { saved }]"
          type="button"
          :disabled="reacting"
          :title="saved ? '取消收藏' : '收藏'"
          :aria-label="saved ? '取消收藏' : '收藏文章'"
          @click="onSave"
        >
          {{ saved ? '♥' : '♡' }}
        </button>
        <button
          :class="['row-action', { liked }]"
          type="button"
          :disabled="reacting"
          :title="liked ? '取消点赞' : '点赞'"
          :aria-label="liked ? '取消点赞' : '点赞文章'"
          @click="onLike"
        >
          ♥ <span>{{ formatNumber(article.likeCount) }}</span>
        </button>
      </slot>
    </div>
  </article>
</template>

<style scoped>
.article-row {
  display: grid;
  grid-template-columns: 116px minmax(0, 1fr) auto;
  gap: 20px;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid var(--line);
  cursor: pointer;
}

.article-row:hover .row-cover img {
  transform: scale(1.05);
}

.row-cover {
  width: 116px;
  height: 87px;
  overflow: hidden;
  background: #dae2dd;
}

.row-cover img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.cover-fallback {
  width: 100%;
  height: 100%;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, var(--blue), var(--lime));
  color: rgba(24, 34, 31, 0.72);
  font-family: 'Playfair Display', Georgia, serif;
}

.cover-fallback.small {
  font-size: 38px;
}

.row-info {
  min-width: 0;
}

.row-info h3 {
  margin: 0 0 7px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 21px;
  line-height: 1.18;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.top-mark {
  display: inline-block;
  margin-right: 8px;
  padding: 2px 6px;
  background: var(--coral-tint);
  color: var(--coral-dark);
  font-family: 'DM Mono', monospace;
  font-size: 9px;
  font-weight: 500;
  letter-spacing: 0.05em;
  vertical-align: middle;
}

.row-info p {
  max-width: 540px;
  margin: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.55;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.row-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 11px;
  margin-top: 11px;
  color: var(--muted);
  font-size: 10px;
}

.row-actions {
  display: flex;
  align-items: center;
  gap: 5px;
}

.row-action {
  min-width: 34px;
  height: 32px;
  padding: 0 7px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: var(--muted);
  background: transparent;
  font-size: 11px;
}

.row-action:hover,
.row-action.liked,
.row-action.saved {
  color: var(--coral-dark);
}

.row-action:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.row-action.saved {
  background: var(--coral-tint);
}

@media (max-width: 600px) {
  .article-row {
    grid-template-columns: 84px minmax(0, 1fr);
    gap: 13px;
  }

  .row-cover {
    width: 84px;
    height: 76px;
  }

  .cover-fallback.small {
    font-size: 30px;
  }

  .row-info h3 {
    font-size: 18px;
  }

  .row-info p {
    -webkit-line-clamp: 1;
  }

  .row-actions {
    grid-column: 2;
    justify-content: flex-start;
    margin-top: -5px;
  }
}
</style>
