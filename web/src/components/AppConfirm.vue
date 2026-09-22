<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()

function onKeydown(event: KeyboardEvent) {
  if (!ui.confirmState.open) return
  if (event.key === 'Escape') ui.answerConfirm(false)
  if (event.key === 'Enter') ui.answerConfirm(true)
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <div v-if="ui.confirmState.open" class="modal-backdrop" @click.self="ui.answerConfirm(false)">
      <div class="modal-panel" role="dialog" aria-modal="true">
        <h3 class="modal-title">{{ ui.confirmState.title }}</h3>
        <p class="modal-text">{{ ui.confirmState.message }}</p>
        <div class="modal-actions">
          <button class="btn btn-outline btn-sm" type="button" @click="ui.answerConfirm(false)">
            {{ ui.confirmState.cancelText }}
          </button>
          <button
            :class="['btn', 'btn-sm', ui.confirmState.danger ? 'btn-coral' : '']"
            type="button"
            @click="ui.answerConfirm(true)"
          >
            {{ ui.confirmState.confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
