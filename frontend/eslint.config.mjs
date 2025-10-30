// @ts-check
import withNuxt from "./.nuxt/eslint.config.mjs";

export default withNuxt({
  rules: {
    "vue/html-self-closing": "off",
    "vue/multi-word-component-names": "off",
    "@typescript-eslint/no-empty-object-type": [
      "error",
      {
        allowInterfaces: "always",
      },
    ],
  },
});
