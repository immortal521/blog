<script setup lang="ts">
import { ActionButton, PanelButton } from "@/components/FloatActionGroup";

const message = useMessage();

const { $ts } = useI18n();

const backTop = () => {
  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};

const { $localePath } = useI18n();

// TODO: 跳转用户设置界面逻辑, 双 Token 刷新机制
// Navigation logic to the user settings page, with a dual-token refresh mechanism
const toUserSetting = () => {
  navigateTo($localePath("/admin"));
};

const copyRSSFeedUrl = async () => {
  await navigator.clipboard.writeText(window.location.origin + "/api/v1/rss");
  message.success($ts("message.rssCopied"), { keepAliveOnHover: true });
};
</script>

<template>
  <FloatActionGroup>
    <ActionButton
      icon="iconamoon:arrow-up-2-fill"
      :title="$ts('tooltip.backToTop')"
      @click="backTop"
    />
    <ActionButton :title="$ts('tooltip.rss')" icon="ion:logo-rss" @click="copyRSSFeedUrl" />
    <PanelButton animation="right" :title="$ts('tooltip.contect')" icon="ion:mail">
      <ContectGroup />
    </PanelButton>
    <PanelButton
      :title="$ts('tooltip.themeOption')"
      animation="down"
      icon="fluent:apps-list-detail-24-filled"
    >
      <ThemeControlCard />
    </PanelButton>
    <ActionButton
      icon="iconamoon:settings-fill"
      :title="$ts('tooltip.userSetting')"
      @click="toUserSetting"
    />
  </FloatActionGroup>
</template>
