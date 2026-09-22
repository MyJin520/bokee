/** 日期 / 数字展示工具 */

export function formatDate(value: string | Date, withTime = false): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleDateString(
    'zh-CN',
    withTime
      ? { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' }
      : { year: 'numeric', month: 'long', day: 'numeric' },
  )
}

export function formatShortDate(value: string | Date): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

/** 归档用：YYYY 年 MM 月 */
export function formatMonth(value: string | Date): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return `${date.getFullYear()} 年 ${String(date.getMonth() + 1).padStart(2, '0')} 月`
}

export function formatNumber(value: number): string {
  return value.toLocaleString('en-US')
}

/** 取标题首字符作为无封面时的占位 */
export function initialOf(text: string): string {
  return (text || '?').trim().slice(0, 1).toUpperCase()
}
