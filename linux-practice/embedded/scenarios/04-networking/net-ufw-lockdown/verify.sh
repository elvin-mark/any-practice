#!/bin/bash
set -eo pipefail

STATUS_OUTPUT=$(ufw status verbose)

if ! echo "$STATUS_OUTPUT" | grep -qi "Status: active"; then
  echo "❌ Error: UFW is not active."
  exit 1
fi

if ! echo "$STATUS_OUTPUT" | grep -qi "Default: deny (incoming)"; then
  echo "❌ Error: Default incoming policy is not set to DENY."
  exit 1
fi

if ! echo "$STATUS_OUTPUT" | grep -qi "allow (outgoing)"; then
  echo "❌ Error: Default outgoing policy is not set to ALLOW."
  exit 1
fi

# Check port 22
if ! echo "$STATUS_OUTPUT" | grep -qE "(22/tcp|22).*ALLOW.*IN"; then
  echo "❌ Error: Port 22/tcp (SSH) is not allowed incoming."
  exit 1
fi

# Check port 80
if ! echo "$STATUS_OUTPUT" | grep -qE "(80/tcp|80).*ALLOW.*IN"; then
  echo "❌ Error: Port 80/tcp (HTTP) is not allowed incoming."
  exit 1
fi

echo "✅ Success: UFW firewall is active with default deny and ports 22, 80 permitted!"
exit 0
