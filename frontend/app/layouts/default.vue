<script setup lang="ts">
useHead({
  bodyAttrs: {
    class: "with-bg",
  },
});
</script>

<template>
  <div class="default-layout">
    <ReadProgress />
    <SiteUtilityBar />
    <BaseHeader />
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
