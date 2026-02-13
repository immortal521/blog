# Blog Server

A blog backend service built with Go + Fiber.

## Tech Stack

- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Dependency Injection**: [Uber FX](https://github.com/uber-go/fx)
- **Database**: PostgreSQL + GORM
- **Cache**: Redis
- **Object Storage**: AWS S3 (Compatible)
- **Logging**: Zap
- **Configuration**: Viper

## Features

- User authentication and authorization (JWT)
- Post management (CRUD, tags, categories)
- Link management (status monitoring)
- Image upload (S3)
- RSS feed
- Statistics
- Scheduled tasks (link status check, view count flush)

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis

### Install Dependencies

```bash
go mod download
```

### Configuration

Create a `config.yml` configuration file:

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
  access_expiration: 24h
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
go run cmd/main.go
```

### Build

```bash
go build -o blog-server cmd/main.go
```

## Project Structure

```
cmd/
  main.go          # Entry point

internal/
  cache/           # Redis cache
  config/          # Configuration management
  database/        # Database connection and migration
  entity/          # Data models
  handler/         # HTTP handlers
  middleware/      # Middleware
  repository/      # Data access layer
  request/         # Request DTOs
  response/        # Response DTOs
  router/          # Router configuration
  scheduler/       # Scheduled tasks
  service/         # Business logic layer
  storage/         # Object storage

pkg/
  errs/            # Error handling
  logger/          # Logger wrapper
  utils/           # Utility functions
  validatorx/      # Validator

templates/        # Email and other templates
```

## API Endpoints

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET /api/posts` - Get posts list
- `GET /api/posts/:id` - Get post details
- `POST /api/posts` - Create post (requires auth)
- `GET /api/links` - Get links list
- `GET /api/rss` - RSS feed
- `GET /api/stats` - Statistics

## License

MIT
