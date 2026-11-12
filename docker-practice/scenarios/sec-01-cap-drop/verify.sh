#!/usr/bin/env bash
set -e

CONTAINER="least-priv-web"
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

# 2. Check CapDrop includes ALL
CAP_DROP=$(docker inspect "$CONTAINER" --format '{{json .HostConfig.CapDrop}}')
if [[ "$CAP_DROP" == *"ALL"* ]] || [[ "$CAP_DROP" == *"all"* ]]; then
  echo "[PASS] HostConfig.CapDrop has ALL capabilities dropped"
else
  echo "[FAIL] HostConfig.CapDrop is '$CAP_DROP' (expected [\"ALL\"])"
  FAILED=1
fi

# 3. Check CapAdd includes NET_BIND_SERVICE
CAP_ADD=$(docker inspect "$CONTAINER" --format '{{json .HostConfig.CapAdd}}')
if [[ "$CAP_ADD" == *"NET_BIND_SERVICE"* ]]; then
  echo "[PASS] HostConfig.CapAdd explicitly adds NET_BIND_SERVICE"
else
  echo "[FAIL] HostConfig.CapAdd is '$CAP_ADD' (missing NET_BIND_SERVICE)"
  FAILED=1
fi

# 4. HTTP response on port 8089
RESP=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8089 || echo "000")
if [ "$RESP" = "200" ]; then
  echo "[PASS] Container serves HTTP traffic successfully with dropped capabilities"
else
  echo "[FAIL] HTTP request failed on port 8089 (HTTP status: $RESP)"
  FAILED=1
fi

exit $FAILED
