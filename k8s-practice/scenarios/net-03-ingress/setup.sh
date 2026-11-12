#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace ecommerce --dry-run=client -o yaml | $KUBECTL apply -f -
$KUBECTL delete ingress ecommerce-ingress -n ecommerce --ignore-not-found=true

# Deploy dummy backend services
$KUBECTL apply -n ecommerce -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: catalog-svc
spec:
  ports:
  - port: 80
    targetPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: checkout-svc
spec:
  ports:
  - port: 80
    targetPort: 80
EOF

echo "Namespace 'ecommerce' and services initialized."
