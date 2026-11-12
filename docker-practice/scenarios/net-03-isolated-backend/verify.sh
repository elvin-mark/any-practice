#!/usr/bin/env bash
set -e

FAILED=0

# 1. Networks exist
for n in frontend-net backend-net; do
  if ! docker network inspect "$n" >/dev/null 2>&1; then
    echo "[FAIL] Network '$n' not found."
    exit 1
  fi
done
echo "[PASS] Networks 'frontend-net' and 'backend-net' exist"

# 2. Containers running
for c in secure-db external-client app-gateway; do
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
echo "[PASS] All three containers are running"

# 3. app-gateway can reach secure-db
if docker exec app-gateway ping -c 1 -W 2 secure-db >/dev/null 2>&1; then
  echo "[PASS] app-gateway successfully connected to secure-db"
else
  echo "[FAIL] app-gateway cannot reach secure-db"
  FAILED=1
fi

# 4. app-gateway can reach external-client
if docker exec app-gateway ping -c 1 -W 2 external-client >/dev/null 2>&1; then
  echo "[PASS] app-gateway successfully connected to external-client"
else
  echo "[FAIL] app-gateway cannot reach external-client"
  FAILED=1
fi

# 5. external-client MUST NOT reach secure-db
if docker exec external-client ping -c 1 -W 2 secure-db >/dev/null 2>&1; then
  echo "[FAIL] Security violation: external-client reached secure-db! Networks are not isolated."
  FAILED=1
else
  echo "[PASS] Network isolation verified: external-client cannot reach secure-db"
fi

exit $FAILED
