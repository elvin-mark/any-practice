#!/usr/bin/env bash
set -e

CONTAINER="prod-worker"
FAILED=0

# 1. Container exists and is running
if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[FAIL] Container '$CONTAINER' does not exist."
  exit 1
fi

STATE=$(docker inspect "$CONTAINER" --format '{{.State.Status}}')
if [ "$STATE" = "running" ]; then
  echo "[PASS] Container '$CONTAINER' is running"
else
  echo "[FAIL] Container '$CONTAINER' status is '$STATE' (expected running)"
  FAILED=1
fi

# 2. Check Memory (256MB = 268435456 bytes)
MEM=$(docker inspect "$CONTAINER" --format '{{.HostConfig.Memory}}')
if [ "$MEM" -eq 268435456 ]; then
  echo "[PASS] Memory limit is configured to 256MB (268435456 bytes)"
else
  echo "[FAIL] Memory limit is $MEM (expected 268435456 bytes / 256m)"
  FAILED=1
fi

# 3. Check CPU limit (0.5 CPUs = NanoCpus: 500000000)
CPUS=$(docker inspect "$CONTAINER" --format '{{.HostConfig.NanoCpus}}')
if [ "$CPUS" -eq 500000000 ]; then
  echo "[PASS] CPU limit is configured to 0.5 CPUs"
else
  echo "[FAIL] CPU limit NanoCpus is $CPUS (expected 500000000)"
  FAILED=1
fi

# 4. Check PIDs limit
PIDS=$(docker inspect "$CONTAINER" --format '{{.HostConfig.PidsLimit}}')
if [ "$PIDS" -eq 100 ]; then
  echo "[PASS] PidsLimit is configured to 100"
else
  echo "[FAIL] PidsLimit is $PIDS (expected 100)"
  FAILED=1
fi

# 5. Check Restart Policy
RESTART=$(docker inspect "$CONTAINER" --format '{{.HostConfig.RestartPolicy.Name}}')
if [ "$RESTART" = "unless-stopped" ]; then
  echo "[PASS] Restart policy is 'unless-stopped'"
else
  echo "[FAIL] Restart policy is '$RESTART' (expected unless-stopped)"
  FAILED=1
fi

exit $FAILED
