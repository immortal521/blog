<script setup lang="ts">
import { generateMdAst } from "@/utils/parseMarkdown";

const props = defineProps({
  markdown: {
    type: String,
    required: true,
  },
  // enableCodeHighlight prop 已移除
});

const renderedVNode = shallowRef<VNode | VNode[] | null>(null);

const renderMarkdown = async () => {
  renderedVNode.value = (await generateMdAst(props.markdown)) as VNode;
};

watch(() => props.markdown, renderMarkdown, { immediate: true });
</script>

<template>
  <div ref="container" class="container">
    <component :is="item" v-for="item in renderedVNode" :key="item" />
  </div>
</template>

<style lang="less" scoped>
.container:deep(img) {
  max-width: 100%;
  height: auto;
  margin: 0 auto;
  border-radius: 5px;
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

:deep(code:not(.hljs)) {
  display: inline-block;
  padding: 0.2em 0.4em;
  font-family: "MapleMono", "Noto Sans SC", monospace;
  font-size: 0.92em;
  color: var(--color-primary-base); /* 可自定义颜色变量 */
  background-color: var(--bg-code);
  border-radius: 4px;
  line-height: 1.4;
  white-space: break-spaces;
  word-break: break-word;
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

/* 列表嵌套支持（缩进一致性） */
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
  text-decoration-line: underline; /* 首先需要有下划线 */
  text-decoration-color: var(--color-primary-base);
  text-underline-offset: 0.5rem;
  text-decoration-thickness: 2px;
}

:deep(pre) {
  position: relative;
}

:deep(pre) .clipboard-button {
  position: absolute;
  right: 1rem;
  top: 1rem;
  cursor: pointer;
  opacity: 0;
  background-color: var(--color-primary-base);
  width: 6rem;
  height: 3rem;
  color: #f0f0f0;
  border-radius: 5px;
  box-shadow: var(--shadow-sm);
  transition: opacity 0.3s ease-in-out;

  &:hover {
    background-color: var(--color-primary-hover);
  }

  &:active {
    background-color: var(--color-primary-active);
  }
}

:deep(pre):hover .clipboard-button {
  opacity: 1;
}
</style>

<style lang="less">
.hljs {
  display: block;
  position: relative;
  background: var(--bg-code);
  color: var(--text-color-base);
  font-family: "MapleMono", monospace;
  font-size: 1.5rem;
  line-height: 1.5;
  padding: 3rem 1.6rem 1.5rem 1.6rem;
  border-radius: var(--radius-card, 0.5rem);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  overflow-x: auto;
  transition: var(--transition-base, 0.3s);
  margin-top: 20px;
}

.code-line {
  height: 1.5rem;
  letter-spacing: 0.05rem;
}

.hljs-comment,
.hljs-quote {
  color: #d4d0ab;
}

.hljs-variable,
.hljs-template-variable,
.hljs-tag,
.hljs-name,
.hljs-selector-id,
.hljs-selector-class,
.hljs-regexp,
.hljs-deletion {
  color: #a9b7c6;
}

/* Orange */
.hljs-number,
.hljs-built_in,
.hljs-literal,
.hljs-type,
.hljs-params,
.hljs-meta,
.hljs-link {
  color: #6897bb;
}

/* Yellow */
.hljs-attribute {
  color: #ffd700;
}

/* Green */
.hljs-string,
.hljs-symbol,
.hljs-bullet,
.hljs-addition {
  color: #6a8759;
}

/* Blue */
.hljs-title,
.hljs-section {
  color: var(--color-primary-base);
}

.hljs-keyword,
.hljs-selector-tag {
  color: #cc7832;
}

.hljs-emphasis {
  font-style: italic;
}

.hljs-strong {
  font-weight: bold;
}

.code-line.numbered-code-line {
  position: relative;
  width: 100%;
  margin-left: 3rem;
  padding-left: 0.25rem;
}

.code-line.numbered-code-line::before {
  content: attr(data-line-number);
  position: absolute;
  display: inline-block;
  left: -3.25rem;
  padding-right: 0.5rem;
  width: 3rem;
  text-align: right;
  color: var(--text-color-muted);
  font-size: 1.5rem;
  line-height: 1.5;
  border-right: 1px solid var(--border-color-disabled);
}
</style>
