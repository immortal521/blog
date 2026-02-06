export const useSidebarStore = defineStore("sidebar", () => {
  const openKeys = ref(new Set<string>());

  return {
    openKeys,
  };
});
