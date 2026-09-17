<script lang="ts" setup>
import { useEditor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import { Markdown } from "@tiptap/markdown";
import { watch } from "vue";

const content = defineModel<string>("content", { default: "" });

const editor = useEditor({
  content: content.value,
  extensions: [StarterKit, Markdown],
  onUpdate: ({ editor }) => {
    content.value = editor.getMarkdown();
  },
  autofocus: true,
});

watch(content, (value) => {
  if (!editor.value) {
    return;
  }

  if (value === editor.value.getMarkdown()) {
    return;
  }

  editor.value.commands.setContent(value, {
    contentType: "markdown",
    emitUpdate: false,
  });
});
</script>

<template>
  <EditorContent :editor="editor" class="editor" />
</template>

<style lang="less">
.editor {
  width: 100%;
  height: 100%;

  .tiptap {
    outline: none;
  }
}
</style>
