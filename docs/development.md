[English](development.md) · [中文](development.zh.md)

# Local Development Environment

Tool versions: Node `24` (see the `node:24` base image in `frontend/Dockerfile`), pnpm `11.15.1` (see `package.json`'s `packageManager`); for the backend see `backend/go.mod`: `go 1.27`.

## Backend

Run from the root directory:

```bash
pnpm dev:backend
```

Or from the `backend` directory:

```bash
go run ./cmd
```

The backend depends on the following services:

- redis
- postgres (a database named `blog` must exist)

The backend program creates the tables automatically.

## Frontend

First install dependencies in the `frontend` directory:

```bash
pnpm i --frozen-lockfile
```

Run from the root directory:

```bash
pnpm dev:frontend
```

Or from the `frontend` directory:

```bash
pnpm dev
```
