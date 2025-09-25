export const useAuthStore = defineStore("auth", () => {
  const accessToken = ref(useCookie("accessToken").value || "");
  const refreshToken = ref(useCookie("refreshToken").value || "");

  const saveToCookies = (accessToken: string, refreshToken: string) => {
    useCookie("accessToken", { httpOnly: true }).value = accessToken;
    useCookie("refreshToken", { httpOnly: true }).value = refreshToken;
  };

  const setTokens = (newAccessToken: string, newRefreshToken: string) => {
    accessToken.value = newAccessToken;
    refreshToken.value = newRefreshToken;
    saveToCookies(newAccessToken, newRefreshToken);
  };

  const clearCookies = () => {
    useCookie("accessToken").value = null;
    useCookie("refreshToken").value = null;
  };

  const clear = () => {
    accessToken.value = "";
    refreshToken.value = "";
    clearCookies();
  };

  return {
    accessToken,
    refreshToken,
    setTokens,
    clear,
  };
});
