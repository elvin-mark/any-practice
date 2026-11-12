#!/usr/bin/env bash
set -e

IMAGE="hardened-app:latest"
FAILED=0

# 1. Image exists
if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
  echo "[FAIL] Image '$IMAGE' not found. Build it with 'docker build -t $IMAGE .'"
  exit 1
fi
echo "[PASS] Docker image '$IMAGE' exists"

# 2. Check Config.User in image
USER_VAL=$(docker image inspect "$IMAGE" --format '{{.Config.User}}')
if [ "$USER_VAL" = "10001" ] || [ "$USER_VAL" = "secuser" ] || [ "$USER_VAL" = "10001:10001" ]; then
  echo "[PASS] Image Config.User configured as '$USER_VAL'"
else
  echo "[FAIL] Image Config.User is '$USER_VAL' (expected 10001 or secuser)"
  FAILED=1
fi

# 3. Test execution and runtime UID
TEST_OUT=$(docker run --rm "$IMAGE" 2>&1 || true)
if echo "$TEST_OUT" | grep -q "Running as UID 10001"; then
  echo "[PASS] Container runtime executed with unprivileged UID 10001"
else
  echo "[FAIL] Container did not run as UID 10001: $TEST_OUT"
  FAILED=1
fi

# 4. Test write permissions in /data
if echo "$TEST_OUT" | grep -q "Writing timestamp..."; then
  echo "[PASS] User 10001 successfully wrote to /data directory"
else
  echo "[FAIL] User 10001 failed to write to /data directory (check chown)"
  FAILED=1
fi

exit $FAILED
