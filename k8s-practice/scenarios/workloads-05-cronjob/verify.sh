#!/usr/bin/env bash

KUBECTL="kubectl"
if [ -n "$TARGET_CONTEXT" ]; then
  KUBECTL="kubectl --context $TARGET_CONTEXT"
fi

FAILED=0

# 1. Check CronJob exists
if ! $KUBECTL get cronjob report-runner -n batch-jobs >/dev/null 2>&1; then
  echo "[FAIL] CronJob 'report-runner' not found in namespace batch-jobs"
  exit 1
fi
echo "[PASS] CronJob 'report-runner' exists"

# 2. Check schedule
SCHEDULE=$($KUBECTL get cronjob report-runner -n batch-jobs -o jsonpath='{.spec.schedule}')
if [ "$SCHEDULE" = "*/15 * * * *" ]; then
  echo "[PASS] CronJob schedule is '*/15 * * * *'"
else
  echo "[FAIL] CronJob schedule is '$SCHEDULE' (expected '*/15 * * * *')"
  FAILED=1
fi

# 3. Check history limits
SUCCESS_LIMIT=$($KUBECTL get cronjob report-runner -n batch-jobs -o jsonpath='{.spec.successfulJobsHistoryLimit}')
FAILED_LIMIT=$($KUBECTL get cronjob report-runner -n batch-jobs -o jsonpath='{.spec.failedJobsHistoryLimit}')

if [ "$SUCCESS_LIMIT" = "3" ]; then
  echo "[PASS] successfulJobsHistoryLimit is 3"
else
  echo "[FAIL] successfulJobsHistoryLimit is '$SUCCESS_LIMIT' (expected 3)"
  FAILED=1
fi

if [ "$FAILED_LIMIT" = "1" ]; then
  echo "[PASS] failedJobsHistoryLimit is 1"
else
  echo "[FAIL] failedJobsHistoryLimit is '$FAILED_LIMIT' (expected 1)"
  FAILED=1
fi

# 4. Check backoffLimit
BACKOFF=$($KUBECTL get cronjob report-runner -n batch-jobs -o jsonpath='{.spec.jobTemplate.spec.backoffLimit}')
if [ "$BACKOFF" = "2" ]; then
  echo "[PASS] backoffLimit is 2"
else
  echo "[FAIL] backoffLimit is '$BACKOFF' (expected 2)"
  FAILED=1
fi

exit $FAILED
