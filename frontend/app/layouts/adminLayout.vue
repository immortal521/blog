<script setup lang="ts">
const sidebar = useTemplateRef("sidebar");

const handleToggle = () => {
  sidebar.value?.toggleOpen();
};
</script>

<template>
  <div class="admin-layout">
    <AdminSidebar ref="sidebar" class="sidebar" />
    <header class="header">
      <button class="toggle" @click="handleToggle">
        {{ sidebar?.isOpen ? "close" : "open" }}
      </button>
    </header>
    <main class="content">
      <slot></slot>
    </main>
  </div>
</template>

<style lang="less" scoped>
.admin-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  background-color: var(--bg-app);
}

.sidebar {
  height: 100vh;
  z-index: 1000;
}

.header {
  position: fixed;
  height: 60px;
  width: 100%;
  background-color: var(--bg-nav-base);
}

.content {
  color: white;
  padding: 10px;
}
</style>
