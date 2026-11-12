# Solution: Deploy a Multi-Container Pod with Shared Volume

### Step 1: Create the Pod Manifest
Save to `multi-app.yaml`:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-app
  namespace: default
spec:
  volumes:
  - name: shared-logs
    emptyDir: {}
  containers:
  - name: producer
    image: busybox:1.36
    command: ["sh", "-c", "while true; do date >> /var/log/app.log; sleep 2; done"]
    volumeMounts:
    - name: shared-logs
      mountPath: /var/log
  - name: consumer
    image: busybox:1.36
    command: ["sh", "-c", "tail -n+1 -F /var/log/app.log"]
    volumeMounts:
    - name: shared-logs
      mountPath: /var/log
```

### Step 2: Apply the Manifest
```bash
kubectl apply -f multi-app.yaml
```

### Step 3: Verify the Running Pod & Sidecar Logs
```bash
kubectl get pod multi-app
kubectl logs multi-app -c consumer
```
