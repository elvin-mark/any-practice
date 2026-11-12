#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check DaemonSet exists
if ! $KUBECTL get daemonset node-monitor -n default >/dev/null 2>&1; then
  echo "[FAIL] DaemonSet 'node-monitor' does not exist in default namespace"
  exit 1
fi
echo "[PASS] DaemonSet 'node-monitor' exists"

# 2. Check desired vs ready pods matches total node count
TOTAL_NODES=$($KUBECTL get nodes --no-headers | wc -l)
DESIRED=$($KUBECTL get daemonset node-monitor -n default -o jsonpath='{.status.desiredNumberScheduled}')
READY=$($KUBECTL get daemonset node-monitor -n default -o jsonpath='{.status.numberReady}')

if [ "$DESIRED" -eq "$TOTAL_NODES" ]; then
  echo "[PASS] DaemonSet scheduled on all $TOTAL_NODES cluster nodes (including control-plane)"
else
  echo "[FAIL] DaemonSet desiredNumberScheduled is $DESIRED, but cluster has $TOTAL_NODES nodes"
  FAILED=1
fi

if [ "$READY" -eq "$TOTAL_NODES" ]; then
  echo "[PASS] All $READY DaemonSet pods are in Ready status"
else
  echo "[FAIL] Ready pods count ($READY) does not match total nodes ($TOTAL_NODES)"
  FAILED=1
fi

exit $FAILED
