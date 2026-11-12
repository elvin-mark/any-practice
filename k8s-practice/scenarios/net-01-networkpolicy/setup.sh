#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace secure-zone --dry-run=client -o yaml | $KUBECTL apply -f -
$KUBECTL delete networkpolicy db-access-policy -n secure-zone --ignore-not-found=true

# Deploy test pods
$KUBECTL apply -n secure-zone -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: db
  labels:
    role: db
spec:
  containers:
  - name: nginx
    image: nginx:1.25-alpine
    ports:
    - containerPort: 80
---
apiVersion: v1
kind: Pod
metadata:
  name: backend
  labels:
    role: backend
spec:
  containers:
  - name: client
    image: busybox:1.36
    command: ["sh", "-c", "sleep 3600"]
---
apiVersion: v1
kind: Pod
metadata:
  name: frontend
  labels:
    role: frontend
spec:
  containers:
  - name: client
    image: busybox:1.36
    command: ["sh", "-c", "sleep 3600"]
EOF

echo "Waiting for test pods in 'secure-zone' to be ready..."
$KUBECTL wait --for=condition=ready pod -l role=db -n secure-zone --timeout=45s 2>/dev/null || true
$KUBECTL wait --for=condition=ready pod -l role=backend -n secure-zone --timeout=45s 2>/dev/null || true
$KUBECTL wait --for=condition=ready pod -l role=frontend -n secure-zone --timeout=45s 2>/dev/null || true
echo "Test environment in 'secure-zone' ready."
