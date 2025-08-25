export default defineSitemapEventHandler(() => {
	try {
		const blogs = [{ id: 1 }];

		const routeList = [];

		const blogRoutes = blogs.map((blog) => ({
			loc: `/blog/${blog.id}`,
			_i18nTransform: true,
			lastmod: "2023-01-01",
		}));

		routeList.push(...blogRoutes);

		return routeList ?? [];
	} catch (error) {
		console.error("获取sitemap失败", error);
		return [];
	}
});
