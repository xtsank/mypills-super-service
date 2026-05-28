#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

log() {
  printf '%s %s\n' "$(date +%H:%M:%S)" "$*"
}

set -a
if [ -f .env ]; then
  . ./.env
fi
set +a

profiles=(--profile postgres --profile mongo)

cleanup() {
  log "Stopping databases"
  sudo docker compose "${profiles[@]}" stop db mongo
}
trap cleanup EXIT

log "Starting databases"
sudo docker compose "${profiles[@]}" up -d db mongo

log "Waiting for postgres"
ready=false
for _ in {1..30}; do
  if sudo docker compose exec -T db pg_isready -U "${DB_USER:-user}" >/dev/null 2>&1; then
    ready=true
    break
  fi
  sleep 1
done

if [ "$ready" != "true" ]; then
  log "Postgres did not become ready"
  exit 1
fi

log "Waiting for mongo"
ready=false
for _ in {1..30}; do
  if sudo docker compose exec -T mongo mongosh --quiet --username "${DB_USER:-user}" --password "${DB_PASSWORD:-password}" --authenticationDatabase admin --eval "db.adminCommand({ ping: 1 })" >/dev/null 2>&1; then
    ready=true
    break
  fi
  sleep 1
done

if [ "$ready" != "true" ]; then
  log "Mongo did not become ready"
  exit 1
fi

log "Running migrator"
MIGRATION_DIRECTION="${MIGRATION_DIRECTION:-pg2mongo}" \
MIGRATION_CLEANUP="${MIGRATION_CLEANUP:-true}" \
POSTGRES_HOST=localhost \
MONGO_HOST=localhost \
go run ./src/cmd/migrator
