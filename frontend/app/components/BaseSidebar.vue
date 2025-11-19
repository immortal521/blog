<script setup lang="ts">
interface Props {
  menuItems?: MenuItem[];
  width?: number | string;
  widthInFold?: number | string;
}

const { menuItems = [], width = 240, widthInFold = 60 } = defineProps<Props>();

const foldWidth = computed(() => `${widthInFold}px`);
const iconWidth = computed(() => {
  if (typeof widthInFold === "number") return `${widthInFold - 36}px`;
  else return `${parseInt(widthInFold) - 36}px`;
});

const sidebarWidth = computed(() => `${width}px`);

const isOpen = ref(true);
const isFold = ref(false);

const open = () => {
  isOpen.value = true;
};

const close = () => {
  isOpen.value = false;
};

const fold = () => {
  isFold.value = true;
};

const unfold = () => {
  isFold.value = false;
};

const toggleOpen = () => {
  if (isOpen.value) close();
  else open();
};

const toggleFold = () => {
  if (isFold.value) unfold();
  else fold();
};

defineExpose({
  open,
  close,
  fold,
  unfold,
  isOpen,
  toggleOpen,
  toggleFold,
});
</script>

<template>
  <div class="sidebar-wrapper" :class="{ fold: isFold, close: !isOpen }">
    <div class="sidebar-content">
      <slot name="header" />
      <ul class="menu">
        <li v-for="item in menuItems" :key="item.label" class="menu-item">
          <NuxtLinkLocale :to="item.to">
            <Icon :name="item.icon" size="18" class="icon" />
            <Transition name="menu-item-lable">
              <span v-show="!isFold" class="label">
                {{ item.label }}
              </span>
            </Transition>
          </NuxtLinkLocale>
        </li>
      </ul>
    </div>
  </div>
</template>

<style lang="less" scoped>
.sidebar-wrapper {
  background-color: var(--bg-sidebar);
  max-width: v-bind(sidebarWidth);
  width: 100%;
  transition: max-width 0.3s ease-in-out;
  overflow-x: hidden;
  box-shadow: var(--shadow-md);
}

.sidebar-content {
  opacity: 1;

  transition: opacity 0.3s ease-in-out;
}

.menu {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 10px;

  &-item {
    color: var(--text-color-base);
    height: 40px;
    padding: 0 8px;
    border-radius: 5px;

    &:hover {
      background-color: var(--bg-sidebar-item-hover);
    }

    a {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .label {
      flex-shrink: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .icon {
      flex-shrink: 0;
      display: inline-block;
      width: v-bind(iconWidth);
    }
  }
}

.menu-item-lable-enter-active,
.menu-item-lable-leave-active {
  transition: opacity 0.4s;
}
.menu-item-lable-enter-from,
.menu-item-lable-leave-to {
  opacity: 0;
}

.fold {
  max-width: v-bind(foldWidth);
}

.close {
  max-width: 0;

  .sidebar-content {
    opacity: 0;
  }
}
</style>
