import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, getUserInfo } from '@/api/user'
import type { UserInfo } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userName = computed(() => userInfo.value?.name || '')

  /** 登录 */
  async function login(userName: string, password: string) {
    const res = await loginApi(userName, password)
    token.value = res.token
    localStorage.setItem('token', res.token)
    localStorage.setItem('userName', res.name)
    // 登录成功后获取用户完整信息
    try {
      userInfo.value = await getUserInfo()
    } catch {
      // 获取用户信息失败不影响登录状态
    }
    return res
  }

  /** 登出 */
  async function doLogout() {
    try {
      await logoutApi()
    } catch {
      // 后端登出失败也清理本地状态
    }
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userName')
  }

  /** 恢复登录态（页面刷新时调用） */
  async function restoreSession() {
    if (!token.value) return
    try {
      userInfo.value = await getUserInfo()
    } catch {
      // Token 过期或无效，清理本地状态
      token.value = ''
      userInfo.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('userName')
    }
  }

  return { token, userInfo, isLoggedIn, userName, login, doLogout, restoreSession }
})
