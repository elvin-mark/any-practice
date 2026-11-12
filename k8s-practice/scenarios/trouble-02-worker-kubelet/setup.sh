#!/usr/bin/env bash
set -e

WORKER_CONTAINER=$(docker ps --filter "name=worker" --format "{{.Names}}" 2>/dev/null | head -n 1 || true)
if [ -n "$WORKER_CONTAINER" ]; then
  docker exec "$WORKER_CONTAINER" systemctl stop kubelet 2>/dev/null || true
  echo "Stopped kubelet on worker node container: $WORKER_CONTAINER"
else
  echo "Note: If running rootless, ensure worker node kubelet can be managed."
fi
