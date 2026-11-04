#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
compose_file=$repo_root/docker-compose.test.yaml
run_id=${TEST_RUN_ID:-test}

case "$run_id" in
  *[!a-z0-9-]* | -*) echo "TEST_RUN_ID must contain lowercase letters, digits, or hyphens" >&2; exit 2 ;;
esac
project=mypills-lab2-$run_id

dc() {
  docker compose -p "$project" -f "$compose_file" "$@"
}

stand_up() {
  if [ "${TEST_APP_BUILD:-0}" = 1 ] || ! docker image inspect "${project}-app:latest" >/dev/null 2>&1; then
    dc up -d --build --wait db app
  else
    dc up -d --no-build --wait db app
  fi
}

stand_reset() {
  dc down -v --remove-orphans
  stand_up
}

run_cleanup() {
  status=$?
  trap - EXIT INT TERM
  if ! stand_reset; then
    echo "Failed to restore the test stand" >&2
    status=1
  fi
  exit "$status"
}

case "${1:-}" in
  up)
    stand_up
    ;;
  reset)
    stand_reset
    ;;
  down)
    dc down -v --remove-orphans
    ;;
  ps)
    dc ps
    ;;
  db)
    shift
    dc exec -T db psql -X -U mypills_test -d mypills_test "$@"
    ;;
  run)
    shift
    if [ "$#" -eq 0 ]; then
      echo "Usage: $0 run COMMAND [ARG...]" >&2
      exit 2
    fi
    stand_reset
    trap run_cleanup EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    "$@"
    ;;
  *)
    echo "Usage: $0 {up|reset|down|ps|db|run COMMAND [ARG...]}" >&2
    exit 2
    ;;
esac
