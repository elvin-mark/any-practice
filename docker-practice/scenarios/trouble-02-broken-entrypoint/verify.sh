#!/usr/bin/env bash
set -e

CONTAINER="fixed-worker"
FAILED=0

# 1. Container running
if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[FAIL] Container '$CONTAINER' not found."
  exit 1
fi

STATE=$(docker inspect "$CONTAINER" --format '{{.State.Status}}')
if [ "$STATE" != "running" ]; then
  EXIT_CODE=$(docker inspect "$CONTAINER" --format '{{.State.ExitCode}}')
  echo "[FAIL] Container '$CONTAINER' status is '$STATE' (exit code: $EXIT_CODE)"
  exit 1
fi
echo "[PASS] Container '$CONTAINER' is running"

# 2. Check logs for success token
LOGS=$(docker logs "$CONTAINER" 2>&1 || true)
if echo "$LOGS" | grep -q "ENTRYPOINT_EXECUTION_SUCCESSFUL"; then
  echo "[PASS] Entrypoint executed successfully and printed expected verification log"
else
  echo "[FAIL] Logs do not contain 'ENTRYPOINT_EXECUTION_SUCCESSFUL': $LOGS"
  FAILED=1
fi

exit $FAILED
