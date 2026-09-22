<script setup lang="ts">
import { useRouter } from 'vue-router'
import ArticleEditor from '@/components/ArticleEditor.vue'

const router = useRouter()

// 创建成功后后端不返回新文章 ID，统一进入“我的文章”
function onSuccess() {
  router.push('/my-articles')
}
</script>

<template>
  <div class="editor-page">
    <header class="editor-header">
      <span class="eyebrow">New article</span>
      <h1>写一篇新文章</h1>
      <p>标题、封面、摘要与正文。支持 Markdown，可以先预览再发布。</p>
    </header>
    <ArticleEditor @success="onSuccess" />
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
