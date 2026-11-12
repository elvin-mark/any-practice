#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete storageclass fast-retain --ignore-not-found=true
echo "Cleaned existing fast-retain StorageClass."
