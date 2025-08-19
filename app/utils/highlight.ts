import type { ElementContent } from "hast";
import { createHighlighter } from "shiki";
import type { VNodeChild } from "vue";

export function hastToVNode(hasts: ElementContent[]): VNodeChild[] {
  const result: VNodeChild[] = [];

  for (const hast of hasts) {
    if (hast.type === "text") {
      result.push(hast.value);
    }
    if (hast.type === "element") {
      result.push(h(hast.tagName, { ...hast.properties }, hastToVNode(hast.children)));
    }
  }

  return result;
}

export const highlighter = await createHighlighter({
  themes: ["tokyo-night", "one-light"],
  langs: [
    "bat",
    "c",
    "cmake",
    "cpp",
    "css",
    "dart",
    "desktop",
    "docker",
    "fish",
    "go",
    "groovy",
    "html",
    "http",
    "java",
    "javascript",
    "json",
    "json5",
    "jsonc",
    "js",
    "jsx",
    "kotlin",
    "latex",
    "less",
    "log",
    "lua",
    "markdown",
    "md",
    "mermaid",
    "nginx",
    "nushell",
    "nu",
    "postcss",
    "powershell",
    "prisma",
    "proto",
    "protobuf",
    "py",
    "python",
    "rust",
    "rs",
    "sass",
    "scss",
    "shell",
    "shellscript",
    "sql",
    "stylus",
    "systemd",
    "toml",
    "ts",
    "ts-tags",
    "tsx",
    "typescript",
    "vue",
    "vue-html",
    "vue-vine",
    "wasm",
    "xml",
    "xsl",
    "yaml",
    "yml",
    "zig",
  ],
});
