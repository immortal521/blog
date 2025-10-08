<script setup lang="ts">
interface Props {
  id?: string;
  type?: "text" | "password" | "number" | "email" | "tel";
  placeholder?: string;
  disabled?: boolean;
  readonly?: boolean;
}

const value = defineModel<string | number>();

const {
  id = "",
  type: inputType = "text",
  disabled = false,
  placeholder = "",
  readonly = false,
} = defineProps<Props>();

const isPasswordVisible = ref(false);
const currentType = computed(() => {
  if (inputType === "password" && isPasswordVisible.value) {
    return "text";
  }
  return inputType;
});

const togglePasswordVisibility = () => {
  if (disabled || readonly) return;
  isPasswordVisible.value = !isPasswordVisible.value;
};
</script>

<template>
  <div class="input-wrapper">
    <span v-if="$slots.prefix" class="input-prefix">
      <slot name="prefix" />
    </span>
    <input
      :id
      v-model="value"
      class="base-input-field"
      :type="currentType"
      :disabled
      :placeholder
      :readonly
    />
    <div class="input-suffix-group">
      <button
        v-if="inputType === 'password'"
        type="button"
        :aria-label="isPasswordVisible ? '隐藏密码' : '显示密码'"
        class="password-toggle-button"
        :disabled="disabled"
        @click="togglePasswordVisibility"
      >
        <Icon
          :name="isPasswordVisible ? 'weui:eyes-on-filled' : 'weui:eyes-off-filled'"
          class="icon"
        />
      </button>

      <span v-if="$slots.suffix" class="input-suffix">
        <slot name="suffix"></slot>
      </span>
    </div>
  </div>
</template>

<style lang="less" scoped>
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  color: inherit;
  border-radius: inherit;
  border: 2px solid var(--border-color-base);
  transition: border 0.2s ease-in-out;

  &:focus-within {
    border: 2px solid var(--color-primary-base);
    box-shadow: 0 0 0 1px var(--color-primary-base);
  }

  .input-prefix,
  .input-suffix-group {
    display: flex;
    align-items: center;
    padding: 8px 0; // 统一垂直填充，保证高度
    color: var(--color-text-secondary, #909399);
    font-size: 14px;
  }

  .input-prefix {
    padding-left: 10px;
    padding-right: 8px;
  }
  .base-input-field {
    flex-grow: 1;

    border: none;
    outline: none;

    padding: 8px;
    margin: 0;

    color: inherit;
    background-color: transparent;
  }

  .input-suffix-group {
    padding-right: 10px; // 后缀组与右侧边框的距离
  }

  // 密码切换按钮样式
  .password-toggle-button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0; // 按钮自身不再需要太多 padding
    margin-left: 4px;
    margin-right: 4px; // 按钮与后缀插槽/输入框之间的间距
    font-size: 18px;
    line-height: 1;
    transition: color 0.2s;

    .icon {
      color: var(--color-text-secondary);

      &:hover:not(.password-toggle-button:disabled) {
        color: var(--color-primary-base);
      }
    }
  }
}
</style>
