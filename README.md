# Immortal's Blog

## 部署
暂时无法一键部署

## 开发环境
### 后端
在理论上根目录运行
```bash
pnpm dev:backend
```
或在`backend`目录运行
```bash
go run ./cmd
```

后端需要一些服务：
- redis
- mysql

mysql 需要拥有名为 `blog` 的数据库

后端程序会自动建表

### 前端
需要先在`frontend`目录安装依赖
```bash
pnpm i --frozen-lockfile
```
在根目录运行
```bash
pnpm dev:frontend
```
或在`frontend`目录运行
```bash
pnpm dev
```


