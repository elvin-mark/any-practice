#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace web-apps --dry-run=client -o yaml | $KUBECTL apply -f -
$KUBECTL delete service front-web-svc -n web-apps --ignore-not-found=true

# Deploy sample app
$KUBECTL apply -n web-apps -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: front-web
  namespace: web-apps
spec:
  replicas: 2
  selector:
    matchLabels:
      app: front-web
  template:
    metadata:
      labels:
        app: front-web
    spec:
      containers:
      - name: web
        image: nginx:1.25-alpine
        ports:
        - containerPort: 80
EOF

$KUBECTL wait --for=condition=available deployment/front-web -n web-apps --timeout=45s 2>/dev/null || true
echo "Namespace 'web-apps' and deployment 'front-web' initialized."
