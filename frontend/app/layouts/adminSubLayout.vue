<script setup lang="ts">
const router = useRouter();

const handleBack = () => {
  router.back();
};
</script>

<template>
  <div class="admin-sub-layout">
    <header class="header">
      <button class="back-btn" @click="handleBack">
        <Icon name="eva:arrow-back-fill" size="28" />
      </button>
      <div class="actions">
        <slot name="actions" />
      </div>
    </header>
    <main class="main">
      <slot />
    </main>
  </div>
</template>

<style lang="less" scoped>
.admin-sub-layout {
  width: 100%;
  height: 100dvh;
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
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
}

.actions {
  margin-left: auto;
}

.main {
  height: 100%;
  padding: 65px 15px 15px;
}

.back-btn {
  background: none;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  color: inherit;
}
</style>
