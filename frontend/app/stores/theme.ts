import {
  applyBaseThemeMode,
  applyThemeColorsToCSSVars,
  generateThemeColors,
  getInitialMode,
  getInitialPrimaryColor,
  withViewTransition,
} from "@/utils/theme";

export const useThemeStore = defineStore("theme", () => {
  // 服务端渲染时没有 window，mode 初始化为空（或默认）
  const mode = ref<ThemeMode | null>(null);
  const primaryColor = ref<string>(getInitialPrimaryColor());
  const themeColors = ref<ThemeColors>(generateThemeColors(primaryColor.value, "light")); // 先用 light 占位

  // 仅客户端执行初始化，读取实际初始值并同步
  if (import.meta.client) {
    const initialMode = getInitialMode();
    mode.value = initialMode;

    const newColors = generateThemeColors(primaryColor.value, initialMode);
    themeColors.value = newColors;

    // 同步到 DOM 和 localStorage
    applyBaseThemeMode(initialMode);
    applyThemeColorsToCSSVars(newColors);
    localStorage.setItem("theme-mode", initialMode);
    localStorage.setItem("theme-primary-color", primaryColor.value);
  }

  // 设置主题模式：Light / Dark
  const setMode = (newMode: ThemeMode) => {
    if (newMode === mode.value) return;

    const newColors = generateThemeColors(primaryColor.value, newMode);
    themeColors.value = newColors;
    mode.value = newMode;

    withViewTransition(() => {
      if (import.meta.client) {
        localStorage.setItem("theme-mode", newMode);
        applyBaseThemeMode(newMode);
        applyThemeColorsToCSSVars(newColors);
      }
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
        localStorage.setItem("theme-primary-color", newColor);
        applyThemeColorsToCSSVars(newColors);
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
