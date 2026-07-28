import type { FloatActionItem } from "~/components/FloatActionGroup/types";

export const useFloatActionStore = defineStore("floatAction", () => {
  const _items = ref<FloatActionItem[]>([]);

  const floatActionItems = computed(() => _items.value);

  function add(...items: FloatActionItem[]) {
    const existingIds = new Set(_items.value.map((item) => item.id));

    const newItems = items.filter((item) => !existingIds.has(item.id));

    _items.value.push(...newItems);
  }

  function remove(id: string) {
    const index = _items.value.findIndex((item) => item.id === id);

    if (index !== -1) {
      _items.value.splice(index, 1);
    }
  }

  function update(id: string, patch: Partial<FloatActionItem>) {
    const item = _items.value.find((item) => item.id === id);

    if (item) {
      Object.assign(item, patch);
    }
  }

  function clear() {
    _items.value = [];
  }

  function get(id: string) {
    return _items.value.find((item) => item.id === id);
  }

  return {
    floatActionItems,
    add,
    remove,
    update,
    clear,
    get,
  };
});
