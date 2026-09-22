<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { createArticle, getArticleInfo, updateArticle } from '@/api/article'
import { uploadFiles } from '@/api/file'
import { useUiStore } from '@/stores/ui'
import { renderMarkdown } from '@/utils/markdown'

const props = defineProps<{ articleId?: number }>()
const emit = defineEmits<{ success: [id: number] }>()

const ui = useUiStore()

const isEdit = computed(() => Boolean(props.articleId))
const loading = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const tab = ref<'write' | 'preview'>('write')
const previewHtml = ref('')

const title = ref('')
const summary = ref('')
const cover = ref('')
const content = ref('')
const isTop = ref(false)
const coverBroken = ref(false)
const coverUrlInput = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const errors = ref<{ title?: string; content?: string; cover?: string }>({})

const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
const MAX_SIZE = 32 * 1024 * 1024

async function loadArticle() {
  if (!props.articleId) return
  loading.value = true
  try {
    const data = await getArticleInfo(props.articleId)
    title.value = data.title
    summary.value = data.summary || ''
    cover.value = data.cover || ''
    coverBroken.value = false
    content.value = data.content
    isTop.value = data.isTop
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '文章加载失败')
  } finally {
    loading.value = false
  }
}

function pickCover() {
  fileInput.value?.click()
}

async function onCoverChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  errors.value.cover = undefined
  if (!ALLOWED_TYPES.includes(file.type)) {
    errors.value.cover = '仅支持 JPG / PNG / GIF / WebP 格式'
    return
  }
  if (file.size > MAX_SIZE) {
    errors.value.cover = '图片大小不能超过 32MB'
    return
  }

  uploading.value = true
  try {
    const files = await uploadFiles([file])
    if (files[0]?.url) {
      cover.value = files[0].url
      coverBroken.value = false
      coverUrlInput.value = ''
      ui.toast('封面上传成功')
    } else {
      ui.toastError('上传返回异常，请改用 URL')
    }
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '封面上传失败')
  } finally {
    uploading.value = false
  }
}

function applyCoverUrl() {
  const url = coverUrlInput.value.trim()
  if (url) {
    cover.value = url
    coverBroken.value = false
    ui.toast('已应用封面地址')
  }
}

function onCoverError() {
  coverBroken.value = true
  cover.value = ''
  ui.toastError('封面图片无法加载，请检查地址或重新上传')
}

function removeCover() {
  cover.value = ''
  coverBroken.value = false
  coverUrlInput.value = ''
}

async function switchPreview() {
  if (tab.value === 'preview') {
    previewHtml.value = await renderMarkdown(content.value)
  }
}

function validate(): boolean {
  errors.value = {}
  if (!title.value.trim()) errors.value.title = '请填写文章标题'
  if (!content.value.trim()) errors.value.content = '正文内容不能为空'
  return !errors.value.title && !errors.value.content
}

