<script setup lang="ts">
const localePath = useLocalePath();

const { data } = await useFetch<{
  data: Post[];
}>("/api/v1/posts", {
  method: "get",
});

const posts = computed(() => data.value?.data);
</script>

<template>
  <div>
    <WelcomePanel />
    <main class="main">
      <h1 style="color: var(--text-color-base)">以下内容为测试文章列表所用</h1>
      <NuxtLink
        v-for="item in posts"
        :key="item.id"
        :to="localePath({ name: 'blog-id', params: { id: item.id } })"
      >
        {{ item.title }}
      </NuxtLink>
    </main>
  </div>
</template>

<style lang="less" scoped>
.main {
  background-color: var(--bg-content);
  width: 100%;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>
