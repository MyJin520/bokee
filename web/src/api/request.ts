/**
 * 基于 fetch 的 HTTP 请求工具
 * 自动处理 JWT Token 注入、FormData 上传、响应解析、错误处理与登录态失效跳转
 */

const BASE_URL = ''

interface RequestConfig {
  method?: string
  headers?: Record<string, string>
  body?: unknown
  params?: Record<string, string | number | undefined>
}

/** 业务约定的未认证/无权限码（后端始终以 HTTP 200 返回，错误码在响应体 code 中） */
const CODE_UNAUTHORIZED = 401
const CODE_FORBIDDEN = 403

export class ApiError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}

/** 登录态失效：清理本地令牌并跳转登录页（避免在登录页重复跳转） */
function handleSessionExpired() {
  localStorage.removeItem('token')
  localStorage.removeItem('userName')
  const path = window.location.pathname
  if (path !== '/login' && path !== '/register') {
    const redirect = encodeURIComponent(`${path}${window.location.search}`)
    window.location.href = `/login?redirect=${redirect}`
  }
}

class HttpClient {
  private getToken(): string | null {
    return localStorage.getItem('token')
  }

  private buildUrl(path: string, params?: Record<string, string | number | undefined>): string {
    const url = new URL(`${BASE_URL}${path}`, window.location.origin)

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          url.searchParams.set(key, String(value))
        }
      })
    }

    return url.toString()
  }

  private async request<T>(path: string, config: RequestConfig = {}): Promise<T> {
    const { method = 'GET', body, params, headers = {} } = config

    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const fetchOptions: RequestInit = {
      method,
      headers: { ...headers },
    }

    // FormData（文件上传）由浏览器自动设置 Content-Type（含 boundary），不能手动指定
    if (body && method !== 'GET') {
      if (body instanceof FormData) {
        fetchOptions.body = body
      } else {
        fetchOptions.headers = { 'Content-Type': 'application/json', ...headers }
        fetchOptions.body = JSON.stringify(body)
      }
    }

    let response: Response
    try {
      response = await fetch(this.buildUrl(path, params), fetchOptions)
    } catch {
      throw new ApiError(-1, '网络连接异常，请检查后端服务是否可用')
    }

    let result: { code?: number; data?: T; msg?: string }
    try {
      result = await response.json()
    } catch {
      throw new ApiError(response.status, '服务响应异常，请稍后重试')
    }

    const code = result.code ?? -1
    if (code !== 0) {
      const message = result.msg || '请求失败'
      if (code === CODE_UNAUTHORIZED) {
        handleSessionExpired()
      }
      throw new ApiError(code, code === CODE_FORBIDDEN ? `${message}（当前账号可能缺少角色权限）` : message)
    }

    return result.data as T
  }

  get<T>(path: string, params?: Record<string, string | number | undefined>): Promise<T> {
    return this.request<T>(path, { method: 'GET', params })
  }

  post<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, { method: 'POST', body })
  }

  put<T>(path: string, body?: unknown): Promise<T> {
    return this.request<T>(path, { method: 'PUT', body })
  }

  delete<T>(path: string, params?: Record<string, string | number | undefined>): Promise<T> {
    return this.request<T>(path, { method: 'DELETE', params })
  }
}

export const http = new HttpClient()
