<script setup lang="ts">
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()
</script>

<template>
  <div class="toast-stack" role="status" aria-live="polite">
    <div
      v-for="item in ui.toasts"
      :key="item.id"
      :class="['toast-item', { 'toast-error': item.type === 'error' }]"
      @click="ui.dismissToast(item.id)"
    >
      {{ item.message }}
    </div>
  </div>
</template>
