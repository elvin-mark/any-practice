#!/bin/bash
set -eo pipefail

SUDO_FILE="/etc/sudoers.d/devops"

if [ ! -f "$SUDO_FILE" ]; then
  echo "❌ Error: $SUDO_FILE does not exist."
  exit 1
fi

if ! visudo -cf "$SUDO_FILE" > /dev/null 2>&1; then
  echo "❌ Error: $SUDO_FILE still contains syntax errors:"
  visudo -cf "$SUDO_FILE" || true
  exit 1
fi

# Check that NOPASSWD is used and group is devops
if ! grep -qE '^%devops.*NOPASSWD:' "$SUDO_FILE"; then
  echo "❌ Error: $SUDO_FILE must grant %devops passwordless sudo (NOPASSWD)."
  exit 1
fi

echo "✅ Success: $SUDO_FILE parsed cleanly and grants valid sudo rights to %devops!"
exit 0
