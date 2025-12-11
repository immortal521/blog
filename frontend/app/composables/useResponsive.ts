import { useMediaQuery, useWindowSize } from "@vueuse/core";

export function useResponsive() {
  // SSR 时的初始宽度（来自 cookie）
  const widthCookie = useCookie("windowWidth");
  const ssrWidth = parseInt(widthCookie.value ?? "1080");

  // 客户端自动更新窗口宽度
  const { width: clientWidth } = useWindowSize();

  // 统一管理设备断点
  const isMobile = useMediaQuery("(max-width: 768px)", {
    ssrWidth,
  });

  const isTablet = useMediaQuery("(max-width: 1024px)", {
    ssrWidth,
  });

  const isDesktop = useMediaQuery("(min-width: 1025px)", {
    ssrWidth,
  });

  // 用计算属性统一给组件读取的 width
  const width = computed(() => {
    // SSR 渲染阶段返回 cookie 里的宽度
    if (import.meta.server) return ssrWidth;

    // 客户端自动返回真实宽度
    return clientWidth.value;
  });

  // 客户端时同步写入 cookie（让下一次 SSR 有正确宽度）
  if (import.meta.client) {
    watch(
      clientWidth,
      (val) => {
        widthCookie.value = String(val);
      },
      { immediate: true },
    );
  }

  return {
    width,
    isMobile,
    isTablet,
    isDesktop,
  };
}
