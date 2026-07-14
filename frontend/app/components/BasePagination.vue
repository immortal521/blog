<script setup lang="ts">
interface PaginationProps {
  total: number;
  siblingCount?: number;
  disabled?: boolean;
}

interface PaginationItem {
  key: string;
  type: "page" | "prev" | "next" | "ellipsis";
  label: string;
  page?: number;
  active: boolean;
  disabled: boolean;
}

const { total, siblingCount = 1, disabled = false } = defineProps<PaginationProps>();

const page = defineModel<number>("page", { required: true });
const pageSize = defineModel<number>("pageSize", { required: true });

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(total / pageSize.value));
});

function createPages(): PaginationItem[] {
  const pages: PaginationItem[] = [];

  const current = page.value;
  const total = totalPages.value;

  function addPage(page: number) {
    pages.push({
      key: `page-${page}`,
      type: "page",
      label: String(page),
      page,
      active: page === current,
      disabled: disabled,
    });
  }

  function addEllipsis(key: string) {
    pages.push({
      key,
      type: "ellipsis",
      label: "...",
      active: false,
      disabled: true,
    });
  }
  if (total <= siblingCount * 2 + 5) {
    for (let i = 1; i <= total; i++) {
      addPage(i);
    }
    return pages;
  }

  addPage(1);

  if (current > siblingCount + 2) {
    addEllipsis("left-ellipsis");
  }

  const start = Math.max(2, current - siblingCount);

  const end = Math.min(total - 1, current + siblingCount);

  for (let i = start; i <= end; i++) {
    addPage(i);
  }

  if (current < total - siblingCount - 2) {
    addEllipsis("right-ellipsis");
  }

  addPage(total);

  return pages;
}

const items = computed<PaginationItem[]>(() => {
  const result: PaginationItem[] = [];

  result.push({
    key: "prev",
    type: "prev",
    label: "Previous",
    active: false,
    disabled: disabled || page.value <= 1,
  });

  result.push(...createPages());

  result.push({
    key: "next",
    type: "next",
    label: "Next",
    active: false,
    disabled: disabled || page.value >= totalPages.value,
  });

  return result;
});

function select(item: PaginationItem) {
  if (item.disabled || item.type === "ellipsis") {
    return;
  }

  let p = page.value;

  if (item.type === "page") {
    p = item.page!;
  }

  if (item.type === "prev") {
    p--;
  }

  if (item.type === "next") {
    p++;
  }

  p = Math.min(Math.max(p, 1), totalPages.value);

  if (p === page.value) {
    return;
  }

  page.value = p;
}

defineExpose({
  totalPages,
  items,
  select,
});
</script>

<template>
  <slot :items="items" :page="page" :total-pages="totalPages" :select="select" />
</template>
