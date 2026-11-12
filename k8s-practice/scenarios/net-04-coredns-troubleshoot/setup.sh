#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

# Deploy client pod
$KUBECTL apply -n default -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: dns-client
  namespace: default
spec:
  containers:
  - name: client
    image: busybox:1.36
    command: ["sh", "-c", "sleep 3600"]
EOF

# Simulate DNS failure by scaling down coredns deployment to 0
$KUBECTL scale deployment coredns -n kube-system --replicas=0

echo "Simulated DNS failure: coredns deployment scaled down to 0."
