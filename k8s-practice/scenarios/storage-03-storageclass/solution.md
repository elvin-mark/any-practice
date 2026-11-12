# Solution: Create a Custom StorageClass

### Step 1: Create the StorageClass YAML
Save to `fast-retain.yaml`:
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-retain
provisioner: rancher.io/local-path
reclaimPolicy: Retain
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
```

### Step 2: Apply the StorageClass
```bash
kubectl apply -f fast-retain.yaml
```

### Step 3: Verify StorageClass
```bash
kubectl get sc fast-retain
kubectl describe sc fast-retain
```
