<script setup lang="ts">
interface Props {
  src?: string;
  title?: string;
}

const { src = "", title = "" } = defineProps<Props>();
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
  height: 400px;
  backdrop-filter: blur(10px);
  background-position: center center;
  background-size: cover;
  background-repeat: no-repeat;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: center;
  transition: height 0.3s ease-in-out;
}

.article-info-container {
  width: 100%;
  height: 100%;
  display: flex;
  overflow: hidden;
  justify-content: center;
  align-items: end;
  backdrop-filter: brightness(0.9);
}

.article-info {
  position: relative;
  bottom: 0;
  width: 100%;
  max-width: 800px;
  overflow: hidden;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.title {
  position: relative;
  padding: 0 20px;
  font-size: 4rem;
  line-height: 2;
  width: 100%;
  font-weight: 700;
  color: var(--article-title-color);
  overflow: hidden;
  text-overflow: ellipsis;
  text-wrap: nowrap;
  text-shadow: 0 0 5px rgb(0 0 0 / 50%);

  &::before {
    position: absolute;
    content: "";
    width: 80px;
    height: 8px;
    bottom: 0;
    border-radius: 5px;
    background-color: var(--color-primary-base);
  }
}

@media (width <= 768px) {
  .article-cover {
    height: 280px;
  }

  .title {
    font-size: 3rem;
  }
}
</style>
