<script setup lang="ts">
import { ActionButton, PanelButton } from "@/components/FloatActionGroup";
const backTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};

// TODO: 跳转用户设置界面逻辑, 双 Token 刷新机制
// Navigation logic to the user settings page, with a dual-token refresh mechanism
const toUserSetting = () => {};
</script>

<template>
  <FloatActionGroup>
    <ActionButton
      icon="iconamoon:arrow-up-2-fill"
      :title="$t('tooltip.backToTop')"
      @click="backTop"
    />
    <PanelButton animation="right" :title="$t('tooltip.contect')" icon="ion:mail">
      <ContectGroup />
    </PanelButton>
    <PanelButton
      :title="$t('tooltip.themeOption')"
      animation="down"
      icon="fluent:apps-list-detail-24-filled"
    >
      <ThemeControlCard />
    </PanelButton>
    <ActionButton
      icon="iconamoon:settings-fill"
      :title="$t('tooltip.userSetting')"
      @click="toUserSetting"
    />
  </FloatActionGroup>
</template>

<style lang="less" scoped></style>
