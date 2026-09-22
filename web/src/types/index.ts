/** 通用 API 响应结构 */
export interface ApiResponse<T = unknown> {
  code: number
  data: T
  msg: string
}

/** 分页数据结构 */
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

/** 文章列表项 */
export interface ArticleListItem {
  id: number
  title: string
  summary: string
  cover: string
  viewCount: number
  likeCount: number
  isTop: boolean
  userId: number
  createdAt: string
}

/** 文章作者公开信息 */
export interface AuthorInfo {
  id: number
  name: string
  avatar: string
}

/** 文章详情 */
export interface ArticleInfo {
  id: number
  title: string
  content: string
  summary: string
  cover: string
  viewCount: number
  likeCount: number
  isTop: boolean
  userId: number
  author: AuthorInfo
  createdAt: string
  updatedAt: string
}

/** 登录响应 */
export interface JwtResp {
  token: string
  name: string
  avatar: string
  email: string
  phone: string
}

/** 用户信息 */
export interface UserInfo {
  id: number
  name: string
  phone: string
  email: string
  status: string
  avatar: string
  roles: UserRole[]
  createdAt: string
  updatedAt: string
}

export interface UserRole {
  id: number
  roleName: string
  roleCode: number
}

/** 文件上传记录 */
export interface UploadedFile {
  id: number
  url: string
  ext: string
  size: number
  originalFileName: string
  hash: string
  createdAt: string
  updatedAt: string
}
