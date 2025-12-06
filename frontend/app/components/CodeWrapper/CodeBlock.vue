<script setup lang="ts">
import type { BundledLanguage, SpecialLanguage } from "shiki";
import type { Element, ElementContent } from "hast";
import CodeHeader from "./CodeHeader.vue";

interface Props {
  code: string;
  lang: BundledLanguage | SpecialLanguage;
}

const { code, lang } = defineProps<Props>();

// 计算行数
const lineCount = code.trim().split("\n").length;

// 将 code 转换为 HAST，并添加行号
const root = highlighter.codeToHast(code, {
  lang,
  themes: {
    dark: "tokyo-night",
    light: "one-light",
  },
  colorReplacements: {
    "tokyo-night": { "#1a1b26": "var(--bg-code)" },
    "one-light": { "#fafafa": "var(--bg-code)" },
  },
  transformers: [
    {
      pre(node) {
        const lineNumbersDiv: Element = {
          type: "element",
          tagName: "div",
          properties: { class: "line-numbers" },
          children: Array.from({ length: lineCount }, (_, i) => ({
            type: "element",
            tagName: "div",
            properties: { class: "line-number", "data-line-number": i + 1 },
            children: [],
          })),
        };

        // 保留原有 children，并在前面插入行号
        node.children = [lineNumbersDiv, ...node.children];
      },
    },
  ],
});

// 提取 <pre> 子节点并转换为 Vue VNode
const vnodes = hastToVNode(root.children as ElementContent[]);
const preVNode = vnodes[0] as VNode;

const headerVNode = h(CodeHeader, { code, lang, class: "code-header" });

// 将 header 插入到 preVNode 的 children 前面
if (preVNode.children && Array.isArray(preVNode.children)) {
  preVNode.children.unshift(headerVNode);
} else {
  preVNode.children = [headerVNode];
}
</script>

<template>
  <component :is="preVNode"></component>
</template>

<style lang="less">
.shiki {
  display: flex;
  position: relative;
  font-size: 1.5rem;
  line-height: 1.5;
  padding: 3rem 3rem 1.5rem 0;
  border-radius: var(--radius-card, 0.5rem);
  box-shadow: 0 1px 3px rgb(0 0 0 / 6%);
  overflow: hidden;
  transition: all 0.3s ease-in-out;
  margin-top: 20px;
  font-family: "Maple Mono", "Noto Sans SC", monospace;
  counter-reset: line-number 0;

  & code::-webkit-scrollbar-thumb {
    height: 2px;
    width: 2px;
    background-color: #0000;
    transition: background 0.3s ease-in-out;
  }

  &:hover code::-webkit-scrollbar-thumb {
    background-color: var(--color-primary-base);
  }
}

.shiki code {
  font-family: inherit;
  overflow-x: auto;
  transition: all 0.3s ease-in-out;
  cursor: text;
}

.line-numbers,
.line {
  transition: all 0.3s ease-in-out;
}

.code-header {
  position: absolute;
  width: 100%;
  height: 3rem;
  top: 0;
  left: 0;
  line-height: 3rem;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}

.shiki:hover .code-header {
  opacity: 1;
}

.line-numbers {
  position: sticky;
  left: 0;
  flex-shrink: 0;
  padding-right: 0.5rem;
  width: 3.5rem;
  margin-right: 0.5rem;
  background-color: var(--bg-code);
  text-align: right;
  color: var(--text-color-muted);
  font-size: 1.5rem;
  line-height: 1.5;
  border-right: 1px solid var(--border-color-disabled);
  cursor: default;
  user-select: none;
  counter-reset: line-number 0;
  z-index: 1;
  font-family: "Maple Mono", "Noto Sans SC", monospace;
}

.line-number::before {
  content: counter(line-number);
  counter-increment: line-number;
}

@media (width <= 768px) {
  .code-header {
    opacity: 1;
  }
}
</style>
