<script setup lang="ts">
const { isVisible } = useLoadingBar();

onMounted(() => {
  useThemeStore();
});
</script>

<template>
  <div>
    <NuxtRouteAnnouncer />
    <LoadingBar :is-visible="isVisible" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style lang="less">
.page-enter-active,
.page-leave-active {
  transition: all 0.4s;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
  filter: blur(10px);
}
</style>
