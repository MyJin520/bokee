import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface ToastItem {
  id: number
  message: string
  type: 'info' | 'error'
}

export interface ConfirmOptions {
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  danger?: boolean
}

interface ConfirmState extends ConfirmOptions {
  open: boolean
  resolve?: (value: boolean) => void
}

/**
 * 全局交互反馈：轻提示 Toast 与模态确认框 Confirm
 * 在组件中使用：const ui = useUiStore(); ui.toast('已保存'); await ui.confirm({...})
 */
export const useUiStore = defineStore('ui', () => {
  const toasts = ref<ToastItem[]>([])
  const confirmState = ref<ConfirmState>({ open: false, message: '' })
  /** 首页搜索框聚焦信号：点击头部搜索按钮时自增，首页监听到后滚动并聚焦 */
  const homeSearchSignal = ref(0)
  let toastSeq = 0

  function focusHomeSearch() {
    homeSearchSignal.value++
  }

  function toast(message: string, type: 'info' | 'error' = 'info', duration = 2000) {
    const id = ++toastSeq
    toasts.value.push({ id, message, type })
    window.setTimeout(() => dismissToast(id), duration)
  }

  function toastError(message: string) {
    toast(message, 'error', 2800)
  }

  function dismissToast(id: number) {
    toasts.value = toasts.value.filter((item) => item.id !== id)
  }

  function confirm(options: ConfirmOptions): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.value = {
        open: true,
        title: options.title || '请确认',
        message: options.message,
        confirmText: options.confirmText || '确认',
        cancelText: options.cancelText || '取消',
        danger: options.danger ?? false,
        resolve,
      }
    })
  }

  function answerConfirm(value: boolean) {
    confirmState.value.resolve?.(value)
    confirmState.value = { open: false, message: '' }
  }

  return {
    toasts,
    confirmState,
    homeSearchSignal,
    focusHomeSearch,
    toast,
    toastError,
    dismissToast,
    confirm,
    answerConfirm,
  }
})
