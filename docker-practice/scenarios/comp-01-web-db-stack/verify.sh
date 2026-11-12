#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

FAILED=0

# 1. Compose file exists
if [ ! -f docker-compose.yml ]; then
  echo "[FAIL] docker-compose.yml not found in $DIR"
  exit 1
fi
echo "[PASS] docker-compose.yml exists"

# 2. Check compose running services
SERVICES=$(docker compose ps --services 2>/dev/null || echo "")
if [[ "$SERVICES" == *"redis"* ]] && [[ "$SERVICES" == *"web"* ]]; then
  echo "[PASS] Both 'redis' and 'web' services are active in Compose"
else
  echo "[FAIL] Active services: '$SERVICES' (expected redis and web). Run 'docker compose up -d'"
  exit 1
fi

# 3. Check redis is healthy
REDIS_CONTAINER=$(docker compose ps -q redis 2>/dev/null || echo "")
if [ -z "$REDIS_CONTAINER" ]; then
  echo "[FAIL] Could not find redis container ID"
  exit 1
fi

REDIS_HEALTH=$(docker inspect "$REDIS_CONTAINER" --format '{{.State.Health.Status}}' 2>/dev/null || echo "none")
if [ "$REDIS_HEALTH" = "healthy" ]; then
  echo "[PASS] Redis container health status is 'healthy'"
else
  echo "[FAIL] Redis container health status is '$REDIS_HEALTH' (expected 'healthy')"
  FAILED=1
fi

# 4. Check web HTTP response
RESP=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8088 || echo "000")
if [ "$RESP" = "200" ]; then
  echo "[PASS] Web service responded with HTTP 200 on port 8088"
else
  echo "[FAIL] Web service failed to respond on port 8088 (code: $RESP)"
  FAILED=1
fi

exit $FAILED
