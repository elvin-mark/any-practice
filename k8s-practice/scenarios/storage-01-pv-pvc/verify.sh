#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check PV exists
if ! $KUBECTL get pv data-pv >/dev/null 2>&1; then
  echo "[FAIL] PersistentVolume 'data-pv' does not exist"
  exit 1
fi
echo "[PASS] PersistentVolume 'data-pv' exists"

# 2. Check PVC exists
if ! $KUBECTL get pvc data-pvc -n default >/dev/null 2>&1; then
  echo "[FAIL] PersistentVolumeClaim 'data-pvc' does not exist"
  exit 1
fi
echo "[PASS] PersistentVolumeClaim 'data-pvc' exists"

# 3. Check PVC status is Bound
PVC_STATUS=$($KUBECTL get pvc data-pvc -n default -o jsonpath='{.status.phase}')
if [ "$PVC_STATUS" = "Bound" ]; then
  echo "[PASS] PersistentVolumeClaim 'data-pvc' is Bound"
else
  echo "[FAIL] PersistentVolumeClaim 'data-pvc' status is: $PVC_STATUS (expected Bound)"
  FAILED=1
fi

# 4. Check PVC bound to data-pv
BOUND_VOLUME=$($KUBECTL get pvc data-pvc -n default -o jsonpath='{.spec.volumeName}')
if [ "$BOUND_VOLUME" = "data-pv" ]; then
  echo "[PASS] PVC is bound specifically to volume 'data-pv'"
else
  echo "[FAIL] PVC bound volume is '$BOUND_VOLUME' (expected 'data-pv')"
  FAILED=1
fi

exit $FAILED
