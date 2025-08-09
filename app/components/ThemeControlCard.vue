<script setup lang="ts">
const isShow = ref(true);

const primaryColorList = ["#7c85ff", "#3B82F6", "#7C3AED", "#10B981", "#F59E0B", "	#D9777F"];

const isDark = computed(() => {
  return useThemeStore().mode === "dark";
});

const setThemeMode = (mode: ThemeMode) => {
  useThemeStore().setMode(mode);
};

const setPrimaryColor = (primaryColor: string) => {
  useThemeStore().setPrimaryColor(primaryColor);
};

const colorInput = useTemplateRef<HTMLInputElement>("colorInput");
const colorInputValue = ref("");

const handleColorChange = () => {
  setPrimaryColor(colorInputValue.value);
};

const { locale, locales, setLocale } = useI18n();

const options = computed(() => {
  return locales.value.map((l) => {
    return { label: l.name as string, value: l.code as string };
  });
});

const localeRef = ref(locale.value);

onMounted(() => {
  watch(localeRef, () => {
    setLocale(localeRef.value);
  });
});
</script>

<template>
  <div v-if="isShow" ref="themeControlCard" class="theme-control-card">
    <p class="title">{{ $t("themeControl.themeMode") }}</p>
    <div :class="{ 'buttons-track': true, dark: isDark }">
      <button class="buttons-item" @click="setThemeMode('light')">
        <Icon name="streamline-plump-color:sun" :size="20" />
        {{ $t("themeControl.light") }}
      </button>
      <button class="buttons-item" @click="setThemeMode('dark')">
        <Icon name="streamline-plump-color:moon-stars" :size="20" />
        {{ $t("themeControl.dark") }}
      </button>
    </div>
    <p class="title">{{ $t("themeControl.primaryColor") }}</p>
    <div class="color-picker">
      <button
        v-for="item in primaryColorList"
        :key="item"
        class="color-picker-item"
        :style="{ background: item }"
        @click="setPrimaryColor(item)"
      />
      <button class="color-picker-item" @click="colorInput?.click()">+</button>
      <input
        v-show="false"
        ref="colorInput"
        v-model="colorInputValue"
        type="color"
        @change="handleColorChange"
      />
    </div>
    <p class="title">{{ $t("themeControl.language") }}</p>
    <InputSelect v-model:value="localeRef" :options />
  </div>
</template>

<style lang="less" scoped>
.theme-control-card {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  border-radius: 8px;
  background-color: var(--bg-nav-base);
  border: 1px solid var(--border-color-nav);
  box-shadow: var(--shadow-md);
}

.title {
  font-size: 12px;
  font-family: "MapleMono";
  font-weight: 500;
  color: var(--text-color-muted);
}

.buttons-track {
  width: calc(100% - 10px);
  height: 35px;
  background-color: var(--bg-button-toggle-track);
  margin: 0 auto;
  border-radius: 5px;
  position: relative;
  display: flex;

  &::after {
    content: "";
    position: absolute;
    left: 2px;
    top: 2px;
    width: 50%;
    height: calc(100% - 4px);
    background-color: var(--bg-button-toggle-thumb);
    z-index: -1;
    border-radius: 5px;
    transition: transform 0.3s ease-in-out;
    border: var(--border-color-button-toggle);
    backdrop-filter: blur(2px);
  }

  &.dark::after {
    transform: translateX(calc(100% - 4px));
  }
}

.buttons-item {
  width: 50%;
  background-color: transparent;
  border: none;
  cursor: pointer;
  font-weight: 600;
  color: var(--text-color-base);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 5px;
}

.color-picker {
  display: grid;
  height: 32px;
  width: calc(100% - 10px);
  margin: 0 auto;
  grid-template-columns: repeat(7, 32px);
  justify-content: center;
  gap: 6px;
}

.color-picker-item {
  cursor: pointer;
  border: none;
  box-shadow: var(--shadow-md);
  border-radius: 5px;
  background-color: transparent;
}
</style>
