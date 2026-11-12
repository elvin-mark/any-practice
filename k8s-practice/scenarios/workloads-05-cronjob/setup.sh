#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL create namespace batch-jobs --dry-run=client -o yaml | $KUBECTL apply -f -
$KUBECTL delete cronjob report-runner -n batch-jobs --ignore-not-found=true
echo "Namespace batch-jobs initialized."
