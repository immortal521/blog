// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: "2025-07-15",
	devtools: { enabled: true },

	app: {
		head: {
			title: "Immortal's Blog",
			htmlAttrs: {
				lang: "en",
			},
			charset: "utf-8",
			viewport: "width=device-width, initial-scale=1, maximum-scale=1",
			link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
		},
	},

	typescript: {
		typeCheck: true,
	},

	nitro: {
		routeRules: {
			"/api/**": { proxy: "http://localhost:8000/api/**" },
		},
	},

	vite: {
		css: {
			devSourcemap: true, // 开启 CSS Source Map，可能有助于调试
		},
	},

	css: ["~/assets/styles/main.less"],

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
			redirectOn: "root", // 可选：只在根路径重定向
		},
	},

	modules: [
		// "@nuxt/fonts",
		"@nuxt/eslint",
		"@nuxt/icon",
		"@nuxt/image",
		"@nuxt/scripts",
		"@pinia/nuxt",
		"@nuxtjs/i18n",
	],
});
