#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete pvc data-pvc -n default --ignore-not-found=true --wait=false
$KUBECTL delete pv data-pv --ignore-not-found=true --wait=false
echo "Storage environment cleaned."
