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
  background-color: var(--bg-content);
  backdrop-filter: blur(10px);
  border-radius: 10px 10px 0 0;
  animation: opacity 0.5s ease;
  padding-bottom: 140px;
  box-shadow: 0 -2px 4px rgba(0, 0, 0, 0.8);
}
</style>
