import VueViewer from "v-viewer";
import "viewerjs/dist/viewer.min.css";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(VueViewer, {
    name: "viewer",
    defaultOptions: {
      navbar: 4,
      button: false,
      toolbar: {
        zoomIn: 4,
        zoomOut: 4,
        oneToOne: 4,
        reset: 4,
        prev: 4,
        play: {
          show: 4,
          size: "large",
        },
        next: 4,
        rotateLeft: 4,
        rotateRight: 4,
        flipHorizontal: 4,
        flipVertical: 4,
      },
    },
  });
});
