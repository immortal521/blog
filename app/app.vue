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
