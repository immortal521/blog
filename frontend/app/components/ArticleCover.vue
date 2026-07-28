<script setup lang="ts">
interface Props {
  src?: string;
  title?: string;
}

const { src = "", title = "" } = defineProps<Props>();
onMounted(() => {
  console.log(src);
});
</script>

<template>
  <header class="article-cover" :style="{ backgroundImage: `url(${src})` }">
    <div class="article-info-container">
      <div class="article-info">
        <h1 class="title">
          {{ title }}
        </h1>
      </div>
    </div>
  </header>
</template>

<style lang="less" scoped>
.article-cover {
  position: relative;
  width: 100%;
  height: 450px;
  background-position: center center;
  background-size: cover;
  background-repeat: no-repeat;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  transition: height 0.3s ease-in-out;

  &::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 60%;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
    pointer-events: none;
  }
}

.article-info-container {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: end;
  backdrop-filter: brightness(0.85);
  transition: backdrop-filter 0.3s ease;

  &:hover {
    backdrop-filter: brightness(0.9);
  }
}

.article-info {
  position: relative;
  bottom: 0;
  width: 100%;
  max-width: 800px;
  overflow: hidden;
  margin-bottom: 32px;
  padding: 0 24px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
}

.title {
  position: relative;
  font-size: 3.6rem;
  line-height: 1.4;
  width: 100%;
  font-weight: 700;
  color: var(--article-title-color);
  text-shadow: 0 2px 8px rgb(0 0 0 / 60%);
  word-wrap: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
  animation: title-slide-up 0.8s ease-out;

  &::before {
    position: absolute;
    content: "";
    width: 60px;
    height: 4px;
    bottom: -12px;
    left: 0;
    border-radius: 2px;
    background: linear-gradient(90deg, var(--color-primary-base), var(--color-primary-hover));
    animation: title-underline 1s ease-out 0.5s both;
  }
}

@keyframes title-slide-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes title-underline {
  from {
    width: 0;
    opacity: 0;
  }
  to {
    width: 60px;
    opacity: 1;
  }
}

@media (width <= 768px) {
  .article-cover {
    height: 320px;
  }

  .title {
    font-size: 2.4rem;
  }

  .article-info {
    margin-bottom: 24px;
    padding: 0 16px;
  }
}
</style>
