#!/usr/bin/env bash
set -e

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

# Ensure clean starting state
$KUBECTL create namespace development --dry-run=client -o yaml | $KUBECTL apply -f -
$KUBECTL delete role developer-role -n development --ignore-not-found=true
$KUBECTL delete rolebinding developer-binding -n development --ignore-not-found=true
echo "Namespace 'development' prepared. RBAC roles cleaned."
