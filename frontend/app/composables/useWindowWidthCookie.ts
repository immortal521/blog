export function useWindowWidthCookie(cookieName = "windowWidth") {
  const width = ref(0);

  function updateWidth() {
    width.value = window.innerWidth;
    useCookie(cookieName, {
      expires: new Date(Date.now() + 356 * 24 * 60 * 60 * 1000),
    }).value = width.value.toString();
  }

  onMounted(() => {
    updateWidth();
    window.addEventListener("resize", updateWidth);
  });

  onBeforeUnmount(() => {
    window.removeEventListener("resize", updateWidth);
  });

  return { width };
}
