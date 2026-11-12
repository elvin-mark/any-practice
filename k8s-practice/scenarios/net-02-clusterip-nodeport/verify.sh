#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check Service exists
if ! $KUBECTL get service front-web-svc -n web-apps >/dev/null 2>&1; then
  echo "[FAIL] Service 'front-web-svc' not found in namespace web-apps"
  exit 1
fi
echo "[PASS] Service 'front-web-svc' exists"

# 2. Check Service Type is NodePort
SVC_TYPE=$($KUBECTL get service front-web-svc -n web-apps -o jsonpath='{.spec.type}')
if [ "$SVC_TYPE" = "NodePort" ]; then
  echo "[PASS] Service type is NodePort"
else
  echo "[FAIL] Service type is $SVC_TYPE (expected NodePort)"
  FAILED=1
fi

# 3. Check nodePort is 30080
NODE_PORT=$($KUBECTL get service front-web-svc -n web-apps -o jsonpath='{.spec.ports[0].nodePort}')
if [ "$NODE_PORT" = "30080" ]; then
  echo "[PASS] Service nodePort is 30080"
else
  echo "[FAIL] Service nodePort is $NODE_PORT (expected 30080)"
  FAILED=1
fi

# 4. Check endpoints populated
EP_COUNT=$($KUBECTL get endpoints front-web-svc -n web-apps -o jsonpath='{.subsets[*].addresses[*].ip}' | wc -w)
if [ "$EP_COUNT" -ge 2 ]; then
  echo "[PASS] Service has $EP_COUNT active endpoints"
else
  echo "[FAIL] Service has $EP_COUNT endpoints (expected at least 2)"
  FAILED=1
fi

exit $FAILED
