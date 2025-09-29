#!/bin/bash
set -e

CONFIG_DIR="/etc/blog"
BACKEND_CONFIG="$CONFIG_DIR/backend/config.yml"
NGINX_CONFIG="$CONFIG_DIR/nginx/default.conf"

sudo mkdir -p "$CONFIG_DIR/backend"
sudo mkdir -p "$CONFIG_DIR/nginx"

# 如果配置文件不存在，则创建默认示例
if [ ! -f "$BACKEND_CONFIG" ]; then
  sudo tee "$BACKEND_CONFIG" > /dev/null <<EOF
server:
  port: 8000
database:
  host: postgres
  port: 5432
  user: postgres
  password: postgres
  name: blog
redis:
  host: redis
  port: 6379
EOF
  echo "Created default backend config at $BACKEND_CONFIG"
fi

if [ ! -f "$NGINX_CONFIG" ]; then
  sudo tee "$NGINX_CONFIG" > /dev/null <<EOF
server {
    listen 80;

    server_name localhost;

    location / {
        proxy_pass http://frontend:3000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location /api/ {
        proxy_pass http://backend:8080/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF
  echo "Created default nginx config at $NGINX_CONFIG"
fi

# 写 docker-compose.yml
cat > docker-compose.yml <<'EOF'
services:
  frontend:
    build: ./frontend
    container_name: blog_frontend
    restart: always
    ports:
     - "3000:3000"
    environment:
     - NODE_ENV=production
    depends_on:
     - backend

  backend:
    build: ./backend
    container_name: blog_backend
    restart: always
    ports:
     - "8000:8000"
    environment:
     - REDIS_HOST=redis
     - REDIS_PORT=6379
     - POSTGRES_HOST=postgres
     - POSTGRES_PORT=5432
     - POSTGRES_USER=postgres
     - POSTGRES_PASSWORD=postgres
     - POSTGRES_DB=mydb
    volumes:
     - /etc/blog/backend/config.yml:/etc/backend/config.yaml:ro
    depends_on:
     - redis
     - postgres

  redis:
    image: docker.io/redis:8-alpine
    container_name: redis
    restart: always
    ports:
      - "6379:6379"

  postgres:
    image: docker.io/postgres:18-alpine
    container_name: postgres
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: blog
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  nginx:
    image: docker.io/nginx:1.25-alpine
    container_name: nginx
    restart: always
    ports:
      - "80:80"
    volumes:
      - /etc/blog/nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      - frontend
      - backend

volumes:
  pgdata:
EOF

# 构建并启动容器
sudo podman compose up -d --build

echo "Deployment completed!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"
echo "Nginx: http://localhost"
