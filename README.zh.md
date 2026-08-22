[English](README.md) · [中文](README.zh.md)

# Immortal's Blog

一个使用 Go（后端）+ Nuxt（前端）的全栈博客，配套 PostgreSQL / Redis / RustFS（S3 兼容对象存储），
通过 Nginx 统一对外提供服务。

整套服务收敛在一个 `compose.yml` 中，clone 后准备一份 `.env` 即可一键启动（Podman 或 Docker 均可）；
也支持完全不使用容器的“自行部署”模式，或交给 GitHub Actions 构建镜像后下发。

## 文档

- [容器部署指南](docs/deployment.zh.md) —— rootless / rootful 自动判定、一键部署、HTTPS / RustFS 反代、证书
- [自行部署指南](docs/self-hosting.zh.md) —— `./deploy.sh self` 打包产物 + 步骤，运行方式由你自行决定
- [GitHub Actions 自动部署](docs/ci.zh.md) —— CI 下发流程与 `remote-deploy.sh` 说明
- [本地开发环境](docs/development.zh.md) —— 后端 `go run`、前端 `pnpm dev`

## 目录结构

```
.
├── compose.yml                 # 所有服务的容器编排（手动 clone 后一键部署用，从源码构建）
├── compose.prod.yml            # CI 自动部署用（预构建镜像，服务器不构建）
├── deploy.sh                   # 一键部署 / 自行部署打包脚本（rootless/rootful 自动判定）
├── .env.example                # 配置样例
├── backend/                    # Go 后端（含 Dockerfile）
├── frontend/                   # Nuxt 前端（含 Dockerfile）
├── config/
│   ├── backend_config.yml      # 后端配置模板（${VAR} 占位符）
│   └── nginx.conf              # 历史参考用的 Nginx 配置
└── deploy/
    ├── remote-deploy.sh        # CI 服务器侧部署脚本，与 deploy.sh 共享端口与配置生成逻辑
    ├── certbot.sh              # certbot 封装，支持 certonly / renew / import；CERT_MODE=project|system
    ├── nginx/                  # 容器模式 Nginx 模板（http / https / rustfs.*）
    ├── letsencrypt/            # 证书目录（容器模式），部署时生成
    ├── acme/                   # ACME HTTP-01 挑战目录，部署时生成
    └── runtime/                # 运行期配置（nginx / backend），部署时生成

docs/                          # 部署与开发文档（见上方链接）
```
