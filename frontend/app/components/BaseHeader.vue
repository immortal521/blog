<script setup lang="ts">
interface Props {
  isSticky?: boolean;
}

const { isSticky = false } = defineProps<Props>();
</script>

<template>
  <header ref="header" :class="{ header: true }">
    <div class="header-container" :class="{ 'is-sticky': isSticky }">
      <NavbarWrapper />
    </div>
  </header>
</template>

<style lang="less" scoped>
.header {
  position: fixed;
  width: 100vw;
  height: 100px;
  z-index: 1999;
}

.header-container {
  position: relative;
  max-width: 1000px;
  margin: auto;
  background: var(--glass-gradient), var(--bg-nav-base);
  border: 1px solid var(--border-color-nav);
  backdrop-filter: var(--blur-nav);
  overflow: hidden;
  width: calc(100% - 40px);
  height: 60px;
  top: 25px;
  border-radius: 15px;
  transition:
    top 0.5s ease-in-out,
    width 0.5s ease-in-out,
    max-width 0.5s ease-in-out,
    border-radius 0.5s ease-in-out,
    background 0.3s ease-in-out,
    border-color 0.3s ease-in-out;
  animation: scale-in 0.5s ease-in-out;
  box-shadow: var(--shadow-nav);

  &:not(.is-sticky):hover {
    background: var(--glass-gradient-strong), var(--bg-nav-hover);
  }

  &.is-sticky {
    border-radius: 0 0 15px 15px;
    top: 0;
    max-width: 100%;
    width: 100%;
  }
}

@media (width <= 768px) {
  .header-container {
    max-width: 100%;
    width: 100%;
    top: 0;
    border-radius: 0;

    &.is-sticky {
      border-radius: 0;
    }
  }
}
</style>
