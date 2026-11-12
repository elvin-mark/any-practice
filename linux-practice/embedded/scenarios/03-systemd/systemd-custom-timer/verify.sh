#!/bin/bash
set -eo pipefail

TIMER_FILE="/etc/systemd/system/backup-report.timer"

if [ ! -f "$TIMER_FILE" ]; then
  echo "❌ Error: $TIMER_FILE was not found."
  exit 1
fi

if ! systemctl is-enabled --quiet backup-report.timer 2>/dev/null; then
  echo "❌ Error: backup-report.timer is not enabled."
  exit 1
fi

if ! systemctl is-active --quiet backup-report.timer 2>/dev/null; then
  echo "❌ Error: backup-report.timer is not active (running)."
  exit 1
fi

# Verify OnCalendar directive exists
if ! grep -qi "OnCalendar=" "$TIMER_FILE"; then
  echo "❌ Error: $TIMER_FILE is missing an OnCalendar directive."
  exit 1
fi

echo "✅ Success: backup-report.timer is enabled, active, and scheduled!"
exit 0
