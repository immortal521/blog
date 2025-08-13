<script setup lang="ts">
interface Props {
  src: string;
  alt?: string;
  lazy?: boolean; // 是否懒加载
}

const props = withDefaults(defineProps<Props>(), {
  alt: "",
  lazy: true,
});

type Status = "loading" | "loaded" | "error";
const status = ref<Status>("loading");
const imgRef = ref<HTMLImageElement | null>(null);
let observer: IntersectionObserver | null = null;

const isModalOpen = ref(false);

const errorPlaceholder =
  "data:image/svg+xml;base64," +
  btoa(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200">
          <rect width="200" height="200" fill="#fdd"/>
          <text x="50%" y="50%" font-size="16" text-anchor="middle" dy=".3em">Error</text>
        </svg>`);

const loadImage = () => {
  if (!imgRef.value) return;
  imgRef.value.src = props.src;
};

onMounted(() => {
  if (props.lazy && imgRef.value) {
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            loadImage();
            observer?.disconnect();
          }
        });
      },
      { threshold: 0.15 },
    );
    observer.observe(imgRef.value);
  } else {
    loadImage();
  }
});

onBeforeUnmount(() => {
  observer?.disconnect();
});

const handleLoad = () => {
  status.value = "loaded";
};

const handleError = () => {
  status.value = "error";
};

const handleClick = () => {
  document.body.style.overflow = "hidden";
  isModalOpen.value = true;
};

const closeFullImage = () => {
  document.body.style.overflow = "";
  isModalOpen.value = false;
};
</script>

<template>
  <div class="image-container">
    <div v-if="status === 'loading'" class="loading"></div>
    <template v-else-if="status === 'error'">
      <img :src="errorPlaceholder" alt="error" class="error" />
    </template>
    <img
      v-if="src"
      ref="imgRef"
      :alt="alt"
      @load="handleLoad"
      @error="handleError"
      @click="handleClick"
    />
    <Teleport to="body">
      <div v-if="isModalOpen" class="img-modal-overlay" @click.self="closeFullImage">
        <div class="modal-content">
          <img :src="src" :alt="alt" />
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.image-container {
  width: 100%;
  display: flex;
  position: relative;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  cursor: pointer;
  overflow: hidden;
}

.loading {
  position: absolute;
  width: 100%;
  height: 100%;
  backdrop-filter: blur(10px);
  z-index: 1;
}

.error {
  position: absolute;
  height: 100%;
  max-height: 200px;
  object-fit: contain;
  z-index: 1;
}

.image-container img {
  width: 100%;
  max-width: 100%;
  max-height: 100%;
  display: block; /* 防止 img 下方空隙 */
}

.img-modal-overlay {
  position: fixed;
  width: 100vw;
  height: 100vh;
  background: #00000020;
  backdrop-filter: blur(5px);
}
</style>
