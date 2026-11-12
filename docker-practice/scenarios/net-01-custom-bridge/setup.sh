#!/usr/bin/env bash
set -e

docker rm -f api-server client-app >/dev/null 2>&1 || true
docker network rm custom-bridge >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for net-01-custom-bridge"
