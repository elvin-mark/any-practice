#!/usr/bin/env bash
set -e

docker rm -f init-worker >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for run-03-init-zombies"
