import { parseCookies } from "h3";
import type { ThemeMode } from "~/types/theme";

export default defineNuxtPlugin((nuxtApp) => {
  const event = nuxtApp.ssrContext?.event;
  const cookies = event ? parseCookies(event) : {};

  const mode = (cookies["theme-mode"] as ThemeMode) || "light";
  const primaryColor = (cookies["theme-primary-color"] as string) || "99a2ff";

  const colors = generateThemeColors(primaryColor, mode);

  useHead({
    htmlAttrs: {
      "data-theme": mode,
      style: `--brand-50: ${colors.brand50};
              --brand-100: ${colors.brand100};
              --brand-200: ${colors.brand200};
              --brand-300: ${colors.brand300};
              --brand-400: ${colors.brand400};
              --brand-500: ${colors.brand500};
              --brand-600: ${colors.brand600};
              --brand-700: ${colors.brand700};
              --brand-800: ${colors.brand800};
              --brand-900: ${colors.brand900};
              --brand-primary: ${colors.primary};
              --brand-hover: ${colors.hover};
              --brand-active: ${colors.active};
              --brand-disabled: ${colors.disabled};
              --brand-bg: ${colors.bg};
              --brand-bg-active: ${colors.bgActive};
              --brand-bg-muted: ${colors.bgMuted};
              --text-on-brand: ${colors.onPrimary};
              --text-on-brand-hover: ${colors.onPrimaryHover};
              --text-on-brand-active: ${colors.onPrimaryActive};
              --border-active: ${colors.primary};
              --border-color-card-hover: ${colors.hover};
              --bg-sidebar-item-hover: ${colors.bgActive};
              --focus-ring: 0 0 0 2px ${colors.bgMuted};
              --selection-bg: ${colors.primary};
              --scrollbar-thumb-bg: ${colors.primary};
              --scrollbar-thumb-hover: ${colors.hover};
              --scrollbar-thumb-active: ${colors.active};
      `,
    },
  });
});
