<script setup lang="ts">
import { DesktopNav, MobileNav } from "./index";

const { isMobile } = useResponsive();
</script>

<template>
  <component :is="isMobile ? MobileNav : DesktopNav" />
</template>
