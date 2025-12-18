export default defineNuxtRouteMiddleware((to, from) => {
  const { start } = useLoadingBar();

  if (from.name && to.fullPath === from.fullPath) {
    return;
  }

  start();
});
