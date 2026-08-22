# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct \
    http_proxy=${HTTPS_PROXY} \
    https_proxy=${HTTPS_PROXY} \
    HTTP_PROXY=${HTTPS_PROXY} \
    HTTPS_PROXY=${HTTPS_PROXY}

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/blog-server \
    ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/migration \
    ./cmd/migration

FROM alpine:3.20

WORKDIR /app

ENV http_proxy=${HTTPS_PROXY} https_proxy=${HTTPS_PROXY} HTTP_PROXY=${HTTPS_PROXY} HTTPS_PROXY=${HTTPS_PROXY}

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/blog-server /app/blog-server
COPY --from=builder /out/migration /app/migration

EXPOSE 8000

CMD ["/app/blog-server"]
