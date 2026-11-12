#!/usr/bin/env bash
set -e

docker rm -f least-priv-web >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for sec-01-cap-drop"
