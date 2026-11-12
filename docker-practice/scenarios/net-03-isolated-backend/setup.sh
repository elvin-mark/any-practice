#!/usr/bin/env bash
set -e

docker rm -f secure-db external-client app-gateway >/dev/null 2>&1 || true
docker network rm frontend-net backend-net >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for net-03-isolated-backend"
