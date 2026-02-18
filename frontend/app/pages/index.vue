<script setup lang="ts">
import type { ApiListResponse } from "~/types/api";
import type { Post } from "~/types/post";

const { data } = await useFetch<ApiListResponse<Post>>("/api/v1/posts", {
  method: "get",
});

const posts = computed(() => data.value?.data ?? []);
</script>

<template>
  <div>
    <WelcomePanel />
    <ContentPanel :spacer="false" class="content-panel">
      <div class="post-list">
        <h1 v-if="posts.length > 0" class="title">Article</h1>
        <PostCard
          v-for="(post, index) in posts"
          :key="post.id"
          :index="index"
          :cover="post.cover"
          :url="'/blog/' + post.id"
          :title="post.title"
          :date="post.updatedAt"
          :summary="post.summary"
        />
      </div>
    </ContentPanel>
  </div>
</template>

<style lang="less" scoped>
.content-panel {
  border-radius: 10px 10px 0 0;
  box-shadow: 0 -4px 8px rgb(0 0 0 / 8%);
}

.title {
  position: relative;
  color: var(--text-color-primary);
  font-weight: 500;

  &::after {
    position: absolute;
    content: "";
    width: 80px;
    height: 6px;
    left: 0;
    bottom: 0;
    border-radius: 5px;
    background-color: var(--color-primary-base);
  }
}

.post-list {
  max-width: 800px;
  margin: 0 auto;
  display: grid;
  gap: 10px;
  padding: 10px 20px;
  grid-template-columns: 1fr;
}
</style>
