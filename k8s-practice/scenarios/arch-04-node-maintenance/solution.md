# Solution: Safely Drain and Maintain a Cluster Node

### Step 1: Identify the Worker Node
```bash
kubectl get nodes
# Notice the worker node (e.g. dev-cluster-worker or cka-cluster-1-worker)
```

### Step 2: Drain the Node
```bash
kubectl drain <worker-node-name> --ignore-daemonsets --delete-emptydir-data --force
```

### Step 3: Check Node Status
```bash
kubectl get nodes
# The worker node should show STATUS: Ready,SchedulingDisabled
```

### Step 4: Uncordon the Node after Maintenance
```bash
kubectl uncordon <worker-node-name>
```

### Step 5: Verify Node is Schedulable
```bash
kubectl get nodes
# STATUS should return to Ready (without SchedulingDisabled)
```
