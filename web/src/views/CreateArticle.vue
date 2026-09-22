<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createArticle } from '@/api/article'

const router = useRouter()

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover: '',
})
const loading = ref(false)
const errorMsg = ref('')

async function handleSubmit() {
  errorMsg.value = ''

  if (!form.value.title.trim()) {
    errorMsg.value = '请输入文章标题'
    return
  }
  if (!form.value.content.trim()) {
    errorMsg.value = '请输入文章内容'
    return
  }

  loading.value = true
  try {
    await createArticle({
      title: form.value.title.trim(),
      content: form.value.content,
      summary: form.value.summary.trim() || undefined,
      cover: form.value.cover.trim() || undefined,
    })
    alert('文章发表成功！')
    router.push('/')
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : '发表失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="create-page">
    <div class="page-heading">
      <h1>写文章</h1>
      <router-link to="/" class="back-link">← 返回首页</router-link>
    </div>

    <form class="article-form" @submit.prevent="handleSubmit">
      <div class="form-group">
        <label class="form-label">标题 <span class="required">*</span></label>
        <input
          v-model="form.title"
          type="text"
          class="form-input form-input-title"
          placeholder="请输入文章标题"
        />
      </div>

      <div class="form-group">
        <label class="form-label">摘要</label>
        <textarea
          v-model="form.summary"
          class="form-textarea"
          rows="3"
          placeholder="文章摘要（选填）"
        ></textarea>
      </div>

      <div class="form-group">
        <label class="form-label">封面图 URL</label>
        <input
          v-model="form.cover"
          type="text"
          class="form-input"
          placeholder="封面图片链接（选填）"
        />
      </div>

      <div class="form-group">
        <label class="form-label">正文内容 <span class="required">*</span></label>
        <textarea
          v-model="form.content"
          class="form-textarea form-textarea-content"
          rows="16"
          placeholder="支持 Markdown 或 HTML 格式&#10;例如：&#10;# 标题&#10;## 二级标题&#10;**粗体** *斜体*&#10;- 列表项&#10;1. 编号列表&#10;&#10;---分割线---&#10;&#10;支持直接使用 HTML 标签"
        ></textarea>
      </div>

      <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>

      <div class="form-actions">
        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? '发表中...' : '发表文章' }}
        </button>
        <button type="button" class="cancel-btn" @click="router.push('/')">
          取消
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.create-page {
  max-width: 720px;
  margin: 0 auto;
}

.page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
}

.page-heading h1 {
  font-size: 1.5rem;
  font-weight: 700;
}

.back-link {
  font-size: 0.9375rem;
  font-weight: 500;
}

.article-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text);
}

.required {
  color: var(--color-danger);
}

.form-input,
.form-textarea {
  padding: 10px 14px;
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  font-size: 0.9375rem;
  outline: none;
  transition: border-color 0.2s;
  font-family: inherit;
  resize: vertical;
}

.form-input:focus,
.form-textarea:focus {
  border-color: var(--color-primary);
}

.form-input-title {
  font-size: 1.125rem;
  font-weight: 600;
  padding: 12px 16px;
}

.form-textarea-content {
  min-height: 300px;
  line-height: 1.7;
  font-size: 0.9375rem;
}

.form-error {
  color: var(--color-danger);
  font-size: 0.875rem;
  text-align: center;
}

.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
}

.submit-btn {
  flex: 1;
  padding: 12px;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-sm);
  font-size: 1rem;
  font-weight: 600;
  transition: background 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.cancel-btn {
  padding: 12px 28px;
  background: var(--color-white);
  color: var(--color-text-light);
  border-radius: var(--radius-sm);
  font-size: 1rem;
  font-weight: 500;
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.cancel-btn:hover {
  border-color: var(--color-text-light);
  color: var(--color-text);
}
</style>
