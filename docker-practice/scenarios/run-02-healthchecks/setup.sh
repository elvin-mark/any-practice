#!/usr/bin/env bash
set -e

docker rm -f monitored-web >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for run-02-healthchecks"
