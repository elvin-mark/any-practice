#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace finance --dry-run=client -o yaml | $KUBECTL apply -f -

# Deploy broken app exiting with error
$KUBECTL apply -n finance -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: finance
spec:
  replicas: 2
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: app
        image: busybox:1.36
        command: ["sh", "-c", "echo 'Error: configuration file /etc/config/app.json missing!' >&2; exit 1"]
EOF

echo "Deployed broken order-service in namespace finance."
