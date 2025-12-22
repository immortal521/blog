export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return;

  if (!to.fullPath.includes("admin")) return;

  const authStore = useAuthStore();
  const localePath = useLocalePath();

  const { accessToken, lagout } = authStore;

  if (!accessToken) {
    lagout();
    return navigateTo(localePath("/auth/login"));
  }
});
