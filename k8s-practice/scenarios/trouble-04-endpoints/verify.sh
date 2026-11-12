#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check Service exists
if ! $KUBECTL get service auth-api -n security >/dev/null 2>&1; then
  echo "[FAIL] Service 'auth-api' does not exist in namespace security"
  exit 1
fi
echo "[PASS] Service 'auth-api' exists"

# 2. Check Service selector
SELECTOR=$($KUBECTL get service auth-api -n security -o jsonpath='{.spec.selector.app}')
if [ "$SELECTOR" = "auth-backend" ]; then
  echo "[PASS] Service selector correctly points to 'auth-backend'"
else
  echo "[FAIL] Service selector is '$SELECTOR' (expected 'auth-backend')"
  FAILED=1
fi

# 3. Check endpoints are populated
EP_COUNT=$($KUBECTL get endpoints auth-api -n security -o jsonpath='{.subsets[*].addresses[*].ip}' | wc -w)
if [ "$EP_COUNT" -ge 2 ]; then
  echo "[PASS] Service 'auth-api' has $EP_COUNT active endpoints"
else
  echo "[FAIL] Service 'auth-api' has $EP_COUNT endpoints (expected at least 2)"
  FAILED=1
fi

exit $FAILED
