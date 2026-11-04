#!/bin/sh
set -eu

if [ "$#" -eq 0 ]; then
  echo "Usage: $0 COMMAND [ARG...]" >&2
  exit 2
fi

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

docker_endpoint=${DOCKER_SOCKET:-$(docker context inspect --format '{{.Endpoints.docker.Host}}')}
case "$docker_endpoint" in
  unix://*) docker_socket=${docker_endpoint#unix://} ;;
  /*) docker_socket=$docker_endpoint ;;
  *) echo "A local Docker socket is required: $docker_endpoint" >&2; exit 1 ;;
esac

if [ ! -S "$docker_socket" ]; then
  echo "Docker socket not found: $docker_socket" >&2
  exit 1
fi

test_uid=${TEST_UID:-$(id -u)}
test_gid=${TEST_GID:-$(id -g)}
docker_socket_gid=$(stat -c '%g' "$docker_socket")

if [ "${CI_RUNNER_BUILD:-0}" = 1 ] || ! docker image inspect mypills-lab2-test-runner:local >/dev/null 2>&1; then
  docker compose -f docker-compose.test.yaml build test-runner
fi
docker run --rm --init \
  --user "$test_uid:$test_gid" \
  --group-add "$docker_socket_gid" \
  --mount "type=bind,src=$repo_root,dst=$repo_root" \
  --mount "type=bind,src=$docker_socket,dst=/var/run/docker.sock" \
  --workdir "$repo_root" \
  -e ALLURE_CLI=/opt/allure/node_modules/.bin/allure \
  -e CGO_ENABLED=1 \
  -e HOME=/tmp \
  -e GOCACHE=/tmp/go-build \
  -e TEST_APP_PORT \
  -e TEST_APP_BUILD \
  -e TEST_RUN_ID \
  -e TEST_UID="$test_uid" \
  -e TEST_GID="$test_gid" \
  -e TEST_SEED_FILE \
  mypills-lab2-test-runner:local "$@"
