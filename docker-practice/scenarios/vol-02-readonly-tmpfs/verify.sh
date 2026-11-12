#!/usr/bin/env bash
set -e

CONTAINER="secure-nginx"
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

# 2. ReadonlyRootfs is true
RO=$(docker inspect "$CONTAINER" --format '{{.HostConfig.ReadonlyRootfs}}')
if [ "$RO" = "true" ]; then
  echo "[PASS] Root filesystem is strictly read-only"
else
  echo "[FAIL] ReadonlyRootfs is '$RO' (expected true). Pass --read-only"
  FAILED=1
fi

# 3. Check Tmpfs mounts
TMPFS_MOUNTS=$(docker inspect "$CONTAINER" --format '{{json .HostConfig.Tmpfs}}')
if [[ "$TMPFS_MOUNTS" == *"/var/run"* ]] && [[ "$TMPFS_MOUNTS" == *"/var/cache/nginx"* ]]; then
  echo "[PASS] Tmpfs paths configured correctly: $TMPFS_MOUNTS"
else
  echo "[FAIL] Missing required tmpfs mounts: $TMPFS_MOUNTS"
  FAILED=1
fi

# 4. HTTP response on port 8085
RESP=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8085 || echo "000")
if [ "$RESP" = "200" ]; then
  echo "[PASS] Secure Nginx responded with HTTP 200 on port 8085"
else
  echo "[FAIL] Failed to reach Nginx on port 8085 (HTTP code: $RESP)"
  FAILED=1
fi

exit $FAILED
