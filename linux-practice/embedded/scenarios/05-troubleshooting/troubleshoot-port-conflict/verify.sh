#!/bin/bash
set -eo pipefail

if ! systemctl is-active --quiet web-app.service; then
  echo "❌ Error: web-app.service is not active."
  echo "Current status:"
  systemctl status web-app.service --no-pager || true
  exit 1
fi

if ! curl -sf http://127.0.0.1:80 | grep -q "Production Web Service Active"; then
  echo "❌ Error: Port 80 is listening, but not serving the expected web-app.service response."
  exit 1
fi

echo "✅ Success: Port 80 conflict resolved and web-app.service is serving traffic!"
exit 0
