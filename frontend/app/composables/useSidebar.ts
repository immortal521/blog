export function useSidebar() {
  const open = ref(false);
  const collapsed = ref(false);
  let lastIsDesktop = false;
  const { width } = useClientWidth();
  watch(
    width,
    (w) => {
      const isDesktop = w >= 768;

      if (isDesktop !== lastIsDesktop) {
        open.value = isDesktop;
        lastIsDesktop = isDesktop;
      }
    },
    { immediate: true },
  );

  const toggleOpen = () => {
    open.value = !open.value;
  };

  const toggleCollapsed = () => {
    collapsed.value = !collapsed.value;
  };

  return { open, collapsed, toggleOpen, toggleCollapsed };
}
