export function useActiveToc(toc: MaybeRef<TocItem[]>) {
  const activeId = ref("");
  let observer: IntersectionObserver | null = null;
  const observe = () => {
    observer?.disconnect();
    const elements = toValue(toc)
      .map((item) => document.getElementById(item.id))
      .filter(Boolean);
    observer = new IntersectionObserver(
      (entries) => {
        const visible = entries.find((entry) => entry.isIntersecting);

        if (visible) {
          activeId.value = visible.target.id;
        }
      },
      {
        rootMargin: "-60px 0px -60% 0px",
      },
    );
    elements.forEach((el) => {
      if (!el) return;
      observer!.observe(el);
    });
  };

  onMounted(async () => {
    await nextTick();
    observe();
  });

  watch(
    () => toValue(toc),
    async () => {
      await nextTick();
      observe();
    },
  );

  onUnmounted(() => {
    observer?.disconnect();
  });

  return {
    activeId,
    refresh: observe,
  };
}
