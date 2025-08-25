<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";

const { icon, animation = undefined } = defineProps<{
  icon: string;
  animation?: "up" | "down" | "left" | "right" | "scale";
}>();

const open = ref(false);

const panelContent = useTemplateRef<HTMLDivElement>("panelContent");
const actionButton = useTemplateRef<HTMLButtonElement>("actionButton");

onClickOutside(panelContent, (event) => {
  // 判断点击事件触发位置是否是 action button
  // 如果是则不关闭，防止 panel content 闪烁
  if (actionButton.value?.contains(event.target as Node)) return;
  open.value = false;
});
</script>

<template>
  <div class="panel-button-wrapper">
    <button ref="actionButton" class="action-button" @click="open = !open">
      <Icon :name="icon" size="24" class="icon" />
    </button>

    <Transition :name="animation ?? 'none'">
      <div v-if="open" ref="panelContent" class="panel-content">
        <slot />
      </div>
    </Transition>
  </div>
</template>

<style lang="less" scoped>
.panel-button-wrapper {
  position: relative;
}

.icon {
  color: var(--text-color-base);
}

.action-button {
  background: var(--bg-nav-base);
  border-radius: 8px;
  width: 100%;
  height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px solid var(--border-color-nav);
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.panel-content {
  position: absolute;
  max-width: 300px;
  right: 120%;
  bottom: 0;
  z-index: 10;
}

.down-enter-active,
.down-leave-active {
  transition:
    opacity 0.5s ease-in-out,
    transform 0.5s ease-in-out;
}

.down-enter-from,
.down-leave-to {
  opacity: 0;
  transform: scale(0.3) translateY(100%);
}

.up-enter-active,
.up-leave-active {
  transition:
    opacity 0.5s ease-in-out,
    transform 0.5s ease-in-out;
}

.up-enter-from,
.up-leave-to {
  opacity: 0;
  transform: scale(0.3) translateY(-100%);
}

.left-enter-active,
.left-leave-active {
  transition:
    opacity 0.5s ease-in-out,
    transform 0.5s ease-in-out;
}

.left-enter-from,
.left-leave-to {
  opacity: 0;
  transform: scale(0.3) translateX(-100%);
}

.right-enter-active,
.right-leave-active {
  transition:
    opacity 0.5s ease-in-out,
    transform 0.5s ease-in-out;
}

.right-enter-from,
.right-leave-to {
  opacity: 0;
  transform: scale(0.3) translateX(100%);
}

.scale-enter-active,
.scale-leave-active {
  transition:
    opacity 0.5s ease-in-out,
    transform 0.5s ease-in-out;
}

.scale-enter-from,
.scale-leave-to {
  opacity: 0;
  transform: scale(0.3);
}

.none-enter-active,
.none-leave-active {
  transition: opacity 0.5s ease-in-out;
}

.none-enter-from,
.none-leave-to {
  opacity: 1;
}
</style>
