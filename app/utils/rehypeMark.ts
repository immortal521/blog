import { visit } from "unist-util-visit";
import type { Element, Parent, Text } from "hast";

export default function rehypeMark() {
  return (tree: Parent) => {
    visit(tree, "text", (node: Text, index?: number, parent?: Element | Parent) => {
      if (!parent || index === undefined) return;
      const value = node.value;
      const regex = /==([\s\S]+?)==/g;

      let match: RegExpExecArray | null;
      let lastIndex = 0;
      const children: (Text | Element)[] = [];

      while ((match = regex.exec(value)) !== null) {
        if (match.index > lastIndex) {
          children.push({
            type: "text",
            value: value.slice(lastIndex, match.index),
          });
        }

        children.push({
          type: "element",
          tagName: "mark",
          children: [
            {
              type: "text",
              value: match[1] as string,
            },
          ],
          properties: {},
        });

        lastIndex = regex.lastIndex;
      }

      if (lastIndex < value.length) {
        children.push({
          type: "text",
          value: value.slice(lastIndex),
        });
      }

      // 替换当前 text 节点
      parent.children.splice(index, 1, ...children);
    });
  };
}
