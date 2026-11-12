#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check NetworkPolicy exists
if ! $KUBECTL get networkpolicy db-access-policy -n secure-zone >/dev/null 2>&1; then
  echo "[FAIL] NetworkPolicy 'db-access-policy' does not exist in namespace secure-zone"
  exit 1
fi
echo "[PASS] NetworkPolicy 'db-access-policy' exists"

# 2. Check podSelector targets role=db
TARGET_ROLE=$($KUBECTL get networkpolicy db-access-policy -n secure-zone -o jsonpath='{.spec.podSelector.matchLabels.role}')
if [ "$TARGET_ROLE" = "db" ]; then
  echo "[PASS] NetworkPolicy correctly targets pods with label role=db"
else
  echo "[FAIL] NetworkPolicy podSelector target is '$TARGET_ROLE' (expected 'db')"
  FAILED=1
fi

# 3. Check ingress from role=backend
FROM_ROLE=$($KUBECTL get networkpolicy db-access-policy -n secure-zone -o jsonpath='{.spec.ingress[*].from[*].podSelector.matchLabels.role}')
if echo "$FROM_ROLE" | grep -q "backend"; then
  echo "[PASS] Ingress rule allows traffic from podSelector role=backend"
else
  echo "[FAIL] Ingress rule does not specify podSelector role=backend"
  FAILED=1
fi

# 4. Check port 80 is allowed
PORTS=$($KUBECTL get networkpolicy db-access-policy -n secure-zone -o jsonpath='{.spec.ingress[*].ports[*].port}')
if echo "$PORTS" | grep -q "80"; then
  echo "[PASS] Port 80 is explicitly allowed in ingress rule"
else
  echo "[FAIL] Port 80 not found in ingress rule ports"
  FAILED=1
fi

# 5. Live traffic test if CNI supports enforcement
DB_IP=$($KUBECTL get pod db -n secure-zone -o jsonpath='{.status.podIP}' 2>/dev/null || true)
if [ -n "$DB_IP" ]; then
  # Test backend -> db (should succeed)
  if $KUBECTL exec backend -n secure-zone -- wget -qO- -T 3 "http://$DB_IP" >/dev/null 2>&1; then
    echo "[PASS] Live traffic: backend pod successfully reaches db"
  else
    echo "[FAIL] Live traffic: backend pod cannot reach db"
    FAILED=1
  fi
fi

exit $FAILED
