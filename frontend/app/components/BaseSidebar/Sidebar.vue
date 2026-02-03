<script setup lang="ts">
import type { SidebarNode } from "./types";

interface Props {
  items?: SidebarNode[];
  collapsed?: boolean;
  width?: number | string;
  widthCollapsed?: number | string;
}

const { collapsed = false, width = 220, widthCollapsed = 60, items = [] } = defineProps<Props>();

const open = defineModel<boolean>("open", { required: false, default: true });

const openKeys = ref<Set<string>>(new Set<string>());
</script>

<template>
  <div :class="{ 'sidebar-wrapper': true, collapsed, 'is-hidden': !open }">
    <aside class="sidebar-content">
      <div v-if="$slots.header" class="sidebar-header">
        <slot name="header" />
      </div>
      <template v-for="node in items" :key="node.key">
        <BaseSidebarSection
          v-if="node.type === 'section'"
          :section="node"
          :collapsed="collapsed"
          :open-keys="openKeys"
        />
        <BaseSidebarItem v-else :item="node" :collapsed="collapsed" :open-keys="openKeys" />
      </template>
    </aside>
  </div>
  <transition name="overlay">
    <div v-if="open" class="overlay" @click.stop="open = false"></div>
  </transition>
</template>

<style lang="less" scoped>
.sidebar-wrapper {
  background: var(--bg-sidebar);
  padding: 0 5px;
  height: 100vh;
  width: 100%;
  max-width: calc(v-bind(width) * 1px);
  overflow-x: hidden;
  transition: max-width 0.3s ease-in-out;
}

.sidebar-content {
  opacity: 1;
  transition: opacity 0.3s ease-in-out;
}

.collapsed {
  max-width: calc(v-bind(widthCollapsed) * 1px);
}

.is-hidden {
  max-width: 0;

  .sidebar-content {
    opacity: 0;
  }
}

.sidebar-header {
  height: 50px;
  width: 100%;
}

.nav-list {
  width: 100%;
  padding: 8px;

  .item {
    width: 100%;
    height: 35px;
    padding: 0 5px;
    line-height: 35px;
    border-radius: 8px;
    color: var(--text-color-base);
    user-select: none;
    cursor: pointer;
    position: relative;
    display: flex;
    margin: 5px 0;

    &.is-active {
      background: var(--color-primary-base);
    }

    &:not(.is-active):hover {
      background: var(--color-primary-bg);
    }

    .icon {
      flex-shrink: 0;
      width: 35px;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
    }

    .label {
      font-weight: bold;
      font-family: "Maple Mono", "Noto Sans SC", monospace;
    }
  }
}

.label-enter-active,
.label-leave-active {
  transition: all 0.3s ease-in-out;
}

.label-enter-from,
.label-leave-to {
  opacity: 0;
}

.overlay-enter-active,
.overlay-leave-active {
  transition: all 0.3s ease-in-out;
}

.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

.overlay {
  position: fixed;
  display: none;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgb(0 0 0 / 50%);
}

@media (width <= 768px) {
  .sidebar-wrapper {
    position: fixed;
    border-radius: 0 8px 8px 0;
    z-index: 9999;
  }

  .overlay {
    display: block;
    z-index: 9998;
  }
}
</style>
