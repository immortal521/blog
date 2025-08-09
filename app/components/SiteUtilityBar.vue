<script setup lang="ts">
import { ActionButton, PanelButton } from "@/components/FloatActionGroup";
const backTop = () => {
	window.scrollTo({
		top: 0,
		behavior: "smooth",
	});
};
</script>

<template>
  <FloatActionGroup>
    <ActionButton
      icon="iconamoon:arrow-up-2-fill"
      :title="$t('tooltip.backToTop')"
      @click="backTop"
    />
    <PanelButton :title="$t('tooltip.contect')" icon="ion:mail">
      <ContectGroup />
    </PanelButton>
    <PanelButton
      :title="$t('tooltip.themeOption')"
      animation="down"
      icon="fluent:apps-list-detail-24-filled"
    >
      <ThemeControlCard />
    </PanelButton>
  </FloatActionGroup>
</template>

<style lang="less" scoped></style>
