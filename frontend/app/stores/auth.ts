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

  const lagout = () => {
    accessToken.value = "";
    localStorage.removeItem("accessToken");
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

  return {
    accessToken,
    setAccessToken,
    login,
    lagout,
  };
});
