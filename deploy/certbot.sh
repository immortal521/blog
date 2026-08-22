#!/usr/bin/env bash
# 以部署用户身份运行 certbot：证书存于项目目录 deploy/letsencrypt，
# rootless/rootful 的 nginx 容器都能直接读取，无需改属主或复制；config/work/logs 也都在项目内。
#
# 子命令：certonly …（签发，参数原样传 certbot）| renew（续期）| import <域名> | 其它原样传给 certbot。
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
DEPLOY_USER="$(stat -c '%U' "$PROJECT_DIR")"
DEPLOY_GROUP="$(stat -c '%G' "$PROJECT_DIR")"
# 证书存放模式：project（默认，项目内 deploy/letsencrypt，归部署用户，容器挂载读取）
#              system（/etc/letsencrypt，root 拥有，由系统 certbot 管理，供 manual 模式使用）
CERT_MODE="${CERT_MODE:-project}"
if [ "$CERT_MODE" = "system" ]; then
  CERT_BASE="/etc/letsencrypt"
  ACME_WEBROOT="/var/www/certbot"
  CERTBOT_ARGS=()
else
  CERT_BASE="$PROJECT_DIR/deploy/letsencrypt"
  # ACME HTTP-01 挑战目录：放在项目内、归部署用户所有，避免依赖宿主 /var/www/certbot 的权限/sudo
  ACME_WEBROOT="$PROJECT_DIR/deploy/acme"
  CERTBOT_ARGS=(
    --config-dir "$CERT_BASE/etc"
    --work-dir   "$CERT_BASE/work"
    --logs-dir   "$CERT_BASE/logs"
  )
fi

# 签发所用邮箱：优先 CERTBOT_EMAIL，回退到 .env 的 EMAIL_FROM
CERTBOT_EMAIL="${CERTBOT_EMAIL:-${EMAIL_FROM:-}}"

cmd="${1:-}"
case "$cmd" in
  import)
    [ "$CERT_MODE" = "system" ] && { echo "ERROR: CERT_MODE=system 时无需 import，证书本就在 /etc/letsencrypt" >&2; exit 1; }
    domain="${2:-}"
    [ -n "$domain" ] || { echo "ERROR: 用法: $0 import <域名>" >&2; exit 1; }
    src_live="/etc/letsencrypt/live/$domain"
    src_arch="/etc/letsencrypt/archive/$domain"
    [ -d "$src_live" ] || { echo "ERROR: $src_live 不存在" >&2; exit 1; }
    sudo mkdir -p "$CERT_BASE/etc/live" "$CERT_BASE/etc/archive"
    sudo cp -rL "$src_live" "$CERT_BASE/etc/live/$domain"
    sudo cp -rL "$src_arch" "$CERT_BASE/etc/archive/$domain"
    sudo chown -R "$DEPLOY_USER:$DEPLOY_GROUP" "$CERT_BASE/etc/live/$domain" "$CERT_BASE/etc/archive/$domain"
    echo "已导入 $domain -> $CERT_BASE/etc/{live,archive}/$domain（属主 $DEPLOY_USER）"
    echo "提示：之后续期请改用 '$0 renew'，或直接 '$0 certonly ...' 重新签发以生成项目内的 renewal 配置。"
    exit 0
    ;;
esac

# certonly 时补全非交互签发参数，避免 certbot 进入交互式询问
if [ "$cmd" = "certonly" ]; then
  # 未显式指定 -w 时，默认用项目内挑战目录
  if ! printf ' %s ' "$*" | grep -q -- '-w'; then
    set -- --webroot -w "$ACME_WEBROOT" "$@"
  fi
  if [ -n "$CERTBOT_EMAIL" ] && ! printf ' %s ' "$*" | grep -q -- '--email'; then
    set -- --agree-tos --no-eff-email --non-interactive --email "$CERTBOT_EMAIL" "$@"
  elif ! printf ' %s ' "$*" | grep -q -- '--agree-tos'; then
    set -- --agree-tos --non-interactive "$@"
  fi
fi

if [ "$CERT_MODE" = "system" ]; then
  # 系统模式：证书在 /etc/letsencrypt，由 root 直接运行 certbot
  mkdir -p "$ACME_WEBROOT"
  exec certbot "${CERTBOT_ARGS[@]}" "$@"
fi

# 项目模式：挑战目录归部署用户所有；root 运行时降权到部署用户写入，并修正目录属主，
# 避免 rootful 下 docker/podman daemon 把 bind 源目录建成 root:root 导致 certbot 写入失败。
mkdir -p "$ACME_WEBROOT" "$CERT_BASE/etc"
if [ "$(id -u)" = "0" ]; then
  chown -R "$DEPLOY_USER:$DEPLOY_GROUP" "$CERT_BASE" "$ACME_WEBROOT"
  exec sudo -u "$DEPLOY_USER" env "HOME=$(getent passwd "$DEPLOY_USER" | cut -d: -f6)" "CERTBOT_EMAIL=$CERTBOT_EMAIL" \
    certbot "${CERTBOT_ARGS[@]}" "$@"
else
  exec certbot "${CERTBOT_ARGS[@]}" "$@"
fi
