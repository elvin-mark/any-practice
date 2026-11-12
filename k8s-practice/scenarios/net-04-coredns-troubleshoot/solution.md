# Solution: Troubleshoot Cluster DNS Resolution Failure

### Step 1: Diagnose DNS Failure
Check DNS pods in `kube-system`:
```bash
kubectl get pods -n kube-system -l k8s-app=kube-dns
# Notice no pods running!
```

### Step 2: Inspect CoreDNS Deployment
```bash
kubectl get deployment coredns -n kube-system
# Notice: READY 0/0
```

### Step 3: Scale Deployment to 2 Replicas
```bash
kubectl scale deployment coredns -n kube-system --replicas=2
```

### Step 4: Verify CoreDNS Pods are Running
```bash
kubectl get pods -n kube-system -l k8s-app=kube-dns
```

### Step 5: Test DNS Resolution from Client Pod
```bash
kubectl exec dns-client -n default -- nslookup kubernetes.default
```
Output should successfully return the cluster IP of the Kubernetes service (e.g. `10.96.0.1`).
