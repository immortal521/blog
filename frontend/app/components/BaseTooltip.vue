<script setup lang="ts">
type Placement = "top" | "left" | "right" | "bottom";
interface Props {
  placement?: Placement | "auto";
  maxWidth?: number | string;
}
interface Style {
  left: string;
  top: string;
  maxWidth: string;
}

const open = ref(false);

const { placement = "auto", maxWidth = 320 } = defineProps<Props>();

const targetRef = useTemplateRef<HTMLDivElement>("target");
const tooltipRef = useTemplateRef<HTMLDivElement>("tooltip");

const placements: Placement[] = ["top", "bottom", "left", "right"];
const acturalPlacement = ref<Placement>("top");

let openTimer: ReturnType<typeof setTimeout> | null = null;
let closeTimer: ReturnType<typeof setTimeout> | null = null;

const toCssSize = (value: number | string) => (typeof value === "number" ? `${value}px` : value);

const style = ref<Style>({
  left: "0",
  top: "0",
  maxWidth: toCssSize(maxWidth),
});

const computePos = (
  target: DOMRect,
  tooltipSize: { width: number; height: number },
  placement: Placement,
) => {
  const { width, height } = tooltipSize;
  let left = 0;
  let top = 0;

  switch (placement) {
    case "top":
      left = target.left + target.width / 2 - width / 2;
      top = target.top - height;
      break;
    case "bottom":
      left = target.left + target.width / 2 - width / 2;
      top = target.bottom;
      break;
    case "left":
      left = target.left - width;
      top = target.top + target.height / 2 - height / 2;
      break;
    case "right":
      left = target.right;
      top = target.top + target.height / 2 - height / 2;
      break;
  }

  return {
    left,
    top,
  };
};

const overflowScope = (
  left: number,
  top: number,
  tooltipSize: { width: number; height: number },
) => {
  const vw = window.innerWidth;
  const vh = window.innerHeight;
  const { width, height } = tooltipSize;

  const overLeft = Math.max(0, -left);
  const overRight = Math.max(0, left + width - vw);
  const overTop = Math.max(0, -top);
  const overBottom = Math.max(0, top + height - vh);

  return overLeft + overRight + overTop + overBottom;
};

const pickPlacement = (targetRect: DOMRect, tooltipSize: { width: number; height: number }) => {
  let list = placements;
  let best: Placement = "top";
  let bestScore = Number.MAX_SAFE_INTEGER;
  if (placement !== "auto") {
    best = placement;
    list = [placement, ...placements.filter((p) => p !== placement)];
  }

  for (const p of list) {
    const { left, top } = computePos(targetRect, tooltipSize, p);
    const score = overflowScope(left, top, tooltipSize);
    if (score < bestScore) {
      bestScore = score;
      best = p;
      if (bestScore === 0) break;
    }
  }

  return best;
};

const caculateStyle: () => Style | null = () => {
  const target = targetRef.value;
  const tooltip = tooltipRef.value;
  if (!target || !tooltip) return null;

  const targetRect = target.getBoundingClientRect();
  const tooltipSize = {
    width: tooltip.offsetWidth,
    height: tooltip.offsetHeight,
  };

  const chosen = pickPlacement(targetRect, tooltipSize);
  acturalPlacement.value = chosen;

  const { left, top } = computePos(targetRect, tooltipSize, acturalPlacement.value);

  return {
    left: toCssSize(left),
    top: toCssSize(top),
    maxWidth: toCssSize(maxWidth),
  };
};

const clearTimers = () => {
  if (openTimer) {
    clearTimeout(openTimer);
    openTimer = null;
  }
  if (closeTimer) {
    clearTimeout(closeTimer);
    closeTimer = null;
  }
};

let rafId = 0;
const updatePosition = () => {
  if (!open.value) return;
  if (rafId) return;
  rafId = requestAnimationFrame(() => {
    rafId = 0;
    const pos = caculateStyle();
    if (pos) style.value = pos;
  });
};

const scheduleOpen = () => {
  clearTimers();
  openTimer = setTimeout(async () => {
    open.value = true;
    await nextTick();
    updatePosition();
  }, 80);
};

const scheduleClose = () => {
  clearTimers();
  closeTimer = setTimeout(() => {
    open.value = false;
  }, 80);
};

watchEffect((onCleanup) => {
  if (!open.value) return;

  const handler = () => updatePosition();
  window.addEventListener("scroll", handler, true);
  window.addEventListener("resize", handler);

  onCleanup(() => {
    window.removeEventListener("scroll", handler, true);
    window.removeEventListener("resize", handler);
  });
});

onBeforeUnmount(() => {
  clearTimers();
  if (rafId) cancelAnimationFrame(rafId);
});
</script>

<template>
  <div
    ref="target"
    class="tooltip-target"
    @mouseenter="scheduleOpen"
    @mouseleave="scheduleClose"
    @focusin="scheduleOpen"
    @focusout="scheduleClose"
  >
    <slot />
  </div>
  <Teleport to="body">
    <Transition name="tooltip-fader">
      <div
        v-if="open"
        ref="tooltip"
        class="tooltip"
        :style
        @mouseenter="clearTimers"
        @mouseleave="scheduleClose"
      >
        <div class="tooltip-content">
          content
          <div class="tooltip-arrow"></div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="less" scoped>
.tooltip {
  position: fixed;
  font-size: 1.4rem;
  line-height: 1.2;
  z-index: 9999;
  transform-origin: bottom center;
  padding: 10px;
}

.tooltip-target {
  cursor: default;
  display: inline-block;
}

.tooltip-content {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  background-color: var(--bg-card-base);
  border: 1px solid var(--border-color-default);
  color: var(--text-color-primary);
  box-shadow: var(--shadow-sm);
  position: relative;
  pointer-events: auto;
}

.tooltip-fader-enter-active,
.tooltip-fader-leave-active {
  transition: opacity 0.1s ease;
}

.tooltip-fader-enter-from,
.tooltip-fader-leave-to {
  opacity: 0;
}
</style>
