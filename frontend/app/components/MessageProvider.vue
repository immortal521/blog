<script setup lang="ts">
interface Props {
  duration?: number;
  keepAliveOnHover?: boolean;
  max?: number;
  closable?: boolean;
}

const {
  max,
  duration: defaultDuration = 2000,
  keepAliveOnHover: defaultKeepAliveOnHover = false,
  closable: defaultClosable = false,
} = defineProps<Props>();

interface Message {
  id: number;
  text: string;
  type: MessageType;
  keepAliveOnHover: boolean;
  icon: MessageIconProps | null;
  closable: boolean;
  timer?: ReturnType<typeof setTimeout>;
  remainTime?: number;
  startTime?: number;
}

const messageIconMap: Record<MessageType, MessageIconProps> = {
  default: { name: "" },
  success: { name: "icon-park-solid:success", color: "var(--color-success)" },
  error: { name: "icon-park-solid:error", color: "var(--color-danger)" },
  info: { name: "fluent:info-sparkle-24-filled", color: "var(--color-info)" },
  warning: { name: "typcn:warning", color: "var(--color-warning)" },
};

const defaultOptions: Required<MessageOptions> = {
  duration: defaultDuration,
  keepAliveOnHover: defaultKeepAliveOnHover,
  icon: null,
  closable: defaultClosable,
};

let seed = 0;
const messages = ref<Message[]>([]);

// Create a new message
const create = (type: MessageType, text: string, options?: MessageOptions) => {
  const merged = { ...defaultOptions, ...options };
  const msg: Message = {
    id: ++seed,
    text,
    type,
    keepAliveOnHover: merged.keepAliveOnHover,
    icon: merged.icon,
    closable: merged.closable,
  };

  if (max && messages.value.length >= max) removeOldestMessage();
  messages.value.push(msg);
  startTimer(msg, merged.duration);
};

// Remove the oldest message when max is reached
const removeOldestMessage = () => {
  const removed = messages.value.shift();
  if (!removed) return;
  if (removed.timer) clearTimeout(removed.timer);
};

// Remove a message by ID
const remove = (id: number) => {
  const index = messages.value.findIndex((m) => m.id === id);
  if (index !== -1) messages.value.splice(index, 1);
};

// Destroy a message and clear its timer
const destroy = (id: number) => {
  const index = messages.value.findIndex((m) => m.id === id);
  if (index === -1) return;
  const removed = messages.value.splice(index, 1);
  if (!removed[0]) return;
  if (removed[0].timer) clearTimeout(removed[0].timer);
};

// Start the auto-dismiss timer for a message
const startTimer = (msg: Message, duration: number) => {
  msg.startTime = Date.now();
  msg.remainTime = duration;
  msg.timer = setTimeout(() => {
    remove(msg.id);
  }, duration);
};

// Pause the timer (on hover)
const pauseTimer = (msg: Message) => {
  if (!msg.timer) return;
  clearTimeout(msg.timer);
  msg.remainTime = (msg.remainTime ?? 0) - (Date.now() - (msg.startTime ?? 0));
  msg.timer = undefined;
};

// Resume the timer (on mouse leave)
const resumeTimer = (msg: Message) => {
  if (msg.timer || (msg.remainTime ?? 0) <= 0) return;
  startTimer(msg, msg.remainTime!);
};

// Generate message API for a specific type
const createMessageApi =
  (type: Exclude<MessageType, "default">) => (text: string, options?: MessageOptions) =>
    create(type, text, options);

// Message API object
const messageApi = {
  create, // General create function to send any type (including default)
  success: createMessageApi("success"),
  error: createMessageApi("error"),
  info: createMessageApi("info"),
  warning: createMessageApi("warning"),
};

provide("message", messageApi);

// TransitionGroup hook to fix leave animation jump
const beforeLeave = (el: Element) => {
  if (!(el instanceof HTMLElement)) return;
  const top = el.offsetTop;
  el.style.position = "absolute";
  el.style.top = `${top}px`;
};
</script>

<template>
  <slot />
  <Teleport defer to="body">
    <div class="message-container">
      <TransitionGroup name="message-list" @before-leave="beforeLeave">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message', `message-${msg.type}`]"
          @mouseenter="msg.keepAliveOnHover && pauseTimer(msg)"
          @mouseleave="msg.keepAliveOnHover && resumeTimer(msg)"
        >
          <div class="icon">
            <Icon
              :name="msg.icon?.name || messageIconMap[msg.type].name"
              :style="{ color: msg.icon?.color || messageIconMap[msg.type].color }"
            />
          </div>
          <span class="msg-text">
            {{ msg.text }}
          </span>
          <div v-if="msg.closable" class="close-btn" @click="destroy(msg.id)">
            <Icon name="material-symbols:close" size="20" />
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style lang="less" scoped>
.message-container {
  position: fixed;
  width: 100vw;
  top: 20px;
  pointer-events: none;
  z-index: 10000;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.message {
  position: relative;
  display: flex;
  width: max-content;
  max-width: 600px;
  margin: 5px 0;
  padding: 10px;
  border: 2px solid var(--border-color-base);
  border-radius: 10px;
  background-color: var(--bg-card-base);
  color: var(--text-color-base);
  backdrop-filter: blur(15px);
  overflow-wrap: break-word;
  pointer-events: auto;
  box-shadow: var(--shadow-md);

  .icon,
  .close-btn {
    width: 20px;
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .close-btn {
    cursor: pointer;
  }

  .msg-text {
    max-width: 520px;
    font-size: 1.8rem;
    margin: 0 10px;
  }

  &.message-default {
    border-color: var(--border-color-base);
  }
  &.message-success {
    border-color: var(--color-success);
    .icon {
      color: var(--color-success);
    }
  }
  &.message-info {
    border-color: var(--color-info);
    .icon {
      color: var(--color-info);
    }
  }
  &.message-error {
    border-color: var(--color-danger);
    .icon {
      color: var(--color-danger);
    }
  }
  &.message-warning {
    border-color: var(--color-warning);
    .icon {
      color: var(--color-warning);
    }
  }
}

.message-list-move,
.message-list-enter-active,
.message-list-leave-active {
  transition: all 0.5s ease;
}

.message-list-enter-from,
.message-list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
