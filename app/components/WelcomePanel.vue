<script setup lang="ts">
const handleDownArrowClick = () => {
  const innerHeight = window.innerHeight;
  const scrollableHeight = document.documentElement.scrollHeight - innerHeight;

  if (scrollableHeight > innerHeight) {
    window.scrollTo({
      top: innerHeight,
      behavior: "smooth",
    });
  } else {
    window.scrollTo({
      top: scrollableHeight,
      behavior: "smooth",
    });
  }
};
</script>

<template>
  <div class="welcome">
    <div class="content">
      <div class="introduction">
        <div class="avatar">
          <NuxtImg
            src="https://q1.qlogo.cn/g?b=qq&nk=188191770&s=640"
            alt="avatar"
            class="avatar-image"
          />
        </div>
        <p class="greetings">Cherish every moment.</p>
      </div>
    </div>
    <div class="down-arrow" @click="handleDownArrowClick">
      <Icon name="mingcute:down-line" size="50" />
    </div>
  </div>
</template>

<style lang="less" scoped>
.welcome {
  width: 100%;
  height: 100vh;
  position: relative;
}

.introduction {
  height: calc(100vh - 60px);
  max-width: 600px;
  width: 95%;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-bottom: 20px;
  text-align: center;
}

.avatar {
  height: 130px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fade-in-down 1s ease-in-out;

  .avatar-image {
    width: 100px;
    height: 100px;
    border: 3px solid #999;
    border-radius: 50%;
    object-fit: cover;
    transform: rotate(0deg);
    transition: transform 0.8s ease-in-out;

    &:hover {
      transform: rotate(360deg);
    }
  }
}

.greetings {
  height: 60px;
  width: 100%;
  background-color: var(--bg-nav-base);
  color: var(--text-color-base);
  border: 1px solid var(--border-color-nav);
  border-radius: 10px;
  line-height: 60px;
  margin: 25px 0;
  font-family: var(--font-family-caveat);
  font-size: 2.5rem;
  font-weight: 800;
  animation: fade-in-down 1s ease-in-out;
  box-shadow: var(--shadow-md);
}

.contect {
  display: flex;
  max-width: 210px;
  list-style: none;
  width: max-content;
  gap: 8px; /* 设置格子之间的空隙 */
  justify-content: start; /* 整个grid水平居中 */
  animation: fade-in-down 1s ease-in-out;
  flex-wrap: wrap;
}

.down-arrow {
  position: absolute;
  bottom: 20px;
  left: 50%;
  cursor: pointer;
  transform: translateX(-50%);
  color: var(--bg-nav-base);

  animation: bounce 2.5s infinite ease-in-out;
}

@keyframes bounce {
  0%,
  100% {
    transform: translate(-50%, 0);
  }
  50% {
    transform: translate(-50%, 10px);
  }
}
</style>
