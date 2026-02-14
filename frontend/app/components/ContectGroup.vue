<script setup lang="ts">
interface ContectItem {
  url: string;
  icon: string;
  title: string;
}

const { $ts } = useI18n();

const contectItems = computed<ContectItem[]>(() => [
  {
    url: "https://github.com/immortal521",
    icon: "codicon:github",
    title: $ts("tooltip.github"),
  },
  {
    url: "https://qm.qq.com/cgi-bin/qm/qr?k=82gkG9i3QYouTHUsD9h4mk1UGFYZdlxT",
    icon: "streamline-logos:qq-logo-solid",
    title: $ts("tooltip.addQQFriend"),
  },
]);
</script>

<template>
  <ul class="contect">
    <li v-for="item in contectItems" :key="item.url" :title="item.title" class="contect-item">
      <a :href="item.url" target="_blank">
        <Icon :name="item.icon" size="24" class="icon" />
      </a>
    </li>
  </ul>
</template>

<style lang="less" scoped>
.contect {
  display: flex;
  max-width: 210px;
  list-style: none;
  width: max-content;
  min-height: 40px;
  margin-left: auto;
  align-items: center;
  gap: 8px;
  justify-content: end;
  flex-wrap: wrap;
}

.contect-item {
  width: 35px;
  height: 35px;
  flex-shrink: 0;
  background: var(--glass-gradient), var(--bg-nav-base);
  border-radius: var(--radius-nav);
  border: 1px solid var(--border-color-nav);
  box-shadow: var(--shadow-nav);
  transition: var(--transition-nav);

  &:hover {
    background: var(--glass-gradient-strong), var(--bg-nav-hover);
  }

  &:active {
    background: var(--glass-gradient-strong), var(--bg-nav-active);
  }

  a {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}

.icon {
  color: var(--text-color-primary);
}
</style>
