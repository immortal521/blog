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
      <div class="cover">
        <NuxtImg :src="cover" />
      </div>

      <!-- 文本部分 -->
      <div class="content">
        <h2 class="title">{{ title }}</h2>
        <p>{{ summary }}</p>
      </div>
    </NuxtLinkLocale>
  </div>
</template>

<style lang="less" scoped>
.post-card {
  height: 260px;
  width: 100%;
  margin: 10px 0;
}

.post-card-link {
  display: block;
  position: relative;
  height: 100%;
  border-radius: 10px;
  overflow: hidden;
  width: 100%;
  background: var(--bg-card-base);
  box-shadow: var(--shadow-card-base);

  .cover {
    height: 180px;
    width: 100%;
    overflow: hidden;

    img {
      object-fit: cover;
      width: 100%;
      height: 100%;
    }
  }

  .content {
    position: relative;
    .title {
      position: relative;
      top: -30px;
      left: 20px;
      width: max-content;
      background: var(--bg-card-title);
      border: 1px solid var(--border-color-base);
      border-radius: 5px;
      padding: 10px 20px;
      font-size: 2rem;
    }
  }
}

@media (max-width: 768px) {
  .post-card-link {
    .content {
      .title {
        font-size: 1.5rem;
        top: -20px;
        left: 10px;
        padding: 5px 10px;
      }
    }
  }
}
</style>
