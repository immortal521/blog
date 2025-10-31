<script setup lang="ts">
definePageMeta({
  layout: "admin-layout",
});

const sidebar = useTemplateRef("sidebar");
console.log(sidebar);

const message = useMessage();
</script>

<template>
  <div class="container">
    <BaseSidebar ref="sidebar" :width="200" :width-in-fold="60" style="height: 100vh" />
    <div style="display: flex; height: 30px; gap: 10px; transition: transform 0.3s ease-in-out">
      <button @click="sidebar?.open">Open</button>
      <button @click="sidebar?.close">Close</button>
      <button @click="sidebar?.fold">Fold</button>
      <button @click="sidebar?.unfold">Unfold</button>
      <button @click="message.success('123', { duration: 5000000 })">success</button>
      <button
        @click="message.info('dhflkjahsjdfklhaslkjfhsas', { duration: 5000000, size: 'small' })"
      >
        info
      </button>
      <button @click="message.error('aslkjfhsas', { duration: 5000000, size: 'large' })">
        error
      </button>
      <button
        @click="
          message.warning('dhflkjahsjdfklhaslkjfhsas', {
            duration: 5000000,
            icon: { name: 'fluent:info-sparkle-24-filled' },
            closable: true,
          })
        "
      >
        warning
      </button>
    </div>
  </div>
</template>

<style lang="less" scoped>
.container {
  width: 100vw;
  height: 100vh;
  display: flex;
}
</style>
