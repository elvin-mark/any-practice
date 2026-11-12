#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace security --dry-run=client -o yaml | $KUBECTL apply -f -

# Deploy app with label app=auth-backend
$KUBECTL apply -n security -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-backend
  namespace: security
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-backend
  template:
    metadata:
      labels:
        app: auth-backend
    spec:
      containers:
      - name: auth
        image: nginx:1.25-alpine
        ports:
        - containerPort: 80
---
# Service with intentional mismatched selector: app=auth-server
apiVersion: v1
kind: Service
metadata:
  name: auth-api
  namespace: security
spec:
  selector:
    app: auth-server
  ports:
  - port: 80
    targetPort: 80
EOF

$KUBECTL wait --for=condition=available deployment/auth-backend -n security --timeout=30s 2>/dev/null || true
echo "Setup complete with mismatched service selector in namespace security."
