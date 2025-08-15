import MarkdownIt from "markdown-it";
import { mark } from "@mdit/plugin-mark";
import { sub } from "@mdit/plugin-sub";
import { sup } from "@mdit/plugin-sup";
import hljs from "highlight.js";
import { isVNode, type VNodeChild } from "vue";
import type Token from "markdown-it/lib/token.mjs";
import { BaseImage, CopyButton } from "#components";

const md = MarkdownIt({ html: false, linkify: true }).use(mark).use(sup).use(sub);

export function parseMarkdownToVNode(markdown: string): VNodeChild[] {
  const tokens = md.parse(markdown, {});
  return tokensToVNode(tokens);
}

function hasComponent(vnodes: VNodeChild[], component: Component): boolean {
  if (!vnodes) {
    return false;
  }

  const nodes = Array.isArray(vnodes) ? vnodes : [vnodes];

  for (const node of nodes) {
    if (!isVNode(node)) {
      continue;
    }

    if (node.type === component) {
      return true;
    }
  }

  return false;
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
      case token.type === "paragraph_open": {
        const attrs = Object.fromEntries(token.attrs ?? []);
        if (token.tag) {
          stack.push({ tag: token.tag, children: [], attrs, key: "" });
        }
        break;
      }
      case token.type === "paragraph_close": {
        const { children, attrs } = stack.pop()!;
        if (hasComponent(children, BaseImage)) {
          const key = getNextKey("img");
          pushToParent(h("div", { ...attrs, class: "img-container", key }, children));
        } else {
          const key = getNextKey("p");
          pushToParent(h("p", { ...attrs, key }, children));
        }
        break;
      }
      case token.type.endsWith("_open"): {
        const attrs = Object.fromEntries(token.attrs ?? []);
        // 将开放标签推入栈中
        if (token.tag) {
          const key = getNextKey(token.tag);
          stack.push({ tag: token.tag, children: [], attrs, key: key });
        }
        break;
      }

      case token.type.endsWith("_close"): {
        // 封闭标签，从栈中弹出并创建 VNode
        const { tag, children, attrs, key: stackKey } = stack.pop()!;
        const vnode = h(tag, { ...attrs, key: stackKey }, children);
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
          pushToParent(h("code", { key: key }, token.content));
        }
        break;
      }

      case token.type === "hr": {
        const key = getNextKey("hr");
        pushToParent(h("hr", { key: key }));
        break;
      }

      case token.type === "image": {
        if (token.attrs) {
          const attrs = Object.fromEntries(token.attrs);
          pushToParent(
            h(BaseImage, {
              src: attrs.src ?? "",
              alt: attrs.alt,
              preview: true,
            }),
          );
        }
        break;
      }

      case token.type === "fence": {
        const lang = token.info.trim() || "plaintext";
        if (lang === "mermaid") break;
        const highlighted = hljs.highlight(token.content, {
          language: lang,
        }).value;
        const lineNumbers = highlighted
          .trim()
          .split("\n")
          .map(() => {
            return `<div class="line-number"></div>`;
          })
          .join("");

        const key = getNextKey("pre");

        pushToParent(
          h("pre", { class: "hljs", key: key }, [
            h("div", { class: "code-header" }, [
              h("span", { class: "code-language" }, lang),
              h(CopyButton, { text: token.content, class: "copy-btn" }),
            ]),
            h("div", { class: "line-numbers", innerHTML: lineNumbers }),
            h("code", {
              class: `language-${lang}`,
              innerHTML: highlighted.trim(),
            }),
          ]),
        );
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
