[English](ci.md) · [中文](ci.zh.md)

# 通过 GitHub Actions 自动部署

适用于「服务器上不保留源码仓库、由 CI 构建镜像并下发」的场景。与 [容器一键部署](deployment.zh.md) 互不冲突，可二选一。

整体流程：

- CI（`.github/workflows/deploy_backend.yml` / `deploy_frontend.yml`）只负责构建镜像，不把密钥写进镜像或传输包。
- 镜像通过 `docker save` 导出为 tar，连同一份不含 secrets 的部署配置包（`compose.prod.yml`、nginx 模板、后端配置模板、`.env.example`）一起 SCP 到服务器。
- 服务器 `podman load` 载入镜像；`.env` 由 `setup-env` action 通过 SSH 从加密的 GitHub Secret（`SERVER_ENV`，完整 `.env` 内容）直接写到服务器，密钥只在 SSH 连接中落盘，不进仓库、不进传输包。
- 之后用服务器本机 `.env` 生成运行期配置，再由服务器侧脚本 `deploy/remote-deploy.sh` 拉起对应服务。该脚本与 `deploy.sh` 共享同一套逻辑：rootless Podman 下 `80/443` 自动回落到 `8080/8443`、`envsubst` 生成后端 / Nginx 配置、`nginx -s reload` 重新接入。工作流 SCP 完镜像与配置包后，通过 SSH 调 `bash deploy/remote-deploy.sh backend|frontend` 即可，不再在 workflow 内联重复部署逻辑。

## 服务器前置要求

- 已安装 Podman（rootless，含 `podman compose`）
- `envsubst`（gettext）
- 在仓库 `Settings → Secrets` 中配置 `SERVER_ENV`：值为服务器所需的完整 `.env` 内容（可直接复制 `.env.example` 后填入真实密钥）。`setup-env` action 仅在服务器上 `~/projects/blog/.env` 不存在或为空时通过 SSH 把它写到服务器（已存在则保留，不会被覆盖）。
  - 若不想用 GitHub Secret，也可在服务器上手动 `cp .env.example .env` 并填写；之后 action 会跳过创建、保留你的版本。

其余编排文件（`compose.prod.yml`、nginx 模板、后端配置模板）由 CI 每次下发，无需手动放置，也不含任何密钥。

## 工作流说明

| 工作流 | 触发路径 | 行为 |
| --- | --- | --- |
| `deploy_backend.yml` | `backend/**`、`config/backend_config.yml` | 构建 `blog_backend:latest` → SCP → 服务器 `podman load` → `bash deploy/remote-deploy.sh backend`（重建 `migration` + `backend`） |
| `deploy_frontend.yml` | `frontend/**` | 构建 `blog_frontend:latest` → SCP → 服务器 `podman load` → `bash deploy/remote-deploy.sh frontend`（重建 `frontend`） |

两个工作流相互独立（不同的 `concurrency` 分组），互不阻塞；各自只更新自己的服务，不会牵连另一个。

服务器侧 `deploy/remote-deploy.sh` 的行为要点：

- **Rootless 端口回退**：与 `deploy.sh` 一致，检测到 rootless Podman 时把 `80/443` 改为 `8080/8443`，并写入 `deploy/runtime/compose.ports.yml` 覆盖文件（用 `-f` 挂上，里面是字面量端口）。之所以用覆盖文件而不是直接 `export` 环境变量，是因为 podman-compose 会优先读取项目目录下的 `.env`（`HTTP_PORT=80`）而忽略 shell 导出的回退值——只有显式 `-f` 覆盖文件里的字面量才生效。
- **迁移顺序由脚本显式保证**：`remote-deploy.sh backend` 先 `up -d migration` 并轮询等待其退出码 0，再 `up -d backend`。不使用 compose 的 `depends_on: service_completed_successfully` 条件——该条件在服务器上的 podman-compose 1.0.6 下会把 backend 创建出来却留在其 `Created` 状态不启动。
- **显式重建**：先 `podman rm -f` 目标容器再 `up -d`，规避 `podman compose --force-recreate` 在依赖容器存在时的删除顺序 bug；`frontend` 已不 `depends_on backend`，故 backend 重建不会受 frontend 容器阻塞。
- **Nginx 幂等**：每次都 `up -d nginx` 并 `nginx -s reload`，容器重建后自动重新解析新实例。

## 增量与首次部署

- **增量**：backend 提交只重建 backend（并强制重跑幂等的 `migration`），frontend 提交只重建 frontend；nginx 配置带动态 resolver，重建后自动重新解析新实例，无需手动 reload（脚本也会顺手 reload 一次）。
- **首次部署**：请先在服务器手动拉起整套服务一次，之后每次 push 即为增量：

  ```bash
  cd ~/projects/blog
  bash deploy/remote-deploy.sh all
  ```

  也可直接推送两个分支的改动各触发一次工作流。

## 两种脚本的分工

- **`deploy.sh`（面向用户）**：一键部署，用户可自行决定三种方式——
  - rootless 容器（普通用户 Podman，自动回落 8080/8443）：`./deploy.sh`
  - rootful 容器（`sudo` 运行或 docker root daemon，直接绑 80/443）：`sudo ./deploy.sh`
  - 自行部署（打包 + 步骤，不自动运行）：`./deploy.sh self`
- **`deploy/remote-deploy.sh`（仓库的远程 CI）**：rootless 场景的服务器侧脚本，配合 `.github/workflows/` 用预构建镜像拉起服务；用户 fork 后亦可原样在自己的 GitHub Actions 中复用。端口回落逻辑与 `deploy.sh` 的 rootless 分支一致。

两者共用同一套配置生成逻辑与 `deploy/runtime/` 生成目录。
