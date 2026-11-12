#!/usr/bin/env bash
set -e

CONTAINER="monitored-web"
FAILED=0

# 1. Container running
if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[FAIL] Container '$CONTAINER' not found."
  exit 1
fi

STATE=$(docker inspect "$CONTAINER" --format '{{.State.Status}}')
if [ "$STATE" != "running" ]; then
  echo "[FAIL] Container '$CONTAINER' status is '$STATE' (expected running)"
  exit 1
fi
echo "[PASS] Container '$CONTAINER' is running"

# 2. Check healthcheck configuration exists
TEST_CMD=$(docker inspect "$CONTAINER" --format '{{json .Config.Healthcheck.Test}}')
if [ -z "$TEST_CMD" ] || [ "$TEST_CMD" = "null" ]; then
  echo "[FAIL] No healthcheck configured on '$CONTAINER'"
  exit 1
fi
echo "[PASS] Healthcheck command configured: $TEST_CMD"

# 3. Wait up to 15 seconds for healthy status
HEALTH_STATUS=""
for i in {1..15}; do
  HEALTH_STATUS=$(docker inspect "$CONTAINER" --format '{{.State.Health.Status}}' 2>/dev/null || echo "unknown")
  if [ "$HEALTH_STATUS" = "healthy" ]; then
    break
  fi
  sleep 1
done

if [ "$HEALTH_STATUS" = "healthy" ]; then
  echo "[PASS] Container health status is 'healthy'"
else
  echo "[FAIL] Container health status is '$HEALTH_STATUS' (expected 'healthy')"
  FAILED=1
fi

exit $FAILED