async function submit() {
  if (!validate()) {
    ui.toastError('请先完善标题与正文')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && props.articleId) {
      await updateArticle({
        id: props.articleId,
        title: title.value.trim(),
        summary: summary.value.trim(),
        cover: cover.value.trim(),
        content: content.value,
        isTop: isTop.value,
      })
      ui.toast('文章更新成功')
      emit('success', props.articleId)
    } else {
      const data = await createArticle({
        title: title.value.trim(),
        summary: summary.value.trim(),
        cover: cover.value.trim(),
        content: content.value,
      })
      const newId = Number(data)
      ui.toast('文章发表成功')
      emit('success', Number.isFinite(newId) && newId > 0 ? newId : 0)
    }
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '保存失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

onMounted(loadArticle)
</script>

<template>
  <div v-if="loading" class="state-block">
    <div class="spinner"></div>
    <p>正在载入文章…</p>
  </div>

  <form v-else class="editor" @submit.prevent="submit">
    <label class="field-label"><span class="required">*</span> 文章标题</label>
    <input
      v-model="title"
      class="field-input title-input serif"
      type="text"
      placeholder="给文章起一个值得回看的标题"
      maxlength="200"
    />
    <p v-if="errors.title" class="form-error">{{ errors.title }}</p>

    <label class="field-label">文章摘要</label>
    <textarea
      v-model="summary"
      class="field-area"
      rows="2"
      placeholder="一到两句话概括这篇文章（可留空，列表页会展示）"
      maxlength="300"
    ></textarea>

    <label class="field-label">封面图片</label>
    <div class="cover-editor">
      <div v-if="cover && !coverBroken" class="cover-preview">
        <img :src="cover" alt="封面预览" @error="onCoverError" />
        <button class="text-button" type="button" @click="removeCover">移除封面 ×</button>
      </div>
      <div v-else class="cover-placeholder">暂无封面</div>
      <div class="cover-actions">
        <button class="btn btn-outline btn-sm" type="button" :disabled="uploading" @click="pickCover">
          {{ uploading ? '上传中…' : '上传封面' }}
        </button>
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/gif,image/webp" hidden @change="onCoverChange" />
        <span class="cover-hint">支持 JPG / PNG / GIF / WebP，≤32MB</span>
      </div>
      <div class="cover-url">
        <input v-model="coverUrlInput" class="field-input" type="url" placeholder="或粘贴图片 URL 后点击应用" />
        <button class="text-button" type="button" @click="applyCoverUrl">应用</button>
      </div>
      <p v-if="errors.cover" class="form-error">{{ errors.cover }}</p>
    </div>

    <div class="content-head">
      <label class="field-label" style="margin: 0"><span class="required">*</span> 正文内容（支持 Markdown）</label>
      <div class="editor-tabs">
        <button :class="['tab', { active: tab === 'write' }]" type="button" @click="tab = 'write'">写作</button>
        <button :class="['tab', { active: tab === 'preview' }]" type="button" @click="tab = 'preview'; switchPreview()">预览</button>
      </div>
    </div>

    <textarea
      v-show="tab === 'write'"
      v-model="content"
      class="field-area content-area mono"
      rows="20"
      placeholder="# 在这里开始书写&#10;&#10;支持 Markdown 语法：标题、列表、引用、代码块、图片、链接…"
    ></textarea>
    <div v-show="tab === 'preview'" class="markdown-body preview-pane" v-html="previewHtml"></div>
    <p v-if="errors.content" class="form-error">{{ errors.content }}</p>

    <label v-if="isEdit" class="top-toggle">
      <input v-model="isTop" type="checkbox" />
      <span>设为置顶文章（将在首页优先展示）</span>
    </label>

    <div class="editor-footer">
      <button class="btn btn-coral" type="submit" :disabled="submitting">
        {{ submitting ? '保存中…' : isEdit ? '保存修改' : '发表文章' }}
      </button>
      <router-link class="btn btn-outline" to="/my-articles">取消</router-link>
    </div>
  </form>
</template>

<style scoped>
.editor {
  max-width: 820px;
}

.title-input {
  padding: 14px 0;
  font-size: 26px;
}

.cover-editor {
  margin-top: 8px;
}

.cover-preview {
  width: 260px;
}

.cover-preview img {
  width: 260px;
  height: 150px;
  display: block;
  object-fit: cover;
  border: 1px solid var(--line);
}

.cover-preview .text-button {
  margin-top: 8px;
}

.cover-placeholder {
  width: 260px;
  height: 150px;
  display: grid;
  place-items: center;
  border: 1px dashed var(--line);
  color: var(--muted);
  font-size: 12px;
}

.cover-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.cover-hint {
  color: var(--muted);
  font-size: 11px;
}

.cover-url {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 10px;
  max-width: 520px;
}

.cover-url .text-button {
  flex-shrink: 0;
  white-space: nowrap;
}

.content-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 22px;
}

.editor-tabs {
  display: flex;
  border: 1px solid var(--line);
}

.tab {
  padding: 7px 14px;
  background: transparent;
  color: var(--muted);
  font-size: 11px;
  font-weight: 700;
}

.tab.active {
  background: var(--ink);
  color: #fff;
}

.content-area {
  margin-top: 10px;
  font-size: 13px;
  line-height: 1.8;
}

.preview-pane {
  margin-top: 10px;
  min-height: 300px;
  padding: 22px 24px;
  border: 1px solid var(--line);
  background: var(--surface);
}

.top-toggle {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-top: 20px;
  color: var(--muted);
  font-size: 12px;
  cursor: pointer;
}

.top-toggle input {
  accent-color: var(--coral);
}

.editor-footer {
  display: flex;
  gap: 10px;
  margin-top: 28px;
}

@media (max-width: 600px) {
  .title-input {
    font-size: 21px;
  }

  .cover-url {
    flex-direction: row;
    align-items: center;
  }
}
</style>
