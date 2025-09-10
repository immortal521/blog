<script setup lang="ts">
useHead({
  title: "markdown 渲染示例",
});

const route = useRoute();

const params = computed(() => route.params);

const { data } = await useFetch<{
  data: Post;
}>("/api/v1/posts/" + params.value.id, {
  method: "get",
});

const post = computed<Post>(() => {
  return (
    data.value?.data ?? {
      id: -1,
      title: "",
      content: "",
      publishedAt: "",
      viewCount: 0,
      readTimeMinutes: 0,
      cover: "",
      tags: [],
    }
  );
});
</script>

<template>
  <ContentPanel :spacer="false">
    <!--  TODO: 修复 v-viewer 服务端渲染时报错的问题 -->
    <article class="article">
      <ArticleCover :src="post.cover" :title="post.title" />
      <main class="content">
        <MarkdownRenderer :markdown="post.content"></MarkdownRenderer>
      </main>
    </article>
  </ContentPanel>
</template>

<style lang="less" scoped>
.article {
  width: 100%;
}

.content {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding-top: 30px;
  animation: article-show 1s ease-in-out;
}

@media (max-width: 768px) {
  .content {
    padding: 0 20px;
  }
}

@keyframes article-show {
  0% {
    opacity: 0;
    transform: translateY(32px);
  }
  100% {
    opacity: 1;
  }
}
</style>
