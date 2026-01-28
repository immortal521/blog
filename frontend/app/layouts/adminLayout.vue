<script setup lang="ts">
interface AdminItem extends SidebarItem {
  to: string;
}
const menuItems = ref<AdminItem[]>([
  {
    icon: undefined,
    label: "Dashboard",
    to: "/admin",
    key: "/admin",
  },
  {
    icon: undefined,
    label: "links",
    to: "/admin/links",
    key: "/admin/links",
  },
  {
    icon: undefined,
    label: "post",
    to: "/admin/posts",
    key: "/admin/posts",
  },
]);

const route = useRoute();

const selectedKey = ref(route.fullPath);

const { open, collapsed } = useSidebar();
const { isMobile } = useResponsive();

const { $localePath } = useI18n();

const onItemClicked = (item: SidebarItem) => {
  const { to, key } = item as AdminItem;
  const localeTo = $localePath(to);
  selectedKey.value = key;
  navigateTo(localeTo);
  if (isMobile.value) {
    open.value = false;
  }
};

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
    <BaseSidebar
      ref="sidebar"
      v-model:selected-key="selectedKey"
      v-model:open="open"
      v-model:collapsed="collapsed"
      :items="menuItems"
      @item-clicked="onItemClicked"
    />
    <div class="right">
      <header class="header">
        <button class="menu-btn" @click="handleToggle">
          <Icon name="hugeicons:menu-11" size="24" class="icon" />
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
  width: 100vw;
  height: 100vh;
  display: flex;
  color: var(--text-color-base);
}

.right {
  width: 100%;
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

  .menu-btn {
    background: none;
    text-align: center;
    padding: 10px;
    color: inherit;
  }
}

.main {
  padding: 10px;
}
</style>
