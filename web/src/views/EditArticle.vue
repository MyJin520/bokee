<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticleInfo, updateArticle } from '@/api/article'

const route = useRoute()
const router = useRouter()

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover: '',
})
const loading = ref(true)
const submitting = ref(false)
const errorMsg = ref('')

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    router.push('/')
    return
  }

  try {
    const article = await getArticleInfo(id)
    form.value = {
      title: article.title,
      content: article.content,
      summary: article.summary || '',
      cover: article.cover || '',
    }
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '文章不存在')
    router.push('/')
  } finally {
    loading.value = false
  }
})

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

  submitting.value = true
  try {
    await updateArticle({
      id: Number(route.params.id),
      title: form.value.title.trim(),
      content: form.value.content,
      summary: form.value.summary.trim() || undefined,
      cover: form.value.cover.trim() || undefined,
    })
    alert('文章更新成功！')
    router.push(`/article/${route.params.id}`)
  } catch (e: unknown) {
    errorMsg.value = e instanceof Error ? e.message : '更新失败'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div v-if="loading" class="loading-state">
    <div class="spinner"></div>
    <p>加载中...</p>
  </div>

  <div v-else class="edit-page">
    <div class="page-heading">
      <h1>编辑文章</h1>
      <router-link :to="`/article/${route.params.id}`" class="back-link">← 返回文章</router-link>
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
          placeholder="支持 Markdown 或 HTML 格式"
        ></textarea>
      </div>

      <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>

      <div class="form-actions">
        <button type="submit" class="submit-btn" :disabled="submitting">
          {{ submitting ? '保存中...' : '保存修改' }}
        </button>
        <button
          type="button"
          class="cancel-btn"
          @click="router.push(`/article/${route.params.id}`)"
        >
          取消
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.loading-state {
  text-align: center;
  padding: 80px 0;
  color: var(--color-text-light);
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.edit-page {
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
