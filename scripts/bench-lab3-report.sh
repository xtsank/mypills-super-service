#!/bin/sh
set -eu

cd "$(dirname "$0")/.."
LAB3_RUN_ID=${1:-}
case "$LAB3_RUN_ID" in
  ''|*[!a-z0-9-]*|-*) echo "Usage: make bench-lab3-report RUN_ID=<run-id>" >&2; exit 2 ;;
esac
LAB3_GRAFANA_PORT=${LAB3_GRAFANA_PORT:-3003}
export LAB3_RUN_ID LAB3_GRAFANA_PORT

for file in samples.csv summary.json; do
  if [ ! -r "benchmark-results/$LAB3_RUN_ID/$file" ]; then
    echo "Missing benchmark-results/$LAB3_RUN_ID/$file" >&2
    exit 1
  fi
done

dc() {
  docker compose -p mypills-lab3-report -f docker-compose.lab3-report.yaml "$@"
}

dc up -d --wait --wait-timeout 90 results-db
dc run --rm --no-deps importer
dc up -d --wait --wait-timeout 90 grafana
printf 'Grafana: http://127.0.0.1:%s/d/mypills-lab3/lab3?var-run_id=%s\n' "$LAB3_GRAFANA_PORT" "$LAB3_RUN_ID"
