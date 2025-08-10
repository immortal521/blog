export default defineNuxtRouteMiddleware((to, from) => {
	const { start, done } = useLoadingBar();

	if (from.name && to.fullPath === from.fullPath) {
		return;
	}

	// 在路由开始导航时启动加载条
	start();

	if (import.meta.client) {
		useRuntimeHook("page:finish", () => {
			setTimeout(() => {
				done();
			}, 500);
		});
	}

	//TODO: 路由加载出错时，也可以调用 `done()`
});
