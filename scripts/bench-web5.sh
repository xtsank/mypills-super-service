#!/bin/sh
set -eu
cd "$(dirname "$0")/.."

pool=${1:-}
case "$pool" in
    6|12) ;;
    *) echo "Usage: sh scripts/bench-web5.sh 6|12" >&2; exit 2 ;;
esac

project="mypills-lab5-exp-$pool-$(date +%s)"
WEB_PORT=${WEB_PORT:-18085}
LAB5_DB_MAX_OPEN_CONNS=$pool
LAB5_DB_MAX_IDLE_CONNS=6
export WEB_PORT LAB5_DB_MAX_OPEN_CONNS LAB5_DB_MAX_IDLE_CONNS

warmup=${WARMUP_SECONDS:-30}
duration=${MEASURE_SECONDS:-180}
case "$warmup:$duration" in
    *[!0-9:]*|0:*|*:0) echo "Durations must be positive integers" >&2; exit 2 ;;
esac

out="$PWD/benchmark-results/$project"
mkdir "$PWD/benchmark-results" 2>/dev/null || test -d "$PWD/benchmark-results"
mkdir "$out"
dc() { docker compose -p "$project" -f docker-compose.yaml -f docker-compose.lab5.yaml "$@"; }
trap 'dc stop' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

dc up -d --build --wait --wait-timeout 120
app_port=$(dc exec -T app sh -c 'printf "%s" "$SERVER_PORT"')
dc exec -T db sh -ec 'exec psql -X -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < tests/fixtures/e2e_seed.sql

{
    date -u +%FT%TZ
    docker info --format 'docker_cpus={{.NCPU}} docker_memory_bytes={{.MemTotal}}'
    printf 'project=%s\nmax_open=%s\nmax_idle=6\nwarmup_seconds=%s\nmeasure_seconds=%s\n' "$project" "$pool" "$warmup" "$duration"
    dc images
} > "$out/environment.txt"

set +e
docker run --rm --network "${project}_default" --user "$(id -u):$(id -g)" \
    -e "WARMUP_SECONDS=$warmup" -e "MEASURE_SECONDS=$duration" \
    -e "APP_PORT=$app_port" \
    -v "$PWD/lab5/load.js:/load.js:ro" -v "$out:/results" \
    grafana/k6:1.6.1 run --out json=/results/samples.jsonl --summary-export=/results/k6-summary.json \
    /load.js > "$out/k6.log" 2>&1
rc=$?
set -e
# k6 exports setup_data with the test user's JWT. Keep metrics without the session.
if [ -s "$out/k6-summary.json" ]; then
    python3 - "$out/k6-summary.json" <<'PYTHON'
import json
import sys
from pathlib import Path

summary = Path(sys.argv[1])
data = json.loads(summary.read_text())
data.get("setup_data", {}).pop("token", None)
summary.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n")
PYTHON
fi
cat "$out/k6.log"

if [ -s "$out/samples.jsonl" ]; then
    python3 scripts/summarize-web5.py "$out/samples.jsonl" "$duration" > "$out/summary.md"
    cat "$out/summary.md"
fi
# Give Prometheus a scrape after the final request before saving DB counters.
sleep 5
curl -fsSG "http://127.0.0.1:$WEB_PORT/monitoring/api/datasources/proxy/uid/lab5-prometheus/api/v1/query" \
    --data-urlencode 'query={__name__=~"go_sql_wait_(count|duration_seconds)_total|go_sql_max_open_connections",job="mypills-api"}' \
    > "$out/db-pool.json" || echo "Could not save the Prometheus snapshot" >&2
printf '\nFiles: %s\n' "$out"
exit "$rc"
