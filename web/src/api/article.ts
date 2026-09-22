import { http } from './request'
import type { ArticleInfo, ArticleListItem, PageData } from '@/types'

/** 获取公开文章列表 */
export function getArticleList(params: {
  title?: string
  userName?: string
  summary?: string
  page?: number
  pageSize?: number
}) {
  return http.post<PageData<ArticleListItem>>('/pub/article/list', params)
}

/** 获取文章详情 */
export function getArticleInfo(id: number) {
  return http.get<ArticleInfo>('/pub/article/get_info', { id })
}

/** 获取指定用户的文章列表 */
export function getUserArticleList(userId: number, page = 1, pageSize = 10) {
  return http.post<PageData<ArticleListItem>>('/pub/article/list_by_user', {
    userId,
    page,
    pageSize,
  })
}

/** 创建文章（需登录） */
export function createArticle(data: {
  title: string
  content: string
  summary?: string
  cover?: string
}) {
  return http.post<string>('/pri/article/create', data)
}

/** 更新文章（需登录） */
export function updateArticle(data: {
  id: number
  title?: string
  content?: string
  summary?: string
  cover?: string
  isTop?: boolean
}) {
  return http.put<string>('/pri/article/update', data)
}

/** 删除文章（需登录） */
export function deleteArticle(id: number) {
  return http.delete<string>('/pri/article/delete', { id })
}
