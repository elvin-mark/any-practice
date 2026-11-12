#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check Ingress exists
if ! $KUBECTL get ingress ecommerce-ingress -n ecommerce >/dev/null 2>&1; then
  echo "[FAIL] Ingress 'ecommerce-ingress' not found in namespace ecommerce"
  exit 1
fi
echo "[PASS] Ingress 'ecommerce-ingress' exists"

# 2. Check Host
HOST=$($KUBECTL get ingress ecommerce-ingress -n ecommerce -o jsonpath='{.spec.rules[0].host}')
if [ "$HOST" = "shop.example.com" ]; then
  echo "[PASS] Ingress host matches shop.example.com"
else
  echo "[FAIL] Ingress host is '$HOST' (expected 'shop.example.com')"
  FAILED=1
fi

# 3. Check /catalog path and backend service
CATALOG_SVC=$($KUBECTL get ingress ecommerce-ingress -n ecommerce -o jsonpath='{.spec.rules[0].http.paths[?(@.path=="/catalog")].backend.service.name}')
if [ "$CATALOG_SVC" = "catalog-svc" ]; then
  echo "[PASS] Path /catalog routes to service 'catalog-svc'"
else
  echo "[FAIL] Path /catalog backend service is '$CATALOG_SVC' (expected 'catalog-svc')"
  FAILED=1
fi

# 4. Check /checkout path and backend service
CHECKOUT_SVC=$($KUBECTL get ingress ecommerce-ingress -n ecommerce -o jsonpath='{.spec.rules[0].http.paths[?(@.path=="/checkout")].backend.service.name}')
if [ "$CHECKOUT_SVC" = "checkout-svc" ]; then
  echo "[PASS] Path /checkout routes to service 'checkout-svc'"
else
  echo "[FAIL] Path /checkout backend service is '$CHECKOUT_SVC' (expected 'checkout-svc')"
  FAILED=1
fi

exit $FAILED
