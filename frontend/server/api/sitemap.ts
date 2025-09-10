export default defineSitemapEventHandler(async () => {
  try {
    // 调用后端 API 获取 posts 的 ids
    const { data } = await $fetch<{ data: { ids: number[] } }>("/api/v1/posts/ids");
    const ids = data?.ids ?? [];

    const routeList = ids.map((id) => ({
      loc: `/blog/${id}`,
      _i18nTransform: true,
      lastmod: new Date().toISOString().split("T")[0], // 当前日期，或可自定义
    }));

    return routeList;
  } catch (error) {
    console.error("获取sitemap失败", error);
    return [];
  }
});
