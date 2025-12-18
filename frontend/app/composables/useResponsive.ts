/**
 * useResponsive
 *
 * 提供基于视口宽度的响应式布局状态，兼容 SSR / CSR 场景。
 *
 * @description
 * 基于 useClientWidth，对窗口宽度进行统一管理，
 * 并在此基础上派生出常用的布局判断。
 *
 * - 在 SSR 阶段即可得到“合理”的布局判断结果，避免首屏闪动
 * - 不直接依赖 matchMedia，所有响应式判断基于同一宽度数据源
 *
 * 判定规则：
 * - width < 768  → 移动端
 * - width ≥ 768  → 桌面端
 *
 * 返回值：
 * @returns {Object}
 * @returns {Ref<number>} returns.width 当前视口宽度
 * @returns {Ref<boolean>} returns.isMobile 是否为移动端布局
 * @returns {Ref<boolean>} returns.isDesktop 是否为桌面端布局
 *
 * @example
 * ```ts
 * const { width, isMobile, isDesktop } = useResponsive();
 *
 * if (isDesktop.value) {
 *   // 桌面端布局逻辑
 * }
 * ```
 */ export function useResponsive(): {
  width: Ref<number>;
  isMobile: Ref<boolean>;
  isDesktop: Ref<boolean>;
} {
  const { width } = useClientWidth();
  const isMobile = computed(() => width.value < 768);
  const isDesktop = computed(() => width.value >= 768);

  return {
    width,
    isMobile,
    isDesktop,
  };
}
