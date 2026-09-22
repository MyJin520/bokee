<script setup lang="ts">
</script>

<template>
  <div class="not-found-page">
    <span class="nf-code serif">404</span>
    <h1>这一页还没有被写出来</h1>
    <p>链接可能有误，或文章已经被作者移走。回到首页继续阅读吧。</p>
    <div class="nf-actions">
      <router-link class="btn" to="/">回到首页</router-link>
      <router-link class="btn btn-outline" to="/articles">浏览文章</router-link>
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  min-height: 60vh;
  padding: 90px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.nf-code {
  font-size: 96px;
  line-height: 1;
  color: var(--coral);
}

.not-found-page h1 {
  margin: 18px 0 10px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 32px;
  font-weight: 500;
}

.not-found-page p {
  max-width: 420px;
  margin: 0 0 26px;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.7;
}

.nf-actions {
  display: flex;
  gap: 10px;
}
</style>
