import {
  applyBaseThemeMode,
  applyThemeColorsToCSSVars,
  generateThemeColors,
  getInitialMode,
  getInitialPrimaryColor,
  withViewTransition,
} from "@/utils/theme";

export const useThemeStore = defineStore("theme", () => {
  const mode = ref<ThemeMode>(getInitialMode());
  const primaryColor = ref<string>(getInitialPrimaryColor());
  const themeColors = ref<ThemeColors>(generateThemeColors(primaryColor.value, mode.value)); // 先用 light 占位

  // 设置主题模式：Light / Dark
  const setMode = (newMode: ThemeMode) => {
    if (newMode === mode.value) return;

    const newColors = generateThemeColors(primaryColor.value, newMode);
    themeColors.value = newColors;
    mode.value = newMode;

    withViewTransition(() => {
      applyBaseThemeMode(newMode);
      applyThemeColorsToCSSVars(newColors);
      document.cookie = `theme-mode=${newMode}; path=/`;
    }, newMode === "light");
  };

  // 设置主题颜色
  const setPrimaryColor = (newColor: string) => {
    if (newColor === primaryColor.value) return;

    const newColors = generateThemeColors(newColor, mode.value ?? "light");
    themeColors.value = newColors;
    primaryColor.value = newColor;

    withViewTransition(() => {
      if (import.meta.client) {
        applyThemeColorsToCSSVars(newColors);
        document.cookie = `theme-primary-color=${newColor}; path=/`;
      }
    });
  };

  return {
    mode,
    primaryColor,
    themeColors,
    setMode,
    setPrimaryColor,
  };
});
