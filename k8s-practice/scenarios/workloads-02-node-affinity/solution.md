# Solution: Schedule Workload Using Node Affinity

### Step 1: Label the Target Node
```bash
# Find node name (e.g. dev-cluster-worker or cka-cluster-1-worker)
kubectl get nodes

kubectl label node <worker-node-name> disktype=ssd
```

### Step 2: Create the Deployment Manifest
Save to `ssd-app.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ssd-app
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ssd-app
  template:
    metadata:
      labels:
        app: ssd-app
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: disktype
                operator: In
                values:
                - ssd
      containers:
      - name: nginx
        image: nginx:1.25-alpine
```

### Step 3: Apply and Verify Scheduling
```bash
kubectl apply -f ssd-app.yaml
kubectl get pods -o wide -l app=ssd-app
```
Both pods must be placed on the node that has `disktype=ssd`.
