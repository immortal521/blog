declare type ThemeMode = "light" | "dark";

declare interface ThemeColors {
  base: string;
  hover: string;
  active: string;
  disabled: string;
  onPrimary: string; // 用于按钮文字（黑 or 白）
}
