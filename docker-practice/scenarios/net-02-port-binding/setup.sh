#!/usr/bin/env bash
set -e

docker rm -f internal-cache >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for net-02-port-binding"
