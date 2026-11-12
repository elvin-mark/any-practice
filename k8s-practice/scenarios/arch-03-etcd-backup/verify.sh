#!/usr/bin/env bash

CP_CONTAINER=$(docker ps --filter "name=control-plane" --format "{{.Names}}" | head -n 1)
if [ -z "$CP_CONTAINER" ]; then
  echo "[FAIL] Control plane container not running in docker"
  exit 1
fi

FAILED=0

# 1. Check if snapshot file exists
if docker exec "$CP_CONTAINER" test -f /var/lib/etcd-backup/snapshot.db; then
  echo "[PASS] Snapshot file exists at /var/lib/etcd-backup/snapshot.db"
else
  echo "[FAIL] Snapshot file /var/lib/etcd-backup/snapshot.db does not exist"
  FAILED=1
fi

# 2. Check snapshot integrity via etcdctl snapshot status
STATUS_OUT=$(docker exec "$CP_CONTAINER" sh -c "ETCDCTL_API=3 etcdctl snapshot status /var/lib/etcd-backup/snapshot.db --write-out=table" 2>&1 || true)
if echo "$STATUS_OUT" | grep -iq "TOTAL KEYS"; then
  echo "[PASS] Snapshot database status verified and intact"
else
  echo "[FAIL] Snapshot verification failed: $STATUS_OUT"
  FAILED=1
fi

exit $FAILED
