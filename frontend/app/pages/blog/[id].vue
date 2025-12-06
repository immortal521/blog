<script setup lang="ts">
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
      summary: "",
      publishedAt: "",
      updatedAt: "",
      viewCount: 0,
      readTimeMinutes: 0,
      cover: "",
      tags: [],
    }
  );
});

useHead({
  title: post.value.title,
});
</script>

<template>
  <ContentPanel :spacer="false">
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

@media (width <= 768px) {
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
