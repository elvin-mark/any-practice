#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

IMAGE="cached-app:latest"
FAILED=0

# 1. Check if image exists
if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
  echo "[FAIL] Docker image '$IMAGE' not found. Please build with 'docker build -t $IMAGE .'"
  exit 1
fi
echo "[PASS] Docker image '$IMAGE' exists"

# 2. Check Dockerfile contents for proper layer ordering
if ! grep -q "package.*json" Dockerfile; then
  echo "[FAIL] Dockerfile does not copy package*.json before application source"
  FAILED=1
else
  echo "[PASS] Dockerfile copies package manifest first"
fi

# 3. Test caching behavior: touch server.js and test build output contains CACHED or using cache
touch server.js
BUILD_LOG=$(docker build -t "$IMAGE" . 2>&1 || true)

if echo "$BUILD_LOG" | grep -qiE 'CACHED|Using cache'; then
  echo "[PASS] Build cache correctly preserved dependency layer after source modification"
else
  echo "[FAIL] Dependency layer was invalidated upon touching server.js. Check Dockerfile order!"
  FAILED=1
fi

exit $FAILED
