<script setup lang="ts">
definePageMeta({
  layout: false,
});

// const route = useRoute();

// const id = computed<number>(() => Number(route.params.id) || -1);

const content = ref<string>("");
const title = ref<string>("");
const summary = ref<string>("");

const summarize = async () => {
  summary.value = "";
  const data = await $fetch<{ sessionId: string }>("/api/v1/model/summarize", {
    method: "post",
    body: {
      content: content.value,
    },
  });

  const es = new EventSource("/api/v1/model/summarize/" + data.sessionId);

  es.onmessage = (event) => {
    setTimeout(() => {
      summary.value += event.data;
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
  <div>
    <NuxtLayout name="admin-sub-layout">
      <template #actions> 一些操作 </template>
      <div class="container">
        <ArticleEdit v-model:content="content" class="editor" />
        <div class="settings">
          <div class="card">
            <h2>标题</h2>
            <BaseInput v-model="title" style="border-radius: 10px" />
          </div>
          <div class="card" style="display: flex; flex-direction: column">
            <textarea
              v-model="summary"
              style="
                resize: none;
                background-color: transparent;
                outline: none;
                border: 1px solid var(--border-color-divider);
                border-radius: 8px;
                color: var(--text-color-base);
                font-size: 1.6rem;
                line-height: 1.5;
                padding: 5px;
                height: 300px;
              "
            ></textarea>
            <button
              style="
                height: 30px;
                background-color: transparent;
                margin-top: 5px;
                color: var(--text-color-base);
                border-radius: 8px;
                border: 1px solid var(--border-color-divider);
              "
              @click="summarize"
            >
              AI 生成摘要
            </button>
          </div>
        </div>
      </div>
    </NuxtLayout>
  </div>
</template>

<style lang="less" scoped>
.container {
  display: flex;
  width: 100%;
  height: 100%;
  min-height: 0;
  min-width: 0;
}

.editor {
  min-width: 0;
  height: 100%;
}

.settings {
  flex: 0 0 320px;
  max-width: 40%;
  min-width: 240px;
  height: 100%;
  margin-left: 10px;
}

.card {
  width: 100%;
  padding: 10px;
  border: 1.5px solid var(--border-color-divider);
  margin-bottom: 10px;
  border-radius: 8px;
  box-shadow: var(--shadow-md);
}

@media (width < 768px) {
  .container {
    display: block;
  }

  .editor {
    width: 100%;
    height: 100%;
  }

  .settings {
    width: 100%;
    height: auto;
  }
}
</style>
