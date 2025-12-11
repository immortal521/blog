interface MenuItem {
  icon: IconName | Component;
  label: string;
  to: string;
  children?: MenuItem;
}
