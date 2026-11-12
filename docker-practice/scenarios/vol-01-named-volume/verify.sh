#!/usr/bin/env bash
set -e

CONTAINER="datastore-app"
VOLUME="db-store"
FAILED=0

# 1. Volume exists
if ! docker volume inspect "$VOLUME" >/dev/null 2>&1; then
  echo "[FAIL] Named volume '$VOLUME' does not exist."
  exit 1
fi
echo "[PASS] Named volume '$VOLUME' exists"

# 2. Container running
if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
  echo "[FAIL] Container '$CONTAINER' not found."
  exit 1
fi
echo "[PASS] Container '$CONTAINER' is running"

# 3. Check mount destination and volume source
MOUNT_INFO=$(docker inspect "$CONTAINER" --format '{{range .Mounts}}{{if eq .Name "'$VOLUME'"}}{{.Destination}}{{end}}{{end}}')
if [ "$MOUNT_INFO" = "/var/lib/appdata" ]; then
  echo "[PASS] Volume '$VOLUME' is mounted to '/var/lib/appdata'"
else
  echo "[FAIL] Volume '$VOLUME' mount destination is '$MOUNT_INFO' (expected /var/lib/appdata)"
  FAILED=1
fi

# 4. Check file content
FILE_CONTENT=$(docker exec "$CONTAINER" cat /var/lib/appdata/persistent.txt 2>/dev/null || echo "")
if [ "$FILE_CONTENT" = "DOCKER_PERSISTENCE_VERIFIED" ]; then
  echo "[PASS] Persistent file contains correct verification token"
else
  echo "[FAIL] File /var/lib/appdata/persistent.txt content is '$FILE_CONTENT' (expected DOCKER_PERSISTENCE_VERIFIED)"
  FAILED=1
fi

exit $FAILED
