<script setup lang="ts">
definePageMeta({
  layout: false,
});

const content = ref("");
const summary = ref("");
const loading = ref(false);

let eventSource: EventSource | null = null;

const closeSSE = () => {
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
};

const summarize = async () => {
  if (!content.value.trim()) {
    return;
  }

  closeSSE();

  summary.value = "";
  loading.value = true;

  try {
    const data = await $fetch<{
      sessionId: string;
    }>("/api/v1/model/summarize", {
      method: "POST",
      body: {
        content: content.value,
      },
    });

    console.log("session:", data.sessionId);

    eventSource = new EventSource(`/api/v1/model/summarize/${data.sessionId}`);

    eventSource.onopen = () => {
      console.log("SSE connected");
    };

    eventSource.onmessage = (event) => {
      console.log("message:", event.data);
      summary.value += event.data;
    };

    eventSource.addEventListener("done", () => {
      console.log("SSE done");
      loading.value = false;
      closeSSE();
    });

    eventSource.onerror = (err) => {
      console.error("SSE error:", err);
      loading.value = false;
      closeSSE();
    };
  } catch (err) {
    console.error("summarize error:", err);
    loading.value = false;
  }
};

onBeforeUnmount(() => {
  closeSSE();
});
</script>

<template>
  <div class="container">
    <textarea v-model="content" class="content-input" placeholder="请输入需要摘要的内容" />

    <button class="submit-btn" :disabled="loading" @click="summarize">
      {{ loading ? "生成中..." : "生成摘要" }}
    </button>

    <div class="summary">
      {{ summary }}
    </div>
  </div>
</template>

<style scoped lang="less">
.container {
  width: 100%;
  min-height: 100vh;
  padding: 40px;
  color: var(--text-color-primary);

  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-input {
  width: 100%;
  min-height: 240px;
  padding: 12px;
  resize: vertical;
}

.submit-btn {
  width: 120px;
  height: 40px;
}

.summary {
  white-space: pre-wrap;
  word-break: break-word;
  text-align: left;
  padding: 16px;
  border: 1px solid #ddd;
  min-height: 120px;
}
</style>
