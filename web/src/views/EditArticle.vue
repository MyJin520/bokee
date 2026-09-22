<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleEditor from '@/components/ArticleEditor.vue'

const route = useRoute()
const router = useRouter()
const articleId = computed(() => Number(route.params.id))

function onSuccess(id: number) {
  if (id) router.push(`/article/${id}`)
  else router.push('/my-articles')
}
</script>

<template>
  <div class="editor-page">
    <header class="editor-header">
      <span class="eyebrow">Edit article</span>
      <h1>编辑文章</h1>
      <p>修改会在保存后立即生效，文章详情缓存也会同步刷新。</p>
    </header>
    <ArticleEditor :article-id="articleId" @success="onSuccess" />
  </div>
</template>

<style scoped>
.editor-page {
  padding: 48px 0 80px;
}

.editor-header {
  margin-bottom: 30px;
}

.editor-header h1 {
  margin: 16px 0 10px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(32px, 4.4vw, 50px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.editor-header p {
  max-width: 520px;
  margin: 0;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.7;
}
</style>
