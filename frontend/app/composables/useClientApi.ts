import { ofetch, type FetchError, type FetchOptions } from "ofetch";

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
  async function doRequest<T>(
    url: string,
    options: FetchOptions<"json", unknown> = {},
  ): Promise<T> {
    const accessToken = useAuthStore().accessToken;
    const headers = accessToken
      ? { ...(options.headers || {}), Authorization: `Bearer ${accessToken}` }
      : options.headers;
    return $baseFetch<T>(url, { ...options, headers });
  }
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

    // Wait if a token refresh is in progress
    if (isRefreshing) {
      await new Promise<void>((resolve) => pendingRequests.push(resolve));
    }

    try {
      return await doRequest<T>(url, options);
    } catch (err: unknown) {
      const fetchErr = err as FetchError;
      if (fetchErr?.response?.status !== 401) throw fetchErr;
      await refreshAccessToken();
      return doRequest<T>(url, options).catch((err: FetchError) => {
        if (err?.response?.status === 401) {
          // TODO: layout logic
          console.log(err);
        }
        throw err;
      });
    }
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
