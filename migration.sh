#!/bin/sh

MIGRATION_DSN="host=${DB_HOST:-postgres} port=${DB_PORT:-5432} dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD"

attempt=0
max_attempts=${MAX_WAIT_ATTEMPTS:-30}

until goose -dir "${MIGRATION_DIR:-./migrations}" postgres "${MIGRATION_DSN}" status >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge "$max_attempts" ]; then
    echo "postgres is not available after ${max_attempts} attempts, last error:" >&2
    goose -dir "${MIGRATION_DIR:-./migrations}" postgres "${MIGRATION_DSN}" status >&2 || true
    exit 1
  fi
  echo "waiting for postgres... (${attempt}/${max_attempts})"
  sleep 1
done

exec goose -dir "${MIGRATION_DIR:-./migrations}" postgres "${MIGRATION_DSN}" up -v
