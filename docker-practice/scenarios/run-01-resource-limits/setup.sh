#!/usr/bin/env bash
set -e

# Remove any pre-existing container with this name
docker rm -f prod-worker >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for run-01-resource-limits"
