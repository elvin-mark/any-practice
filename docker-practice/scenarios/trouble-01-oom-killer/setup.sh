#!/usr/bin/env bash
set -e

docker rm -f mem-hungry >/dev/null 2>&1 || true

# Start intentionally constrained container that immediately OOMs
docker run -d \
  --name mem-hungry \
  --memory=16m \
  --memory-swap=16m \
  --label practice=docker \
  alpine:3.20 sh -c 'python3 -c "x = [0] * (50 * 1024 * 1024)" 2>/dev/null || head -c 50M /dev/zero | tail || dd if=/dev/zero of=/dev/shm/fill bs=1M count=30' || true

sleep 2

echo "[PASS] Intentionally flawed OOM container started for trouble-01-oom-killer"
