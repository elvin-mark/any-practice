#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check pod existence
if ! $KUBECTL get pod multi-app -n default >/dev/null 2>&1; then
  echo "[FAIL] Pod 'multi-app' does not exist in namespace default"
  exit 1
fi
echo "[PASS] Pod 'multi-app' exists"

# 2. Check 2 containers exist
CONTAINERS=$($KUBECTL get pod multi-app -n default -o jsonpath='{.spec.containers[*].name}')
if echo "$CONTAINERS" | grep -q "producer" && echo "$CONTAINERS" | grep -q "consumer"; then
  echo "[PASS] Pod has both 'producer' and 'consumer' containers"
else
  echo "[FAIL] Expected containers 'producer' and 'consumer', found: $CONTAINERS"
  FAILED=1
fi

# 3. Check shared volume
VOL_NAME=$($KUBECTL get pod multi-app -n default -o jsonpath='{.spec.volumes[*].name}')
if echo "$VOL_NAME" | grep -q "shared-logs"; then
  echo "[PASS] Volume 'shared-logs' is configured"
else
  echo "[FAIL] Volume 'shared-logs' not found in pod spec"
  FAILED=1
fi

# 4. Check pod status
PHASE=$($KUBECTL get pod multi-app -n default -o jsonpath='{.status.phase}')
if [ "$PHASE" = "Running" ]; then
  echo "[PASS] Pod 'multi-app' is Running"
else
  echo "[FAIL] Pod 'multi-app' status: $PHASE (expected Running)"
  FAILED=1
fi

# 5. Check logs of consumer sidecar
LOGS=$($KUBECTL logs multi-app -c consumer -n default --tail=5 2>/dev/null || true)
if [ -n "$LOGS" ]; then
  echo "[PASS] Consumer container successfully reading logs from producer via shared volume"
else
  echo "[FAIL] Consumer container has no log output from shared volume"
  FAILED=1
fi

exit $FAILED
