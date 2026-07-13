import { Badge } from "#components";
import { createVNode, render } from "vue";

interface BadgeInstance {
  container: HTMLElement;
}

const badgeInstances = new WeakMap<HTMLElement, BadgeInstance>();

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive("badge", {
    mounted(el: HTMLElement, binding) {
      const container = document.createElement("span");

      container.style.position = "absolute";
      container.style.top = "0";
      container.style.right = "0";

      el.style.position = "relative";

      el.appendChild(container);

      const vnode = createVNode(Badge, {
        value: binding.value,
      });

      render(vnode, container);

      badgeInstances.set(el, {
        container,
      });
    },

    updated(el: HTMLElement, binding) {
      const instance = badgeInstances.get(el);

      if (!instance) {
        return;
      }

      const vnode = createVNode(Badge, {
        value: binding.value,
      });

      render(vnode, instance.container);
    },

    unmounted(el: HTMLElement) {
      const instance = badgeInstances.get(el);

      if (!instance) {
        return;
      }

      render(null, instance.container);

      instance.container.remove();

      badgeInstances.delete(el);
    },
  });
});
