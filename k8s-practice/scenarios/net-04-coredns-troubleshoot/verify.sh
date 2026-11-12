#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check CoreDNS replicas
READY_REPLICAS=$($KUBECTL get deployment coredns -n kube-system -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
if [ "$READY_REPLICAS" -ge 2 ]; then
  echo "[PASS] CoreDNS deployment has $READY_REPLICAS ready replicas"
else
  echo "[FAIL] CoreDNS ready replicas is $READY_REPLICAS (expected at least 2)"
  FAILED=1
fi

# 2. Wait up to 10s for dns-client pod if still starting
$KUBECTL wait --for=condition=ready pod/dns-client -n default --timeout=15s 2>/dev/null || true

# 3. Test nslookup resolution from dns-client pod
LOOKUP_OUT=$($KUBECTL exec dns-client -n default -- nslookup kubernetes.default 2>&1 || true)
if echo "$LOOKUP_OUT" | grep -iq "Address:" && ! echo "$LOOKUP_OUT" | grep -iq "can't resolve"; then
  echo "[PASS] Internal DNS lookup for 'kubernetes.default' succeeded"
else
  echo "[FAIL] DNS resolution failed: $LOOKUP_OUT"
  FAILED=1
fi

exit $FAILED
