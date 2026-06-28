#!/bin/bash
# Recria containers de app com KAFKA_BROKERS atualizado (preserva demais env vars).
set -euo pipefail

KEY_FILE="$1"
KAFKA_BROKERS="$2"
AI_HOST="$3"
SOCKET_HOST="$4"
SERVER_HOSTS="$5"

SSH=(ssh -o StrictHostKeyChecking=no -i "$KEY_FILE" -o ConnectTimeout=15)

recreate_container() {
  local host="$1"
  local container="$2"
  local port_map="$3"

  "${SSH[@]}" "ec2-user@${host}" bash -s -- "$container" "$port_map" "$KAFKA_BROKERS" <<'REMOTE'
set -euo pipefail
CONTAINER="$1"
PORT_MAP="$2"
KAFKA_BROKERS="$3"

if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[$(hostname)] Container $CONTAINER não encontrado, pulando."
  exit 0
fi

IMAGE=$(docker inspect -f '{{.Config.Image}}' "$CONTAINER")
ENV_ARGS=$(docker inspect -f '{{range .Config.Env}}-e {{.}} {{end}}' "$CONTAINER")

if echo "$ENV_ARGS" | grep -q 'KAFKA_BROKERS='; then
  ENV_ARGS=$(echo "$ENV_ARGS" | sed -E "s/-e KAFKA_BROKERS=[^ ]+/-e KAFKA_BROKERS=${KAFKA_BROKERS}/")
else
  ENV_ARGS="$ENV_ARGS -e KAFKA_BROKERS=${KAFKA_BROKERS}"
fi

docker stop "$CONTAINER" || true
docker rm "$CONTAINER" || true
# shellcheck disable=SC2086
docker run -d --name "$CONTAINER" --restart unless-stopped -p "$PORT_MAP" $ENV_ARGS "$IMAGE"

echo "[$(hostname)] $CONTAINER recriado com KAFKA_BROKERS=$KAFKA_BROKERS"
REMOTE
}

echo "Sincronizando KAFKA_BROKERS=$KAFKA_BROKERS"

recreate_container "$AI_HOST" "scd-ai-server" "8070:8070"
recreate_container "$SOCKET_HOST" "scd-socket-server" "8765:8765"

IFS=',' read -ra HOSTS <<< "$SERVER_HOSTS"
for host in "${HOSTS[@]}"; do
  [ -n "$host" ] && recreate_container "$host" "scd-server" "3000:3000"
done

echo "Sync concluído."
