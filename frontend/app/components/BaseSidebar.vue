<script setup lang="ts">
interface Props {
  menuItems?: MenuItem[];
  width?: number | string;
  widthInFold?: number | string;
}

const { menuItems = [], width = 240, widthInFold = 60 } = defineProps<Props>();

const foldWidth = computed(() => `${widthInFold}px`);
const sidebarWidth = computed(() => `${width}px`);

console.log(menuItems);
console.log(width);

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

defineExpose({
  open,
  close,
  fold,
  unfold,
});
</script>

<template>
  <div class="sidebar-wrapper" :class="{ fold: isFold, close: !isOpen }"></div>
</template>

<style lang="less" scoped>
.sidebar-wrapper {
  background-color: var(--bg-content);
  max-width: v-bind(sidebarWidth);
  width: 100%;
  transition: max-width 0.3s ease-in-out;
}

.fold {
  max-width: v-bind(foldWidth);
}

.close {
  max-width: 0;
}
</style>
