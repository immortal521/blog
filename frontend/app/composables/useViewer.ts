import type { ShallowRef } from "vue";
import Viewer from "viewerjs";

export function useViewer(
  elementRef: Readonly<ShallowRef<HTMLDivElement | null>>,
  options: Viewer.Options = {},
) {
  onMounted(() => {
    if (elementRef.value) {
      new Viewer(elementRef.value, options);
    }
  });
}
