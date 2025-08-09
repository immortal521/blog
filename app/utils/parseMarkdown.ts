interface CoreMarkdownModules {
  unified: typeof import("unified").unified;
  remarkParse: typeof import("remark-parse").default;
  remarkGfm: typeof import("remark-gfm").default;
  remarkRehype: typeof import("remark-rehype").default;
  rehypeSanitize: typeof import("rehype-sanitize");
  rehypeRaw: typeof import("rehype-raw").default;
  rehypeMark: typeof import("@/utils/rehypeMark").default;
  astNodeToVNode: typeof import("@/utils/astToVNode").default;
}

let coreModulesPromise: Promise<CoreMarkdownModules> | null = null;

const loadCoreMarkdownModules = async (): Promise<CoreMarkdownModules> => {
  if (!coreModulesPromise) {
    coreModulesPromise = Promise.all([
      import("unified"),
      import("remark-parse"),
      import("remark-gfm"),
      import("remark-rehype"),
      import("rehype-sanitize"),
      import("rehype-raw"),
      import("@/utils/rehypeMark"),
      import("@/utils/astToVNode"),
    ]).then(
      ([
        unifiedModule,
        remarkParseModule,
        remarkGfmModule,
        remarkRehypeModule,
        rehypeSanitizeModule,
        rehypeRawModule,
        rehypeMarkModule,
        astNodeToVNodeModule,
      ]) => ({
        unified: unifiedModule.unified,
        remarkParse: remarkParseModule.default,
        remarkGfm: remarkGfmModule.default,
        remarkRehype: remarkRehypeModule.default,
        rehypeSanitize: rehypeSanitizeModule,
        rehypeRaw: rehypeRawModule.default,
        rehypeMark: rehypeMarkModule.default,
        astNodeToVNode: astNodeToVNodeModule.default,
      }),
    );
  }
  return coreModulesPromise;
};

const detectOptions = (md: string) => ({
  highlight: /(```[\s\S]*?```|`[^`\n]*?`)/.test(md),
  math: /\$\$[^$]*\$\$|[^$]\$[^$\n]*\$/.test(md),
});

export const generateMdAst = async (markdown: string) => {
  if (!markdown) return;

  const { highlight, math } = detectOptions(markdown);
  const {
    unified,
    remarkParse,
    remarkGfm,
    remarkRehype,
    rehypeSanitize,
    rehypeRaw,
    rehypeMark,
    astNodeToVNode,
  } = await loadCoreMarkdownModules();

  const processor = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeRaw)
    .use(rehypeSanitize.default, {
      ...rehypeSanitize.defaultSchema,
      tagNames: [
        ...(rehypeSanitize.defaultSchema.tagNames || []),
        "u", // 自定义标签扩展
      ],
    })
    .use(rehypeMark);

  if (highlight) {
    const [{ default: highlight }, { default: codeLines }] = await Promise.all([
      import("rehype-highlight"),
      import("rehype-highlight-code-lines"),
    ]);
    processor.use(highlight).use(codeLines, { showLineNumbers: true });
  }

  if (math) {
    const [{ default: katex }, { default: math }] = await Promise.all([
      import("rehype-katex"),
      import("remark-math"),
    ]);
    processor.use(math).use(katex);
  }

  try {
    const ast = await processor.run(processor.parse(markdown));
    return astNodeToVNode(ast);
  } catch (err) {
    console.warn("Markdown 解析失败：", err);
    return h("span", {}, "Markdown 解析失败");
  }
};
