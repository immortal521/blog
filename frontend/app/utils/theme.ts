import tinycolor from "tinycolor2";
import { nextTick } from "vue";
import type { ThemeColors, ThemeMode } from "~/types/theme";

function soften(color: tinycolor.Instance): tinycolor.Instance {
  const { h, s, l } = color.toHsl();
  if (l < 0.25 || l > 0.7) return color;
  if (s < 0.55) return color;
  return tinycolor({ h, s: 0.4, l });
}

export function generateThemeColors(baseColor: string, mode: ThemeMode): ThemeColors {
  const base = soften(tinycolor(baseColor));
  const hex = base.toHexString();

  const hover = mode === "light" ? tinycolor(hex).darken(5) : tinycolor(hex).lighten(5);

  const active = mode === "light" ? tinycolor(hex).darken(10) : tinycolor(hex).lighten(10);

  const disabled = tinycolor.mix(hex, "#b0b0b0", 60);

  const whiteContrast = tinycolor.readability(hex, "#ffffff");

  const onPrimary = whiteContrast >= 4.5 ? "#ffffff" : "#1a1a1e";

  return {
    base: hex,
    hover: hover.toHexString(),
    active: active.toHexString(),
    disabled: disabled.toHexString(),
    onPrimary,
  };
}

export function applyThemeColorsToCSSVars(colors: ThemeColors) {
  const root = document.documentElement;

  root.style.setProperty("--color-primary-base", colors.base);
  root.style.setProperty("--color-primary-hover", colors.hover);
  root.style.setProperty("--color-primary-active", colors.active);
  root.style.setProperty("--color-primary-disabled", colors.disabled);
  root.style.setProperty("--color-on-primary", colors.onPrimary);
}

export function applyBaseThemeMode(mode: ThemeMode) {
  const root = document.documentElement;
  root.setAttribute("data-theme", mode);
}

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
