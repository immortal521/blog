<script setup lang="ts">
import { useRouter } from "vue-router";
import { onClickOutside, useMediaQuery } from "@vueuse/core";
const { t } = useI18n();
const localePath = useLocalePath();

const navbarItems = computed<MenuItem[]>(() => {
  return [
    {
      icon: "typcn:home",
      label: t("navbar.home"),
      to: localePath("index"),
    },
    {
      icon: "fluent:apps-list-detail-24-filled",
      label: t("navbar.about"),
      to: localePath("about"),
    },
    {
      icon: "ri:link",
      label: t("navbar.links"),
      to: localePath("links"),
    },
  ];
});

const router = useRouter();

const handleLogoClicked = () => {
  const currentPath = router.currentRoute.value.fullPath;
  const targetPath = "/"; // 您要导航到的目标路径

  if (targetPath !== currentPath) {
    console.log(localePath(targetPath));
    navigateTo(localePath(targetPath));
  }
};

const dropdownOpen = ref(false);

const isMobile = useMediaQuery("(max-width: 768px)");

onMounted(() => {
  watch(
    dropdownOpen,
    () => {
      if (dropdownOpen.value) {
        document.body.style.overflow = "hidden";
      } else {
        document.body.style.overflow = "";
      }
    },
    {
      immediate: true,
    },
  );
  watch(isMobile, () => {
    if (!isMobile.value) {
      document.body.style.overflow = "";
    }
  });
});

const menuBtn = useTemplateRef("menuBtn");
const dropdown = useTemplateRef("dropdown");

onClickOutside(dropdown, (event) => {
  const target = event.target as Node | null;
  if (menuBtn.value?.contains(target)) return;

  dropdownOpen.value = false;
});
</script>

<template>
  <nav id="navbar" class="navbar">
    <div class="logo" @click="handleLogoClicked">
      <img src="http://q1.qlogo.cn/g?b=qq&nk=188191770&s=100" alt="logo" />
    </div>
    <button v-if="isMobile" ref="menuBtn" class="menu-btn" @click="dropdownOpen = !dropdownOpen">
      <Icon :name="dropdownOpen ? 'material-symbols:close' : 'hugeicons:menu-11'" size="28" />
    </button>
    <Teleport to="body">
      <Transition name="dropdown">
        <div v-if="isMobile && dropdownOpen" id="dropdown" ref="dropdown" class="dropdown">
          <NavMenu :items="navbarItems" :vertical="isMobile" @select="dropdownOpen = false" />
        </div>
      </Transition>
    </Teleport>
    <NavMenu v-if="!isMobile" :items="navbarItems" />
  </nav>
</template>

<style lang="less" scoped>
.navbar {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.logo {
  height: 45px;
  width: 45px;
  display: flex;
  position: relative;
  align-items: center;
  justify-content: center;
  user-select: none;
  cursor: pointer;
  border-radius: 50%;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.menu-btn {
  position: relative;
  right: -15px;
  background-color: transparent;
  border: none;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 50px;
  color: var(--text-color-base);
}

.dropdown {
  position: fixed;
  top: 60px;
  right: 0;
  width: 100vw;
  height: calc(100vh - 60px);
  background-color: var(--bg-nav-active);
  border-top: 2px solid var(--border-color-nav);
  box-shadow: var(--shadow-md);
  backdrop-filter: blur(10px);
  padding-top: 20px;
  z-index: 9999;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition:
    opacity 0.3s ease-in-out,
    transform 0.3s ease-in-out;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
