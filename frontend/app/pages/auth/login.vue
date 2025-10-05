<script setup lang="ts">
definePageMeta({
  layout: "auth-layout",
});
</script>

<template>
  <div class="main">
    <div class="login-form-container">
      <h2>欢迎回来</h2>
      <form class="login-form" @submit.prevent="">
        <div class="input-box">
          <label for="email">邮箱</label>
          <input id="email" type="text" placeholder="邮箱" />
        </div>
        <div class="input-box">
          <label for="password">密码</label>
          <input id="password" type="password" placeholder="密码" />
        </div>
        <button class="login-button" type="submit">登录</button>
      </form>
    </div>
  </div>
</template>

<style lang="less" scoped>
.main {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.login-form-container {
  min-width: 350px;
  max-width: 90vw;
  padding: 20px;
  background-color: var(--bg-card-base);
  box-shadow: var(--shadow-card-base);
  border-radius: 10px;
  color: var(--text-color-base);
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
}

.input-box {
  display: flex;
  margin: 10px 0;
  align-items: center;

  label {
    width: 6rem;
    font-size: 1.75rem;
  }

  input {
    display: inline-block;
    width: 100%;
    height: 40px;
    padding: 10px 12px;
    margin-left: 1rem;
    color: var(--text-color-base);
    outline: none;
    font-weight: 700;
    border-radius: 5px;
    border: 1px solid var(--border-color-card);
    background-color: var(--bg-card-base);

    &:focus {
      border: 1px solid var(--color-primary-base);
    }
  }
}

.login-button {
  height: 40px;
  margin-top: 20px;
  border-radius: 10px;
  color: var(--color-on-primary);
  font-size: 1.75rem;
  background-color: var(--color-primary-base);

  &:hover {
    background-color: var(--color-primary-hover);
  }

  &:active {
    background-color: var(--color-primary-active);
  }
}
</style>
