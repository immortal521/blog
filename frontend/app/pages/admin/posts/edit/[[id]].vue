<script setup lang="ts">
definePageMeta({
  layout: "admin-sub-layout",
});

const route = useRoute();

const id = computed<number>(() => Number(route.params.id) || -1);
</script>

<template>
  <div>{{ id }}</div>
</template>

<style lang="scss" scoped></style>
