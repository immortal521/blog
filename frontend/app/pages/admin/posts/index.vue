<script setup lang="ts"></script>

<template>
  <div>
    <div class="title">
      <h2 class="content">文章管理</h2>
      <span class="count">{{ `共  文章` }}</span>
      <button class="create">新建文章</button>
    </div>
  </div>
</template>

<style lang="less" scoped>
.title {
  display: flex;
  height: 40px;
  line-height: 40px;

  .count {
    margin-left: 8px;
    height: 30px;
    line-height: 30px;
    font-size: 14px;
    background: var(--glass-gradient), var(--bg-card-base);
    padding: 0 8px;
    border-radius: 8px;
  }

  .create {
    margin-left: auto;
  }
}
</style>
