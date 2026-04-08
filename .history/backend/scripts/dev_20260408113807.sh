#!/bin/bash
# Быстрый запуск для разработки
set -euo pipefail

echo "=== Hakathon Backend — Dev Setup ==="

# Проверяем .env
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo "✓ .env created. Edit if needed."
fi

# Загружаем переменные
set -a
source .env
set +a

echo ""
echo "1) Starting PostgreSQL..."
docker compose up -d postgres

echo ""
echo "2) Waiting for database..."
until docker compose exec -T postgres pg_isready -U "${POSTGRES_USER}" -q 2>/dev/null; do
    sleep 1
done
echo "   ✓ Database ready"

echo ""
echo "3) Running backend..."
echo "   DATABASE_URL=${DATABASE_URL}"
echo "   SERVER_PORT=${SERVER_PORT}"
echo "   GIN_MODE=${GIN_MODE}"
echo ""

go run ./cmd/server