#!/usr/bin/env bash
set -e

docker rm -f no-escalate >/dev/null 2>&1 || true

echo "[PASS] Environment prepared for sec-02-no-new-privs"
