<script setup lang="ts">
import type { Post } from "~/types/post";
import type { ApiResponse } from "~/types/api";

const route = useRoute();

const params = computed(() => route.params);

const { data } = await useFetch<ApiResponse<Post>>(() => `/api/v1/posts/${params.value.id}`, {
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

useHead(() => ({ title: post.value.title }));

const rendered = computed(() => {
  if (!post.value.content) {
    return {
      content: [],
      toc: [],
    };
  }
  return parseMarkdownToVNode(post.value.content, { toc: true });
});

const content = computed(() => rendered.value.content);
const toc = computed(() => rendered.value.toc ?? []);
const { activeId } = useActiveToc(toc);
</script>

<template>
  <ContentPanel :spacer="false">
    <ArticleCover :src="post.cover" :title="post.title" />
    <article class="article">
      <main class="content">
        <MarkdownRenderer :content="content" />
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
  gap: 24px;
  padding: 0 20px;
  padding-top: 30px;
}

.content {
  width: 100%;
  max-width: 800px;
  background-color: var(--bg-card-base);
  animation: article-show 1s ease-in-out;
  box-shadow: var(--shadow-lg);
  padding: 40px;
  margin-top: 8px;
  border: 1px solid var(--border-color-card);
  border-radius: 12px;
  position: relative;
  overflow: hidden;
}

.toc-container {
  position: sticky;
  top: 80px;
  width: 240px;
  height: max-content;
  background-color: var(--bg-card-base);
  box-shadow: var(--shadow-md);
  padding: 16px;
  margin-top: 8px;
  border: 1px solid var(--border-color-card);
  border-radius: 12px;
  transition:
    box-shadow 0.3s ease,
    transform 0.3s ease;

  &:hover {
    box-shadow: var(--shadow-lg);
    transform: translateY(-2px);
  }
}

@media (width <= 1200px) {
  .toc-container {
    display: none;
  }

  .article {
    padding: 0;
  }

  .content {
    padding: 0 16px;
    margin-top: 0;
  }
}
</style>
