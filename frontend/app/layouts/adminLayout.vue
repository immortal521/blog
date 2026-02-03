<script setup lang="ts">
import type { SidebarNode } from "~/components/BaseSidebar/types";

const { ts } = useI18n();

const menuItems = ref<SidebarNode[]>([
  {
    type: "link",
    icon: "duo-icons:dashboard",
    label: ts("admin.sidebar.dashboard"),
    to: "/admin",
    key: "/admin",
    exact: true,
  },
  {
    type: "section",
    key: "content",
    label: ts("admin.sidebar.content"),
    items: [
      {
        type: "link",
        icon: "ri:link",
        label: ts("admin.sidebar.links"),
        to: "/admin/links",
        key: "/admin/links",
      },
      {
        type: "link",
        icon: "material-symbols:post-rounded",
        label: ts("admin.sidebar.posts"),
        to: "/admin/posts",
        key: "/admin/posts",
      },
    ],
  },

  {
    type: "section",
    key: "settings",
    label: ts("admin.sidebar.settings"),
    items: [
      {
        type: "link",
        icon: "mdi:account-cog-outline",
        label: ts("admin.sidebar.profile"),
        to: "/admin/profile",
        key: "/admin/profile",
      },
      {
        type: "link",
        icon: "mdi:cog-outline",
        label: ts("admin.sidebar.system"),
        to: "/admin/settings",
        key: "/admin/settings",
      },
    ],
  },
]);

const route = useRoute();

const selectedKey = ref(route.fullPath);

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
    <BaseSidebar
      ref="sidebar"
      v-model:selected-key="selectedKey"
      v-model:open="open"
      v-model:collapsed="collapsed"
      :items="menuItems"
    >
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
  width: 100vw;
  height: 100vh;
  display: flex;
  color: var(--text-color-base);
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
