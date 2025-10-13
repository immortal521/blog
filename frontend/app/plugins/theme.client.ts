export default defineNuxtPlugin(() => {
  const mode = useCookie<ThemeMode>("theme-mode");
  if (!mode.value) {
    mode.value =
      window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light";
  }
  applyBaseThemeMode(mode.value);
});
