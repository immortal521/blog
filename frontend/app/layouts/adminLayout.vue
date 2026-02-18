<script setup lang="ts">
const { menuItems } = useAdminMenu();

const { open, collapsed } = useSidebar();
const { isMobile } = useResponsive();

const handleToggle = () => {
  if (isMobile.value) {
    collapsed.value = false;
    open.value = !open.value;
  } else {
    collapsed.value = !collapsed.value;
  }
};
</script>

<template>
  <div class="admin-layout">
    <BaseSidebar ref="sidebar" v-model:open="open" v-model:collapsed="collapsed" :items="menuItems">
      <template #header>
        <div class="logo">
          <h2>Admin</h2>
        </div>
      </template>
    </BaseSidebar>
    <div class="right">
      <header class="header">
        <button class="menu-btn" @click="handleToggle">
          <Icon v-if="!isMobile && !collapsed" name="icon-park-outline:menu-fold" size="28" />
          <Icon v-else-if="!isMobile && collapsed" name="icon-park-outline:menu-unfold" size="28" />
          <Icon v-else name="hugeicons:menu-11" size="24" class="icon" />
        </button>
      </header>
      <main class="main">
        <slot />
      </main>
    </div>
  </div>
</template>

<style lang="less" scoped>
.admin-layout {
  width: 100%;
  height: 100vh;
  display: flex;
  color: var(--text-color-primary);
}

.right {
  width: 100%;
}

.logo {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--text-color-primary);
}

.header {
  height: 50px;
  width: 100%;
  border-bottom: 1px solid var(--border-color-default);
  background: var(--glass-gradient), var(--bg-nav-base);
  backdrop-filter: var(--filter-blur-sm);
  display: flex;
  align-items: center;
  padding: 0 10px;

  .menu-btn {
    background: none;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 10px;
    color: inherit;
  }
}

.main {
  padding: 10px;
}
</style>
