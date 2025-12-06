<script setup lang="ts"></script>

<template>
  <div class="card">
    <slot />
  </div>
</template>

<style lang="less" scoped>
.card {
  background: var(--bg-card-base);
  border: 2px solid var(--border-color-card);
  box-shadow: var(--shadow-card-base);
  border-radius: var(--radius-card);
  transition:
    box-shadow 0.2s ease-in-out,
    background 0.2s ease-in-out,
    border 0.2s ease-in-out;

  &:hover {
    box-shadow: var(--shadow-card-hover);
    background: var(--bg-card-hover);
    border: 2px solid var(--color-primary-bg);
  }
}
</style>
