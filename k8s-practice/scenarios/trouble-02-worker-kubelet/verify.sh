#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# Check that all nodes in the cluster report Ready
NOT_READY_NODES=$($KUBECTL get nodes --no-headers | grep -iv "Ready" | grep -v "SchedulingDisabled" || true)

if [ -z "$NOT_READY_NODES" ]; then
  echo "[PASS] All nodes in the cluster report Ready status"
else
  echo "[FAIL] The following node(s) are not in Ready status:"
  echo "$NOT_READY_NODES"
  FAILED=1
fi

exit $FAILED
