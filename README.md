[English](README.md) · [中文](README.zh.md)

# Immortal's Blog

A full-stack blog built with Go (backend) + Nuxt (frontend), backed by PostgreSQL / Redis / RustFS (an S3-compatible object store), and served through a single Nginx entry point.

The whole stack is consolidated into one `compose.yml`. After cloning, prepare a single `.env` and you can launch with one command — Podman or Docker both work. A fully container-free "self-hosting" mode is also supported, or you can let GitHub Actions build the images and push them to the server.

## Documentation

- [Container deployment guide](docs/deployment.md) — automatic rootless / rootful detection, one-command deploy, HTTPS / RustFS reverse proxy, certificates
- [Self-hosting guide](docs/self-hosting.md) — `./deploy.sh self` packaging + steps; you decide how to run it
- [GitHub Actions auto-deploy](docs/ci.md) — CI push flow and `remote-deploy.sh` explained
- [Local development](docs/development.md) — backend `go run`, frontend `pnpm dev`

## Directory structure

```
.
├── compose.yml                 # Container orchestration for all services (one-command deploy after manual clone; builds from source)
├── compose.prod.yml            # Used by CI auto-deploy (prebuilt images; the server does not build)
├── deploy.sh                   # One-command / self-hosting packaging script (auto-detects rootless/rootful)
├── .env.example                # Sample configuration
├── backend/                    # Go backend (includes Dockerfile)
├── frontend/                   # Nuxt frontend (includes Dockerfile)
├── config/
│   ├── backend_config.yml      # Backend config template (${VAR} placeholders)
│   └── nginx.conf              # Historical reference Nginx config
└── deploy/
    ├── remote-deploy.sh        # Server-side CI deploy script; shares port and config-generation logic with deploy.sh
    ├── certbot.sh              # certbot wrapper supporting certonly / renew / import; CERT_MODE=project|system
    ├── nginx/                  # Container-mode Nginx templates (http / https / rustfs.*)
    ├── letsencrypt/            # Certificate directory (container mode), generated at deploy time
    ├── acme/                   # ACME HTTP-01 challenge directory, generated at deploy time
    └── runtime/                # Runtime config (nginx / backend), generated at deploy time

docs/                          # Deployment and development docs (see links above)
```
