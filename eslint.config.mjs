import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const subProjectPath = path.resolve(__dirname, "frontend/eslint.config.mjs");

// 转成 file:// URL
const subProjectURL = pathToFileURL(subProjectPath).href;

const config = await import(subProjectURL);
export default config.default;
