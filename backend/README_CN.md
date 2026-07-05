# Blog Server

[English](./README.md)

一个使用 Go 构建的博客后端服务。

## 技术栈

- **Web 框架**: [Echo v5](https://github.com/labstack/echo)
- **依赖注入**: [Uber Fx](https://github.com/uber-go/fx)
- **ORM**: [Ent](https://entgo.io/ent)
- **数据库**: PostgreSQL
- **缓存**: Redis
- **对象存储**: AWS S3 (兼容)
- **日志**: Zap
- **配置**: Viper
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt)
- **参数校验**: [go-playground/validator](https://github.com/go-playground/validator)

## 功能特性

- 用户认证与授权 (JWT + RBAC)
- 文章管理 (CRUD、标签、分类、浏览量)
- 友链管理 (状态监控)
- 图片上传 (S3 兼容)
- RSS 订阅
- 定时任务 (友链状态检查、浏览量持久化)
- LLM 集成 (文章摘要生成)

## 快速开始

### 环境要求

- Go 1.26+
- PostgreSQL
- Redis

### 安装依赖

```bash
go mod download
```

### 配置

复制示例配置并修改：

```bash
cp config.yml.example config.yml
```

或直接创建 `config.yml`：

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

### 运行

```bash
go run ./cmd/server
```

### 构建

```bash
go build -o blog-server ./cmd/server
```

### 数据库迁移

```bash
go run ./cmd/migration
```

## 项目结构

```
cmd/
  server/main.go        # 应用入口
  migration/main.go     # 数据库迁移工具

config/                 # 配置结构体、加载器、校验
handler/                # HTTP 处理器与路由注册
service/                # 业务逻辑层
repository/             # 数据访问层
mapper/                 # Ent <-> Entity 映射
entity/                 # 领域实体 (纯结构体)
request/                # 请求 DTO (含参数校验)
response/               # 响应 DTO

middleware/             # HTTP 中间件 (认证、日志、请求体限制)
authz/                  # 授权 (RBAC + 资源归属检查)

datastore/              # 数据库客户端与事务管理
ent/                    # Ent ORM 模式定义与生成代码
cache/                  # Redis 客户端

storage/                # S3 兼容对象存储
scheduler/              # 后台任务调度器

pkg/
  errx/                 # 自定义错误类型与错误码
  jwt/                  # JWT 令牌生成与解析
  validatorx/           # 参数校验封装
  txmgr/                # 事务管理器接口

contextx/               # Context 工具 (用户信息传递)
logger/                 # Zap 结构化日志
utils/                  # 工具函数
templates/              # 邮件模板
```

## API 接口

### 认证

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/refresh` - 刷新访问令牌
- `POST /api/auth/captcha` - 发送验证码邮件

### 文章

- `GET /api/posts` - 获取文章列表 (公开)
- `GET /api/posts/:id` - 获取文章详情 (公开)
- `POST /api/posts` - 创建文章 (管理员)
- `PUT /api/posts/:id` - 更新文章 (管理员/作者)
- `DELETE /api/posts/:id` - 删除文章 (管理员/作者)

### 友链

- `GET /api/links` - 获取友链列表 (公开)
- `POST /api/links` - 创建友链 (管理员)
- `PUT /api/links/:id` - 更新友链 (管理员/作者)
- `DELETE /api/links/:id` - 删除友链 (管理员/作者)

### 其他

- `GET /api/rss` - RSS 订阅
- `POST /api/upload` - 上传图片 (需认证)

## License

MIT
