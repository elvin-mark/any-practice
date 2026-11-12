#!/usr/bin/env bash
set -e

CONTAINER="mem-hungry"
FAILED=0

# 1. Container exists
if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[FAIL] Container '$CONTAINER' not found."
  exit 1
fi

# 2. Container running
STATE=$(docker inspect "$CONTAINER" --format '{{.State.Status}}')
if [ "$STATE" = "running" ]; then
  echo "[PASS] Container '$CONTAINER' is running"
else
  echo "[FAIL] Container '$CONTAINER' is not running (status: $STATE)"
  FAILED=1
fi

# 3. OOMKilled == false
OOM=$(docker inspect "$CONTAINER" --format '{{.State.OOMKilled}}')
if [ "$OOM" = "false" ]; then
  echo "[PASS] Container is not killed by OOM killer"
else
  echo "[FAIL] Container was terminated by OOM killer"
  FAILED=1
fi

# 4. Check memory limit >= 128MB
MEM=$(docker inspect "$CONTAINER" --format '{{.HostConfig.Memory}}')
if [ "$MEM" -ge 134217728 ]; then
  echo "[PASS] Container configured with adequate memory limit ($MEM bytes)"
else
  echo "[FAIL] Memory limit is $MEM bytes (must be at least 128MB/256MB)"
  FAILED=1
fi

exit $FAILED
