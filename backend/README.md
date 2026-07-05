# Blog Server

[中文文档](./README_CN.md)

A blog backend service built with Go.

## Tech Stack

- **Web Framework**: [Echo v5](https://github.com/labstack/echo)
- **Dependency Injection**: [Uber Fx](https://github.com/uber-go/fx)
- **ORM**: [Ent](https://entgo.io/ent)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Object Storage**: AWS S3 (Compatible)
- **Logging**: Zap
- **Configuration**: Viper
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)

## Features

- User authentication and authorization (JWT + RBAC)
- Post management (CRUD, tags, categories, view count)
- Link management (status monitoring)
- Image upload (S3 compatible)
- RSS feed
- Scheduled tasks (link status check, view count flush)
- LLM integration (post summary generation)

## Quick Start

### Prerequisites

- Go 1.26+
- PostgreSQL
- Redis

### Install Dependencies

```bash
go mod download
```

### Configuration

Copy the example config and modify:

```bash
cp config.yml.example config.yml
```

Or create `config.yml` directly:

```yaml
app:
  name: blog-server
  version: 1.0.0
  environment: development
  domain: localhost
  cors_origins:
    - http://localhost:3000

server:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  name: blog

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret
  access_expiration: 14m
  refresh_expiration: 168h

email:
  host: smtp.example.com
  port: 587
  username: your_email
  password: your_password
  from: noreply@example.com

rustfs:
  region: auto
  access_key_id: your_access_key
  secret_access_key: your_secret_key
  endpoint: http://localhost:9000
```

### Run

```bash
go run ./cmd/server
```

### Build

```bash
go build -o blog-server ./cmd/server
```

### Database Migration

```bash
go run ./cmd/migration
```

## Project Structure

```
cmd/
  server/main.go        # Application entry point
  migration/main.go     # Database migration runner

config/                 # Configuration structs, loader, validation
handler/                # HTTP handlers and route registration
service/                # Business logic layer
repository/             # Data access layer
mapper/                 # Ent <-> Entity mapping
entity/                 # Domain entities (plain structs)
request/                # Request DTOs with validation
response/               # Response DTOs

middleware/             # HTTP middleware (auth, logger, body limit)
authz/                  # Authorization (RBAC + ownership check)

datastore/              # Database client and transaction management
ent/                    # Ent ORM schemas and generated code
cache/                  # Redis client

storage/                # S3-compatible object storage
scheduler/              # Background job scheduler

pkg/
  errx/                 # Custom error types with error codes
  jwt/                  # JWT token generation and parsing
  validatorx/           # Validation wrapper
  txmgr/                # Transaction manager interface

contextx/               # Context helpers (user info propagation)
logger/                 # Zap structured logger
utils/                  # Utility functions
templates/              # Email templates
```

## API Endpoints

### Auth

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/captcha` - Send captcha email

### Posts

- `GET /api/posts` - Get posts list (public)
- `GET /api/posts/:id` - Get post details (public)
- `POST /api/posts` - Create post (admin)
- `PUT /api/posts/:id` - Update post (admin/owner)
- `DELETE /api/posts/:id` - Delete post (admin/owner)

### Links

- `GET /api/links` - Get links list (public)
- `POST /api/links` - Create link (admin)
- `PUT /api/links/:id` - Update link (admin/owner)
- `DELETE /api/links/:id` - Delete link (admin/owner)

### Other

- `GET /api/rss` - RSS feed
- `POST /api/upload` - Upload image (authenticated)

## License

MIT
