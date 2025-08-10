<script setup lang="ts">
const { isVisible = true } = defineProps<{ isVisible?: boolean }>();
</script>

<template>
  <Transition name="loading-bar">
    <div v-if="isVisible" class="loading">
      <span class="bar" />
    </div>
  </Transition>
</template>

<style lang="less" scoped>
.loading {
  position: absolute;
  top: 0;
  z-index: 10000;
  height: 3px;
  width: 100%;
  animation: change-bar 2.25s infinite 0.75s;
  will-change: transform;

  .bar {
    position: absolute;
    width: 100%;
    height: 0;
    text-indent: -9999px;
    background: #786bcf;
    animation: change-bar 2.25s infinite;

    &:before {
      position: absolute;
      height: 3px;
      background-color: inherit;
      content: "";
      left: 0;
      right: 0;
      width: 100%;
      transform: scaleX(0);
      animation: moving-bar 0.75s infinite;
    }
  }
}

.loading-bar-enter-active,
.loading-bar-leave-active {
  transition: all 0.4s;
}
.loading-bar-enter-from,
.loading-bar-leave-to {
  opacity: 0;
}

@keyframes change-bar {
  0% {
    background-color: #786bcf;
  }
  33.3% {
    background-color: #786bcf;
  }
  33.33% {
    background-color: #78e0f5;
  }
  66.6% {
    background-color: #78e0f5;
  }
  66.66% {
    background-color: #fc6fb1;
  }
  99.9% {
    background-color: #fc6fb1;
  }
}

@keyframes moving-bar {
  0% {
    transform: scaleX(0);
  }
  99.9% {
    transform: scaleX(1);
  }
  100% {
    transform: scaleX(0);
  }
}
</style>
