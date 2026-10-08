import { http } from './request'
import type { ActionType, PageData, TargetType, UserActionListItem } from '@/types'

/**
 * 用户操作（点赞 / 收藏 / 关注）API
 * 对应后端 /pri/user_action 模块，均需登录
 */

/** 创建用户操作（点赞、收藏文章、关注作者） */
export function createUserAction(data: {
  userId: number
  targetId: number
  actionType: ActionType
  targetType: TargetType
}) {
  return http.post<unknown>('/pri/user_action/create', data)
}

/** 删除用户操作（需传入操作记录 ID） */
export function deleteUserAction(actionId: number) {
  return http.delete<unknown>('/pri/user_action/delete', { id: actionId })
}

/** 分页获取当前用户的某类操作列表 */
export function listUserActions(params: {
  actionType: ActionType
  targetType: TargetType
  page?: number
  pageSize?: number
}) {
  return http.post<PageData<UserActionListItem>>('/pri/user_action/list', params)
}
