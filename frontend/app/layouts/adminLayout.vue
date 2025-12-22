<script setup lang="ts">
interface AdminItem extends SidebarItem {
  to: string;
}
const menuItems = ref<AdminItem[]>([
  {
    icon: undefined,
    label: "Dashboard",
    to: "/admin",
    key: "admin-index",
  },
  {
    icon: undefined,
    label: "links",
    to: "/admin/links",
    key: "admin-links",
  },
  {
    icon: undefined,
    label: "post",
    to: "/admin/posts",
    key: "admin-posts",
  },
]);
const selectedKey = ref(menuItems.value[0]?.key);

const { open } = useSidebar();
const { isMobile } = useResponsive();

const localePath = useLocalePath();

const onItemClicked = (item: SidebarItem) => {
  const { to, key } = item as AdminItem;
  const localeTo = localePath(to);
  selectedKey.value = key;
  navigateTo(localeTo);
  if (isMobile.value) {
    open.value = false;
  }
};
</script>

<template>
  <div class="admin-layout">
    <BaseSidebar
      ref="sidebar"
      v-model:selected-key="selectedKey"
      v-model:open="open"
      :items="menuItems"
      @item-clicked="onItemClicked"
    />
    <div class="right">
      <header class="header">
        <button v-if="isMobile" class="menu-btn" @click="open = !open">
          <Icon name="hugeicons:menu-11" size="24" />
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
}

.right {
  width: 100%;
}

.header {
  height: 50px;
  width: 100%;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-nav-base);
  backdrop-filter: blur(var(--nav-blur));
  display: flex;
  align-items: center;
  padding: 0 10px;

  .menu-btn {
    background: none;
    text-align: center;
    padding: 10px;
  }
}

.main {
  padding: 10px;
}
</style>
