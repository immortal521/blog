<script setup lang="ts">
import { useMediaQuery } from "@vueuse/core";
import { DesktopNav, MobileNav } from "./index";

const width = useCookie("windowWidth");

const isMobileQuery = useMediaQuery("(max-width: 768px)", {
  ssrWidth: parseInt(width.value ?? "1080"),
});
</script>

<template>
  <component :is="isMobileQuery ? MobileNav : DesktopNav" />
</template>
