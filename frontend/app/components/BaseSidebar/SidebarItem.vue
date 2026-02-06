<script setup lang="ts">
import type { SidebarItem } from "./types";
interface Props {
  item: SidebarItem;
  collapsed?: boolean;
  openKeys: Set<string>;
}

const { item, openKeys, collapsed = false } = defineProps<Props>();

const emit = defineEmits<{
  (e: "toggle", key: string): void;
}>();

const isOpen = computed(() => {
  return item.type === "group" && openKeys.has(item.key);
});

const route = useRoute();

const isLinkMatch = (to: string, exact?: boolean) => {
  if (exact) return route.path === to;
  return route.path.startsWith(to);
};

const hasActiveChild = (items: SidebarItem[]): boolean => {
  return items.some((item) => {
    if (item.type === "link") {
      return isLinkMatch(item.to, item.exact);
    }
    if (item.type === "group") {
      return hasActiveChild(item?.children ?? []);
    }
    return false;
  });
};

const isLinkActive = computed(() => {
  return item.type === "link" && isLinkMatch(item.to, item.exact);
});

const isGroupActive = computed(() => {
  if (item.type !== "group") return false;
  return hasActiveChild(item.children ?? []);
});

const onLinkClick = (e: MouseEvent) => {
  if (item.type === "link" && item.disabled) e.preventDefault();
};

const onGroupClick = () => {
  if (item.type !== "group") return;
  if (item.disabled) return;
  emit("toggle", item.key);
};

const onActiveClick = () => {
  if (item.type !== "action") return;
  if (item.disabled) return;
  item.action();
};

const onBeforeEnter = (el: Element) => {
  const e = el as HTMLElement;
  e.style.height = "0";
  e.style.opacity = "0";
  e.style.overflow = "hidden";
};

const onEnter = (el: Element) => {
  const e = el as HTMLElement;
  void e.offsetHeight;
  e.style.height = `${el.scrollHeight}px`;
  e.style.opacity = "1";
};

const onAfterEnter = (el: Element) => {
  const e = el as HTMLElement;
  e.style.height = "auto";
  e.style.overflow = "";
  e.style.opacity = "";
};

const onBeforeLeave = (el: Element) => {
  const e = el as HTMLElement;
  e.style.height = `${el.scrollHeight}px`;
  e.style.opacity = "1";
  e.style.overflow = "hidden";
};

const onLeave = (el: Element) => {
  const e = el as HTMLElement;
  void e.offsetHeight;
  e.style.height = "0";
  e.style.opacity = "0";
};

const onAfterLeave = (el: Element) => {
  const e = el as HTMLElement;
  e.style.height = "";
  e.style.opacity = "";
  e.style.overflow = "";
};
</script>

<template>
  <template v-if="!item.hidden">
    <div v-if="item.type === 'divider'" class="row divider">
      <span v-if="item.label && !collapsed" class="label">{{ item.label }}</span>
      <div class="divider-line"></div>
    </div>
    <div v-else-if="item.type === 'link'" class="row" :class="{ active: isLinkActive, collapsed }">
      <NuxtLinkLocale
        :to="item.to"
        :class="{ btn: true, disabled: item.disabled }"
        :aria-disabled="item.disabled"
        :tabindex="item.disabled ? -1 : 0"
        @click="onLinkClick"
      >
        <div class="icon">
          <Icon v-if="item.icon" :name="item.icon" size="24" />
        </div>
        <span class="label">{{ item.label }}</span>
      </NuxtLinkLocale>
    </div>
    <div v-else-if="item.type === 'action'" class="row" :class="{ collapsed }">
      <button class="btn" :disabled="item.disabled" @click="onActiveClick">
        <div class="icon">
          <Icon v-if="item.icon" :name="item.icon" />
        </div>
        <span class="label">{{ item.label }}</span>
      </button>
    </div>
    <div
      v-else-if="item.type === 'group'"
      class="row"
      :class="{ 'active-group': isGroupActive, collapsed }"
    >
      <button
        v-if="!collapsed"
        class="btn"
        type="button"
        :aria-expanded="isOpen"
        :disabled="item.disabled"
        @click="onGroupClick"
      >
        <div class="icon">
          <Icon v-if="item.icon" :name="item.icon" size="24" />
        </div>
        <span class="label">{{ item.label }}</span>
      </button>
      <Transition
        name="group-collapse"
        @before-enter="onBeforeEnter"
        @enter="onEnter"
        @after-enter="onAfterEnter"
        @before-leave="onBeforeLeave"
        @leave="onLeave"
        @after-leave="onAfterLeave"
      >
        <div v-show="isOpen" class="children">
          <SidebarItem
            v-for="ch in item.children"
            :key="ch.key"
            :item="ch"
            :collapsed="collapsed"
            :open-keys="openKeys"
            @toggle="emit('toggle', $event)"
          />
        </div>
      </Transition>
    </div>
  </template>
</template>

<style lang="less" scoped>
.row {
  margin: 2px 0;
  overflow: hidden;
}

.btn {
  width: 100%;
  display: flex;
  align-items: center;
  background: transparent;
  min-height: 35px;
  font-family: "Maple Mono", monospace;
  text-decoration: none;
  color: var(--text-color-base);
  padding: 2px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;

  .label {
    font-weight: inherit;
    font-size: 1.6rem;
  }
}

.icon {
  width: 30px;
  margin-right: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.active > .btn {
  background: var(--color-primary-bg-active);
}

.active-group > .btn {
  font-weight: 800;
}

.disabled,
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.row > .btn:hover {
  background: var(--color-primary-bg);
}

.row > .btn:active {
  background: var(--color-primary-bg-active);
}

.divider {
  width: 100%;
  padding: 0 4px;
  margin: 8px 0;
  display: flex;
  align-items: center;

  .label {
    padding-right: 5px;
    font-size: 1.2rem;
    color: var(--text-color-muted);
  }

  .divider-line {
    width: 100%;
    height: 1px;
    border-radius: 5px;
    background: var(--border-color-divider);
  }
}

.row.collapsed > .btn {
  flex-direction: column;
  background: transparent;
  min-width: 0;

  .label {
    font-size: 1.2rem;
    white-space: normal;
    overflow-wrap: anywhere;
    min-width: 0;
    text-align: center;
  }
}

.row > .btn.disabled:hover,
.row > .btn:disabled:hover,
.row > .btn.disabled:active,
.row > .btn:disabled:active {
  background: transparent;
}

.row.collapsed > .btn .icon {
  width: 100%;
  height: 35px;
  margin-right: 0;
  border-radius: 24px;
}

.row.collapsed > .btn:hover,
.row.collapsed > .btn:active {
  background: transparent;
}

.row > .children {
  padding-left: 12px;
}

.row.collapsed > .children {
  padding-left: 0;
}

.row.collapsed > .btn:hover .icon {
  background: var(--color-primary-bg);
}

.row.collapsed > .btn:active .icon {
  background: var(--color-primary-bg-active);
}

.row.collapsed.active > .btn .icon {
  background-color: var(--color-primary-bg-active);
}

.row.collapsed > .btn.disabled:hover .icon,
.row.collapsed > .btn:disabled:hover .icon,
.row.collapsed > .btn.disabled:active .icon,
.row.collapsed > .btn:disabled:active .icon {
  background: transparent;
}

.group-collapse-enter-active,
.group-collapse-leave-active {
  transition:
    height 0.3s ease-in-out,
    opacity 0.2s ease-in-out;
  overflow: hidden;
}
</style>
