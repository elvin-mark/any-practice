#!/bin/bash
set -eo pipefail

TARGET="/usr/local/bin/backup-tool"

if [ ! -f "$TARGET" ]; then
  echo "❌ Error: $TARGET does not exist."
  exit 1
fi

# Check SUID bit
if [ -u "$TARGET" ]; then
  echo "❌ Error: $TARGET still has the SUID bit set."
  exit 1
fi

# Check world-writable
PERMS=$(stat -c "%a" "$TARGET")
# Check if last digit has write bit (2, 3, 6, 7)
if [ -w "$TARGET" ] && [ "$(stat -c "%A" "$TARGET" | cut -c 9)" = "w" ]; then
  echo "❌ Error: $TARGET is still world-writable ($PERMS)."
  exit 1
fi

echo "✅ Success: SUID bit revoked and file permissions secured properly ($PERMS)!"
exit 0
