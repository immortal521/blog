<script setup lang="ts"></script>

<template>
  <div class="card">
    <slot />
  </div>
</template>

<style lang="less" scoped>
.card {
  position: relative;
  background: var(--bg-card-base);
  border: 1px solid var(--border-color-card);
  box-shadow: var(--shadow-card-base);
  border-radius: var(--radius-card);
  transition:
    box-shadow 0.3s ease,
    background 0.3s ease,
    border-color 0.3s ease,
    transform 0.3s ease;

  &::before {
    content: "";
    position: absolute;
    inset: 0;
    border-radius: inherit;
    background: var(--glass-highlight);
    pointer-events: none;
  }

  &:hover {
    box-shadow: var(--shadow-card-hover);
    background: var(--bg-card-hover);
    border-color: var(--border-color-card-hover);
    transform: translateY(-1px);
  }
}
</style>
