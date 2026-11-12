#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete pod storage-pod -n default --ignore-not-found=true --wait=false

# Setup PV and PVC
$KUBECTL apply -f - <<EOF
apiVersion: v1
kind: PersistentVolume
metadata:
  name: store-pv
spec:
  capacity:
    storage: 500Mi
  accessModes:
    - ReadWriteOnce
  storageClassName: local-storage
  hostPath:
    path: /tmp/cka-store
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: storage-claim
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: local-storage
  resources:
    requests:
      storage: 500Mi
EOF

echo "Prepared storage-claim PVC in default namespace."
