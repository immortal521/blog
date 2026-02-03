<script setup lang="ts">
import type { SidebarSection, SidebarItemEmits } from "./types";
interface Props {
  section: SidebarSection;
  collapsed?: boolean;
  openKeys: Set<string>;
}

const emit = defineEmits<SidebarItemEmits>();

const { section, collapsed = false, openKeys } = defineProps<Props>();
</script>

<template>
  <section v-if="!section.hidden" class="section">
    <div v-if="section.label && !collapsed" class="section-title">
      <span>{{ section.label }}</span>
    </div>
    <div class="section-items">
      <BaseSidebarItem
        v-for="it in section.items"
        :key="it.key"
        :item="it"
        :collapsed="collapsed"
        :open-keys="openKeys"
        @toggle="emit('toggle', $event)"
      />
    </div>
  </section>
</template>

<style lang="less" scoped>
.section-title {
  font-size: 1.2rem;
  color: var(--text-color-muted);
  margin: 10px 0;
}

.section-items {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
</style>
