#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check if node has taint tier=critical:NoSchedule
TAINTED_NODES=$($KUBECTL get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.taints}{"\n"}{end}' | grep "tier" | grep "critical" | awk '{print $1}')
if [ -n "$TAINTED_NODES" ]; then
  echo "[PASS] Found node with taint tier=critical:NoSchedule: $TAINTED_NODES"
else
  echo "[FAIL] No node found with taint tier=critical:NoSchedule"
  FAILED=1
fi

# 2. Check pod existence
if ! $KUBECTL get pod toleration-pod -n default >/dev/null 2>&1; then
  echo "[FAIL] Pod 'toleration-pod' does not exist"
  exit 1
fi
echo "[PASS] Pod 'toleration-pod' exists"

# 3. Check toleration in pod spec
TOLERATIONS=$($KUBECTL get pod toleration-pod -n default -o jsonpath='{.spec.tolerations}')
if echo "$TOLERATIONS" | grep -q "tier" && echo "$TOLERATIONS" | grep -q "critical"; then
  echo "[PASS] Pod spec has toleration for key 'tier' and value 'critical'"
else
  echo "[FAIL] Pod spec missing matching toleration"
  FAILED=1
fi

# 4. Check pod status
PHASE=$($KUBECTL get pod toleration-pod -n default -o jsonpath='{.status.phase}')
if [ "$PHASE" = "Running" ]; then
  echo "[PASS] Pod 'toleration-pod' is in Running status"
else
  echo "[FAIL] Pod 'toleration-pod' status is: $PHASE"
  FAILED=1
fi

exit $FAILED
