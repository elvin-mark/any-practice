#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

$KUBECTL delete daemonset node-monitor -n default --ignore-not-found=true --wait=false
echo "Environment prepared for DaemonSet scenario."
