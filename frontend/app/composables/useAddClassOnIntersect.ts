import type { Ref } from "vue";

export function useAddClassOnIntersect(
  target: Ref<HTMLElement | null>,
  className: string,
  options?: IntersectionObserverInit,
) {
  let observer: IntersectionObserver | null = null;

  onMounted(() => {
    if (!target.value) return;
    observer = new IntersectionObserver(([entry]) => {
      if (entry?.isIntersecting) {
        target.value!.classList.add(className);
        observer?.disconnect();
      }
    }, options ?? {});
    observer.observe(target.value);
  });

  onBeforeUnmount(() => observer?.disconnect());
}
