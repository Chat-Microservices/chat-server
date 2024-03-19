#!/bin/sh
source /root/.env

wait_for_port() {
    host="$1"
    port="$2"
    timeout="${3:-15}"

    echo "Waiting for $host:$port to be available..."
    timeout $timeout sh -c "while ! nc -z $host $port; do sleep 1; done"
}

# Ожидание доступности порта
wait_for_port "pg-chat-server" "5432"

# накатываем миграции
if ! goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v; then
  echo "Migration failed"
  exit 1
fi

# запускаем приложение
./chat_server -config-path=/root/.env