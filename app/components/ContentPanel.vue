<script setup lang="ts"></script>

<template>
  <div class="container">
    <BaseSpacer height="80" />
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
