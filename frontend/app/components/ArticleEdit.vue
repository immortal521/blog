<script setup lang="ts">
const content = defineModel<string>("content", {
  default: "",
});

watch(content, () => {
  console.log(JSON.stringify(content.value));
});
</script>

<template>
  <div class="article-edit">
    <div class="edit">
      <textarea v-model="content"></textarea>
    </div>
    <div class="preview">
      <MarkdownRenderer :markdown="content" />
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
}

.edit {
  width: 50%;
  border-radius: 8px;

  textarea {
    height: 100%;
    width: 100%;
    color: var(--text-color-base);
    border-radius: 8px 0 0 8px;
    padding: 10px;
    resize: none;
    outline: none;
    border: none;
    border-right: var(--border-color-divider) 1px solid;
    background-color: transparent;
  }
}

.preview {
  width: 50%;
  padding: 10px;
}
</style>
