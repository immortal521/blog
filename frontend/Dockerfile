# =========================
# 构建阶段 (Builder)
# =========================
FROM node:24 AS builder

# 设置工作目录
WORKDIR /app

# 安装 pnpm
RUN npm install -g pnpm
ENV PNPM_STORE_PATH=/tmp/.pnpm-store

# 复制依赖文件并安装依赖
COPY package*.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# 复制源码
COPY . .

# 修复可执行权限（部分工具需要）
RUN chmod +x /app/node_modules/.bin/*

# 构建 Nuxt SSR
ENV NUXT_TELEMETRY_DISABLED=1

RUN pnpm run build

# =========================
# 运行阶段 (Runtime)
# =========================
FROM node:24-alpine

WORKDIR /app

# 复制构建产物和依赖文件
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package.json /app/
COPY --from=builder /app/pnpm-lock.yaml /app/

# 安装 pnpm 并仅安装生产依赖
RUN npm install -g pnpm
ENV PNPM_STORE_PATH=/tmp/.pnpm-store
RUN pnpm install --prod --frozen-lockfile --ignore-scripts

# 暴露 Nuxt SSR 默认端口
EXPOSE 3000

# 启动 Nuxt SSR
CMD ["node", ".output/server/index.mjs"]
