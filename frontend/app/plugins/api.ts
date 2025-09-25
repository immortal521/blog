import { useAuthStore } from "~/stores/auth";

let isRefreshing = false;
const pendingQueue: Array<() => void> = [];

export default defineNuxtPlugin((nuxtApp) => {
  const api: typeof $fetch = $fetch.create({
    baseURL: import.meta.server ? "http://localhost:8000/api/v1" : "/api/v1",

    async onRequest({ options }) {
      const authStore = useAuthStore();

      const token = import.meta.server ? useCookie("accessToken").value : authStore.accessToken;

      if (token) {
        options.headers.set("Authorization", `Bearer ${token}`);
      }
    },

    async onResponseError({ response, request, options }) {
      if (response.status !== 401) throw response;

      const authStore = useAuthStore();
      let newAccessToken: string | undefined;

      if (import.meta.server) {
        try {
          const data = await $fetch<{ accessToken: string; refreshToken: string }>(
            "/api/v1/auth/refresh",
            {
              method: "POST" as const,
            },
          );
          nuxtApp.runWithContext(() => {
            useCookie("accessToken", { httpOnly: true }).value = data.accessToken;
            useCookie("refreshToken", { httpOnly: true }).value = data.refreshToken;
          });
          newAccessToken = data.accessToken;
        } catch (err) {
          console.error("SSR Token Refresh Failed:", err);
          nuxtApp.runWithContext(() => {
            useCookie("accessToken").value = undefined;
            useCookie("refreshToken").value = undefined;
          });
          throw response;
        }
      } else {
        if (isRefreshing) {
          await new Promise<void>((resolve) => pendingQueue.push(resolve));
        } else {
          isRefreshing = true;

          try {
            const data = await $fetch<{ accessToken: string; refreshToken: string }>(
              "/api/v1/auth/refresh",
              {
                method: "POST" as const,
              },
            );

            authStore.setTokens(data.accessToken, data.refreshToken);
            pendingQueue.forEach((resolve) => resolve());
            pendingQueue.splice(0);
          } catch (err) {
            pendingQueue.splice(0);
            console.error("CSR Token Refresh Failed:", err);
            throw err;
          } finally {
            isRefreshing = false;
          }
        }
        newAccessToken = authStore.accessToken;
      }
      if (!newAccessToken) {
        throw response;
      }

      const newHeaders = new Headers(options.headers || {});
      newHeaders.set("Authorization", `Bearer ${newAccessToken}`);

      return api(request, {
        ...options,
        headers: newHeaders,
        method: options.method as typeof options.method &
          ("GET" | "HEAD" | "PATCH" | "POST" | "PUT" | "DELETE" | "CONNECT" | "OPTIONS" | "TRACE"),
      });
    },
  });

  return {
    provide: {
      api,
    },
  };
});
