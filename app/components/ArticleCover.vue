<script setup lang="ts">
interface Props {
  src?: string;
  title?: string;
}

const { src = "", title = "" } = defineProps<Props>();
</script>

<template>
  <header class="article-cover" :style="{ backgroundImage: `url(${src})` }">
    <div class="cover-content">
      <h1 class="title">{{ title }}</h1>
    </div>
  </header>
</template>

<style lang="less" scoped>
.article-cover {
  position: relative;
  width: 100vw;
  height: 400px;
  backdrop-filter: blur(10px);
  background-position: center center;
  background-size: cover;
  background-repeat: no-repeat;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  transition: height 0.3s ease-in-out;
}

.cover-content {
  width: 100%;
  height: 100%;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  color: var(--text-color-base);
}

.title {
  text-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
}

@media (max-width: 768px) {
  .article-cover {
    height: 280px;
  }
}
</style>
