<script setup lang="ts">
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
</script>

<template>
  <div class="container">
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
  </div>
</template>

<style lang="less" scoped>
.container {
  width: 98vw;
  height: 90vh;
  display: flex;
  margin: 0 auto;
  padding-top: 100px;
  flex-direction: column;
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
