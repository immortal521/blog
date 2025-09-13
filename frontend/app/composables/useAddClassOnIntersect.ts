import { onMounted, onBeforeUnmount, nextTick, type Ref } from "vue";

type ComponentWithEl = { $el?: HTMLElement };
type Target = HTMLElement | ComponentWithEl | null;
type TargetOrArray = Target | Target[];

/** 类型守卫：判断是否是组件实例 */
function isComponent(target: Target): target is ComponentWithEl {
  return target !== null && typeof target === "object" && "$el" in target;
}

/** 安全提取 DOM 元素 */
function getDomElement(target: Target): HTMLElement | null {
  if (!target) return null;
  if (isComponent(target)) return target.$el ?? null;
  return target;
}
/**
 * useAddClassOnIntersect
 * 当元素进入视口时添加 class
 */
export function useAddClassOnIntersect(
  targets: Ref<TargetOrArray>,
  className: string,
  options?: IntersectionObserverInit,
) {
  let observers: IntersectionObserver[] = [];

  onMounted(async () => {
    await nextTick(); // 确保 DOM 渲染完成

    const list: Target[] = Array.isArray(targets.value) ? targets.value : [targets.value];

    list.forEach((item) => {
      const el = getDomElement(item);
      if (!el) return;

      const observer = new IntersectionObserver(([entry]) => {
        if (entry?.isIntersecting) {
          el.classList.add(className);
          observer.disconnect();
        }
      }, options ?? {});

      observer.observe(el);
      observers.push(observer);
    });
  });

  onBeforeUnmount(() => observers.forEach((obs) => obs.disconnect()));
}
