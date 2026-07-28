<script setup lang="ts">
const { text } = defineProps<{ text: string }>();

const { $ts } = useI18n();

const copied = ref(false);

const buttonLabel = computed(() => (copied.value ? $ts("copied") : $ts("copy")));

const copy = async () => {
  try {
    await navigator.clipboard.writeText(text);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 1500);
  } catch (error) {
    console.error(error);
  }
};
</script>

<template>
  <button class="btn" @click="copy">{{ buttonLabel }}</button>
</template>

<style lang="less" scoped>
.btn {
  background: var(--bg-card-base);
  padding: 0 1rem;
  color: var(--text-primary);
  border-radius: 5px;
  font-size: 1.2rem;
  font-family: var(--font-family-base);
  border: 1px solid var(--brand-primary);

  &:hover {
    border: 1px solid var(--brand-hover);
    background-color: var(--brand-hover);
  }
}
</style>
