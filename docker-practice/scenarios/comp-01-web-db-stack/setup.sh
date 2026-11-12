#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

# Down any existing compose stack
if [ -f docker-compose.yml ]; then
  docker compose down -v >/dev/null 2>&1 || true
  rm -f docker-compose.yml
fi

echo "[PASS] Environment prepared for comp-01-web-db-stack in $DIR"
