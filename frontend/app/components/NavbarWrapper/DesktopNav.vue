<script setup lang="ts">
import { navLinks } from "./navLinks";

const { $ts } = useI18n();

const items = computed<MenuItem[]>(() => {
  return navLinks.map((item) => {
    return {
      key: item.to,
      icon: item.icon,
      label: $ts(item.labelKey),
      to: item.to,
    };
  });
});
</script>

<template>
  <nav id="navbar" class="navbar">
    <BaseLogo />
    <ul :class="{ 'nav-menu': true }">
      <li v-for="item in items" :key="item.to" class="nav-item">
        <NuxtLinkLocale :to="item.to">
          <div class="icon">
            <Icon :name="item.icon" size="18" />
          </div>
          <span>
            {{ item.label }}
          </span>
        </NuxtLinkLocale>
      </li>
    </ul>
  </nav>
</template>

<style lang="less" scoped>
.navbar {
  width: 100%;
  height: 100%;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.nav-menu {
  display: flex;
  height: 100%;
  align-items: center;
  width: max-content;
  list-style: none;
  gap: 10px;

  .nav-item {
    height: 100%;

    a {
      width: 100%;
      height: 100%;
      display: flex;
      color: var(--text-color-primary);
      position: relative;
      align-items: center;
      transition: color 0.3s ease-in-out;
      padding: 0 2px;

      &::after {
        content: "";
        position: absolute;
        height: 6px;
        border-radius: 10px;
        width: 0;
        background-color: var(--color-primary-base);
        bottom: 10px;
        transition: width 0.3s ease-in-out;
      }

      &:active {
        color: var(--color-primary-active);
      }

      &:hover {
        color: var(--color-primary-base);
      }

      &:hover::after {
        width: 100%;
      }

      :deep(.svg-icon) {
        transition: color 0.3s ease-in-out;
      }

      &:hover:deep(.svg-icon) {
        color: var(--color-primary-base);
      }

      &:active::after {
        background-color: var(--color-primary-active);
      }

      .icon {
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }
  }
}
</style>
