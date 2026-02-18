export type MessageType = "success" | "error" | "info" | "warning" | "default";

export type MessageSizeType = "small" | "medium" | "large";

export interface MessageOptions {
  duration?: number;
  keepAliveOnHover?: boolean;
  icon?: MessageIconProps | null;
  closable?: boolean;
  size?: MessageSizeType;
}

export interface MessageIconProps {
  name: string;
  color?: string;
  size?: number;
}

export interface MessageApi {
  create: (type: MessageType, msg: string, options?: MessageOptions) => void;
  success: (msg: string, options?: MessageOptions) => void;
  error: (msg: string, options?: MessageOptions) => void;
  info: (msg: string, options?: MessageOptions) => void;
  warning: (msg: string, options?: MessageOptions) => void;
}
