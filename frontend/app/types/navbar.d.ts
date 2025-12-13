interface MenuItem extends SidebarItem {
  to: string;
  children?: MenuItem;
}

interface SidebarItem {
  key: string;
  label: string;
  icon: IconName | Component;
}
