#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete pod toleration-pod -n default --ignore-not-found=true --wait=false

# Remove tier taint from all nodes if any
NODES=$($KUBECTL get nodes --no-headers | awk '{print $1}')
for n in $NODES; do
  $KUBECTL taint node "$n" tier:NoSchedule- 2>/dev/null || true
done

echo "Cleaned previous toleration-pod and taints."
