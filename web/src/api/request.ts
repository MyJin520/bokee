/**
 * 基于 fetch 的 HTTP 请求工具
 * 自动处理 JWT Token 注入、响应解析、错误处理
 */

const BASE_URL = ''

interface RequestConfig {
  method?: string
  headers?: Record<string, string>
  body?: unknown
  params?: Record<string, string | number | undefined>
}

class HttpClient {
  private getToken(): string | null {
    return localStorage.getItem('token')
  }

  private buildUrl(path: string, params?: Record<string, string | number | undefined>): string {
    const url = new URL(`${BASE_URL}${path}`, window.location.origin)

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
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
      headers: {
        'Content-Type': 'application/json',
        ...headers,
      },
    }

    if (body && method !== 'GET') {
      fetchOptions.body = JSON.stringify(body)
    }

    const url = this.buildUrl(path, params)
    const response = await fetch(url, fetchOptions)
    const result = await response.json()

    if (result.code !== 0) {
      throw new Error(result.msg || '请求失败')
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
