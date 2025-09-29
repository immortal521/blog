# 构建阶段
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/

# 运行阶段
FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/main .
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh
EXPOSE 8000

CMD ["/usr/local/bin/wait-for-it.sh", "redis:6379", "--", "wait-for-it.sh", "postgres:5432", "--", "./main"]

