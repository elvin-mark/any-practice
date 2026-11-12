#!/usr/bin/env bash
set -e

docker rm -f secure-nginx >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for vol-02-readonly-tmpfs"
