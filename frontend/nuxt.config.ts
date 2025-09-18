// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",

  devServer: {
    host: "0.0.0.0",
    port: 3000,
  },
  // 启用开发工具和时间线调试
  devtools: {
    enabled: true,
    timeline: {
      enabled: true,
    },
  },

  // 应用信息配置
  app: {
    head: {
      title: "Immortal's Blog",
      htmlAttrs: { lang: "en" },
      charset: "utf-8",
      viewport: "width=device-width, initial-scale=1, maximum-scale=1",
      link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
    },
    pageTransition: { name: "page", mode: "out-in" },
  },

  // 站点基本信息
  site: {
    url: "https://blog.immortel.top",
    name: "Immortal's Blog",
  },

  // Sitemap 配置
  sitemap: {
    sources: ["/api/sitemap"],
    // cacheMaxAgeSeconds: 6 * 60 * 60, // 6小时缓存，可按需开启
    autoLastmod: true, // 自动生成最后修改时间，方便爬虫
  },

  // TypeScript 配置
  typescript: {
    typeCheck: true,
  },

  // Nitro 配置
  nitro: {
    routeRules: {
      "/api/v1/**": { proxy: "http://localhost:8000/api/v1/**" }, // API 代理转发
    },
  },

  fonts: {
    families: [
      { name: "Open Sans", provider: "google" },
      {
        name: "Noto Sans SC",
        provider: "google",
      },
      { name: "Caveat", provider: "google" },
    ],
  },

  // Vite 配置
  vite: {
    css: {
      devSourcemap: true, // 开启 CSS Source Map，方便调试
    },
  },

  // 全局样式
  css: ["~/assets/styles/main.less"],

  // 国际化配置 (i18n)
  i18n: {
    strategy: "prefix_except_default",
    defaultLocale: "zh",
    locales: [
      { code: "en", name: "English", file: "en.json" },
      { code: "zh", name: "简体中文", file: "zh.json" },
      { code: "ja", name: "日本語", file: "ja.json" },
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_redirected",
      redirectOn: "root", // 只在根路径重定向
    },
  },

  image: {},

  // 使用的模块
  modules: [
    "@nuxt/fonts",
    "@nuxt/eslint",
    "@nuxt/icon",
    "@nuxt/scripts",
    "@pinia/nuxt",
    "@nuxtjs/i18n",
    "@nuxtjs/sitemap",
    "@nuxt/image",
    "motion-v/nuxt",
  ],
});
