<script setup lang="ts">
const router = useRouter();
const { $localePath } = useI18n();
const handleLogoClicked = () => {
  const currentPath = router.currentRoute.value.fullPath;
  const targetPath = "/";
  if (targetPath !== currentPath) {
    navigateTo($localePath(targetPath));
  }
};
</script>

<template>
  <div class="logo" @click="handleLogoClicked">
    <NuxtImg src="https://q1.qlogo.cn/g?b=qq&nk=188191770&s=640" alt="logo" class="avatar-image" />
  </div>
</template>

<style lang="less" scoped>
.logo {
  height: 45px;
  width: 45px;
  display: flex;
  position: relative;
  align-items: center;
  justify-content: center;
  user-select: none;
  cursor: pointer;
  border-radius: 50%;
  overflow: hidden;

  .avatar-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
