<script setup lang="ts">
import { useMediaQuery, onClickOutside } from "@vueuse/core";

definePageMeta({
  layout: "admin-layout",
});

const menuItems = ref<MenuItem[]>([
  { icon: undefined, label: "Dashboard", to: "/admin" },
  { icon: undefined, label: "links", to: "/admin/links" },
  { icon: undefined, label: "post", to: "/admin/posts" },
]);

const width = useCookie("windowWidth");

// 初始化宽度，SSR时用cookie，客户端用真实窗口宽度
const windowWidth = ref(parseInt(width.value ?? "1080"));

const isMobileQuery = useMediaQuery("(max-width: 768px)", {
  ssrWidth: windowWidth.value,
});

const navIsOpen = ref(false);

const toggleNav = () => {
  navIsOpen.value = !navIsOpen.value;
};

// 点击外部关闭导航
const navCard = useTemplateRef("nav-card");

onClickOutside(navCard, () => {
  if (isMobileQuery.value) {
    navIsOpen.value = false;
  }
});
</script>

<template>
  <div>
    <header class="header">
      <button @click="toggleNav">toggle</button>
    </header>
    <div class="wrapper">
      <NavCard
        ref="nav-card"
        :menu-items="menuItems"
        :class="{ 'nav-aside': true, collapsed: isMobileQuery && !navIsOpen }"
      />
      <main class="main" :class="{ mobile: isMobileQuery }"></main>
    </div>
  </div>
</template>

<style lang="less" scoped>
.header {
  height: 60px;
  background: var(--bg-content);
  border-bottom: 1px solid var(--border-color-base);
  box-shadow: var(--shadow-md);
}

.wrapper {
  height: calc(100vh - 60px);
  overflow: hidden;
  padding: 20px;
  width: 100vw;
  position: relative;
  display: flex;
}

.nav-aside {
  transition: transform 0.3s ease-in-out;

  &.collapsed {
    transform: translateX(-220px);
  }
}

.main {
  width: 100%;
  margin-left: 220px;
  height: 100%;
  background: red;

  &.mobile {
    margin-left: 0;
  }
}
</style>
