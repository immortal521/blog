import { h, type VNode } from "vue";
import type {
	Root,
	Element,
	Text as HastText,
	Node as HastNode,
	RootContent,
} from "hast";
// import ClipboardButton from "@comp/ClipboardButton/ClipboardButton.vue";

type ElementReplacer = (
	node: Element,
	props: Record<string, unknown>,
	children: (VNode | string)[],
) => VNode | null;

/**
 * 递归提取 VNode 或 VNode 数组中所有字符串内容
 */
// function extractStringsFromVNodes(vnodes: VNode | VNode[] | null | undefined): string[] {
//   if (!vnodes) return [];
//
//   const vnodeArray = Array.isArray(vnodes)
//     ? vnodes.filter(isVNode)
//     : isVNode(vnodes)
//       ? [vnodes]
//       : [];
//
//   const strings: string[] = [];
//
//   for (const vnode of vnodeArray) {
//     // 纯文本节点
//     if (vnode.type === Text || typeof vnode.children === "string") {
//       strings.push(String(vnode.children));
//     }
//
//     // children 是数组时，递归提取
//     if (Array.isArray(vnode.children)) {
//       strings.push(...extractStringsFromVNodes(vnode.children as VNode[]));
//     } else if (typeof vnode.children === "object" && vnode.children !== null) {
//       // 插槽函数情况
//       const slotObj = vnode.children as Record<string, unknown>;
//       if (typeof slotObj.default === "function") {
//         const slotVNodes = slotObj.default();
//         strings.push(...extractStringsFromVNodes(slotVNodes.filter(isVNode)));
//       } else {
//         // 其它对象结构，尝试扁平处理
//         const nestedChildren = Object.values(slotObj).flat().filter(isVNode);
//         if (nestedChildren.length) {
//           strings.push(...extractStringsFromVNodes(nestedChildren));
//         }
//       }
//     }
//   }
//
//   return strings;
// }

const elementReplacers: Record<string, ElementReplacer> = {
	// img: (_node, props) =>
	//   h(
	//     NImage,
	//     {
	//       src: props.src as string,
	//       lazy: true,
	//     },
	//     {
	//       error: () => h(Icon, { icon: "image", fill: "var(--text-color-muted)", size: 100 }),
	//       placeholder: () => h(Icon, { icon: "image", fill: "var(--text-color-muted)", size: 100 }),
	//     },
	//   ),
	// pre: (_node, props, children) => {
	//   if (typeof children === "string") {
	//     return h("pre", props, children);
	//   }
	//   // 提取文本，添加复制按钮
	//   const textContent = extractStringsFromVNodes(children as VNode[]).join("");
	//   children.push(h(ClipboardButton, { text: textContent }));
	//   return h("pre", props, children);
	// },
};

/**
 * 将 HAST 的 properties 转为 Vue VNode 所需 props，处理 data 属性的 camelCase 转 kebab-case
 */
function transformProperties(
	properties: Element["properties"] = {},
): Record<string, unknown> {
	return Object.entries(properties).reduce(
		(acc, [key, value]) => {
			let propKey = key;

			if (
				key.startsWith("data") &&
				key.length > 4 &&
				key[4]?.toUpperCase() === key[4]
			) {
				propKey =
					"data" +
					key
						.slice(4)
						.replace(/([A-Z])/g, "-$1")
						.toLowerCase();
			}

			if (typeof value === "boolean") {
				acc[propKey] = value;
			} else if (Array.isArray(value)) {
				acc[propKey] = value.join(" ");
			} else if (value != null) {
				acc[propKey] = String(value);
			}

			return acc;
		},
		{} as Record<string, unknown>,
	);
}

/**
 * 递归将 HAST AST 节点转换成 Vue VNode
 */
function astNodeToVNode(node: HastNode | RootContent): VNode | VNode[] | null {
	switch (node.type) {
		case "root":
			return (node as Root).children
				.map(astNodeToVNode)
				.filter(Boolean) as VNode[];

		case "element": {
			const el = node as Element;
			const children = el.children.map(astNodeToVNode).filter(Boolean) as (
				| VNode
				| string
			)[];
			const props = transformProperties(el.properties);

			const replacer = elementReplacers[el.tagName];
			if (replacer) {
				return replacer(el, props, children);
			}
			return h(el.tagName, props, children);
		}

		case "text":
			return h("span", {}, (node as HastText).value);

		case "comment":
			return null;

		default:
			console.warn(`Unknown HAST node type: ${node.type}`);
			return null;
	}
}

export default astNodeToVNode;
