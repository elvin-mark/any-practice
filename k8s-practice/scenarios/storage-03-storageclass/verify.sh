#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check SC exists
if ! $KUBECTL get storageclass fast-retain >/dev/null 2>&1; then
  echo "[FAIL] StorageClass 'fast-retain' does not exist"
  exit 1
fi
echo "[PASS] StorageClass 'fast-retain' exists"

# 2. Check provisioner
PROVISIONER=$($KUBECTL get storageclass fast-retain -o jsonpath='{.provisioner}')
if [ "$PROVISIONER" = "rancher.io/local-path" ]; then
  echo "[PASS] Provisioner is rancher.io/local-path"
else
  echo "[FAIL] Provisioner is '$PROVISIONER' (expected rancher.io/local-path)"
  FAILED=1
fi

# 3. Check reclaimPolicy
RECLAIM=$($KUBECTL get storageclass fast-retain -o jsonpath='{.reclaimPolicy}')
if [ "$RECLAIM" = "Retain" ]; then
  echo "[PASS] ReclaimPolicy is Retain"
else
  echo "[FAIL] ReclaimPolicy is '$RECLAIM' (expected Retain)"
  FAILED=1
fi

# 4. Check volumeBindingMode
BINDING=$($KUBECTL get storageclass fast-retain -o jsonpath='{.volumeBindingMode}')
if [ "$BINDING" = "WaitForFirstConsumer" ]; then
  echo "[PASS] VolumeBindingMode is WaitForFirstConsumer"
else
  echo "[FAIL] VolumeBindingMode is '$BINDING' (expected WaitForFirstConsumer)"
  FAILED=1
fi

# 5. Check allowVolumeExpansion
EXPANSION=$($KUBECTL get storageclass fast-retain -o jsonpath='{.allowVolumeExpansion}')
if [ "$EXPANSION" = "true" ]; then
  echo "[PASS] AllowVolumeExpansion is true"
else
  echo "[FAIL] AllowVolumeExpansion is '$EXPANSION' (expected true)"
  FAILED=1
fi

exit $FAILED
