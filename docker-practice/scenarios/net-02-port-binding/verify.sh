#!/usr/bin/env bash
set -e

CONTAINER="internal-cache"
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

# 2. Check HostIp binding
HOST_IP=$(docker inspect "$CONTAINER" --format '{{range $p, $conf := .NetworkSettings.Ports}}{{range $conf}}{{.HostIp}}{{end}}{{end}}')
HOST_PORT=$(docker inspect "$CONTAINER" --format '{{range $p, $conf := .NetworkSettings.Ports}}{{range $conf}}{{.HostPort}}{{end}}{{end}}')

if [ "$HOST_IP" = "127.0.0.1" ]; then
  echo "[PASS] Port is bound strictly to loopback interface 127.0.0.1"
else
  echo "[FAIL] Port HostIp is '$HOST_IP' (expected 127.0.0.1). Do not bind to 0.0.0.0!"
  FAILED=1
fi

if [ "$HOST_PORT" = "6379" ]; then
  echo "[PASS] HostPort is 6379"
else
  echo "[FAIL] HostPort is '$HOST_PORT' (expected 6379)"
  FAILED=1
fi

# 3. Test redis ping on loopback
PING_RES=$(docker run --rm --network host redis:alpine redis-cli -h 127.0.0.1 -p 6379 ping 2>/dev/null || echo "")
if [ "$PING_RES" = "PONG" ]; then
  echo "[PASS] Redis responded with PONG on 127.0.0.1:6379"
else
  echo "[FAIL] Redis ping failed on 127.0.0.1:6379: $PING_RES"
  FAILED=1
fi

exit $FAILED
