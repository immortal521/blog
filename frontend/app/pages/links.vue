<script setup lang="ts">
const { t } = useI18n();
useHead({
  title: t("page.links"),
});

const { data } = await useFetch<{
  data: FriendLink[];
}>("/api/v1/links", {
  method: "get",
});

const links = computed(() => data.value?.data);

const friendLinkFormisShow = ref(false);
const handleShow = () => {
  friendLinkFormisShow.value = true;
};

const linkCard = useTemplateRef<HTMLElement>("link-card");

useAddClassOnIntersect(linkCard, "show");
</script>

<template>
  <ContentPanel style="min-height: 100vh">
    <article class="content">
      <h1 class="title">{{ t("friendLink.title") }}</h1>
      <div class="submit-link-tips">
        申请友链需要满足以下条件:
        <ul>
          <li>请先在您的站点添加本博客的友链信息，再提交申请。</li>
          <li>HTTPS 访问稳定，正常加载。</li>
          <li>友链信息不包含广告、色情、政治等违法信息。</li>
          <li>若长期无法访问、域名失效或出现违规内容，将在不另行通知的情况下移除链接。</li>
          <li>提交申请后，请耐心等待审核，审核通过后会在友链列表中显示。</li>
        </ul>
        <button class="submit-link-button" @click="handleShow">提交友链</button>
      </div>
      <li class="links">
        <ul v-for="link in links" :key="link.id">
          <BaseCard class="item" ref="link-card">
            <a :href="link.url" target="_blank">
              <div class="avatar">
                <img :src="link.avatar" />
              </div>
              <div class="description-container">
                <p class="name">{{ link.name.trim() }}</p>
                <p class="description">{{ link.description }}</p>
              </div>
            </a>
          </BaseCard>
        </ul>
      </li>
    </article>
    <FriendLinkForm v-model:show="friendLinkFormisShow" />
  </ContentPanel>
</template>

<style lang="less" scoped>
.content {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.title {
  text-align: center;
  color: var(--text-color-base);
  display: block;
  height: 130px;
  line-height: 130px;
  font-size: 30px;
  font-weight: 600;
}

.submit-link-tips {
  position: relative;
  background: var(--bg-card-base);
  border-radius: 10px;
  box-shadow: var(--shadow-card-base);
  padding: 20px;
  color: var(--text-color-base);

  li {
    list-style: none;
    padding-left: 2rem;
  }

  .submit-link-button {
    display: block;
    position: relative;
    margin-left: auto;
    width: 100px;
    padding: 0 1.5rem;
    height: 30px;
    background: var(--color-primary-base);
    color: #e0e0e0;
    border-radius: 5px;

    &:hover {
      background-color: var(--color-primary-hover);
    }

    &:active {
      background-color: var(--color-primary-active);
    }
  }
}

.links {
  width: 100%;
  display: grid; /* 启用 Flexbox */
  gap: 20px; /* 卡片之间的间距 */
  margin-top: 20px;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

.item {
  a {
    display: block;
    height: 100%;
    width: 100%;
    padding: 20px;
  }
}

.show {
  animation: fade-in-up 1s ease-in-out;
}

.avatar {
  width: 100%;
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;

  img {
    height: 60px;
    width: 60px;
    object-fit: cover;
    border-radius: 50%;
    border: 2px solid #b0b0b0;
  }
}

.name {
  text-align: center;
  font-size: 1.25em;
  line-height: 2em;
}

.description-container {
  height: 110px;
}

.description {
  overflow: hidden; /* 隐藏溢出内容 */
  text-overflow: ellipsis; /* 显示省略号 */
  display: -webkit-box; /* Webkit 浏览器专用，将元素作为弹性盒子显示 */
  -webkit-box-orient: vertical; /* 垂直方向排列内容 */
  line-clamp: 3;
  -webkit-line-clamp: 3; /* 限制显示为 3 行 */
  color: var(--text-color-muted);
  font-size: 0.925em;
}
</style>
