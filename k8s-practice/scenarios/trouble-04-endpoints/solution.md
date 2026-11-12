# Solution: Troubleshoot Service Missing Endpoints

### Step 1: Inspect Service Endpoints
```bash
kubectl get endpoints auth-api -n security
# Output shows: <none>
```

### Step 2: Compare Selectors and Pod Labels
```bash
# Check service selector:
kubectl describe svc auth-api -n security | grep Selector
# Output: Selector: app=auth-server

# Check actual pod labels:
kubectl get pods -n security --show-labels
# Output: LABELS: app=auth-backend
```
Notice the mismatch: `app=auth-server` vs `app=auth-backend`.

### Step 3: Update the Service Selector
```bash
kubectl set selector service auth-api -n security "app=auth-backend"
```
Or edit with `kubectl edit svc auth-api -n security` and change `app: auth-server` to `app: auth-backend`.

### Step 4: Verify Endpoints
```bash
kubectl get endpoints auth-api -n security
# Output now lists the pod IPs: e.g. 10.244.1.5:80, 10.244.1.6:80
```
