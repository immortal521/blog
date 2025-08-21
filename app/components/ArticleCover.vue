<script setup lang="ts">
interface Props {
  src?: string;
  title?: string;
}

const { src = "", title = "" } = defineProps<Props>();
</script>

<template>
  <header class="article-cover">
    <NuxtImg v-if="src" :src class="cover" />
    <div class="cover-content">
      <h1 class="title">{{ title }}</h1>
    </div>
  </header>
</template>

<style lang="less" scoped>
.article-cover {
  position: relative;
  width: 100vw;
  height: 500px;
  background-color: var(--bg-article-cover);
  backdrop-filter: blur(10px);
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
}

.cover {
  width: 100%;
  object-fit: cover;
}

.cover-content {
  width: 100%;
  height: 100%;
  position: absolute;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  color: var(--text-color-base);
  top: 0;
  left: 0;
  z-index: 1;
}

.title {
  text-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
}
</style>
