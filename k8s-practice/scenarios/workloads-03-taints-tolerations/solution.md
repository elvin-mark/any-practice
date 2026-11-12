# Solution: Configure Taints and Pod Tolerations

### Step 1: Add Taint to the Worker Node
```bash
# Identify worker node
kubectl get nodes

kubectl taint nodes <worker-node-name> tier=critical:NoSchedule
```

### Step 2: Create Pod with Toleration
Create `toleration-pod.yaml`:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: toleration-pod
  namespace: default
spec:
  tolerations:
  - key: "tier"
    operator: "Equal"
    value: "critical"
    effect: "NoSchedule"
  containers:
  - name: nginx
    image: nginx:1.25-alpine
```

### Step 3: Apply and Verify
```bash
kubectl apply -f toleration-pod.yaml
kubectl get pod toleration-pod -o wide
```
Verify the pod is scheduled on the tainted node and transitions to `Running`.
