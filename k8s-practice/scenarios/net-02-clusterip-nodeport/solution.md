# Solution: Expose Deployment via NodePort Service

### Method 1: Using YAML Manifest
Save to `front-svc.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: front-web-svc
  namespace: web-apps
spec:
  type: NodePort
  selector:
    app: front-web
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080
    protocol: TCP
```
Apply:
```bash
kubectl apply -f front-svc.yaml
```

### Method 2: Imperative Command with Edit
```bash
kubectl expose deployment front-web \
  --name=front-web-svc \
  --namespace=web-apps \
  --type=NodePort \
  --port=80 \
  --target-port=80 \
  --dry-run=client -o yaml > front-svc.yaml
```
Edit `front-svc.yaml` to add `nodePort: 30080` under `ports:`, then run `kubectl apply -f front-svc.yaml`.

### Verification
```bash
kubectl get svc front-web-svc -n web-apps
kubectl get ep front-web-svc -n web-apps
```
