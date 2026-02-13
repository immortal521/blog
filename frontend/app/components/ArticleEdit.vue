<script setup lang="ts">
const content = defineModel<string>("content", {
  default: "",
});

const editorRef = useTemplateRef<HTMLDivElement>("editor");

const { isMobile } = useResponsive();

type Mode = "preview" | "edit" | "both";
const mode = ref<Mode>("both");
const userTouchedMode = ref(false);

watch(
  isMobile,
  (m) => {
    if (!userTouchedMode.value) {
      mode.value = m ? "edit" : "both";
    }
  },
  { immediate: true },
);

const isFullscreen = ref(false);

function syncFullscreenState() {
  isFullscreen.value = !!document.fullscreenElement;
}

onMounted(() => {
  document.addEventListener("fullscreenchange", syncFullscreenState);
  syncFullscreenState();
});

onBeforeUnmount(() => {
  document.removeEventListener("fullscreenchange", syncFullscreenState);
});

async function toggleFullscreen() {
  try {
    if (document.fullscreenElement) {
      await document.exitFullscreen();
    } else {
      await editorRef.value?.requestFullscreen();
    }
  } catch (e) {
    console.warn(e);
  }
}

function toggleMode() {
  userTouchedMode.value = true;

  if (isMobile.value) {
    mode.value = mode.value === "edit" ? "preview" : "edit";
    return;
  }

  // PC 三态循环（从 both 开始更符合桌面习惯）
  if (mode.value === "both") mode.value = "edit";
  else if (mode.value === "edit") mode.value = "preview";
  else mode.value = "both";
}
</script>

<template>
  <div ref="editor" class="article-edit">
    <div class="toolbar">
      <div class="tools-left"></div>
      <div class="tools-right">
        <button class="btn" @click="toggleMode">
          <Icon v-if="mode === 'edit'" name="mingcute:eye-line" size="18" />

          <Icon v-else-if="mode === 'preview'" name="mingcute:edit-2-line" size="18" />

          <Icon v-else name="mingcute:layout-grid-line" size="18" />
        </button>
        <button class="btn" @click="toggleFullscreen">
          <Icon v-if="!isFullscreen" name="mingcute:fullscreen-fill" size="18" />
          <Icon v-else name="mingcute:fullscreen-exit-fill" size="18" />
        </button>
      </div>
    </div>
    <div class="main">
      <div v-if="mode === 'edit' || mode === 'both'" class="edit">
        <textarea v-model="content" :class="{ 'split-border': mode === 'both' }"></textarea>
      </div>
      <div v-if="mode === 'preview' || mode === 'both'" class="preview">
        <MarkdownRenderer :markdown="content" />
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.article-edit {
  border: 1.5px solid var(--border-color-divider);
  border-radius: 8px;
  box-shadow: var(--shadow-md);
  width: 96%;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: transparent;
}

.toolbar {
  height: 40px;
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  padding: 5px 15px;
  border-bottom: 1px solid var(--border-color-divider);

  .tools-left {
    display: flex;
    gap: 6px;
    margin-right: auto;
    min-width: 0;
  }

  .tools-right {
    display: flex;
    gap: 6px;
    margin-left: auto;
    min-width: 0;
  }

  .btn {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    color: var(--text-color-base);
    background-color: transparent;
    padding: 6px;
    border-radius: 6px;
    border: none;
    cursor: pointer;
    user-select: none;

    &:hover {
      background-color: var(--border-color-base);
    }

    &:active {
      background-color: var(--border-color-divider);
      transform: scale(0.97);
    }
  }
}

.main {
  width: 100%;
  flex: 1;
  min-height: 0;
  display: flex;
  padding: 8px 0;
  gap: 0;
}

.edit,
.preview {
  flex: 1;
  min-width: 0;
  height: 100%;
}

.edit {
  display: flex;
}

.edit textarea {
  flex: 1;
  width: 100%;
  height: 100%;
  color: var(--text-color-base);
  padding: 10px 20px;
  resize: none;
  font-size: 1.5rem;
  line-height: 1.5;
  outline: none;
  border: none;
  border-radius: 8px;
  background-color: transparent;
  min-width: 0;
  overflow: auto;
}

.edit textarea.split-border {
  border-right: 1px solid var(--border-color-divider);
  border-radius: 8px 0 0 8px;
}

.preview {
  padding: 10px 14px;
  overflow: auto;
}

.preview.split {
  border-radius: 0 8px 8px 0;
}
</style>
