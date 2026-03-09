<script setup lang="ts">
const { spacer = true } = defineProps<{ spacer?: boolean }>();
</script>

<template>
  <div class="container">
    <BaseSpacer v-if="spacer" height="80" />
    <div class="blur-layer-container">
      <div class="blur-layer"></div>
    </div>
    <slot />
  </div>
</template>

<style lang="less" scoped>
.container {
  position: relative;
  background: var(--glass-gradient), var(--bg-content);
  animation: opacity 0.5s ease;
  padding-bottom: 140px;
}

.blur-layer-container {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: -1;
}

.blur-layer {
  position: sticky;
  top: 0;
  width: 100%;
  height: 100vh;
  backdrop-filter: var(--glass-blur);
  backface-visibility: hidden;
  transform: translateZ(0);
}
</style>
