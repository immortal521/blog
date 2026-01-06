type ScrollTarget = HTMLElement | Window | Document | Ref<HTMLElement | Window | Document>;
type ScrollCallback = (isOverThreshold: boolean) => void;

export function useScrollHandler(
  target: ScrollTarget,
  threshold: number,
  callback: ScrollCallback,
): void;
export function useScrollHandler(
  threshold: number,
  callback: (isOverThreshold: boolean) => void,
): void;

export function useScrollHandler(
  ...args: [ScrollTarget, number, ScrollCallback] | [number, ScrollCallback]
) {
  let target: ScrollTarget;
  let threshold: number;
  let callback: (isOverThreshold: boolean) => void;
  if (typeof args[0] === "number") {
    target = window;
    threshold = args[0];
    callback = args[1] as ScrollCallback;
  } else {
    target = args[0];
    threshold = args[1] as number;
    callback = args[2] as ScrollCallback;
  }
  let lastState: boolean | null = null;

  const getScrollTop = (el: HTMLElement | Window | Document) => {
    if (el instanceof Window) return el.scrollY;

    if (el instanceof Document) return el.documentElement.scrollTop;

    return el.scrollTop;
  };

  const handler = () => {
    const el = unref(target);
    if (!el) return;

    const over = getScrollTop(el) >= threshold;
    if (over !== lastState) {
      lastState = over;
      callback(over);
    }
  };

  let cleanup: (() => void) | undefined;

  const register = () => {
    const el = unref(target);
    if (!el) return;

    el.addEventListener("scroll", handler, { passive: true });
    cleanup = () => el.removeEventListener("scroll", handler);
    handler();
  };

  onMounted(register);

  watch(
    () => unref(target),
    () => {
      cleanup?.();
      register();
    },
  );

  onUnmounted(() => {
    cleanup?.();
  });
}
