<script setup lang="ts">
useHead({
  bodyAttrs: {
    class: "with-bg",
  },
});

const isSticky = ref(false);
useScrollHandler(100, (over) => (isSticky.value = over));
</script>

<template>
  <div class="default-layout">
    <ReadProgress />
    <SiteUtilityBar />
    <BaseHeader :is-sticky="isSticky" />
    <main class="main">
      <slot />
    </main>
    <BaseFooter />
  </div>
</template>

<style lang="less" scoped>
.main {
  margin: 0;
}
</style>
