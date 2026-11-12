#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check if at least one node has label disktype=ssd
LABELED_NODES=$($KUBECTL get nodes -l disktype=ssd --no-headers | awk '{print $1}')
if [ -n "$LABELED_NODES" ]; then
  echo "[PASS] Found node(s) with label disktype=ssd: $LABELED_NODES"
else
  echo "[FAIL] No node found with label disktype=ssd"
  FAILED=1
fi

# 2. Check deployment exists
if ! $KUBECTL get deployment ssd-app -n default >/dev/null 2>&1; then
  echo "[FAIL] Deployment 'ssd-app' does not exist"
  exit 1
fi
echo "[PASS] Deployment 'ssd-app' exists"

# 3. Check nodeAffinity in deployment spec
AFFINITY=$($KUBECTL get deployment ssd-app -n default -o jsonpath='{.spec.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[*].matchExpressions[*].key}')
if echo "$AFFINITY" | grep -q "disktype"; then
  echo "[PASS] Deployment has required NodeAffinity for 'disktype'"
else
  echo "[FAIL] Deployment missing required NodeAffinity for 'disktype'"
  FAILED=1
fi

# 4. Check all pods are running on the labeled node
POD_NODES=$($KUBECTL get pods -l app=ssd-app -n default -o jsonpath='{.items[*].spec.nodeName}')
for pn in $POD_NODES; do
  if echo "$LABELED_NODES" | grep -q "$pn"; then
    echo "[PASS] Pod scheduled on correct labeled node: $pn"
  else
    echo "[FAIL] Pod scheduled on non-labeled node: $pn"
    FAILED=1
  fi
done

exit $FAILED
