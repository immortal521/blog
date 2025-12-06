<script setup lang="ts">
const isSticky = ref(false);

const header = useTemplateRef("header");

const scrollHandler = () => {
  if (!header.value) return;

  const scrollTop = document.documentElement.scrollTop;

  if (scrollTop > 100) {
    isSticky.value = true;
  } else {
    isSticky.value = false;
  }
};

onMounted(() => {
  window.addEventListener("scroll", scrollHandler);
});

onUnmounted(() => {
  window.removeEventListener("scroll", scrollHandler);
});
</script>

<template>
  <header ref="header" :class="{ header, sticky: isSticky }">
    <NavbarWrapper />
  </header>
</template>

<style lang="less" scoped>
.header {
  position: fixed;
  width: calc(100vw - 40px);
  height: 60px;
  left: 20px;
  top: 20px;
  border-radius: 15px;
  overflow: hidden;
  background-color: var(--bg-nav-base);
  border: 1px solid var(--border-color-nav);
  backdrop-filter: blur(var(--nav-blur));
  transition:
    left 0.8s ease-in-out,
    top 0.8s ease-in-out,
    width 0.8s ease-in-out,
    border-radius 0.8s ease-in-out,
    background-color 0.3s ease-in-out,
    border-color 0.3s ease-in-out;
  animation: scale-in 0.5s ease-in-out;
  box-shadow: var(--shadow-nav);
  z-index: 1999;

  &:hover {
    background-color: var(--bg-nav-hover);
  }

  @media (width <= 768px) {
    width: 100%;
    left: 0;
    top: 0;
    border-radius: 0;
  }
}

.sticky {
  border-radius: 0;
  width: 100%;
  left: 0;
  top: 0;
}
</style>
