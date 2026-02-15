import type { SidebarNode } from "~/components/BaseSidebar/types";

export function useAdminMenu() {
  const { $ts, $localePath } = useI18n();

  const menuItems = computed<SidebarNode[]>(() => [
    {
      type: "link",
      icon: "duo-icons:dashboard",
      label: $ts("admin.sidebar.dashboard"),
      to: "/admin",
      key: "/admin",
      exact: true,
    },
    {
      type: "section",
      key: "content",
      label: $ts("admin.sidebar.content"),
      items: [
        {
          type: "group",
          icon: "material-symbols:post-rounded",
          label: $ts("admin.sidebar.posts"),
          key: "/admin/posts",
          children: [
            {
              type: "link",
              icon: "mdi:format-list-bulleted",
              label: $ts("admin.sidebar.postsList"),
              to: "/admin/posts",
              key: "/admin/posts",
              exact: true,
            },
            {
              type: "link",
              icon: "mdi:plus-box-outline",
              label: $ts("admin.sidebar.postCreate"),
              to: "/admin/posts/edit",
              key: "/admin/posts/edit",
            },
            {
              type: "link",
              icon: "mdi:trash-can-outline",
              label: $ts("admin.sidebar.trash"),
              to: "/admin/posts/trash",
              key: "/admin/posts/trash",
            },
          ],
        },
        {
          type: "link",
          icon: "ri:link",
          label: $ts("admin.sidebar.links"),
          to: "/admin/links",
          key: "/admin/links",
        },
        {
          type: "link",
          icon: "mdi:folder-outline",
          label: $ts("admin.sidebar.categories"),
          to: "/admin/categories",
          key: "/admin/categories",
        },
        {
          type: "link",
          icon: "mdi:tag-outline",
          label: $ts("admin.sidebar.tags"),
          to: "/admin/tags",
          key: "/admin/tags",
        },
      ],
    },
    {
      type: "section",
      key: "settings",
      label: $ts("admin.sidebar.settings"),
      items: [
        {
          type: "group",
          icon: "mdi:account-cog-outline",
          label: $ts("admin.sidebar.profile"),
          key: "/admin/profile",
          children: [
            {
              type: "link",
              label: $ts("admin.sidebar.profile"),
              to: "/admin/profile",
              key: "/admin/profile",
              exact: true,
            },
          ],
        },
        {
          type: "group",
          icon: "mdi:cog-outline",
          label: $ts("admin.sidebar.system"),
          key: "/admin/settings",
          children: [
            {
              type: "link",
              label: $ts("admin.sidebar.system"),
              to: "/admin/settings",
              key: "/admin/settings",
              exact: true,
            },
          ],
        },
      ],
    },
    {
      type: "action",
      label: $ts("admin.sidebar.back"),
      icon: "material-symbols:text-select-move-back-word-rounded",
      key: "admin.sidebar.back",
      action: () => {
        navigateTo($localePath("/"));
      },
    },
  ]);

  return { menuItems };
}
