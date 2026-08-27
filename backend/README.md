# Blog Server

A blog backend service written in Go. It serves published posts and links, handles
JWT-based user authentication, generates RSS feeds, and streams LLM-generated post
summaries over Server-Sent Events.

The service is built around [Echo v5](https://github.com/labstack/echo) for HTTP,
[Uber Fx](https://github.com/uber-go/fx) for dependency injection,
[Ent](https://entgo.io/) as the ORM on PostgreSQL, and Redis for caching and
counters. Scheduled background jobs handle link health checks and view-count
persistence.

## Features

- **Authentication & authorization**
  - Email + password registration and login with JWT access/refresh tokens.
  - Refresh token is returned as an `HttpOnly` cookie; access token is returned in the JSON body.
  - Captcha codes sent by email (SMTP) before registration.
  - RBAC with two roles (`admin`, `reader`). Only `admin` can create/update/delete posts and links; `reader` can only read.
- **Posts**
  - Public list / detail / metadata endpoints.
  - Admin CRUD with tags and categories (many-to-many).
  - View counts are incremented in Redis on read and flushed to PostgreSQL hourly.
  - `read_time_minutes` is derived automatically from content length.
- **Links (友链)**
  - Public list of enabled links.
  - Public "apply for a link" endpoint that creates a disabled link pending review.
  - Link availability is checked hourly by a scheduler job (HTTPS `GET`, status `normal`/`abnormal`).
- **RSS**
  - Paginated RSS 2.0 feed with `atom:link` (self/next/prev/first/last) navigation.
  - Complete archive feed (all published posts).
- **LLM summarization**
  - `POST /api/v1/model/summarize` opens a session; `GET /api/v1/model/summarize/:sessionId` streams the summary as SSE (`data:` lines, `event: done`, `event: error`).
  - Calls GitHub Models / Azure AI inference (`models.inference.ai.azure.com`) with the configured API key.
- **Infrastructure**
  - Structured Zap logging (console in dev, JSON in production).
  - TOML configuration with environment-variable overrides and hot-reload (file change triggers a graceful restart).
  - Graceful HTTP shutdown with a configurable timeout.

## Tech Stack

| Concern            | Library / Tool                                   |
| ------------------ | ----------------------------------------------- |
| Web framework      | Echo v5                                         |
| Dependency inject. | Uber Fx                                         |
| ORM                | Ent v0.14                                       |
| Database           | PostgreSQL (`lib/pq`)                            |
| Cache / counters   | Redis (`redis/go-redis`)                        |
| Object storage     | AWS S3 SDK v2 (S3-compatible; implementation present but not wired into the running server — see below) |
| Logging            | Uber Zap                                        |
| Config             | Viper (TOML) + `fsnotify` hot-reload            |
| JWT                | `golang-jwt/jwt/v5`                             |
| Validation         | `go-playground/validator/v10`                   |
| Email              | `gopkg.in/gomail.v2`                            |

## Requirements

- Go 1.27+ (declared `go 1.27.0` in `go.mod`)
- PostgreSQL
- Redis

## Configuration

Configuration is loaded from a TOML file (default `config.toml` in the working
directory) via Viper. Values can be overridden by environment variables using
the `APP` prefix with `.` replaced by `_` (e.g. `APP_SERVER_PORT`, `APP_JWT_SECRET`,
`APP_DATABASE_PASSWORD`). In production (`app.environment = "production"`) a JWT
secret and a database password are required.

There is **no committed example file** (the repo's `.gitignore` excludes
`config.toml`). Copy the values below into a `config.toml` at the project root.

```toml
[app]
name = "my-app"
version = "1.0.1"
environment = "development"      # "development" | "production"
debug = true
domain = "http://192.168.31.75:3000"   # used to build RSS item URLs
cors_origins = ["*"]             # currently not consumed by any middleware

[server]
host = "0.0.0.0"
port = 8000
read_timeout = "30s"
write_timeout = "30s"
idle_timeout = "60s"
max_header_bytes = 1048576
graceful_shutdown = "10s"

[database]
host = "localhost"
port = 5432
user = "immortal"
password = "immortal"
name = "blog"
ssl_mode = "disable"
max_open_conns = 25
max_idle_conns = 5
conn_max_lifetime = "1h"
conn_max_idle_time = "30m"
timeout = "5s"

[redis]
host = "localhost"
port = 6379
password = ""
db = 0
pool_size = 10
min_idle_conns = 2
dial_timeout = "5s"
read_timeout = "3s"
write_timeout = "3s"
pool_timeout = "0s"
idle_timeout = "0s"
idle_check_frequency = "0s"

[jwt]
secret = "change-me"
access_expiration = "14m"
refresh_expiration = "168h"
issuer = "my-app"

[log]
level = "info"            # zap level; "debug" in development
format = ""               # "json" forces JSON output (otherwise dev console)
file_path = ""
max_size = 100
max_backups = 3
max_age = 30
compress = true

[email]
host = "mail.example.com"
port = 587
username = "noreply@example.com"
password = "your_password"
from = "noreply@example.com"

[llm]
apikey = "your-github-models-or-azure-key"

[rustfs]                  # S3-compatible storage (see "Status / Not yet exposed")
region = "us-east-1"
access_key_id = "rustfsadmin"
secret_access_key = "rustfsadmin"
endpoint = "http://localhost:9000/"
```

### Hot reload

When the config file changes on disk, the `config` module triggers an Fx
graceful shutdown (`fx.Shutdowner`), so the process restarts and picks up the
new configuration.

## Running

```bash
# Download dependencies
go mod download

# Run the HTTP server (loads config.toml, starts scheduled jobs)
go run ./cmd/server

# Build a static binary
go build -o blog-server ./cmd/server
./blog-server

# Create / migrate the database schema (Ent schema create; drops stale columns/indexes)
go run ./cmd/migration
```

### Docker

`Dockerfile` builds both the `blog-server` and `migration` binaries (CGO disabled)
and exposes port `8000`. Mount or bake in a `config.toml` and run `./migration`
once before `./blog-server`.

## API

All routes are mounted under `/api/v1`. Protected routes require an
`Authorization: Bearer <accessToken>` header. Responses use a common envelope:

```json
{ "code": 0, "msg": "success", "data": { } }
```

Errors return the same envelope with a non-zero `code` and an HTTP status mapped
from it (e.g. `401` for unauthorized/forbidden, `400` for invalid params,
`404` for not found, `409` for conflicts, `500` for internal errors). Messages
are in Chinese for client display.

### Auth

| Method | Path                         | Auth | Description                                  |
| ------ | ---------------------------- | ---- | -------------------------------------------- |
| POST   | `/api/v1/auth/captcha`       | —    | Send a captcha email.                        |
| POST   | `/api/v1/auth/register`      | —    | Register with email + captcha.               |
| POST   | `/api/v1/auth/login`         | —    | Login, returns access token; sets refresh cookie. |
| POST   | `/api/v1/auth/logout`        | —    | Clears the refresh-token cookie.             |
| POST   | `/api/v1/auth/refresh`       | cookie | Rotates access + refresh tokens from cookie. |

Request bodies:

```jsonc
// POST /api/v1/auth/captcha
{ "email": "user@example.com", "type": "Register" }   // type ∈ {Register, PasswordReset, ChangeEmail}

// POST /api/v1/auth/register
{ "email": "user@example.com", "password": "secret123", "passwordConfirm": "secret123", "captcha": "AB12CD" }

// POST /api/v1/auth/login
{ "email": "user@example.com", "password": "secret123" }
```

`login` / `register` responses include `accessToken`, `uuid`, `username`,
`role`, and `avatar`; the `refreshToken` is set as an `HttpOnly` cookie and is
not returned in the JSON body.

### Posts

| Method | Path                              | Auth            | Description                          |
| ------ | --------------------------------- | --------------- | ------------------------------------ |
| GET    | `/api/v1/posts`                   | —               | Paginated published posts (`page`, `pageSize`). |
| GET    | `/api/v1/posts/meta`              | —               | Post id + updated_at (for sitemaps). |
| GET    | `/api/v1/posts/:id`              | —               | Single published post (increments view count). |
| POST   | `/api/v1/posts`                  | admin           | Create a post.                       |
| GET    | `/api/v1/admin/posts`            | admin           | Admin list (filters: `status`, `keyword`). |
| GET    | `/api/v1/admin/posts/:id`        | admin           | Admin post detail.                   |
| PUT    | `/api/v1/admin/posts/:id`        | admin or owner  | Update a post.                       |
| DELETE | `/api/v1/admin/posts/:id`        | admin or owner  | Soft-delete a post.                  |

`POST /api/v1/posts` body (UserID is taken from the JWT, never from the client):

```jsonc
{
  "title": "My post",
  "summary": "Optional summary",
  "cover": "https://...",          // optional
  "content": "Post body...",
  "status": "published",           // draft | published | archived
  "categoryIDs": [1, 2],           // optional, existing category ids
  "tags": [3, 4]                   // optional, existing tag ids
}
```

`status` values: `draft`, `published`, `archived`. Post tags and categories are
referenced by **existing** ids; there are no API endpoints to create tags,
categories, or users.

### Links

| Method | Path                               | Auth | Description                              |
| ------ | ---------------------------------- | ---- | ---------------------------------------- |
| GET    | `/api/v1/links`                    | —    | List enabled links (public).             |
| POST   | `/api/v1/links/apply-link`         | —    | Submit a link application (created disabled). |

`POST /api/v1/links/apply-link` body:

```jsonc
{ "name": "Example", "url": "https://example.com", "description": "…", "avatar": "https://…" }
```

Link status (`normal`/`abnormal`) is computed only by the hourly scheduler job;
there is no HTTP endpoint to approve, edit, or delete links.

### RSS

| Method | Path                       | Auth | Description                                  |
| ------ | -------------------------- | ---- | -------------------------------------------- |
| GET    | `/api/v1/rss`              | —    | RSS 2.0 feed; optional `?page=N` (page size 10). |
| GET    | `/api/v1/rss/complete`     | —    | RSS 2.0 feed with all published posts.       |

### LLM summarization

| Method | Path                                      | Auth | Description                              |
| ------ | ----------------------------------------- | ---- | ---------------------------------------- |
| POST   | `/api/v1/model/summarize`                 | —    | Create a summary session, returns `sessionId`. |
| GET    | `/api/v1/model/summarize/:sessionId`      | —    | SSE stream of the generated summary.     |

`POST /api/v1/model/summarize` body: `{ "content": "…article text…" }`.

The stream emits `data: <text>\n\n` lines, terminates with `event: done\ndata: [DONE]`,
and reports failures via `event: error\ndata: <message>`.

## Project Structure

```
cmd/
  server/main.go        # Fx app wiring + HTTP server lifecycle
  migration/main.go     # Ent schema migration runner

config/                 # Config structs, Viper loader, validation, hot-reload
handler/                # Echo handlers + route registration (per-resource groups)
middleware/             # JWT auth, request logger, body-size limit
authz/                  # RBAC + ownership (ABAC) authorizer
contextx/               # Request-scoped user/request-id helpers
service/                # Business logic (auth, post, link, rss, mail, model)
repository/             # Data-access layer over Ent (post, link, user)
mapper/                 # Ent <-> domain entity conversions
entity/                 # Plain domain structs (Post, User, Link, RSS, ...)
request/                # Inbound DTOs with validation tags
response/               # Outbound DTOs / response envelopes
datastore/              # Ent client + transaction manager (txmgr.TxManager)
ent/                    # Ent schema definitions + generated code
cache/                  # Redis client (Store / AtomicStore / PatternScanner)
storage/                # S3-compatible Storage interface + implementation (not wired)
scheduler/              # Background jobs (view-count flush, link status check)
pkg/
  errx/                 # Error codes, AppError, HTTP-status mapping
  jwt/                  # Token generation / parsing
  validatorx/          # validator.v10 wrapper
  txmgr/                # Transaction manager interface
logger/                 # Zap-backed structured logger
templates/              # Embedded email templates (captcha.html)
utils/                  # Password hashing, random strings, Cloudflare IP fetch
```

### Request flow

1. `cmd/server/main.go` builds an Fx container with all modules and starts Echo.
2. Global middleware: request logger (attaches request id + per-request logger)
   and a 10 MiB body-size limit.
3. Routes are registered under `/api/v1`. Protected routes attach
   `middleware.AuthMiddleware`, which parses the Bearer token, resolves the user
   into `contextx.User`, and stores it on the request context.
4. Handlers bind + validate the request, call a `service`, and write a
   `response` envelope. Domain errors are wrapped in `errx.AppError` and rendered
   by the central `ErrorHandler`.
5. Services use repositories (Ent) and the transaction manager for writes. Posts
   and links authorize through `authz.Authorizer` (RBAC + ownership check).

### Scheduled jobs

Started on app start, cancelled on stop:

- **View-count flush** — every hour, reads `blog:post:view_count:*` keys from
  Redis (`SCAN` + `GETDEL`) and applies bulk increments to PostgreSQL.
- **Link status check** — every hour, `GET`s each HTTPS link and updates its
  status to `normal`/`abnormal` when it changes.

## Status / Not yet exposed

The following exist in the codebase but are **not** reachable over HTTP:

- **S3 object storage**: `storage.S3Storage` implements `Upload`/`Download`/
  `Delete`/`Copy`/`Exists`, but it is not registered in the Fx graph and no
  handler uses it. There is **no `/api/upload` endpoint**.
- **Links admin CRUD**: no endpoints to create/update/delete/approve links;
  link lifecycle is partially driven by the scheduler.
- **Tags, categories, users**: no endpoints to create or manage them; posts only
  reference pre-existing ids.
- **Comments**: a `Comment` table is defined in the Ent schema, but there is no
  repository, service, or handler for it.
- **CORS**: `app.cors_origins` is defined in config but no CORS middleware is
  registered, so cross-origin requests are not handled by the server.
- `service.PostService.GetPostsWithContent` is a stub that always returns
  `(nil, nil)`.
