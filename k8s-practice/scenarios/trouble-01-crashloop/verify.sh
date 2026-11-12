#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check deployment exists
if ! $KUBECTL get deployment order-service -n finance >/dev/null 2>&1; then
  echo "[FAIL] Deployment 'order-service' not found in namespace finance"
  exit 1
fi
echo "[PASS] Deployment 'order-service' exists"

# 2. Check ready replicas
READY=$($KUBECTL get deployment order-service -n finance -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
if [ "$READY" -ge 2 ]; then
  echo "[PASS] Deployment 'order-service' has $READY/2 ready replicas"
else
  echo "[FAIL] Deployment 'order-service' ready replicas: $READY (expected 2)"
  FAILED=1
fi

# 3. Check pods are Running without crash loop
CRASH_COUNT=$($KUBECTL get pods -n finance -l app=order-service --no-headers | grep -iE 'CrashLoopBackOff|Error' | wc -l)
if [ "$CRASH_COUNT" -eq 0 ]; then
  echo "[PASS] Zero pods in CrashLoopBackOff or Error state"
else
  echo "[FAIL] Found $CRASH_COUNT pods still crashing in namespace finance"
  FAILED=1
fi

exit $FAILED
