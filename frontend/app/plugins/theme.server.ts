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
      style: `--color-primary-base: ${colors.base};
              --color-primary-hover: ${colors.hover};
              --color-primary-active: ${colors.active};
              --color-primary-disabled: ${colors.disabled};
              --color-on-primary: ${colors.onPrimary};
      `,
    },
  });
});
