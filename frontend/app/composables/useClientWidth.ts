import { useDebounceFn } from "@vueuse/core";

export function useClientWidth() {
  const widthCookie = useCookie<number>("client_width", {
    path: "/",
    maxAge: 356 * 24 * 60 * 60,
  });

  const getWidthFromUA = () => {
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
