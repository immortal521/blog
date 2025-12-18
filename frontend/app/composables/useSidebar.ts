export function useSidebar() {
  const open = ref(false);
  const collapsed = ref(false);
  let lastIsDesktop = false;

  const { isDesktop } = useResponsive();

  watch(
    isDesktop,
    (value) => {
      if (value !== lastIsDesktop) {
        open.value = value;
        lastIsDesktop = value;
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

  return {
    open,
    collapsed,
    toggleOpen,
    toggleCollapsed,
  };
}
