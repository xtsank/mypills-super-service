#!/bin/sh
set -eu

cd "$(dirname "$0")/.."
LAB3_PROFILE=${1:-full}
case "$LAB3_PROFILE" in
  full|smoke) ;;
  *) echo "PROFILE must be full or smoke" >&2; exit 2 ;;
esac

LAB3_RUN_ID="$(date -u +%Y%m%dT%H%M%SZ | tr '[:upper:]' '[:lower:]')-$$"
if [ -z "${LAB3_CPU:-}" ]; then
  LAB3_CPU=$(awk '/^Cpus_allowed_list:/ { split($2, cpu, /[-,]/); print cpu[1] }' /proc/self/status 2>/dev/null || true)
  LAB3_CPU=${LAB3_CPU:-0}
fi
case "$LAB3_CPU" in
  *[!0-9]*) echo "LAB3_CPU must be one CPU index" >&2; exit 2 ;;
esac
LAB3_UID=$(id -u)
LAB3_GID=$(id -g)
export LAB3_PROFILE LAB3_RUN_ID LAB3_CPU LAB3_UID LAB3_GID

results="benchmark-results/$LAB3_RUN_ID"
mkdir -p benchmark-results
mkdir "$results"

dc() {
  docker compose -p "mypills-lab3-$LAB3_RUN_ID" -f docker-compose.lab3.yaml "$@"
}

cleanup() {
  status=$?
  trap - EXIT INT TERM
  dc down -v || status=1
  exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

printf 'Run %s (%s, CPU %s). Results: %s\n' "$LAB3_RUN_ID" "$LAB3_PROFILE" "$LAB3_CPU" "$results"
dc build benchmark
dc up -d --wait --wait-timeout 90 db
dc run --rm --no-deps benchmark
