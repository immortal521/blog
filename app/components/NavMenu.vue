<script setup lang="ts">
const { items, vertical: vertical = false } = defineProps<{
  items: MenuItem[];
  vertical?: boolean;
}>();

const emit = defineEmits<{
  (e: "select"): void;
}>();
</script>

<template>
  <ul :class="{ 'nav-menu': true, vertical: vertical }">
    <li v-for="item in items" :key="item.to" class="nav-item" @click="emit('select')">
      <NuxtLink :to="item.to">
        <div class="icon">
          <Icon :name="item.icon" size="18" />
        </div>
        <span>
          {{ item.label }}
        </span>
      </NuxtLink>
    </li>
  </ul>
</template>

<style lang="less" scoped>
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
      color: var(--text-color-base);
      position: relative;
      align-items: center;
      transition: color 0.3s ease-in-out;
      padding: 0 2px;

      &::after {
        content: "";
        position: absolute;
        height: 6px;
        width: 0;
        border-radius: 10px;
        background-color: var(--color-primary-base);
        bottom: 10px;
        transition: width 0.3s ease-in-out;
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

      &:active {
        color: var(--color-primary-active);
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

.vertical {
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
  height: auto;
  flex-direction: column;
  gap: 10px;

  .nav-item {
    width: 100%;
    height: 40px;

    a {
      width: 100%;
      height: 100%;
      display: flex;
      color: var(--text-color-base);
      position: relative;
      align-items: center;
      transition:
        color 0.3s ease-in-out,
        border-bottom 0.3s ease-in-out,
        background-color 0.3s ease-in-out;
      padding: 0 8px;
      border-radius: 10px 10px 0 0;
      border-bottom: 1px solid var(--text-color-muted);

      &::after {
        display: none;
      }

      &:hover {
        border-bottom: 1px solid var(--color-primary-base);
        color: var(--color-primary-base);
        background-color: var(--bg-nav-base);
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

      &:active {
        color: var(--color-primary-active);
      }
      &:active::after {
        background-color: var(--color-primary-active);
      }

      .icon {
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        margin-right: 10px;
      }
    }
  }
}
</style>
