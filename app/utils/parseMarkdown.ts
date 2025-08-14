import MarkdownIt from "markdown-it";
import { mark } from "@mdit/plugin-mark";
import { sub } from "@mdit/plugin-sub";
import { sup } from "@mdit/plugin-sup";
import hljs from "highlight.js";
import type { VNodeChild } from "vue";
import type Token from "markdown-it/lib/token.mjs";

const md = MarkdownIt({ html: false, linkify: true })
	.use(mark)
	.use(sup)
	.use(sub);

export function parseMarkdownToVNode(markdown: string): VNodeChild[] {
	const tokens = md.parse(markdown, {});
	return tokensToVNode(tokens);
}

function tokensToVNode(tokens: Token[]): VNodeChild[] {
	interface Stack {
		tag: string;
		children: VNodeChild[];
		attrs: { [k: string]: string };
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

	for (const token of tokens) {
		switch (true) {
			case token.type.endsWith("_open"): {
				const attrs = Object.fromEntries(token.attrs ?? []);
				// 将开放标签推入栈中
				if (token.tag) {
					stack.push({ tag: token.tag, children: [], attrs });
				}
				break;
			}

			case token.type.endsWith("_close"): {
				// 封闭标签，从栈中弹出并创建 VNode
				const { tag, children, attrs } = stack.pop()!;
				const vnode = h(tag, { ...attrs }, children);
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
					pushToParent(h("code", {}, token.content));
				}
				break;
			}

			case token.type === "hr": {
				break;
			}

			case token.type === "image": {
				if (token.attrs) {
					const attrs = Object.fromEntries(token.attrs);
					pushToParent(h("img", { src: attrs.src, alt: attrs.alt }));
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
					});

				pushToParent(
					h("pre", { class: "hljs" }, [
						h("div", { class: "line-numbers", innerHTML: lineNumbers }),
						h("code", {
							class: `language-${lang}`,
							innerHTML: highlighted.trim(),
						}),
					]),
				);
				break;
			}

			case token.type === "html_block": {
				pushToParent(h("div", { innerHTML: token.content }));
				break;
			}

			case token.type === "html_inline": {
				pushToParent(h("span", { innerHTML: token.content }));
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
