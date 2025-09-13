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
  let observer: IntersectionObserver | null = null;

  onMounted(async () => {
    await nextTick(); // 确保 DOM 渲染完成

    // 将 targets 转为数组并提取真实 DOM 元素
    const list: HTMLElement[] = (Array.isArray(targets.value) ? targets.value : [targets.value])
      .map(getDomElement)
      .filter((el): el is HTMLElement => el !== null); // TS 类型守卫

    if (list.length === 0) return;

    observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add(className);
          observer?.unobserve(entry.target); // 只处理一次
        }
      });
    }, options ?? {});

    list.forEach((el) => observer?.observe(el));
  });

  onBeforeUnmount(() => observer?.disconnect());
}
