type MessageType = "success" | "error" | "info" | "warning" | "default";

declare interface MessageOptions {
  duration?: number;
  keepAliveOnHover?: boolean;
  icon?: MessageIconProps | null;
  closable?: boolean;
}

interface MessageIconProps {
  name: string;
  color?: string;
  size?: number;
}

interface MessageApi {
  create: (type: MessageType, msg: string, options?: MessageOptions) => void;
  success: (msg: string, options?: MessageOptions) => void;
  error: (msg: string, options?: MessageOptions) => void;
  info: (msg: string, options?: MessageOptions) => void;
  warning: (msg: string, options?: MessageOptions) => void;
}
