import MarkdownIt from "markdown-it";
import { mark } from "@mdit/plugin-mark";
import { sub } from "@mdit/plugin-sub";
import { sup } from "@mdit/plugin-sup";
import hljs from "highlight.js";
import type { VNodeChild } from "vue";
import type Token from "markdown-it/lib/token.mjs";

const md = MarkdownIt({ html: true, linkify: true })
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
				// 将开放标签推入栈中
				if (token.tag) {
					stack.push({ tag: token.tag, children: [] });
				}
				break;
			}

			case token.type.endsWith("_close"): {
				// 封闭标签，从栈中弹出并创建 VNode
				const { tag, children } = stack.pop()!;
				const vnode = h(tag, {}, children);
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
				console.log(highlighted.trim().split("\n"));
				let lineCount = 1;
				const code = highlighted
					.trim()
					.split("\n")
					.map((item) => {
						return `<span class="code-line numbered-code-line" data-line-number="${lineCount++}">${item}</span>`;
					})
					.join("\n");

				pushToParent(
					h("pre", { class: "hljs" }, [
						h("code", { class: `language-${lang}`, innerHTML: code }),
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
