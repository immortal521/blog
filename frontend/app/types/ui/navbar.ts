export interface MenuItem {
  key: string;
  label: string;
  icon: string;
  to: string;
  children?: MenuItem;
}
