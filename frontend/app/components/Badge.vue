<script setup lang="ts">
interface Props {
  value: string | number;
}

const { value } = defineProps<Props>();
</script>

<template>
  <span class="badge">{{ value }}</span>
</template>

<style lang="less" scoped>
.badge {
  position: absolute;
  top: -6px;

  min-width: 18px;
  height: 18px;

  padding: 0 5px;

  border-radius: 999px;

  background: #f56c6c;
  color: white;

  font-size: 12px;

  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
