<script setup lang="ts">
import BaseImage from "~/components/BaseImage.vue";
</script>

<template>
  <div>
    <WelcomePanel />
    <main>
      <BaseImage src="https://gitee.com/Immortal112/image/raw/master/background/28.webp" lazy class="img"/>
    </main>
  </div>
</template>

<style lang="less" scoped>
.img {
  width: 200px;
  border-radius: 20px;
}
</style>
