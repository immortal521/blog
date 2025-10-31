<script setup lang="ts">
const localePath = useLocalePath();

const { data } = await useFetch<{
  data: Post[];
}>("/api/v1/posts", {
  method: "get",
});

const posts = computed(() => data.value?.data ?? []);
</script>

<template>
  <div>
    <WelcomePanel />
    <ContentPanel :spacer="false">
      <div class="post-list">
        <PostCard
          v-for="(post, index) in posts"
          :key="post.id"
          :index="index"
          :cover="post.cover"
          :url="localePath('/blog/' + post.id)"
          :title="post.title"
          :date="post.updatedAt"
          :summary="post.summary"
        />
      </div>
    </ContentPanel>
  </div>
</template>

<style lang="less" scoped>
.post-list {
  max-width: 800px;
  margin: 0 auto;
  display: grid;
  gap: 10px;
  padding: 10px 20px;
  grid-template-columns: 1fr;
}
</style>
