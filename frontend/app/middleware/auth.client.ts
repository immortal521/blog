export default defineNuxtRouteMiddleware(async (to) => {
  if (!to.path.includes("/admin")) return;

  const auth = useAuthStore();

  if (auth.accessToken) return;

  try {
    await auth.refresh();
  } catch {
    auth.logout();

    return navigateTo("/auth/login");
  }
});
