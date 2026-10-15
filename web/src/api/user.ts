import { http } from './request'
import type { AuthorInfo, JwtResp, PageData, UserInfo } from '@/types'

/** 用户登录 */
export function login(userName: string, password: string) {
  return http.post<JwtResp>('/pub/user/login', { name: userName, password })
}

/** 用户注册 */
export function register(data: {
  name: string
  password: string
  phone?: string
  email?: string
}) {
  return http.post<string>('/pub/user/create', data)
}

/** 获取用户信息（需登录）；传入 id 时查询指定用户，未传时返回当前登录用户 */
export function getUserInfo(id?: number) {
  return http.get<UserInfo>('/pri/user/get_info', id ? { id } : undefined)
}

/** 更新用户信息（需登录） */
export function updateUserInfo(data: {
  name?: string
  phone?: string
  email?: string
  avatar?: string
}) {
  return http.put<string>('/pri/user/update', data)
}

/** 用户登出（需登录） */
export function logout() {
  return http.get<string>('/pri/user/logout')
}

/** 修改密码（需登录，普通用户需提供旧密码） */
export function changePassword(data: { oldPassword: string; newPassword: string }) {
  return http.post<string>('/pri/user/forget_password', data)
}

/**
 * 查询其他用户（需登录，普通用户即可调用）
 * 仅支持用户名模糊搜索，返回脱敏的用户简要信息
 */
export function findUsers(params: { name?: string; page?: number; pageSize?: number }) {
  // 后端以 JSON 请求体绑定参数，缺省字段需显式传空串，避免请求体解析失败
  return http.post<PageData<AuthorInfo>>('/pri/user/find_users', {
    name: params.name ?? '',
    page: params.page ?? 1,
    pageSize: params.pageSize ?? 12,
  })
}

/**
 * 分页获取用户列表（需登录，且当前角色需拥有 /pri/user/list 私有路由权限）
 * 用户名 / 手机号 / 邮箱为模糊搜索，状态为精确匹配
 */
export function getUserList(params: {
  name?: string
  phone?: string
  email?: string
  status?: string
  page?: number
  pageSize?: number
}) {
  // 后端以 JSON 请求体绑定参数，缺省字段需显式传空串，避免请求体解析失败
  return http.post<PageData<UserInfo>>('/pri/user/list', {
    name: params.name ?? '',
    phone: params.phone ?? '',
    email: params.email ?? '',
    status: params.status ?? '',
    page: params.page ?? 1,
    pageSize: params.pageSize ?? 10,
  })
}
