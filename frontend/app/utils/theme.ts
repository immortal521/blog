import tinycolor from "tinycolor2";
import { nextTick } from "vue";
import type { ThemeColors, ThemeMode } from "~/types/theme";

interface BrandScale {
  brand50: string;
  brand100: string;
  brand200: string;
  brand300: string;
  brand400: string;
  brand500: string;
  brand600: string;
  brand700: string;
  brand800: string;
  brand900: string;
}

function generateBrandScale(baseColor: string): BrandScale {
  const base = tinycolor(baseColor);
  const { h, s, l } = base.toHsl();

  return {
    brand50: tinycolor({ h, s: Math.max(s * 0.3, 0.1), l: 0.97 }).toHexString(),
    brand100: tinycolor({ h, s: Math.max(s * 0.4, 0.15), l: 0.93 }).toHexString(),
    brand200: tinycolor({ h, s: Math.max(s * 0.5, 0.2), l: 0.86 }).toHexString(),
    brand300: tinycolor({ h, s: Math.max(s * 0.6, 0.25), l: 0.76 }).toHexString(),
    brand400: tinycolor({ h, s: Math.max(s * 0.8, 0.3), l: 0.64 }).toHexString(),
    brand500: base.toHexString(),
    brand600: tinycolor({ h, s: Math.min(s * 1.1, 1), l: Math.max(l - 0.15, 0.2) }).toHexString(),
    brand700: tinycolor({ h, s: Math.min(s * 1.2, 1), l: Math.max(l - 0.25, 0.15) }).toHexString(),
    brand800: tinycolor({ h, s: Math.min(s * 1.3, 1), l: Math.max(l - 0.35, 0.1) }).toHexString(),
    brand900: tinycolor({ h, s: Math.min(s * 1.4, 1), l: Math.max(l - 0.45, 0.05) }).toHexString(),
  };
}

export function generateThemeColors(baseColor: string, mode: ThemeMode): ThemeColors {
  const scale = generateBrandScale(baseColor);

  const hover =
    mode === "light" ? tinycolor(scale.brand500).darken(5) : tinycolor(scale.brand500).lighten(5);
  const active =
    mode === "light" ? tinycolor(scale.brand500).darken(10) : tinycolor(scale.brand500).lighten(10);
  const disabled = tinycolor.mix(scale.brand500, "#b0b0b0", 60);

  const bgActive =
    mode === "light"
      ? tinycolor(scale.brand200).lighten(5).toHexString()
      : tinycolor(scale.brand800).lighten(5).toHexString();

  const bgMuted =
    mode === "light"
      ? tinycolor(scale.brand100).lighten(3).toHexString()
      : tinycolor(scale.brand900).lighten(3).toHexString();

  const whiteContrast = tinycolor.readability(scale.brand500, "#ffffff");
  const onPrimary = whiteContrast >= 4.5 ? "#ffffff" : "#1a1a1e";

  const hoverContrast = tinycolor.readability(hover, "#ffffff");
  const onPrimaryHover = hoverContrast >= 4.5 ? "#ffffff" : "#1a1a1e";

  const activeContrast = tinycolor.readability(active, "#ffffff");
  const onPrimaryActive = activeContrast >= 4.5 ? "#ffffff" : "#1a1a1e";

  return {
    ...scale,
    primary: scale.brand500,
    hover: hover.toHexString(),
    active: active.toHexString(),
    disabled: disabled.toHexString(),
    bg: scale.brand50,
    bgActive,
    bgMuted,
    onPrimary,
    onPrimaryHover,
    onPrimaryActive,
  };
}

export function applyThemeColorsToCSSVars(colors: ThemeColors) {
  const root = document.documentElement;

  // Brand scale
  root.style.setProperty("--brand-50", colors.brand50);
  root.style.setProperty("--brand-100", colors.brand100);
  root.style.setProperty("--brand-200", colors.brand200);
  root.style.setProperty("--brand-300", colors.brand300);
  root.style.setProperty("--brand-400", colors.brand400);
  root.style.setProperty("--brand-500", colors.brand500);
  root.style.setProperty("--brand-600", colors.brand600);
  root.style.setProperty("--brand-700", colors.brand700);
  root.style.setProperty("--brand-800", colors.brand800);
  root.style.setProperty("--brand-900", colors.brand900);

  // Semantic brand tokens
  root.style.setProperty("--brand-primary", colors.primary);
  root.style.setProperty("--brand-hover", colors.hover);
  root.style.setProperty("--brand-active", colors.active);
  root.style.setProperty("--brand-disabled", colors.disabled);

  // Brand backgrounds
  root.style.setProperty("--brand-bg", colors.bg);
  root.style.setProperty("--brand-bg-active", colors.bgActive);
  root.style.setProperty("--brand-bg-muted", colors.bgMuted);

  // Text on brand
  root.style.setProperty("--text-on-brand", colors.onPrimary);
  root.style.setProperty("--text-on-brand-hover", colors.onPrimaryHover);
  root.style.setProperty("--text-on-brand-active", colors.onPrimaryActive);

  // Derived tokens
  root.style.setProperty("--border-active", colors.primary);
  root.style.setProperty("--border-color-card-hover", colors.hover);
  root.style.setProperty("--bg-sidebar-item-hover", colors.bgActive);
  root.style.setProperty("--focus-ring", `0 0 0 2px ${colors.bgMuted}`);
  root.style.setProperty("--selection-bg", colors.primary);
  root.style.setProperty("--scrollbar-thumb-bg", colors.primary);
  root.style.setProperty("--scrollbar-thumb-hover", colors.hover);
  root.style.setProperty("--scrollbar-thumb-active", colors.active);
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
