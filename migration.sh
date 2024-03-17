#!/bin/bash
source .env

# Ожидание доступности PostgreSQL
while ! nc -z pg-chat-server 5432; do
  >&2 echo "PostgreSQL недоступен - ожидание..."
  sleep 2
done


sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v
