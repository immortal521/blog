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
  background: var(--bg-content);
  animation: opacity 0.5s ease;
  padding-bottom: 140px;
  min-height: 100vh;
}

.blur-layer-container {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: -1;
}

.blur-layer {
  position: sticky;
  top: 0;
  width: 100%;
  height: 100%;
  max-height: 100vh;
  backface-visibility: hidden;
  transform: translateZ(0);
  will-change: transform, opacity;
  background: radial-gradient(ellipse at 50% 0%, var(--color-primary-bg-muted) 0%, transparent 70%);
  opacity: 0.3;
}
</style>
