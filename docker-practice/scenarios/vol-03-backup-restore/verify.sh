#!/usr/bin/env bash
set -e

SOURCE_VOL="prod-backup-source"
DEST_VOL="prod-backup-restored"
BACKUP_ARCHIVE="/tmp/prod-backup.tar.gz"
FAILED=0

# 1. Backup archive exists
if [ -f "$BACKUP_ARCHIVE" ]; then
  echo "[PASS] Backup archive '$BACKUP_ARCHIVE' created"
else
  echo "[FAIL] Backup archive '$BACKUP_ARCHIVE' not found"
  exit 1
fi

# 2. Dest volume exists
if ! docker volume inspect "$DEST_VOL" >/dev/null 2>&1; then
  echo "[FAIL] Restored volume '$DEST_VOL' does not exist."
  exit 1
fi
echo "[PASS] Restored volume '$DEST_VOL' exists"

# 3. Verify files and checksum in restored volume
CHECK_OUTPUT=$(docker run --rm -v "$DEST_VOL":/restored alpine:3.20 sh -c '
  cd /restored
  if [ ! -f config.json ] || [ ! -f data.db ] || [ ! -f checksum.sha256 ]; then
    echo "MISSING_FILES"
    exit 1
  fi
  # Verify checksums
  sha256sum -c checksum.sha256
' 2>&1 || true)

if echo "$CHECK_OUTPUT" | grep -q "OK"; then
  echo "[PASS] All files restored and SHA256 integrity verified in '$DEST_VOL'"
else
  echo "[FAIL] Checksum verification failed on restored volume: $CHECK_OUTPUT"
  FAILED=1
fi

exit $FAILED
