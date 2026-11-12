# Solution: Troubleshoot Application Stuck in CrashLoopBackOff

### Step 1: Check Pod Status
```bash
kubectl get pods -n finance
# Notice STATUS: CrashLoopBackOff
```

### Step 2: Read Container Logs
```bash
kubectl logs -n finance -l app=order-service --tail=20
# Output shows: Error: configuration file /etc/config/app.json missing!
```

### Step 3: Edit the Deployment to Fix the Command
```bash
kubectl edit deployment order-service -n finance
```
Replace the crashing `command:` in `spec.template.spec.containers[0]` with a working command, e.g.:
```yaml
command: ["sh", "-c", "echo 'Order service started successfully'; sleep 3600"]
```
Alternatively, apply an updated manifest:
```bash
kubectl set image deployment/order-service app=nginx:1.25-alpine -n finance
kubectl set command deployment/order-service -n finance --
```

### Step 4: Verify All Replicas are Ready
```bash
kubectl rollout status deployment/order-service -n finance
kubectl get pods -n finance
```
