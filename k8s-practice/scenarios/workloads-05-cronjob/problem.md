# Task: Create a Batch CronJob with Limits

Create a recurring batch CronJob in namespace `batch-jobs`:

1. Ensure the namespace `batch-jobs` exists.
2. Create a `CronJob` named `report-runner` with:
   - Schedule: `*/15 * * * *` (runs every 15 minutes)
   - Image: `busybox:1.36`
   - Command: `sh -c 'echo "Generating periodic usage report"; date'`
   - Successful Jobs History Limit: `3`
   - Failed Jobs History Limit: `1`
   - Job Backoff Limit: `2`
   - Active Deadline Seconds: `60`
3. Verify that the CronJob is registered with the correct schedule.
