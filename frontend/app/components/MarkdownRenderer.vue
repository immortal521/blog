<script setup lang="ts">
import "viewerjs/dist/viewer.css";
import { parseMarkdownToVNode } from "@/utils/parseMarkdown";
import type { VNodeChild } from "vue";

const props = defineProps({
  markdown: {
    type: String,
    required: true,
  },
});

const renderedVNode = shallowRef<VNodeChild>([]);
const containerRef = useTemplateRef("container");

const { update } = useViewer(containerRef);

const renderMarkdown = () => {
  const { content } = parseMarkdownToVNode(props.markdown, { toc: true });
  renderedVNode.value = content;
};

watch(
  () => props.markdown,
  async () => {
    renderMarkdown();
    await update();
  },
  { immediate: true },
);
</script>

<template>
  <div ref="container" class="container">
    <component :is="item" v-for="item in renderedVNode" :key="item" />
  </div>
</template>

<style lang="less" scoped>
.container:deep(.img) {
  max-width: 100%;
  margin: 0 auto;
  border-radius: 5px;
  cursor: pointer;
  box-shadow: var(--shadow-sm);
}

:deep(em) {
  font-style: italic;
}

:deep(p),
:deep(blockquote),
:deep(ul),
:deep(ol),
:deep(dl),
:deep(table) {
  margin: 0.8rem 0;
}

:deep(h1),
:deep(h2),
:deep(h3),
:deep(h4),
:deep(h5),
:deep(h6) {
  font-weight: 700;
  color: var(--color-header);
  line-height: 1.5;
  margin-top: 2rem;
  margin-bottom: 0.8rem;
}

:deep(h1) {
  font-size: 3.2rem;
  padding-bottom: 0.48rem;
}

:deep(h2) {
  font-size: 2.6rem;
  padding-bottom: 0.32rem;
}

:deep(h3) {
  font-size: 2.2rem;
}

:deep(h4) {
  font-size: 1.8rem;
  font-weight: 600;
}

:deep(h5) {
  font-size: 1.7rem;
  font-weight: 600;
}

:deep(h6) {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-color-muted);
}

:deep(p) {
  color: var(--text-color-base);
  overflow-wrap: break-word;
}

:deep(:not(pre.shiki) > code) {
  display: inline-block;
  padding: 0.2em 0.4em;
  font-family: "Maple Mono", "Noto Sans SC", monospace;
  font-size: 0.92em;
  color: var(--color-primary-base); /* 可自定义颜色变量 */
  background-color: var(--bg-code);
  border-radius: 4px;
  line-height: 1.4;
  white-space: break-spaces;
  overflow-wrap: break-word;
}

:deep(table) {
  width: 100%;
  border-collapse: collapse;
  border-radius: 8px;
  background-color: var(--bg-content);
  overflow: hidden;
  font-size: var(--font-size-table-row);
  color: var(--text-color-base);
}

:deep(thead) {
  background-color: var(--table-head-bg);
  color: var(--color-header);
  font-weight: bold;
  font-size: var(--font-size-table-header);
}

:deep(thead) th {
  padding: 0.75em 1em;
  text-align: left;
  border-bottom: 1px solid var(--table-head-border);
}

:deep(tbody) tr {
  border-bottom: 1px solid var(--border-table);
  font-size: var(--font-size-table-row);
}

:deep(tbody) td {
  padding: 0.65em 1em;
}

:deep(tbody) tr:nth-child(odd) {
  background-color: var(--table-row-even-bg);
}

:deep(tbody) tr:hover {
  background-color: var(--table-row-hover-bg);
}

:deep(blockquote) {
  border-radius: 5px;
  padding: 10px 16px;
  background-color: var(--bg-card-base);
  position: relative;
  border-left: none;

  &::before {
    display: block;
    position: absolute;
    content: "";
    width: 4px;
    left: 0;
    top: 0;
    height: 100%;
    background-color: var(--color-primary-base);
    border-radius: 2px;
  }
}

:deep(ul),
:deep(ol) {
  margin: 1em 0;
  padding-left: 1.5em;
  color: var(--text-color-base);
  font-size: 1.5rem;
  line-height: 1.75;
}

:deep(ul li),
:deep(ol li) {
  margin: 0.3rem 0;
  padding-left: 0.25rem;
  position: relative;
}

/* 无序列表圆点样式 */
:deep(ul li::marker) {
  color: var(--text-color-muted);
}

/* 有序列表数字样式 */
:deep(ol li::marker) {
  font-weight: bold;
  color: var(--text-color-muted);
}

:deep(ul ul),
:deep(ul ol),
:deep(ol ul),
:deep(ol ol) {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
  padding-left: 1.5em;
}

:deep(.contains-task-list) {
  padding-left: 0;
  list-style: none;
}

:deep(mark) {
  padding: 0 0.5rem;
  color: #fff;
  background-color: var(--color-primary-base);
  border-radius: 10px;
}

:deep(u) {
  text-decoration-line: underline;
  text-decoration-color: var(--color-primary-base);
  text-underline-offset: 0.5rem;
  text-decoration-thickness: 2px;
}
</style>
