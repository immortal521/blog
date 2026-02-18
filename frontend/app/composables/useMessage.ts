import type { MessageApi } from "~/types/ui/message";

export function useMessage(): MessageApi {
  const api = inject<MessageApi>("message");
  if (!api) {
    throw new Error("useMessage must be used within a MessageProvider");
  }

  return api;
}
