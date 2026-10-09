import { http } from './request'
import type { UploadedFile } from '@/types'

/**
 * 上传文件（支持多文件）
 * 后端接口：POST /pri/file/uploads（私有路由，需登录），multipart 字段名为 files
 * 允许类型：jpg/png/gif/webp/pdf/docx/xlsx，单文件最大 32MB
 */
export function uploadFiles(files: File[] | FileList) {
  const formData = new FormData()
  Array.from(files).forEach((file) => formData.append('files', file))
  return http.post<UploadedFile[]>('/pri/file/uploads', formData)
}
