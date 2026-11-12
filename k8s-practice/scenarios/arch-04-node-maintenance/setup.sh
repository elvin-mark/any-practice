#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

# Ensure all nodes are uncordoned first
NODES=$($KUBECTL get nodes --no-headers | awk '{print $1}')
for n in $NODES; do
  $KUBECTL uncordon "$n" 2>/dev/null || true
done

# Create test deployment
$KUBECTL apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: maintenance-test
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: maintenance-test
  template:
    metadata:
      labels:
        app: maintenance-test
    spec:
      containers:
      - name: web
        image: nginx:1.25-alpine
EOF

echo "Deployed test workload and ensured nodes are in uncordoned state."
