#!/usr/bin/env bash
set -e

# Find control-plane container name
CP_CONTAINER=$(docker ps --filter "name=control-plane" --format "{{.Names}}" | head -n 1)

if [ -n "$CP_CONTAINER" ]; then
  docker exec "$CP_CONTAINER" rm -f /etc/kubernetes/manifests/my-static-web.yaml 2>/dev/null || true
fi
echo "Cleaned any prior static pod manifests from control-plane node."
