#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check Role existence
if $KUBECTL get role developer-role -n development >/dev/null 2>&1; then
  echo "[PASS] Role 'developer-role' exists in namespace development"
else
  echo "[FAIL] Role 'developer-role' does not exist in namespace development"
  FAILED=1
fi

# 2. Check RoleBinding existence
if $KUBECTL get rolebinding developer-binding -n development >/dev/null 2>&1; then
  echo "[PASS] RoleBinding 'developer-binding' exists in namespace development"
else
  echo "[FAIL] RoleBinding 'developer-binding' does not exist in namespace development"
  FAILED=1
fi

# 3. Test permissions via auth can-i
CAN_GET_PODS=$($KUBECTL auth can-i get pods -n development --as=developer 2>/dev/null || echo "no")
if [ "$CAN_GET_PODS" = "yes" ]; then
  echo "[PASS] User 'developer' can get pods in development"
else
  echo "[FAIL] User 'developer' cannot get pods in development"
  FAILED=1
fi

CAN_CREATE_SVC=$($KUBECTL auth can-i create services -n development --as=developer 2>/dev/null || echo "no")
if [ "$CAN_CREATE_SVC" = "yes" ]; then
  echo "[PASS] User 'developer' can create services in development"
else
  echo "[FAIL] User 'developer' cannot create services in development"
  FAILED=1
fi

CAN_KUBE_SYSTEM=$($KUBECTL auth can-i list pods -n kube-system --as=developer 2>/dev/null || true)
if [ "$CAN_KUBE_SYSTEM" = "no" ]; then
  echo "[PASS] User 'developer' correctly cannot list pods in kube-system"
else
  echo "[FAIL] User 'developer' has unauthorized access to kube-system pods"
  FAILED=1
fi

exit $FAILED
