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

export function parseMarkdownToVNode(markdown: string): VNodeChild[] {
  const tokens = md.parse(markdown, {});
  return tokensToVNode(tokens);
}

function tokensToVNode(tokens: Token[]): VNodeChild[] {
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

  const keyCounters = new Map<string, number>();

  const getNextKey = (type: string): string => {
    const currentCount = keyCounters.get(type) ?? 0;
    keyCounters.set(type, currentCount + 1);
    return `${type}-${currentCount}`;
  };

  for (const token of tokens) {
    switch (true) {
      case token.type.endsWith("_open"): {
        const attrs = Object.fromEntries(token.attrs ?? []);
        // 将开放标签推入栈中
        if (token.tag) {
          const key = getNextKey(token.tag);
          stack.push({ tag: token.tag, children: [], attrs, key });
        }
        break;
      }

      case token.type.endsWith("_close"): {
        // 封闭标签，从栈中弹出并创建 VNode
        const { tag, children, attrs, key } = stack.pop()!;
        const vnode = h(tag, { ...attrs, key }, children);
        pushToParent(vnode);
        break;
      }

      case token.type === "inline": {
        // 递归处理内联内容，并将其结果直接添加到当前父节点的 children 中
        if (token.children) {
          const inlineVNodes = tokensToVNode(token.children);
          for (const node of inlineVNodes) {
            pushToParent(node);
          }
        }
        break;
      }

      case token.type === "text": {
        if (token.content) {
          pushToParent(token.content);
        }
        break;
      }

      case token.type === "code_inline": {
        if (token.content) {
          const key = getNextKey("code");
          pushToParent(h("code", { key }, token.content));
        }
        break;
      }

      case token.type === "checkbox_input": {
        if (token.attrs) {
          const attrs = Object.fromEntries(token.attrs);
          pushToParent(h("input", { ...attrs }));
        }
        break;
      }

      case token.type === "hr": {
        const key = getNextKey("hr");
        pushToParent(h("hr", { key }));
        break;
      }

      case token.type === "image": {
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
        break;
      }

      case token.type === "fence": {
        let lang = token.info.trim() as BundledLanguage | SpecialLanguage;
        if (!highlighter.getLoadedLanguages().includes(lang)) {
          lang = "plaintext";
        }
        const code = token.content.trim();
        const key = getNextKey("pre");
        pushToParent(h(CodeWrapper, { lang, code, key }));
        break;
      }

      // 可以根据需要添加更多 token 类型
      default: {
        console.warn(`未处理的 token 类型: ${token.type}`);
        console.warn(token);
        break;
      }
    }
  }

  return result;
}
