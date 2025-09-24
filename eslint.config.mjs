import path from "node:path";
import { fileURLToPath } from "node:url";

// 解析子项目路径
const __dirname = path.dirname(fileURLToPath(import.meta.url));
const subProject = path.resolve(__dirname, "frontend/eslint.config.mjs");

// 直接 re-export 子项目的配置
const config = await import(subProject);
export default config.default;
