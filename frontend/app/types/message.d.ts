type MessageType = "success" | "error" | "info" | "warning" | "default";

type MessageSizeType = "small" | "medium" | "large";

interface MessageOptions {
  duration?: number;
  keepAliveOnHover?: boolean;
  icon?: MessageIconProps | null;
  closable?: boolean;
  size?: MessageSizeType;
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
