<script setup lang="ts">
interface Props {
  index: number;
  title: string;
  summary?: string;
  cover?: string;
  date?: string;
  author?: string;
  url: string;
  variant?: "horizontal" | "vertical" | "mini";
}

const {
  index,
  title,
  summary = "",
  cover = "",
  date = "1970-01-01",
  author = "",
  variant = "horizontal",
} = defineProps<Props>();

console.log(index, title, summary, cover, date, author, variant);
</script>

<template>
  <div class="post-card">
    <NuxtLinkLocale :to="url" class="post-card-link">
      <div class="post-cover" :style="{ backgroundImage: `url(${cover})` }"></div>

      <!-- 文本部分 -->
      <div class="content">
        <h2>{{ title }}</h2>
        <p>{{ summary }}</p>
      </div>
    </NuxtLinkLocale>
  </div>
</template>

<style lang="less" scoped>
.post-card {
  height: 200px;
  width: 100%;
  background: var(--bg-card-base);
  border-radius: 10px;
  box-shadow: var(--shadow-card-base);
  margin: 10px 0;
}

.post-card-link {
  display: block;
  height: 100%;
  width: 100%;
}
</style>
