export type ThemeMode = "light" | "dark";

export interface ThemeColors {
  // Brand scale (50-900)
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

  // Semantic brand tokens
  primary: string;
  hover: string;
  active: string;
  disabled: string;

  // Brand backgrounds
  bg: string;
  bgActive: string;
  bgMuted: string;

  // On-brand text colors
  onPrimary: string;
  onPrimaryHover: string;
  onPrimaryActive: string;
}
