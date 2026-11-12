#!/usr/bin/env bash
set -e

NET="custom-bridge"
FAILED=0

# 1. Network exists
if ! docker network inspect "$NET" >/dev/null 2>&1; then
  echo "[FAIL] Network '$NET' not found."
  exit 1
fi
echo "[PASS] Network '$NET' exists"

# 2. Subnet check
SUBNET=$(docker network inspect "$NET" --format '{{(index .IPAM.Config 0).Subnet}}')
if [ "$SUBNET" = "172.28.0.0/16" ]; then
  echo "[PASS] Subnet is configured as 172.28.0.0/16"
else
  echo "[FAIL] Subnet is '$SUBNET' (expected 172.28.0.0/16)"
  FAILED=1
fi

# 3. Containers running
for c in api-server client-app; do
  if ! docker inspect "$c" >/dev/null 2>&1; then
    echo "[FAIL] Container '$c' not found."
    exit 1
  fi
  st=$(docker inspect "$c" --format '{{.State.Status}}')
  if [ "$st" != "running" ]; then
    echo "[FAIL] Container '$c' status is '$st' (expected running)"
    FAILED=1
  fi
done
echo "[PASS] Both api-server and client-app are running"

# 4. Test DNS resolution and HTTP connectivity from client-app to api-server
if docker exec client-app wget -q -O - http://api-server | grep -q "Welcome to nginx"; then
  echo "[PASS] DNS resolution and HTTP request from client-app to http://api-server succeeded"
else
  echo "[FAIL] client-app could not resolve or reach api-server via DNS"
  FAILED=1
fi

exit $FAILED
