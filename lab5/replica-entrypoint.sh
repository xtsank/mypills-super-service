#!/bin/sh
set -eu

: "${PGDATA:?}"
: "${PRIMARY_HOST:?}"
: "${REPLICATION_USER:?}"
: "${PGPASSWORD:?}"

mkdir -p "$PGDATA"
chown -R postgres:postgres "$PGDATA"
chmod 700 "$PGDATA"

if [ ! -f "$PGDATA/PG_VERSION" ]; then
    if [ -n "$(ls -A "$PGDATA")" ]; then
        echo "Replica data directory is incomplete; inspect or replace only the replica volume." >&2
        exit 1
    fi

    echo "Creating PostgreSQL replica from ${PRIMARY_HOST}" >&2
    su-exec postgres pg_basebackup \
        -h "$PRIMARY_HOST" -p 5432 -U "$REPLICATION_USER" \
        -D "$PGDATA" -R -X stream -S mypills_lab5 -w
fi

exec docker-entrypoint.sh postgres -c hot_standby=on
