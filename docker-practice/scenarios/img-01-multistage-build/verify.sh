#!/usr/bin/env bash
set -e

IMAGE="go-microservice:optimized"
FAILED=0

# 1. Check if image exists
if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
  echo "[FAIL] Docker image '$IMAGE' not found. Did you run 'docker build -t $IMAGE .'?"
  exit 1
fi
echo "[PASS] Docker image '$IMAGE' exists"

# 2. Check image size < 35MB (36700160 bytes)
SIZE_BYTES=$(docker image inspect "$IMAGE" --format '{{.Size}}')
SIZE_MB=$((SIZE_BYTES / 1024 / 1024))
if [ "$SIZE_BYTES" -le 36700160 ]; then
  echo "[PASS] Image size is ${SIZE_MB}MB (< 35MB requirement)"
else
  echo "[FAIL] Image size is ${SIZE_MB}MB (exceeds 35MB limit). Use multi-stage build!"
  FAILED=1
fi

# 3. Test running the container
CONTAINER_NAME="test-go-ms-$RANDOM"
docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

CONTAINER_ID=$(docker run -d --rm --name "$CONTAINER_NAME" -p 18080:8080 "$IMAGE" 2>/dev/null || true)
if [ -z "$CONTAINER_ID" ]; then
  echo "[FAIL] Failed to start container from '$IMAGE'"
  exit 1
fi

# Cleanup container on exit
trap 'docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true' EXIT

sleep 2

# 4. Check HTTP response
RESP=$(curl -s http://127.0.0.1:18080/health || echo "FAILED")
if [[ "$RESP" == *"OK"* ]]; then
  echo "[PASS] Microservice responds to HTTP health check"
else
  echo "[FAIL] Microservice failed to respond on port 8080: $RESP"
  FAILED=1
fi

# 5. Check non-root user
UID_VAL=$(docker exec "$CONTAINER_NAME" id -u 2>/dev/null || echo "0")
if [ "$UID_VAL" != "0" ]; then
  echo "[PASS] Container executes as unprivileged UID $UID_VAL"
else
  echo "[FAIL] Container is running as root (UID 0). Use 'USER 10001' in Dockerfile."
  FAILED=1
fi

exit $FAILED
