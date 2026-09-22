import { ref } from 'vue'

/**
 * 文章点赞 / 收藏状态
 * 后端当前未提供点赞、收藏接口，与 bokeeUI 原型一致：
 * 状态保存在浏览器 localStorage（按文章 ID 记录），点赞数在前端乐观增减
 */
const LIKED_KEY = 'bokee:liked'
const SAVED_KEY = 'bokee:saved'

function readStorage(key: string): number[] {
  try {
    const raw = localStorage.getItem(key)
    const parsed = raw ? JSON.parse(raw) : []
    return Array.isArray(parsed) ? parsed.filter((n) => typeof n === 'number') : []
  } catch {
    return []
  }
}

const likedIds = ref<number[]>(readStorage(LIKED_KEY))
const savedIds = ref<number[]>(readStorage(SAVED_KEY))

function persist() {
  localStorage.setItem(LIKED_KEY, JSON.stringify(likedIds.value))
  localStorage.setItem(SAVED_KEY, JSON.stringify(savedIds.value))
}

export function useReactions() {
  function isLiked(id: number) {
    return likedIds.value.includes(id)
  }

  function isSaved(id: number) {
    return savedIds.value.includes(id)
  }

  /** 切换点赞，返回切换后的状态 */
  function toggleLiked(id: number): boolean {
    if (isLiked(id)) {
      likedIds.value = likedIds.value.filter((item) => item !== id)
    } else {
      likedIds.value = [...likedIds.value, id]
    }
    persist()
    return isLiked(id)
  }

  /** 切换收藏，返回切换后的状态 */
  function toggleSaved(id: number): boolean {
    if (isSaved(id)) {
      savedIds.value = savedIds.value.filter((item) => item !== id)
    } else {
      savedIds.value = [...savedIds.value, id]
    }
    persist()
    return isSaved(id)
  }

  return { likedIds, savedIds, isLiked, isSaved, toggleLiked, toggleSaved }
}
