#!/usr/bin/env bash
set -e

CONTAINER="no-escalate"
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

# 2. Check SecurityOpt
SEC_OPTS=$(docker inspect "$CONTAINER" --format '{{json .HostConfig.SecurityOpt}}')
if [[ "$SEC_OPTS" == *"no-new-privileges"* ]]; then
  echo "[PASS] SecurityOpt contains 'no-new-privileges': $SEC_OPTS"
else
  echo "[FAIL] SecurityOpt does not contain 'no-new-privileges': $SEC_OPTS"
  FAILED=1
fi

exit $FAILED
