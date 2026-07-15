<script setup lang="ts" generic="T extends Record<string, any>">
import type { VNode } from "vue";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface Column<RowType extends Record<string, any>> {
  key: keyof RowType | string;
  title: string;

  width?: string;

  align?: "left" | "center" | "right";

  formatter?: (row: RowType) => string | number;

  render?: (row: RowType, index: number) => VNode;
}

const {
  data,
  columns,
  rowKey = () => "id",
  virtual = false,
  height = 500,
  rowHeight = 48,
  overscan = 5,
} = defineProps<{
  data: T[];

  columns: Column<T>[];

  rowKey?: keyof T | ((row: T) => string | number);

  virtual?: boolean;

  height?: number;

  rowHeight?: number;

  overscan?: number;
}>();

const getRowKey = (row: T, index: number): string | number => {
  if (typeof rowKey === "function") {
    return rowKey(row);
  }

  return (row[rowKey] as string | number) ?? index;
};

function renderCell(column: Column<T>, row: T, index: number) {
  if (column.render) {
    return column.render(row, index);
  }

  if (column.formatter) {
    return column.formatter(row);
  }

  return row[column.key];
}

const scrollTop = ref(0);

const startIndex = computed(() => {
  if (!virtual) {
    return 0;
  }

  return Math.max(0, Math.floor(scrollTop.value / rowHeight) - overscan);
});

const endIndex = computed(() => {
  if (!virtual) {
    return data.length;
  }

  return Math.min(data.length, Math.ceil((scrollTop.value + height) / rowHeight) + overscan);
});

const visibleRows = computed(() => data.slice(startIndex.value, endIndex.value));

const offsetTop = computed(() => startIndex.value * rowHeight);

const totalHeight = computed(() => data.length * rowHeight);

function onScroll(e: Event) {
  scrollTop.value = (e.target as HTMLElement).scrollTop;
}
</script>

<template>
  <div class="my-table">
    <!-- header -->

    <div class="table-header">
      <div class="table-row">
        <div
          v-for="column in columns"
          :key="column.key"
          class="table-cell"
          :class="`align-${column.align ?? 'left'}`"
          :style="{
            width: column.width,
            flex: column.width ? 'none' : 1,
          }"
        >
          {{ column.title }}
        </div>
      </div>
    </div>

    <!-- body -->

    <div
      class="table-body"
      :style="{
        height: virtual ? `${height}px` : undefined,
        overflow: virtual ? 'auto' : undefined,
      }"
      @scroll="onScroll"
    >
      <!-- virtual -->

      <template v-if="virtual">
        <div
          :style="{
            height: `${totalHeight}px`,
            position: 'relative',
          }"
        >
          <div
            :style="{
              transform: `translateY(${offsetTop}px)`,
            }"
          >
            <div
              v-for="(row, index) in visibleRows"
              :key="getRowKey(row, startIndex + index)"
              class="table-row"
              :style="{
                height: `${rowHeight}px`,
              }"
            >
              <div
                v-for="column in columns"
                :key="column.key"
                class="table-cell"
                :class="`align-${column.align ?? 'left'}`"
                :style="{
                  width: column.width,
                  flex: column.width ? 'none' : 1,
                }"
              >
                <component :is="renderCell(column, row, startIndex + index)" v-if="column.render" />

                <template v-else>
                  {{ renderCell(column, row, startIndex + index) }}
                </template>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- normal -->

      <template v-else>
        <div v-for="(row, index) in data" :key="getRowKey(row, index)" class="table-row">
          <div
            v-for="column in columns"
            :key="column.key"
            class="table-cell"
            :class="`align-${column.align ?? 'left'}`"
            :style="{
              width: column.width,
              flex: column.width ? 'none' : 1,
            }"
          >
            <component :is="renderCell(column, row, index)" v-if="column.render" />

            <template v-else>
              {{ renderCell(column, row, index) }}
            </template>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style lang="less" scoped>
.my-table {
  width: 100%;
  border-collapse: collapse;

  .table-row {
    display: flex;
    width: 100%;
    align-items: center;
    border-bottom: 1px solid #f0f0f0;

    &:hover {
      background-color: #fafafa;
    }
  }

  .table-header .table-row {
    background-color: #fafafa;
    font-weight: 600;
  }

  .table-cell {
    padding: 12px 16px;
    box-sizing: border-box;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap; /* 视情况决定是否开启单行截断 */

    // 真正起作用的对齐样式
    &.align-left {
      text-align: left;
      justify-content: flex-start;
    }
    &.align-center {
      text-align: center;
      justify-content: center;
    }
    &.align-right {
      text-align: right;
      justify-content: flex-end;
    }
  }
}
</style>
