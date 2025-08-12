export default defineNuxtRouteMiddleware(() => {
	const headers = useRequestHeaders();
	const ua = headers["user-agent"] || "";

	// SSR 阶段初步判断
	const isMobileSSR = /Mobile|Android|iP(hone|od)|IEMobile|BlackBerry/i.test(
		ua,
	);

	useState("device", () => ({
		isMobile: isMobileSSR,
		width: isMobileSSR ? 768 : 1080,
	}));
});
