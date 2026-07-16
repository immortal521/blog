import type { ApiResponse } from "~/types/api";
import type { User } from "~/types/user";

export type AuthData = {
  accessToken: string;
} & User;

export const useAuthStore = defineStore("auth", () => {
  const { apiFetch } = useClientApi();

  const accessToken = ref(import.meta.client ? localStorage.getItem("accessToken") || "" : "");

  const setAccessToken = (newAccessToken: string) => {
    accessToken.value = newAccessToken;
    localStorage.setItem("accessToken", newAccessToken);
  };

  const logout = async () => {
    await apiFetch("/auth/logout", {
      method: "POST",
    });
    accessToken.value = "";
    localStorage.removeItem("accessToken");
    navigateTo("/auth/login");
  };

  const login = async (email: string, password: string) => {
    const { code, data } = await apiFetch<ApiResponse<AuthData>>("/auth/login", {
      method: "POST",
      body: { email, password },
    });
    if (code !== 200) {
      // TODO: login failed logic
      console.warn(code);
    }
    setAccessToken(data.accessToken);

    // TODO: login successed logic

    navigateTo("/admin");
  };

  async function refresh() {
    const res = await $fetch<ApiResponse<{ accessToken: string }>>("/api/v1/auth/refresh", {
      method: "POST",
      credentials: "include",
    });

    accessToken.value = res.data.accessToken;
    localStorage.setItem("accessToken", res.data.accessToken);
  }

  return {
    accessToken,
    setAccessToken,
    refresh,
    login,
    logout,
  };
});
