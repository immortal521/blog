<script setup lang="ts">
import "viewerjs/dist/viewer.css";
const props = defineProps<{
  src: string;
}>();

const el = useTemplateRef("el");

const { update } = useViewer(el);

onMounted(() => {
  update();
});
</script>

<template>
  <div ref="el">
    <img :src="props.src" class="img" />
  </div>
</template>

<style lang="less" scoped>
.img {
  object-fit: cover;
  width: 100%;
  height: 100%;
}
</style>
