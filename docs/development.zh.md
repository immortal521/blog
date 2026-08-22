[English](development.md) · [中文](development.zh.md)

# 本地开发环境

工具版本：Node `24`（见 `frontend/Dockerfile` 的 `node:24` 基础镜像）、pnpm `11.15.1`（见 `package.json` 的 `packageManager`）；后端见 `backend/go.mod`：`go 1.27`。

## 后端

在根目录运行：

```bash
pnpm dev:backend
```

或在 `backend` 目录运行：

```bash
go run ./cmd
```

后端依赖以下服务：

- redis
- postgres（需存在名为 `blog` 的数据库）

后端程序会自动建表。

## 前端

先在 `frontend` 目录安装依赖：

```bash
pnpm i --frozen-lockfile
```

在根目录运行：

```bash
pnpm dev:frontend
```

或在 `frontend` 目录运行：

```bash
pnpm dev
```
