<script setup lang="ts">
const router = useRouter();

const handleBack = () => {
  router.back();
};
</script>

<template>
  <div class="admin-sub-layout">
    <header class="header">
      <button @click="handleBack">1</button>
    </header>
    <main class="main">
      <slot />
    </main>
  </div>
</template>

<style lang="less" scoped>
.admin-sub-layout {
  width: 100vw;
  height: 100vh;
  color: var(--text-color-base);
}

.header {
  height: 50px;
  width: 100%;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-nav-base);
  backdrop-filter: var(--filter-blur-sm);
  display: flex;
  align-items: center;
  padding: 0 10px;
}

.main {
  padding: 15px;
}
</style>
