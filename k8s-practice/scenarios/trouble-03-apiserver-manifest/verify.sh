#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Find infra-agent pod in kube-system
POD_NAME=$($KUBECTL get pods -n kube-system -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | grep '^infra-agent' | head -n 1 || true)

if [ -z "$POD_NAME" ]; then
  echo "[FAIL] Pod starting with name 'infra-agent' not found in namespace kube-system"
  exit 1
fi
echo "[PASS] Found pod '$POD_NAME'"

# 2. Check pod status is Running
PHASE=$($KUBECTL get pod "$POD_NAME" -n kube-system -o jsonpath='{.status.phase}' 2>/dev/null || echo "Unknown")
if [ "$PHASE" = "Running" ]; then
  echo "[PASS] Pod '$POD_NAME' is Running"
else
  echo "[FAIL] Pod '$POD_NAME' status is: $PHASE (expected Running)"
  FAILED=1
fi

exit $FAILED
