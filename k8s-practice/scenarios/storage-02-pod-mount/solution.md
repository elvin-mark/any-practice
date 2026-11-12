# Solution: Mount Persistent Storage into a Pod

### Step 1: Create the Pod Manifest
Save to `storage-pod.yaml`:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: storage-pod
  namespace: default
spec:
  volumes:
  - name: my-storage
    persistentVolumeClaim:
      claimName: storage-claim
  containers:
  - name: busybox
    image: busybox:1.36
    command: ["sh", "-c", "echo 'cka-practice-success' > /opt/store/status.txt; sleep 3600"]
    volumeMounts:
    - name: my-storage
      mountPath: /opt/store
```

### Step 2: Apply the Manifest
```bash
kubectl apply -f storage-pod.yaml
```

### Step 3: Verify the File
```bash
kubectl wait --for=condition=ready pod/storage-pod --timeout=30s
kubectl exec storage-pod -- cat /opt/store/status.txt
# Output: cka-practice-success
```
