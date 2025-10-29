import type { MessageOptions } from "~/components/MessageProvider.vue";

export interface MessageApi {
  success: (msg: string, options?: MessageOptions) => void;
  error: (msg: string, options?: MessageOptions) => void;
  info: (msg: string, options?: MessageOptions) => void;
  warning: (msg: string, options?: MessageOptions) => void;
}

export function useMessage(): MessageApi {
  const api = inject<MessageApi>("message");
  if (!api) {
    throw new Error("useMessage");
  }

  return api;
}
