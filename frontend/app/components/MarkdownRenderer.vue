<script setup lang="ts">
import "viewerjs/dist/viewer.css";
import { parseMarkdownToVNode, type TocItem } from "@/utils/parseMarkdown";
import type { VNodeChild } from "vue";

interface Props {
  markdown: string;
  toc?: boolean;
}

interface Emits {
  tocChange: [toc: TocItem[]];
}

const { markdown, toc = false } = defineProps<Props>();
const emits = defineEmits<Emits>();

const renderedVNode = shallowRef<VNodeChild>([]);
const containerRef = useTemplateRef("container");

const { update } = useViewer(containerRef);

const renderMarkdown = () => {
  const { content, toc: tocItems } = parseMarkdownToVNode(markdown, { toc: toc });

  if (toc && tocItems) {
    emits("tocChange", tocItems);
  }

  renderedVNode.value = content;
};

watch(
  () => markdown,
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
  border-radius: 8px;
  cursor: pointer;
  box-shadow: var(--shadow-md);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;

  &:hover {
    transform: scale(1.02);
    box-shadow: var(--shadow-lg);
  }
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
  margin: 1.2rem 0;
}

:deep(h1),
:deep(h2),
:deep(h3),
:deep(h4),
:deep(h5),
:deep(h6) {
  scroll-margin-top: var(--header-height);
  font-weight: 700;
  color: var(--color-header);
  line-height: 1.4;
  margin-top: 2.4rem;
  margin-bottom: 1rem;
  position: relative;
}

:deep(h1) {
  font-size: 3.2rem;
  padding-bottom: 0.6rem;

  &::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    width: 80px;
    height: 4px;
    background: linear-gradient(90deg, var(--color-primary-base), transparent);
    border-radius: 2px;
  }
}

:deep(h2) {
  font-size: 2.6rem;
  padding-bottom: 0.4rem;

  &::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    width: 60px;
    height: 3px;
    background: linear-gradient(90deg, var(--color-primary-base), transparent);
    border-radius: 2px;
  }
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
  color: var(--text-color-secondary);
}

:deep(p) {
  color: var(--text-color-primary);
  overflow-wrap: break-word;
  line-height: 1.8;
}

:deep(:not(pre.shiki) > code) {
  display: inline;
  padding: 0.2em 0.45em;
  font-family: "Maple Mono", "Noto Sans SC", monospace;
  font-size: 0.86em;
  font-weight: 500;
  color: var(--inline-code-color);
  background-color: var(--inline-code-bg);
  border-radius: 4px;
  line-height: 1.45;
  white-space: break-spaces;
  overflow-wrap: break-word;
}

:deep(table) {
  width: 100%;
  border-collapse: collapse;
  border-radius: 8px;
  background: var(--bg-content);
  overflow: hidden;
  font-size: var(--font-size-table-row);
  color: var(--text-color-primary);
  box-shadow: var(--shadow-sm);
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
  border-bottom: 2px solid var(--table-head-border);
}

:deep(tbody) tr {
  border-bottom: 1px solid var(--border-table);
  font-size: var(--font-size-table-row);
  transition: background-color 0.2s ease;
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
  border-radius: 8px;
  padding: 16px 20px;
  background: var(--bg-card-base);
  position: relative;
  border-left: none;
  margin: 1.5rem 0;
  box-shadow: var(--shadow-sm);

  &::before {
    display: block;
    position: absolute;
    content: "";
    width: 4px;
    left: 0;
    top: 0;
    height: 100%;
    background: linear-gradient(180deg, var(--color-primary-base), var(--color-primary-hover));
    border-radius: 2px;
  }
}

:deep(ul),
:deep(ol) {
  margin: 1.2em 0;
  padding-left: 1.5em;
  color: var(--text-color-primary);
  font-size: 1.5rem;
  line-height: 1.75;
}

:deep(ul li),
:deep(ol li) {
  margin: 0.4rem 0;
  padding-left: 0.25rem;
  position: relative;
}

// 无序列表圆点样式
:deep(ul li::marker) {
  color: var(--color-primary-base);
}

// 有序列表数字样式
:deep(ol li::marker) {
  font-weight: bold;
  color: var(--color-primary-base);
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
  color: var(--text-on-brand);
  background: linear-gradient(135deg, var(--color-primary-base), var(--color-primary-hover));
  border-radius: 4px;
}

:deep(u) {
  text-decoration-line: underline;
  text-decoration-color: var(--color-primary-base);
  text-underline-offset: 0.5rem;
  text-decoration-thickness: 2px;
}

:deep(a) {
  color: var(--color-primary-base);
  text-decoration: none;
  position: relative;
  transition: color 0.2s ease;

  &::after {
    content: "";
    position: absolute;
    bottom: -2px;
    left: 0;
    width: 0;
    height: 2px;
    background: var(--color-primary-base);
    transition: width 0.2s ease;
  }

  &:hover::after {
    width: 100%;
  }
}

:deep(hr) {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border-color-default), transparent);
  margin: 2rem 0;
}
</style>
