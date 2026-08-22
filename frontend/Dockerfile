FROM node:24 AS builder

WORKDIR /app

ENV http_proxy=${HTTPS_PROXY} https_proxy=${HTTPS_PROXY} HTTP_PROXY=${HTTPS_PROXY} HTTPS_PROXY=${HTTPS_PROXY}

RUN npm install -g pnpm
ENV PNPM_STORE_PATH=/tmp/.pnpm-store

COPY package*.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY . .

RUN chmod +x /app/node_modules/.bin/*

ENV NUXT_TELEMETRY_DISABLED=1

RUN pnpm run build

FROM node:24-alpine

WORKDIR /app

ENV NUXT_TELEMETRY_DISABLED=1

COPY --from=builder /app/.output /app/.output

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
