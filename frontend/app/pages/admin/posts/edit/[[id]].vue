<script setup lang="ts">
definePageMeta({
  layout: false,
});

// const route = useRoute();

// const id = computed<number>(() => Number(route.params.id) || -1);

const content = ref<string>("");
const title = ref<string>("");
// const summary = ref<string>("");

// const summarize = async () => {
//   summary.value = "";
//   const data = await $fetch<{ sessionId: string }>("/api/v1/model/summarize", {
//     method: "post",
//     body: {
//       content: content.value,
//     },
//   });
//
//   const es = new EventSource("/api/v1/model/summarize/" + data.sessionId);
//
//   es.onmessage = (event) => {
//     summary.value += event.data;
//   };
//
//   es.addEventListener("done", () => {
//     es.close();
//   });
//
//   es.addEventListener("error", () => {
//     es.close();
//   });
// };
</script>

<template>
  <div>
    <NuxtLayout name="admin-sub-layout">
      <template #actions>
        <div class="editor-header">
          <BaseInput v-model="title" class="title-input" placeholder="输入文章标题" />
        </div>
      </template>
      <div class="editor-body">
        <ArticleEdit v-model:content="content" class="editor" />
      </div>
      <div class="editor-footer">
        <div class="detail">
          <div>行数：0 字符数：0</div>
        </div>
        <div class="actions">
          <button class="btn trash">存入草稿箱</button>
          <button class="btn publish">发布</button>
        </div>
      </div>
    </NuxtLayout>
  </div>
</template>

<style lang="less" scoped>
.editor-body {
  display: flex;
  width: 100%;
  height: calc(100% - 60px);
  min-height: 0;
  min-width: 0;
}

.editor {
  width: 100%;
  height: 100%;
}

.editor-header {
  display: flex;
  width: 100%;
  height: 100%;
  justify-content: center;
  margin: 0 auto;
  max-width: 1200px;
}

.title-input {
  height: 100%;
  border: none;
  font-size: 2rem;
  width: 100%;

  &:focus-within {
    border: none;
    box-shadow: none;
  }
}

.editor-footer {
  width: 100%;
  height: 50px;
  margin-top: 10px;
  border: 1px solid var(--border-color-divider);
  border-radius: 8px;
  box-shadow: var(--shadow-sm);
  display: flex;
  padding: 0 8px;

  .detail {
    flex-grow: 1;
    line-height: 50px;
    font-size: 1.4rem;
  }

  .actions {
    flex-shrink: 0;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 8px;
    gap: 5px;

    .btn {
      height: 100%;
      padding: 8px;
      line-height: 1;
      border-radius: 8px;
      color: var(--text-color-base);
      background-color: transparent;
      box-shadow: var(--shadow-sm);
    }

    .trash {
      border: 1px solid var(--color-primary-base);

      &:hover {
        border: 1px solid var(--color-primary-hover);
      }

      &:active {
        border: 1px solid var(--color-primary-active);
      }
    }

    .publish {
      background-color: var(--color-primary-base);

      &:hover {
        background-color: var(--color-primary-hover);
      }

      &:active {
        background-color: var(--color-primary-active);
      }
    }
  }
}
</style>
