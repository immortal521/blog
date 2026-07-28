interface BaseItem {
  id: string;
  icon: string;
  label: string;
}

interface FloatButton extends BaseItem {
  type: "button";
  onClick: () => void;
}

interface FloatPanel extends BaseItem {
  type: "panel";
  component: ShowModelComponent;
}

type ShowModelComponent = Component & {
  new (): {
    $props: {
      readonly show: boolean;
      readonly "onUpdate:show"?: (value: boolean) => void;
    };
  };
};

export type FloatActionItem = FloatButton | FloatPanel;
