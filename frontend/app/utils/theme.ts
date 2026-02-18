import tinycolor from "tinycolor2";
import { nextTick } from "vue";
import type { ThemeColors, ThemeMode } from "~/types/theme";

/**
 * 计算一组主题颜色，包括 hover/active/disabled
 */
export function generateThemeColors(baseColor: string, mode: ThemeMode): ThemeColors {
  const base = tinycolor(baseColor);

  const hover = mode === "light" ? tinycolor(baseColor).darken(5) : tinycolor(baseColor).lighten(5);

  const active =
    mode === "light" ? tinycolor(baseColor).darken(10) : tinycolor(baseColor).lighten(10);

  const disabled = tinycolor.mix(base, "#cccccc", 60);

  const whiteContrast = tinycolor.readability(base, "#ffffff");

  const onPrimary = whiteContrast >= 4.5 ? "#ffffff" : "#000000";

  return {
    base: base.toHexString(),
    hover: hover.toHexString(),
    active: active.toHexString(),
    disabled: disabled.toHexString(),
    onPrimary,
  };
}

/**
 * 把主题颜色同步到 CSS 变量
 */
export function applyThemeColorsToCSSVars(colors: ThemeColors) {
  const root = document.documentElement;

  root.style.setProperty("--color-primary-base", colors.base);
  root.style.setProperty("--color-primary-hover", colors.hover);
  root.style.setProperty("--color-primary-active", colors.active);
  root.style.setProperty("--color-primary-disabled", colors.disabled);
  root.style.setProperty("--color-on-primary", colors.onPrimary);
}

/**
 * 设置整体明暗模式
 */
export function applyBaseThemeMode(mode: ThemeMode) {
  const root = document.documentElement;
  root.setAttribute("data-theme", mode);
}

/**
 * 如果浏览器支持 View Transition API，则对当前文档应用视图过渡。
 * 如果浏览器不支持视图过渡，则立即执行提供的函数。
 */
export function withViewTransition(applyFn: () => void, direction: boolean = true) {
  if (typeof document !== "undefined" && document.startViewTransition) {
    const transition = document.startViewTransition(async () => {
      applyFn();
      await nextTick();
    });

    transition.ready
      .then(() => {
        const innerHeight = window.innerHeight;
        const innerWidth = window.innerWidth;
        const radius = Math.sqrt(innerHeight ** 2 + innerWidth ** 2);

        const clipPath = [`circle(0 at 100% 100%)`, `circle(${radius}px at 100% 100%)`];

        document.documentElement.animate(
          {
            clipPath: direction ? clipPath : [...clipPath].reverse(),
          },
          {
            duration: 400,
            easing: "ease-in",
            fill: "both",
            pseudoElement: direction
              ? "::view-transition-new(root)"
              : "::view-transition-old(root)",
          },
        );
      })
      .catch(console.warn);
  } else {
    applyFn();
  }
}

export function getInitialMode(): ThemeMode {
  if (import.meta.server) return "light";

  const themeMode = useCookie<ThemeMode>("theme-mode", {
    default: () => "light",
  });

  return themeMode.value;
}

export function getInitialPrimaryColor(): string {
  if (import.meta.server) return "#99a2ff";

  const color = useCookie<string>("theme-primary-color");

  return color.value ? decodeURIComponent(color.value) : "#99a2ff";
}
