<script setup lang="ts">
const { spacer = true } = defineProps<{ spacer?: boolean }>();
</script>

<template>
  <div class="container">
    <BaseSpacer v-if="spacer" height="80" />
    <slot />
  </div>
</template>

<style lang="less" scoped>
.container {
  min-height: 100vh;
  background-color: var(--bg-content);
  backdrop-filter: blur(10px);
  animation: opacity 0.5s ease;
  padding-bottom: 140px;
}
</style>
