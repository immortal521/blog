import Viewer from "viewerjs";
import type { ShallowRef } from "vue";

export function useViewer(
  elementRef: Readonly<ShallowRef<HTMLDivElement | null>>,
  options: Viewer.Options = {},
) {
  const viewer = shallowRef<Viewer | null>(null);

  const init = () => {
    if (!elementRef.value) return;
    viewer.value?.destroy();
    viewer.value = new Viewer(elementRef.value, options);
  };

  const update = async () => {
    await nextTick();
    viewer.value?.update();
  };

  onMounted(() => {
    init();
  });

  onBeforeUnmount(() => {
    viewer.value?.destroy();
    viewer.value = null;
  });

  return { viewer, init, update };
}
