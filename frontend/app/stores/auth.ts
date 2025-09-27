export const useAuthStore = defineStore("auth", () => {
  const accessToken = ref(localStorage.getItem("accessToken") || "");

  const setAccessToken = (newAccessToken: string) => {
    accessToken.value = newAccessToken;
    localStorage.setItem("accessToken", newAccessToken);
  };

  const clear = () => {
    accessToken.value = "";
    localStorage.removeItem("accessToken");
  };

  return {
    accessToken,
    setAccessToken,
    clear,
  };
});
