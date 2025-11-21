<script setup lang="ts">
const { t } = useI18n();

const sidebar = useTemplateRef("sidebar");

const menuItems: MenuItem[] = [
  {
    icon: "ri:link",
    label: t("navbar.links"),
    to: "/admin/links",
  },
];

const isOpen = computed(() => sidebar.value?.isOpen);

defineExpose({
  toggleOpen: () => sidebar.value?.toggleOpen(),
  isOpen,
});
</script>

<template>
  <BaseSidebar ref="sidebar" :menu-items="menuItems">
    <template #header>
      <div class="header">
        <BaseLogo />
        <button>x</button>
      </div>
    </template>
  </BaseSidebar>
</template>

<style lang="less" scoped>
.header {
  width: 100%;
  height: 60px;
  display: flex;
  position: relative;
  justify-content: center;
  align-items: center;
}
</style>
