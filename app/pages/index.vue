<script setup lang="ts">
const localePath = useLocalePath();
</script>

<template>
  <div>
    <WelcomePanel />
    <main class="main">
      <NuxtLink :to="localePath('/blog/1')"> 前往 Markdown 渲染示例 </NuxtLink>
    </main>
  </div>
</template>

<style lang="less" scoped>
.main {
  background-color: var(--bg-content);
  width: 100%;
  height: 200px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
