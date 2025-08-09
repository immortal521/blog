<script setup lang="ts">
// import { getLinksApi } from "@/api/link";
const { t } = useI18n();

const { data } = await useFetch<{
  data: FriendLink[];
}>("/api/links", {
  method: "get",
});

const links = computed(() => data.value?.data);
</script>

<template>
  <ContentPanel>
    <article class="content">
      <h1 class="title">{{ t("friendLink.title") }}</h1>
      <li class="links">
        <ul v-for="link in links" :key="link.id">
          <BaseCard class="item">
            <a :href="link.url" target="_blank">
              <div class="avatar">
                <img :src="link.avatar" />
              </div>
              <div class="description-container">
                <p class="name">{{ link.name }}</p>
                <p class="description">{{ link.description }}</p>
              </div>
            </a>
          </BaseCard>
        </ul>
      </li>
    </article>
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

.links {
  width: 100%;
  display: grid; /* 启用 Flexbox */
  gap: 20px; /* 卡片之间的间距 */
  padding: 20px;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
}

.item {
  animation: fade-in-up 1s ease-in-out;

  a {
    display: block;
    height: 100%;
    width: 100%;
    padding: 20px;
  }
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
