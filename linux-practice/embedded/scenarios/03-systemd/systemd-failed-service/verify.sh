#!/bin/bash
set -eo pipefail

if ! systemctl is-active --quiet data-worker.service; then
  echo "❌ Error: data-worker.service is not active/running."
  echo "Current status:"
  systemctl status data-worker.service --no-pager || true
  exit 1
fi

echo "✅ Success: data-worker.service is active and healthy!"
exit 0
