<script setup lang="ts">
const content = defineModel<string>("content", {
  default: "",
});
</script>

<template>
  <div class="article-edit">
    <div class="toolbar">
      <div class="tools"></div>
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
  margin: 0 2%;
  border: 1.5px solid var(--border-color-divider);
  border-radius: 8px;
  box-shadow: var(--shadow-md);
  width: 96%;
  min-height: 400px;
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

  .tools {
    width: calc(100% - 4px);
    height: 30px;
    background-color: var(--color-primary-bg-muted);
    border-radius: 8px;
  }
}

.main {
  width: 100%;
  flex: 1;
  min-height: 0;
  display: flex;
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
  padding: 10px;
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
