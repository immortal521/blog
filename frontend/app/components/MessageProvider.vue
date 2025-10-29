<script setup lang="ts">
type MessageType = "success" | "error" | "info" | "warning";

interface Props {
  duration?: number;
  keepAliveOnHover?: boolean;
}

const { duration: defaultDuration = 2000, keepAliveOnHover: defaultKeepAliveOnHover = false } =
  defineProps<Props>();

interface Message {
  id: number;
  text: string;
  type: MessageType;
  keepAliveOnHover: boolean;
  timer?: ReturnType<typeof setTimeout>;
  remainTime?: number;
  startTime?: number;
}

export interface MessageOptions {
  duration?: number;
  keepAliveOnHover?: boolean;
}

const defaultOptions: Required<MessageOptions> = {
  duration: defaultDuration,
  keepAliveOnHover: defaultKeepAliveOnHover,
};

let seed = 0;
const messages = ref<Message[]>([]);

const create = (type: MessageType, text: string, options?: MessageOptions) => {
  const merged = { ...defaultOptions, ...options };
  const id = ++seed;
  const msg: Message = { id, text, type, keepAliveOnHover: merged.keepAliveOnHover };
  messages.value.push(msg);
  startTimer(msg, merged.duration);
};

const remove = (id: number) => {
  const index = messages.value.findIndex((m) => m.id === id);
  if (index !== -1) messages.value.splice(index, 1);
};

const startTimer = (msg: Message, duration: number) => {
  msg.startTime = Date.now();
  msg.timer = setTimeout(() => {
    remove(msg.id);
  }, duration);
  msg.remainTime = duration;
};

const pauseTimer = (msg: Message) => {
  if (!msg.timer) return;
  clearTimeout(msg.timer);
  msg.remainTime = msg.remainTime! - (Date.now() - (msg.startTime ?? 0));
  msg.timer = undefined;
};

const resumeTimer = (msg: Message) => {
  if (msg.timer || msg.remainTime! <= 0) return;
  startTimer(msg, msg.remainTime!);
};

const messageApi = {
  success: (msg: string, options?: MessageOptions) => create("success", msg, options),
  error: (msg: string, options?: MessageOptions) => create("error", msg, options),
  info: (msg: string, options?: MessageOptions) => create("info", msg, options),
  warning: (msg: string, options?: MessageOptions) => create("warning", msg, options),
};

provide("message", messageApi);
</script>

<template>
  <slot />
  <Teleport to="body">
    <div class="message-container">
      <TransitionGroup name="message-list" tag="div">
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="message"
          @mouseenter="msg.keepAliveOnHover && pauseTimer(msg)"
          @mouseleave="msg.keepAliveOnHover && resumeTimer(msg)"
        >
          {{ msg.text }}
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

  [name="message-list"] {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
}

.message {
  position: relative;
  background-color: var(--bg-card-base);
  color: var(--text-color-base);
  width: 90%;
  max-width: 600px;
  border: 1px solid var(--border-color-base);
  margin: 5px 0;
  padding: 10px 20px;
  border-radius: 10px;
  backdrop-filter: blur(15px);
  overflow-wrap: break-word;
  pointer-events: auto;
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

.message-list-leave-active {
  position: absolute;
}
</style>
