# Solution: Create PersistentVolume and PersistentVolumeClaim

### Step 1: Create the PersistentVolume
Save to `pv.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: data-pv
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: manual
  hostPath:
    path: /mnt/data
```
Apply:
```bash
kubectl apply -f pv.yaml
```

### Step 2: Create the PersistentVolumeClaim
Save to `pvc.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: manual
  resources:
    requests:
      storage: 1Gi
```
Apply:
```bash
kubectl apply -f pvc.yaml
```

### Step 3: Check Status
```bash
kubectl get pv data-pv
kubectl get pvc data-pvc
```
Verify that both report `STATUS: Bound`.
