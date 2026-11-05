#!/usr/bin/env bash
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

export TEST_RUN_ID=${TEST_RUN_ID:-wireshark-$(date +%s)-$$}
export TEST_APP_PORT=${TEST_APP_PORT:-18081}
export TEST_SEED_FILE=e2e_seed.sql
export E2E_BASE_URL=http://127.0.0.1:$TEST_APP_PORT
capture_file=${1:-docs/lr-testing/e2e-http.pcapng}
capture_pid=

cleanup() {
  local status=$?
  trap - EXIT INT TERM
  if [[ -n $capture_pid ]] && kill -0 "$capture_pid" 2>/dev/null; then
    kill -INT "$capture_pid"
    wait "$capture_pid" || status=1
  fi
  ./scripts/test-stand.sh down || status=1
  exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

command -v dumpcap >/dev/null || { echo 'Wireshark dumpcap is required' >&2; exit 1; }
command -v curl >/dev/null || { echo 'curl is required' >&2; exit 1; }
command -v jq >/dev/null || { echo 'jq is required' >&2; exit 1; }

./scripts/test-stand.sh reset
mkdir -p "$(dirname "$capture_file")"
echo "Capturing HTTP traffic on lo:$TEST_APP_PORT to $capture_file"
dumpcap -q -i lo -f "tcp port $TEST_APP_PORT" -w "$capture_file" &
capture_pid=$!
sleep 1
kill -0 "$capture_pid" 2>/dev/null || { echo 'Wireshark capture did not start' >&2; exit 1; }

./scripts/replay-e2e-http.sh
# Give libpcap time to flush the final HTTP responses before stopping capture.
sleep 2
kill -INT "$capture_pid"
wait "$capture_pid"
capture_pid=
echo "Capture saved to $capture_file"
