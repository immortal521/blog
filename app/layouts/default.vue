<script setup lang="ts"></script>

<template>
  <div>
    <SiteUtilityBar />
    <BaseHeader />
    <main class="main">
      <slot />
    </main>
    <BaseFooter />
  </div>
</template>

<style lang="less" scoped>
.main {
  margin: 0;
}
</style>
