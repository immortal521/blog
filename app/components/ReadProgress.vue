<script setup lang="ts">
const width = ref("0%");
const handleScroll = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight;
  const clientHeight = document.documentElement.clientHeight || document.body.clientHeight;
  const progress = (scrollTop / (scrollHeight - clientHeight)) * 100;
  width.value = `${progress}%`;
};

onMounted(() => {
  handleScroll();
  window.addEventListener("scroll", handleScroll);
});

onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});
</script>

<template>
  <div class="read-progress" :style="{ width }"></div>
</template>

<style lang="less" scoped>
.read-progress {
  position: fixed;
  top: 0;
  left: 0;
  width: 0;
  height: 2px;
  background-color: var(--color-primary-base);
  z-index: 9999;
}
</style>
