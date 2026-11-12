#!/usr/bin/env bash
set -e

CONTAINER="init-worker"
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

# 2. HostConfig.Init == true
INIT_FLAG=$(docker inspect "$CONTAINER" --format '{{.HostConfig.Init}}')
if [ "$INIT_FLAG" = "true" ]; then
  echo "[PASS] Container HostConfig.Init is true (--init enabled)"
else
  echo "[FAIL] Container HostConfig.Init is '$INIT_FLAG' (expected true). Launch with --init"
  FAILED=1
fi

# 3. Check PID 1 process inside container
PID1_CMD=$(docker exec "$CONTAINER" ps -o pid,comm | awk '$1 == 1 {print $2}')
if [[ "$PID1_CMD" == *"init"* ]] || [[ "$PID1_CMD" == *"tini"* ]]; then
  echo "[PASS] PID 1 inside container is init reaper ($PID1_CMD)"
else
  echo "[FAIL] PID 1 inside container is '$PID1_CMD' (expected docker-init / tini)"
  FAILED=1
fi

exit $FAILED
