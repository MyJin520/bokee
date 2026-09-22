import { http } from './request'
import type { JwtResp, UserInfo } from '@/types'

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

/** 获取当前用户信息（需登录） */
export function getUserInfo() {
  return http.get<UserInfo>('/pri/user/get_info')
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
