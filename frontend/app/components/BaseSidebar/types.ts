export type ItemType = "link" | "group" | "action" | "divider";
export type NodeType = ItemType | "section";

export interface BaseNode<T extends NodeType> {
  type: T;
  key: string;
}

export interface BaseItem<T extends NodeType> extends BaseNode<T> {
  title?: string;
  icon?: string;
  hidden?: boolean;
  disabled?: boolean;
}

export interface SidebarLinkItem extends BaseItem<"link"> {
  to: string;
  exact?: boolean;
  badge?: string;
  label: string;
}

export interface SidebarGroupItem extends BaseItem<"group"> {
  children: SidebarItem[];
  label: string;
}

export interface SidebarActionItem extends BaseItem<"action"> {
  action: () => void;
  label: string;
}

export interface SidebarDividerItem extends BaseItem<"divider"> {
  label?: string;
}

export type SidebarItem =
  | SidebarLinkItem
  | SidebarGroupItem
  | SidebarActionItem
  | SidebarDividerItem;

export type SidebarSection = {
  key: string;
  type: "section";
  title?: string;
  label?: string;
  items: SidebarItem[];
  hidden?: boolean;

  collapsible?: boolean;
  defaultOpen?: boolean;
};

export type SidebarNode = SidebarSection | SidebarItem;

export interface SidebarItemEmits {
  (e: "toggle", key: string): void;
}
