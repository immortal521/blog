[English](ci.md) · [中文](ci.zh.md)

# Auto-deploy via GitHub Actions

Suitable for the scenario where "the server keeps no source repo; CI builds the images and pushes them down." It does not conflict with the [container one-command deploy](deployment.md) — pick either one.

The overall flow:

- CI (`.github/workflows/deploy_backend.yml` / `deploy_frontend.yml`) only builds images; it does not write secrets into the image or the transfer package.
- The images are exported to a tar via `docker save`, and together with a deployment config package that contains no secrets (`compose.prod.yml`, nginx templates, backend config template, `.env.example`) are SCP'd to the server.
- The server `podman load`s the images; `.env` is written to the server by the `setup-env` action via SSH from an encrypted GitHub Secret (`SERVER_ENV`, the full `.env` content). The secret only lands on disk inside the SSH connection — it does not enter the repo nor the transfer package.
- Then the server-local `.env` is used to generate runtime config, and the server-side script `deploy/remote-deploy.sh` brings up the corresponding services. This script shares the same logic with `deploy.sh`: under rootless Podman, `80/443` automatically fall back to `8080/8443`; `envsubst` generates the backend / Nginx config; `nginx -s reload` reconnects. After the workflow SCPs the images and config package, it just calls `bash deploy/remote-deploy.sh backend|frontend` over SSH — it does not repeat the deploy logic inline in the workflow.

## Server prerequisites

- Podman installed (rootless, including `podman compose`)
- `envsubst` (gettext)
- Configure `SERVER_ENV` in the repo's `Settings → Secrets`: the value is the full `.env` content the server needs (you can copy `.env.example` and fill in the real secrets). The `setup-env` action only writes it to the server via SSH when `~/projects/blog/.env` does not exist or is empty on the server (if it already exists, it is preserved and not overwritten).
  - If you do not want to use a GitHub Secret, you can also manually `cp .env.example .env` on the server and fill it in; afterwards the action skips creation and keeps your version.

The remaining orchestration files (`compose.prod.yml`, nginx templates, backend config template) are pushed by CI every time — no manual placement needed, and they contain no secrets.

## Workflow reference

| Workflow | Trigger path | Behavior |
| --- | --- | --- |
| `deploy_backend.yml` | `backend/**`, `config/backend_config.yml` | Build `blog_backend:latest` → SCP → server `podman load` → `bash deploy/remote-deploy.sh backend` (rebuild `migration` + `backend`) |
| `deploy_frontend.yml` | `frontend/**` | Build `blog_frontend:latest` → SCP → server `podman load` → `bash deploy/remote-deploy.sh frontend` (rebuild `frontend`) |

The two workflows are independent (different `concurrency` groups) and do not block each other; each only updates its own service and does not affect the other.

Key behaviors of the server-side `deploy/remote-deploy.sh`:

- **Rootless port fallback**: consistent with `deploy.sh` — when rootless Podman is detected, change `80/443` to `8080/8443`, and write a `deploy/runtime/compose.ports.yml` override file (mounted with `-f`, containing literal ports). The reason for an override file rather than directly `export`ing an environment variable is that podman-compose prioritizes the project-directory `.env` (`HTTP_PORT=80`) over shell-exported fallback values — only the literal values in an explicit `-f` override file take effect.
- **Migration order explicitly guaranteed by the script**: `remote-deploy.sh backend` first `up -d migration` and polls until its exit code is 0, then `up -d backend`. It does not use compose's `depends_on: service_completed_successfully` condition — under podman-compose 1.0.6 on the server, that condition creates `backend` but leaves it in its `Created` state without starting.
- **Explicit rebuild**: first `podman rm -f` the target container, then `up -d`, working around a podman-compose `--force-recreate` deletion-ordering bug when dependency containers exist; `frontend` no longer `depends_on backend`, so a backend rebuild is not blocked by the frontend container.
- **Idempotent Nginx**: every time `up -d nginx` and `nginx -s reload`; after a container rebuild it automatically re-resolves the new instance.

## Incremental vs. first deploy

- **Incremental**: a backend commit only rebuilds backend (and force-reruns the idempotent `migration`); a frontend commit only rebuilds frontend; the nginx config has a dynamic resolver, so after a rebuild it automatically re-resolves the new instance without a manual reload (the script also reloads once anyway).
- **First deploy**: first manually bring up the whole stack once on the server; afterwards every push is incremental:

  ```bash
  cd ~/projects/blog
  bash deploy/remote-deploy.sh all
  ```

  You can also just push changes to both branches, each triggering a workflow once.

## Division of labor between the two scripts

- **`deploy.sh` (user-facing)**: one-command deploy, where the user chooses one of three modes —
  - rootless container (plain-user Podman, auto fallback 8080/8443): `./deploy.sh`
  - rootful container (`sudo` or Docker root daemon, binds 80/443 directly): `sudo ./deploy.sh`
  - self-hosting (package + steps, does not auto-run): `./deploy.sh self`
- **`deploy/remote-deploy.sh` (the repo's remote CI)**: the server-side script for the rootless scenario, used with `.github/workflows/` to bring up services with prebuilt images; after a user forks, they can reuse it as-is in their own GitHub Actions. The port-fallback logic matches the rootless branch of `deploy.sh`.

Both share the same config-generation logic and the `deploy/runtime/` generation directory.
