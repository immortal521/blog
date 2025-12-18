import { useDebounceFn } from "@vueuse/core";

/**
 * 获取并维护客户端窗口宽度、兼容 SSR CSR。
 *
 * @description
 * 用于在 Nuxt 应用中获取客户端窗口宽度，
 * 并通过 Cookie 在服务端渲染与客户端渲染之间同步状态。
 *
 * - 根据屏幕宽度在 JS 层面决定布局或组件渲染
 * - SSR 场景下区分移动端 / 桌面端渲染逻辑
 * - 减少首屏 hydration 前后的布局抖动
 *
 * @returns {Ref<number>} returns.width 当前窗口宽度
 *
 * @example
 * ```ts
 * const { width } = useClientWidth();
 *
 * const isMobile = computed(() => width.value < 768);
 * ```
 */
export function useClientWidth(): { width: Ref<number> } {
  const widthCookie = useCookie<number>("client_width", {
    path: "/",
    maxAge: 356 * 24 * 60 * 60,
  });

  /**
   * 在无法访问 window（即 SSR 环境）时，
   * 通过 User-Agent 推测一个“合理的”客户端宽度
   */
  const getWidthFromUA = () => {
    if (import.meta.client) return window.innerWidth;
    const headers = useRequestHeaders(["user-agent"]);
    const ua = headers["user-agent"] || "";

    const isMobile = /Mobile|Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      ua,
    );

    return isMobile ? 375 : 1024;
  };

  const clientWidth = widthCookie.value ? widthCookie.value : getWidthFromUA();

  const width = ref(clientWidth);

  const syncWidthToCookie = useDebounceFn((value: number) => {
    widthCookie.value = value;
  }, 500);

  const update = () => {
    const next = window.innerWidth;
    width.value = next;
    syncWidthToCookie(window.innerWidth);
  };

  onMounted(() => {
    window.addEventListener("resize", update);
    update();
  });

  onBeforeUnmount(() => {
    window.removeEventListener("resize", update);
  });

  return {
    width,
  };
}
