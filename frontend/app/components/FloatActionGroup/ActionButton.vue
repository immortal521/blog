<script setup lang="ts">
defineEmits<{
  (e: "click", event: MouseEvent): void;
}>();

const { height = 30, icon } = defineProps<{
  height?: number | string;
  icon: string;
}>();
</script>

<template>
  <button :style="{ height: `${height}px` }" class="action-button" @click="$emit('click', $event)">
    <Icon :name="icon" size="24" class="icon" />
  </button>
</template>

<style lang="less" scoped>
.action-button {
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--bg-nav-base);
  border: 1px solid var(--border-color-nav);
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  cursor: pointer;
}

.icon {
  color: var(--text-color-base);
}
</style>
