import { ofetch, type FetchOptions } from "ofetch";

/**
 * Flag indicating whether the Access Token is currently being refreshed
 */
let isRefreshing = false;

/**
 * Queue of pending requests while a token refresh is in progress
 */
const pendingRequests: (() => void)[] = [];

/**
 * Base ofetch instance with default baseURL and JSON headers
 */
const $baseFetch = ofetch.create({
  baseURL: "/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * Client-side API composable that provides a fetch function with automatic
 * token handling and refresh logic.
 */
export function useClientApi() {
  /**
   * Client-only API request function. Automatically attaches the Access Token
   * and supports token refresh.
   *
   * **Note: This function must only be called on the client.**
   * Calling it during SSR will throw an error.
   *
   * @template T - Type of the returned data
   * @param {string} url - The URL to request (relative or absolute)
   * @param {FetchOptions<"json", unknown>} [options={}] - Optional request configuration
   * @throws {Error} Throws "apiFetch must be called on the client" if called during SSR
   * @returns {Promise<T>} Returns a promise resolving to the data of type T
   *
   * @example
   * import { useClientApi } from '~/composables/useClientApi'
   *
   * const { apiFetch } = useClientApi()
   *
   * onMounted(async () => {
   *   const user = await apiFetch<User>('/api/protected/user', { method: 'get' })
   *   console.log(user)
   * })
   */
  async function apiFetch<T>(url: string, options: FetchOptions<"json", unknown> = {}): Promise<T> {
    if (import.meta.server) {
      throw new Error("apiFetch must be called on the client");
    }

    // Access the store only inside the function to ensure Pinia is active
    const accessToken = storeToRefs(useAuthStore()).accessToken;

    // Wait if a token refresh is in progress
    if (isRefreshing) {
      await new Promise<void>((resolve) => pendingRequests.push(resolve));
    }

    // Refresh token if missing or expired
    if (!accessToken.value || isTokenExpired(accessToken.value)) {
      await refreshAccessToken();
    }

    // Add Authorization header
    options.headers = {
      ...(options.headers || {}),
      Authorization: `Bearer ${accessToken.value}`,
    };

    const data = await $baseFetch<T>(url, options);
    return data;
  }

  /**
   * Refreshes the Access Token.
   *
   * If a refresh is already in progress, this function will return immediately.
   * Once the refresh completes, all pending requests will resume.
   *
   * @async
   */
  async function refreshAccessToken() {
    if (isRefreshing) return;
    isRefreshing = true;
    try {
      const res = await $baseFetch<{ accessToken: string }>("/auth/refresh", {
        method: "POST",
        credentials: "include",
      });
      // Update the token in Pinia store
      useAuthStore().setAccessToken(res.accessToken);
    } finally {
      isRefreshing = false;
      // Resolve all pending requests
      pendingRequests.forEach((resolve) => resolve());
      pendingRequests.splice(0);
    }
  }

  /**
   * Checks whether a JWT Access Token is expired.
   *
   * @param {string} token - The JWT Access Token
   * @returns {boolean} Returns true if the token is invalid or expired, otherwise false
   */
  function isTokenExpired(token: string): boolean {
    try {
      const parts = token.split(".");
      if (parts.length !== 3) return true;
      const payloadStr = parts[1];
      if (!payloadStr) return true;
      const payload = JSON.parse(atob(payloadStr));
      return Date.now() >= payload.exp * 1000;
    } catch {
      return true;
    }
  }

  /**
   * Convenience wrapper for GET requests
   */
  function get<T>(url: string, options: FetchOptions<"json", unknown> = {}) {
    return apiFetch<T>(url, { ...options, method: "get" });
  }

  /**
   * Convenience wrapper for POST requests
   */
  function post<Req extends Record<string, unknown>, Res>(
    url: string,
    body?: Req,
    options: FetchOptions<"json", unknown> = {},
  ) {
    return apiFetch<Res>(url, { ...options, method: "post", body });
  }

  /**
   * Convenience wrapper for PUT requests
   */
  function put<Req extends Record<string, unknown>, Res>(
    url: string,
    body?: Req,
    options: FetchOptions<"json", unknown> = {},
  ) {
    return apiFetch<Res>(url, { ...options, method: "put", body });
  }

  /**
   * Convenience wrapper for DELETE requests
   */
  function del<T>(url: string, options: FetchOptions<"json", unknown> = {}) {
    return apiFetch<T>(url, { ...options, method: "delete" });
  }

  return { apiFetch, get, post, put, delete: del };
}
