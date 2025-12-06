<script setup lang="ts">
interface SpacerProps {
  height?: number | string;
}

const { height = 0 } = defineProps<SpacerProps>();
</script>

<template>
  <div class="spacer" :style="{ height: `${height}px` }"></div>
</template>
