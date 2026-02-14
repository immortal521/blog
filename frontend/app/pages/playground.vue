<script setup lang="ts">
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
</script>

<template>
  <div class="container">
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
            background: var(--glass-gradient), var(--bg-card-base);
            color: var(--text-color-primary);
            padding: 10px;
          "
        ></textarea>
      </div>
      <button style="width: 80px; height: 40px; border-radius: 8px" @click="summarize">
        generate
      </button>
      <h2 style="color: var(--text-color-primary)">摘要：</h2>
      <div
        style="
          background: var(--glass-gradient), var(--bg-card-base);
          color: var(--text-color-primary);
        "
      >
        <TransitionGroup name="msgs">
          <span v-for="msg in messages" :key="msg">{{ msg }}</span>
        </TransitionGroup>
      </div>
      <hr />
      <br />
      <br />
      <ArticleEdit v-model:content="article" />
    </main>
  </div>
</template>

<style lang="less" scoped>
.container {
  width: 100vw;
  height: 100vh;
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
  border-bottom: 1px solid var(--border-color-default);
  background: var(--glass-gradient), var(--bg-nav-base);
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
