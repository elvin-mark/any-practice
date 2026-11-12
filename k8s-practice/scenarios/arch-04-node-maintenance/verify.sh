#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check all nodes are schedulable (uncordoned)
SCHED_DISABLED=$($KUBECTL get nodes --no-headers | grep -i "SchedulingDisabled" || true)
if [ -z "$SCHED_DISABLED" ]; then
  echo "[PASS] All nodes are uncordoned and schedulable"
else
  echo "[FAIL] A node is still marked SchedulingDisabled: $SCHED_DISABLED"
  FAILED=1
fi

# 2. Check deployment replicas are ready
READY=$($KUBECTL get deployment maintenance-test -n default -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
if [ "$READY" -ge 2 ]; then
  echo "[PASS] Deployment 'maintenance-test' has $READY/2 ready replicas"
else
  echo "[FAIL] Deployment 'maintenance-test' ready replicas: $READY (expected 2)"
  FAILED=1
fi

exit $FAILED
