#!/usr/bin/env bash
set -e

docker rm -f datastore-app >/dev/null 2>&1 || true
docker volume rm -f db-store >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for vol-01-named-volume"
