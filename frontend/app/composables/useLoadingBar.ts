const count = ref(0);
const isVisible = computed(() => count.value > 0);

export function useLoadingBar() {
  return {
    isVisible,
    start() {
      count.value++;
    },
    done() {
      if (count.value > 0) {
        count.value--;
      }
    },
    reset() {
      count.value = 0;
    },
  };
}
