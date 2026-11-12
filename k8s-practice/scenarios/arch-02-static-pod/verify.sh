#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check if control-plane has the manifest file
CP_CONTAINER=$(docker ps --filter "name=control-plane" --format "{{.Names}}" | head -n 1)
if [ -z "$CP_CONTAINER" ]; then
  echo "[FAIL] Control plane container not running in docker"
  exit 1
fi

if docker exec "$CP_CONTAINER" test -f /etc/kubernetes/manifests/my-static-web.yaml; then
  echo "[PASS] Manifest /etc/kubernetes/manifests/my-static-web.yaml exists on control-plane"
else
  echo "[FAIL] Manifest /etc/kubernetes/manifests/my-static-web.yaml not found on control-plane"
  FAILED=1
fi

# 2. Check if Pod exists and is Running
POD_NAME=$($KUBECTL get pods -n default -o jsonpath='{.items[*].metadata.name}' | tr ' ' '\n' | grep '^my-static-web' | head -n 1 || true)

if [ -n "$POD_NAME" ]; then
  echo "[PASS] Pod '$POD_NAME' discovered via kubectl"
  STATUS=$($KUBECTL get pod "$POD_NAME" -n default -o jsonpath='{.status.phase}' 2>/dev/null || echo "Unknown")
  if [ "$STATUS" = "Running" ]; then
    echo "[PASS] Pod '$POD_NAME' is in Running status"
  else
    echo "[FAIL] Pod '$POD_NAME' is in state: $STATUS (expected Running)"
    FAILED=1
  fi
else
  echo "[FAIL] No pod starting with name 'my-static-web' found in default namespace"
  FAILED=1
fi

exit $FAILED
