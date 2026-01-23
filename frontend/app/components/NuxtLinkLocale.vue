<script setup lang="ts">
import type { RouteLocationNamedRaw, RouteLocationResolvedGeneric } from "vue-router";

const props = defineProps<{
  to: string | RouteLocationNamedRaw | RouteLocationResolvedGeneric;
}>();

const { $localePath } = useI18n();

const resolvedTo = computed(() => $localePath(props.to));
</script>

<template>
  <NuxtLink :to="resolvedTo">
    <slot />
  </NuxtLink>
</template>
