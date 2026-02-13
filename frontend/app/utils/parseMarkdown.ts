import MarkdownIt from "markdown-it";
import { mark } from "@mdit/plugin-mark";
import { sub } from "@mdit/plugin-sub";
import { sup } from "@mdit/plugin-sup";
import type { VNodeChild } from "vue";
import type Token from "markdown-it/lib/token.mjs";
import { CodeWrapper, NuxtImg } from "#components";
import type { BundledLanguage, SpecialLanguage } from "shiki";
import { tasklist } from "@mdit/plugin-tasklist";

const md = MarkdownIt({ html: false, linkify: true }).use(mark).use(sup).use(sub).use(tasklist);

interface ParseResult {
  content: VNodeChild[];
  toc?: TocItem[];
}

interface TocItem {
  id: string;
  level: number;
  text: string;
}

interface Options {
  toc?: boolean;
}

interface NormalizedOptions {
  toc: boolean;
}

function slugifyHeading(text: string): string {
  return text
    .trim()
    .toLowerCase()
    .replace(/[`~!@#$%^&*()+=[\]{};:'"\\|,.<>/?，。；：“”‘’、】【、？！…—]/g, " ")
    .replace(/[^\w\u4e00-\u9fa5\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");
}

function uniqueSlug(text: string, slugCounters: Map<string, number>) {
  const n = slugCounters.get(text) ?? 0;
  slugCounters.set(text, n + 1);
  return n === 0 ? text : `${text}-${n}`;
}

function ensureHeadingID(text: string, level: number, slugCounters: Map<string, number>) {
  const slug = slugifyHeading(text);
  if (slug) return uniqueSlug(slug, slugCounters);
  return uniqueSlug(`heading-${level}`, slugCounters);
}

function getInlinePlaintext(inline: Token | undefined): string {
  if (!inline || inline.type !== "inline") return "";
  const parts: string[] = [];
  for (const c of inline.children ?? []) {
    if (c.type === "text") parts.push(c.content);
  }
  return parts.join("");
}

export function parseMarkdownToVNode(markdown: string, options?: Options): ParseResult {
  const tokens = md.parse(markdown, {});
  const normalizedOptions: NormalizedOptions = {
    toc: options?.toc ?? false,
  };
  const keyCounters = new Map<string, number>();
  const toc: TocItem[] = [];
  const content = tokensToVNode(tokens, normalizedOptions, keyCounters, toc);
  console.log(toc);

  return { content, toc };
}

function tokensToVNode(
  tokens: Token[],
  options: NormalizedOptions,
  keyCounters: Map<string, number>,
  toc: TocItem[],
): VNodeChild[] {
  interface Stack {
    tag: string;
    children: VNodeChild[];
    attrs: { [k: string]: string };
    key: string;
  }
  const stack: Stack[] = [];
  const result: VNodeChild[] = [];

  const pushToParent = (node: VNodeChild) => {
    if (stack.length > 0) {
      (stack[stack.length - 1] as Stack).children.push(node);
    } else {
      result.push(node);
    }
  };

  const getNextKey = (type: string): string => {
    const currentCount = keyCounters.get(type) ?? 0;
    keyCounters.set(type, currentCount + 1);
    return `${type}-${currentCount}`;
  };

  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i]!;
    if (options.toc && token.type === "heading_open") {
      const level = Number(token.tag.slice(1) ?? 0);
      const next = tokens[i + 1];
      const text = next?.type === "inline" ? getInlinePlaintext(next) : "";

      if (level >= 1 && level <= 6) {
        const id = ensureHeadingID(text, level, keyCounters);
        token.attrs = token.attrs ?? [];
        if (!token.attrs.some((attr) => attr[0] === "id")) token.attrs.push(["id", id]);
        toc.push({ id, level, text });
      }
    }

    if (token.tag && token.type.endsWith("_open")) {
      const attrs = Object.fromEntries(token.attrs ?? []);
      // 将开放标签推入栈中
      const key = getNextKey(token.tag);
      stack.push({ tag: token.tag, children: [], attrs, key });
      continue;
    }

    if (token.type.endsWith("_close")) {
      // 封闭标签，从栈中弹出并创建 VNode
      const top = stack.pop();
      if (!top) continue;
      const { tag, children, attrs, key } = top;
      const vnode = h(tag, { ...attrs, key }, children);
      pushToParent(vnode);
      continue;
    }

    if (token.type === "inline") {
      // 递归处理内联内容，并将其结果直接添加到当前父节点的 children 中
      if (token.children) {
        const inlineVNodes = tokensToVNode(token.children, options, keyCounters, toc);
        for (const node of inlineVNodes) {
          pushToParent(node);
        }
      }
      continue;
    }

    if (token.type === "text" && token.content) {
      pushToParent(token.content);
      continue;
    }

    if (token.type === "code_inline" && token.content) {
      const key = getNextKey("code");
      pushToParent(h("code", { key }, token.content));
      continue;
    }

    if (token.type === "checkbox_input" && token.attrs) {
      const attrs = Object.fromEntries(token.attrs);
      pushToParent(h("input", { ...attrs }));
      continue;
    }

    if (token.type === "hr") {
      const key = getNextKey("hr");
      pushToParent(h("hr", { key }));
      continue;
    }

    if (token.type === "image") {
      const key = getNextKey("img");
      if (token.attrs) {
        const attrs = Object.fromEntries(token.attrs);
        pushToParent(
          h(NuxtImg, {
            key,
            src: attrs.src ?? "",
            alt: attrs.alt,
            class: "img",
            loading: "lazy",
          }),
        );
      }
      continue;
    }

    if (token.type === "fence") {
      let lang = token.info.trim() as BundledLanguage | SpecialLanguage;
      if (!highlighter.getLoadedLanguages().includes(lang)) {
        lang = "plaintext";
      }
      const code = token.content.trim();
      const key = getNextKey("pre");
      pushToParent(h(CodeWrapper, { lang, code, key }));
      continue;
    }

    if (token.type === "softbreak") {
      const key = getNextKey("br");
      pushToParent(h("br", { key }));
      continue;
    }

    console.warn(`未处理的 token 类型: ${token.type}`);
  }

  return result;
}
