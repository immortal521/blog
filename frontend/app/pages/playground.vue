<script setup lang="ts">
import type {
  SidebarActionItem,
  SidebarDividerItem,
  SidebarGroupItem,
  SidebarLinkItem,
  SidebarSection,
} from "~/components/BaseSidebar/types";

definePageMeta({
  layout: "admin-layout",
});
const messages = ref<string[]>([]);
const article = ref<string>(""); // 用户输入的文章

const summarize = async () => {
  messages.value = [];
  const data = await $fetch<{ sessionId: string }>("/api/v1/model/summarize", {
    method: "post",
    body: {
      content: article.value,
    },
  });

  const es = new EventSource("/api/v1/model/summarize/" + data.sessionId);

  es.onmessage = (event) => {
    setTimeout(() => {
      messages.value.push(event.data);
      console.log(event.data);
    }, 1000);
  };

  es.addEventListener("done", () => {
    console.log("生成完成");
    es.close();
  });

  es.addEventListener("error", (e) => {
    console.error("SSE 错误:", e);
    es.close();
  });
};

const open = ref(true);

onMounted(() => {
  open.value = window.innerWidth >= 768;
});

const link: SidebarLinkItem = {
  to: "1",
  type: "link",
  key: "",
  icon: "ri:link",
  label: "Link",
};

const divider: SidebarDividerItem = {
  type: "divider",
  key: "",
  title: "Divider",
};

const action: SidebarActionItem = {
  action: function (): void {
    collapsed.value = !collapsed.value;
  },
  label: "Action",
  type: "action",
  key: "",
};

const group: SidebarGroupItem = {
  children: [
    {
      type: "link",
      label: "playground",
      to: "/playground",
      key: "/playground",
    },
  ],
  type: "group",
  key: "group",
  label: "Group",
};

const section = ref<SidebarSection>({
  type: "section",
  items: [link, divider, action, group],
  key: "",
  label: "section",
});

const openKeys = ref(new Set<string>());
openKeys.value.add(group.key);

const collapsed = ref(false);
</script>

<template>
  <div class="container">
    <div class="sidebar" :class="{ collapsed }">
      <BaseSidebarItem :item="link" :open-keys="openKeys" :collapsed />
      <BaseSidebarItem :item="group" :open-keys="openKeys" :collapsed />
      <BaseSidebarItem :item="action" :open-keys="openKeys" :collapsed />
      <BaseSidebarItem :item="divider" :open-keys="openKeys" :collapsed />
      <BaseSidebarSection :section="section" :open-keys="openKeys" :collapsed />
    </div>
    <div class="right">
      <header class="header">
        <button @click="open = !open">{{ open }}</button>
      </header>
      <main class="main">
        <h1>生成摘要</h1>
        <div class="article-box">
          <textarea
            v-model="article"
            style="
              height: 200px;
              width: 100%;
              background-color: var(--bg-card-base);
              color: var(--text-color-base);
              padding: 10px;
            "
          ></textarea>
        </div>
        <button style="width: 80px; height: 40px; border-radius: 8px" @click="summarize">
          generate
        </button>
        <h2 style="color: var(--text-color-base)">摘要：</h2>
        <div style="background-color: var(--bg-card-base); color: var(--text-color-base)">
          <TransitionGroup name="msgs">
            <span v-for="msg in messages" :key="msg">{{ msg }}</span>
          </TransitionGroup>
        </div>
        <hr />
      </main>
    </div>
  </div>
</template>

<style lang="less" scoped>
.container {
  width: 100vw;
  height: 100vh;
  display: flex;
}

.sidebar {
  width: 240px;
  height: 100vh;
  transition: width 0.3s ease;
}

.collapsed {
  width: 60px;
}

.right {
  width: 100%;
}

.header {
  height: 60px;
  width: 100%;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-nav-base);
  backdrop-filter: var(--filter-blur-sm);
}

.main {
  padding: 10px;
}

.msgs-enter-active,
.msgs-leave-active {
  transition: all 0.5s ease;
}

.msgs-enter-from,
.msgs-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
