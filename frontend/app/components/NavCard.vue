<script setup lang="ts">
interface Props {
  menuItems: MenuItem[];
}

const { menuItems } = defineProps<Props>();

const route = useRoute();
</script>

<template>
  <aside class="aside">
    <ul>
      <li v-for="item in menuItems" :key="item.label" :class="{ active: item.to === route.path }">
        {{ item.label }}
      </li>
    </ul>
  </aside>
</template>

<style lang="less" scoped>
.aside {
  height: max-content;
  width: 200px;
  border-radius: 8px;
  box-shadow: var(--shadow-md);
  position: absolute;
  top: 20px;
  left: 10px;
  background: var(--glass-gradient), var(--bg-content);
  color: var(--text-color-primary);
  overflow: hidden;

  ul {
    list-style: none;

    li {
      padding: 10px;
      border-radius: 8px;
    }

    .active {
      background: var(--color-primary-base);
    }
  }
}
</style>
