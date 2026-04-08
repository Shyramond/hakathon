#!/bin/sh
# Ожидает доступность PostgreSQL перед запуском приложения
set -e

HOST="${DB_HOST:-localhost}"
PORT="${DB_PORT:-5432}"
MAX_RETRIES="${MAX_RETRIES:-30}"
RETRY_INTERVAL="${RETRY_INTERVAL:-2}"

echo "Waiting for PostgreSQL at ${HOST}:${PORT}..."

retries=0
until pg_isready -h "$HOST" -p "$PORT" -q 2>/dev/null; do
    retries=$((retries + 1))
    if [ "$retries" -ge "$MAX_RETRIES" ]; then
        echo "ERROR: PostgreSQL not available after ${MAX_RETRIES} attempts"
        exit 1
    fi
    echo "  Attempt ${retries}/${MAX_RETRIES} — retrying in ${RETRY_INTERVAL}s..."
    sleep "$RETRY_INTERVAL"
done

echo "PostgreSQL is ready!"
exec "$@"