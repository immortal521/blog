<script setup lang="ts">
interface Props {
  toc: TocItem[];
  activeId: string;
}

const { toc } = defineProps<Props>();

const handleTocClick = (id: string) => {
  const target = document.getElementById(id);

  if (!target) return;

  target.scrollIntoView({
    behavior: "smooth",
    block: "start",
  });

  history.pushState(null, "", `#${id}`);
};
</script>

<template>
  <ul class="toc">
    <li v-for="item in toc" :key="item.id" :class="{ active: activeId === item.id }">
      <a :href="`#${item.id}`" @click.prevent="handleTocClick(item.id)">{{ item.text }}</a>
    </li>
  </ul>
</template>

<style lang="less" scoped>
.toc {
  max-height: 500px;
  width: 100%;
  list-style: none;
}
</style>
