export interface CollapseTransitionOptions {
  getHeight?: (el: Element) => string;
}

export const useCollapseTransition = (options: CollapseTransitionOptions = {}) => {
  const getHeight = options.getHeight ?? ((el) => `${el.scrollHeight}px`);

  const onBeforeEnter = (el: Element) => {
    const e = el as HTMLElement;
    e.style.height = "0";
  };

  const onEnter = (el: Element) => {
    const e = el as HTMLElement;
    void e.offsetHeight;
    e.style.height = getHeight(e);
  };

  const onAfterEnter = (el: Element) => {
    const e = el as HTMLElement;
    e.style.height = "";
  };

  const onBeforeLeave = (el: Element) => {
    const e = el as HTMLElement;
    e.style.height = getHeight(e);
  };

  const onLeave = (el: Element) => {
    const e = el as HTMLElement;
    void e.offsetHeight;
    e.style.height = "0";
  };

  const onAfterLeave = (el: Element) => {
    const e = el as HTMLElement;
    e.style.height = "0";
  };

  return {
    onBeforeEnter,
    onEnter,
    onAfterEnter,
    onBeforeLeave,
    onLeave,
    onAfterLeave,
  };
};
