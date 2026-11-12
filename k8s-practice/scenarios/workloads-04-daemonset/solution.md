# Solution: Deploy a Cluster-Wide DaemonSet with Tolerations

### Step 1: Create the DaemonSet YAML
Save to `node-monitor.yaml`:
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-monitor
  namespace: default
spec:
  selector:
    matchLabels:
      app: node-monitor
  template:
    metadata:
      labels:
        app: node-monitor
    spec:
      tolerations:
      - operator: "Exists"
      containers:
      - name: monitor
        image: busybox:1.36
        command: ["sh", "-c", "while true; do echo monitoring; sleep 30; done"]
```
*(Note: `operator: "Exists"` without a key tolerates all taints, including `node-role.kubernetes.io/control-plane:NoSchedule`)*.

### Step 2: Apply the Manifest
```bash
kubectl apply -f node-monitor.yaml
```

### Step 3: Verify Pod Scheduling
```bash
kubectl get daemonset node-monitor
kubectl get pods -o wide -l app=node-monitor
```
Ensure one pod is running on each node in the cluster.
