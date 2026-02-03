interface MenuItem extends SidebarItem {
  key: string;
  label: string;
  icon: IconName | Component;
  to: string;
  children?: MenuItem;
}
