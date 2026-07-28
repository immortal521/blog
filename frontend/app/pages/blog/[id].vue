<script setup lang="ts">
import type { Post } from "~/types/post";
import type { ApiResponse } from "~/types/api";

const route = useRoute();

const params = computed(() => route.params);

const { data } = await useFetch<ApiResponse<Post>>("/api/v1/posts/" + params.value.id, {
  method: "get",
});

const post = computed<Post>(() => {
  return (
    data.value?.data ?? {
      id: -1,
      title: "",
      content: "",
      summary: "",
      publishedAt: "",
      updatedAt: "",
      viewCount: 0,
      readTimeMinutes: 0,
      cover: "",
      tags: [],
      categories: [],
    }
  );
});

useHead({
  title: post.value.title,
});

const toc = ref<TocItem[]>([]);

const handleTocChange = (value: TocItem[]) => {
  toc.value = value;
};

const { activeId } = useActiveToc(toc);
</script>

<template>
  <ContentPanel :spacer="false">
    <ArticleCover :src="post.cover" :title="post.title" />
    <article class="article">
      <main class="content">
        <MarkdownRenderer :markdown="post.content" :toc="true" @toc-change="handleTocChange" />
      </main>
      <div class="toc-container">
        <ArticleToc :toc="toc" :active-id="activeId" />
      </div>
    </article>
  </ContentPanel>
</template>

<style lang="less" scoped>
.article {
  position: relative;
  width: 100%;
  display: flex;
  justify-content: center;
  align-self: flex-start;
}

.content {
  width: 100%;
  max-width: 800px;
  padding-top: 30px;
  margin-right: 20px;
  background-color: var(--bg-card-base);
  animation: article-show 1s ease-in-out;
  box-shadow: var(--shadow-md);
  padding: 10px;
  margin-top: 4px;
  border: 1px solid var(--border-color-card);
  border-radius: 8px;
}

.toc-container {
  position: sticky;
  top: 70px;
  width: 220px;
  height: max-content;
  background-color: var(--bg-card-base);
  box-shadow: var(--shadow-md);
  padding: 10px;
  margin-top: 4px;
  border: 1px solid var(--border-color-card);
  border-radius: 8px;
}

@media (width <= 1200px) {
  .toc-container {
    display: none;
  }
}
</style>
