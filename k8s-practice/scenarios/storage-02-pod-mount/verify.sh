#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check pod exists
if ! $KUBECTL get pod storage-pod -n default >/dev/null 2>&1; then
  echo "[FAIL] Pod 'storage-pod' does not exist in namespace default"
  exit 1
fi
echo "[PASS] Pod 'storage-pod' exists"

# 2. Check PVC mount in pod spec
PVC_USED=$($KUBECTL get pod storage-pod -n default -o jsonpath='{.spec.volumes[*].persistentVolumeClaim.claimName}')
if echo "$PVC_USED" | grep -q "storage-claim"; then
  echo "[PASS] Pod references PersistentVolumeClaim 'storage-claim'"
else
  echo "[FAIL] Pod does not reference 'storage-claim' PVC (found: '$PVC_USED')"
  FAILED=1
fi

MOUNT_PATH=$($KUBECTL get pod storage-pod -n default -o jsonpath='{.spec.containers[0].volumeMounts[*].mountPath}')
if echo "$MOUNT_PATH" | grep -q "/opt/store"; then
  echo "[PASS] Volume mounted at '/opt/store'"
else
  echo "[FAIL] Volume not mounted at '/opt/store' (found: '$MOUNT_PATH')"
  FAILED=1
fi

# 3. Check status is Running
PHASE=$($KUBECTL get pod storage-pod -n default -o jsonpath='{.status.phase}')
if [ "$PHASE" = "Running" ]; then
  echo "[PASS] Pod 'storage-pod' is in Running status"
else
  echo "[FAIL] Pod 'storage-pod' status: $PHASE"
  FAILED=1
fi

# 4. Check file contents inside mount
CONTENT=$($KUBECTL exec storage-pod -n default -- cat /opt/store/status.txt 2>/dev/null || true)
if echo "$CONTENT" | grep -q "cka-practice-success"; then
  echo "[PASS] Verified data written to persistent storage: $CONTENT"
else
  echo "[FAIL] Expected data not found in /opt/store/status.txt (read: '$CONTENT')"
  FAILED=1
fi

exit $FAILED
