import { ref, type Ref } from 'vue'
import { createUserAction, deleteUserAction, listUserActions } from '@/api/userAction'
import { useUserStore } from '@/stores/user'
import type { ActionType, TargetType, UserActionListItem } from '@/types'

/**
 * 文章点赞 / 收藏 / 作者关注状态（后端 user_action 模块驱动）
 *
 * 后端未提供「查询单个目标操作状态」的接口，因此以列表接口（actionType + targetType）
 * 拉取当前用户的全部操作，建立 targetId -> actionId 映射，供各页面同步状态。
 * 点赞/收藏/关注操作后，再刷新对应类型的映射以拿到新操作的 actionId（删除时需要）。
 */

const PAGE_SIZE = 100 // 后端单页上限 100
const MAX_PAGES = 20 // 最多拉取 2000 条，超出后静默截断（防止异常数据导致请求风暴）

/** 文章ID -> actionId */
const likedActions = ref<Map<number, number>>(new Map())
/** 文章ID -> actionId */
const savedActions = ref<Map<number, number>>(new Map())
/** 作者用户ID -> actionId */
const followedActions = ref<Map<number, number>>(new Map())

let loadedForUserId: number | null = null
let loadingPromise: Promise<void> | null = null

/** 分页拉全某一类型的操作列表 */
async function fetchAllActions(actionType: ActionType, targetType: TargetType): Promise<UserActionListItem[]> {
  const all: UserActionListItem[] = []
  let page = 1
  while (page <= MAX_PAGES) {
    const res = await listUserActions({ actionType, targetType, page, pageSize: PAGE_SIZE })
    all.push(...res.list)
    if (all.length >= res.total || res.list.length === 0) break
    page += 1
  }
  return all
}

/** 列表项 -> targetId => actionId 映射（响应为扁平结构，无需按类型取目标） */
function toActionMap(list: UserActionListItem[]): Map<number, number> {
  const map = new Map<number, number>()
  for (const item of list) {
    map.set(item.targetId, item.actionId)
  }
  return map
}

export function useReactions() {
  const userStore = useUserStore()

  function isLiked(articleId: number) {
    return likedActions.value.has(articleId)
  }

  function isSaved(articleId: number) {
    return savedActions.value.has(articleId)
  }

  function isFollowed(authorUserId: number) {
    return followedActions.value.has(authorUserId)
  }

  /** 拉取当前用户点赞/收藏/关注状态（已加载过则跳过；未登录时清空） */
  async function ensureReactions() {
    const userId = userStore.userInfo?.id
    if (!userId) {
      if (loadedForUserId !== null) resetReactions()
      return
    }
    if (loadedForUserId === userId) return
    if (loadingPromise) return loadingPromise

    loadingPromise = (async () => {
      const [likes, saves, follows] = await Promise.all([
        fetchAllActions('like', 'article'),
        fetchAllActions('bookmark', 'article'),
        fetchAllActions('follow', 'author'),
      ])
      likedActions.value = toActionMap(likes)
      savedActions.value = toActionMap(saves)
      followedActions.value = toActionMap(follows)
      loadedForUserId = userId
    })().finally(() => {
      loadingPromise = null
    })
    return loadingPromise
  }

  function resetReactions() {
    likedActions.value = new Map()
    savedActions.value = new Map()
    followedActions.value = new Map()
    loadedForUserId = null
  }

  function requireUserId(): number {
    const userId = userStore.userInfo?.id
    if (!userId) throw new Error('请先登录后再操作')
    return userId
  }

  /**
   * 切换某个操作：已存在则删除（返回 false），不存在则创建（返回 true）
   * 创建接口不返回 actionId，创建后刷新该类型映射以补齐
   */
  async function toggleAction(
    map: Ref<Map<number, number>>,
    actionType: ActionType,
    targetType: TargetType,
    targetId: number,
  ): Promise<boolean> {
    const userId = requireUserId()
    const existing = map.value.get(targetId)

    if (existing) {
      await deleteUserAction(existing)
      map.value.delete(targetId)
      return false
    }

    await createUserAction({ userId, targetId, actionType, targetType })
    // 先用占位 ID 保持本地状态，随后刷新拿到真实 actionId；刷新失败时占位保留，下次刷新校正
    map.value.set(targetId, -1)
    try {
      map.value = toActionMap(await fetchAllActions(actionType, targetType))
    } catch {
      // 忽略：保持占位状态，后续 ensure/toggle 会校正
    }
    return true
  }

  async function toggleLiked(articleId: number): Promise<boolean> {
    return toggleAction(likedActions, 'like', 'article', articleId)
  }

  async function toggleSaved(articleId: number): Promise<boolean> {
    return toggleAction(savedActions, 'bookmark', 'article', articleId)
  }

  async function toggleFollow(authorUserId: number): Promise<boolean> {
    return toggleAction(followedActions, 'follow', 'author', authorUserId)
  }

  return {
    likedActions,
    savedActions,
    followedActions,
    isLiked,
    isSaved,
    isFollowed,
    ensureReactions,
    resetReactions,
    toggleLiked,
    toggleSaved,
    toggleFollow,
  }
}
