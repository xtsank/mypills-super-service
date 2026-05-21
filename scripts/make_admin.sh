#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ -f "${ROOT_DIR}/.env" ]; then
  set -a
  . "${ROOT_DIR}/.env"
  set +a
fi

if [ -z "${DB_USER:-}" ] || [ -z "${DB_PASSWORD:-}" ] || [ -z "${DB_NAME:-}" ] || [ -z "${DB_HOST:-}" ] || [ -z "${DB_PORT:-}" ]; then
  echo "DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT must be set" >&2
  exit 1
fi

if [ -z "${1:-}" ]; then
  echo "Usage: $(basename "$0") <login>" >&2
  exit 1
fi

LOGIN="$1"
SAFE_LOGIN=$(printf "%s" "$LOGIN" | sed "s/'/''/g")

sudo docker compose exec -T db psql \
  -U "$DB_USER" \
  -d "$DB_NAME" \
  -v ON_ERROR_STOP=1 \
  -c "UPDATE Users SET is_admin = true WHERE login = '$SAFE_LOGIN';" \
  -c "SELECT id, login, is_admin FROM Users WHERE login = '$SAFE_LOGIN';"
