#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete deployment ssd-app -n default --ignore-not-found=true --wait=false

# Remove existing disktype label from all nodes to ensure clean test
NODES=$($KUBECTL get nodes --no-headers | awk '{print $1}')
for n in $NODES; do
  $KUBECTL label node "$n" disktype- 2>/dev/null || true
done

echo "Cleaned disktype label and ssd-app deployment."
