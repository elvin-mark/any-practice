# Solution: Create a Batch CronJob with Limits

### Step 1: Create the CronJob YAML
Save to `report-runner.yaml`:
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: report-runner
  namespace: batch-jobs
spec:
  schedule: "*/15 * * * *"
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
  jobTemplate:
    spec:
      backoffLimit: 2
      activeDeadlineSeconds: 60
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: runner
            image: busybox:1.36
            command: ["sh", "-c", "echo 'Generating periodic usage report'; date"]
```

### Step 2: Apply the CronJob
```bash
kubectl apply -f report-runner.yaml
```

### Step 3: Verify the CronJob
```bash
kubectl get cronjob report-runner -n batch-jobs
```
