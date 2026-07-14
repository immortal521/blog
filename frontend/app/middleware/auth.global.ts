export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) return;

  if (!to.fullPath.includes("admin")) return;

  const authStore = useAuthStore();
  const { $localePath } = useI18n();

  const { accessToken, logout } = authStore;

  if (!accessToken) {
    logout();
    return navigateTo($localePath("/auth/login"));
  }
});
