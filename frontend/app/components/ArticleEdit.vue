<script setup lang="ts">
const content = defineModel<string>("content", {
  default: "",
});

const editorRef = useTemplateRef<HTMLDivElement>("editor");

const isFullscreen = ref(false);

const fullscreen = () => {
  if (isFullscreen.value) {
    document.exitFullscreen();
    isFullscreen.value = false;
  } else {
    editorRef.value?.requestFullscreen();
    isFullscreen.value = true;
  }
};
</script>

<template>
  <div ref="editor" class="article-edit">
    <div class="toolbar">
      <div class="tools-left"></div>
      <div class="tools-right">
        <button class="btn" @click="fullscreen">
          <Icon v-if="!isFullscreen" name="mingcute:fullscreen-fill" size="18" />
          <Icon v-else name="mingcute:fullscreen-exit-fill" size="18" />
        </button>
      </div>
    </div>
    <div class="main">
      <div class="edit">
        <textarea v-model="content"></textarea>
      </div>
      <div class="preview">
        <MarkdownRenderer :markdown="content" />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.article-edit {
  border: 1.5px solid var(--border-color-divider);
  border-radius: 8px;
  box-shadow: var(--shadow-md);
  width: 96%;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.toolbar {
  height: 40px;
  flex: 0 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 5px 15px;
  border-bottom: 1px solid var(--border-color-divider);

  .tools-left {
    margin-right: auto;
  }

  .tools-right {
    margin-left: auto;
  }

  .btn {
    display: flex;
    justify-content: center;
    align-items: center;
    color: var(--text-color-base);
    background-color: transparent;
    padding: 4px;
    border-radius: 5px;

    &:hover {
      background-color: var(--border-color-base);
    }

    &:active {
      background-color: var(--border-color-divider);
      transform: scale(0.95);
    }
  }
}

.main {
  width: 100%;
  flex: 1;
  min-height: 0;
  display: flex;
  padding: 8px 0;
}

.edit {
  width: 50%;
  min-width: 0;
  display: flex;
}

.edit textarea {
  flex: 1;
  width: 100%;
  color: var(--text-color-base);
  border-radius: 8px 0 0 8px;
  padding: 10px 20px;
  resize: none;
  font-size: 1.5rem;
  line-height: 1.5;
  outline: none;
  border: none;
  border-right: var(--border-color-divider) 1px solid;
  background-color: transparent;
}

.preview {
  width: 50%;
  padding: 10px;
  min-width: 0;
  overflow: auto; /* 预览区内容多时滚动 */
}
</style>
