<script setup lang="ts">
import type { BundledLanguage, SpecialLanguage } from "shiki";

interface Props {
  code: string;
  lang: BundledLanguage | SpecialLanguage;
}

const { code, lang } = defineProps<Props>();
</script>

<template>
  <div class="code-header">
    <span class="code-lang">{{ lang }}</span>
    <CopyButton :text="code" class="copy-btn" />
  </div>
</template>

<style lang="less" scoped>
.code-header {
  display: flex;
  align-items: center;
  padding: 0 1.5rem;
}

.copy-btn {
  margin-left: auto;
}
</style>
